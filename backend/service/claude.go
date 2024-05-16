package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"gopkg.in/yaml.v3"
)

const defaultRegion = "us-east-1"
const defaultModelID = "anthropic.claude-3-sonnet-20240229-v1:0"

var brc *bedrockruntime.Client

var (
	runConditionPromptTemplate string = `
	Instruct: 你是AWS IT专家,你会根据我提供的context和condition进行判断,并且根据condition标注具体的指标,比如EC2实例id, S3 bucketname
	以中文和YAML格式返回, 返回yaml的例子,(不要包含任何markdown标记,纯文本):
	result:
		pass : true 
		reason: |
			your answer

	context: 
	{{.context}}

	condition:
	{{.condition}}

	`

	imageTextExtracPromptTemplate string = `
	Instruct:  你是AWS IT专家,你会根据我提供的context来提取图片里面的信息,如果是EC2, S3, ALB请带上具体的信息,比如实例id,ARN
	
	使用中文和JSON格式返回,返回的例子
	{"text":###text exctract from image"}

	context: 
	{{.context}}
	`

	projectSummaryPromptTemplate string = `
	Instruct: you are senior SRE , I give some metrics , you need analtyics and summary IT infrastructure issue or other thing
you need provide detail information in summary and you need return use json and need include all raw metrics


	health level:
	Very Good , all metrics pass
	Good,  95% mestrics pass 
	Not Good 90% mestrics pass
	Bad , less 90% pass 
	Very Bad less 80% pass 

	metrics example:

	[{
            "id": 1,
            "project": "Project1",
            "catalog": "EC2",
            "item_desc": "EC2 实例的CPU用量",
            "item_condition": "如果CPU大于45%则不通过",
            "dashboard_url": "http://localhost:3100/d/tmsOtSxZk/amazon-ec2?orgId=1\u0026viewPanel=2",
            "status": false,
            "status_desc": "根据提供的上下文信息,其中列出了三个 EC2 实例的 ID 及其最大 CPU 利用率。\n根据条件\"如果 CPU 大于 45% 则不通过\"的要求,其中一个实例 i-0cf8b09fb079046e3 的最大 CPU 利用率为 47.5%,超过了 45% 的阈值,因此不通过验证。",
            "screen": "static/9033d403-8c12-4d05-9017-4d9dbec3c29c.png",
            "check_date": "2024-04-25T21:33:17+08:00"
		},
		{
            "id": 3,
            "project": "Project1",
            "catalog": "EC2",
            "item_desc": "出口网络流量",
            "item_condition": "网络流量不大于2k",
            "dashboard_url": "http://localhost:3100/d/tmsOtSxZk/amazon-ec2?orgId=1\u0026viewPanel=18",
            "status": false,
            "status_desc": "根据提供的上下文信息,有一个 EC2 实例 ID i-0cf8b09fb079046e3 在 14:00 左右出现了明显的流量峰值。由于条件是网络流量不大于 2k,因此该实例在该时间段内的网络流量很可能超过了 2k。因此,无法通过条件检查。\n\n具体检查项目:\n  - EC2 实例 ID: \n    - i-0222ef85c03b5b8c5\n    - i-0333da98ab126702b\n    - i-0cf8b09fb079046e3",
            "screen": "static/43ee21f9-1b3e-4b4d-9d1e-e4d50c9fca4d.png",
            "check_date": "2024-04-24T14:44:10+08:00"
			}]

	input metrics
	{{.metrics}}

	output is json , not include any markdown code, your summary, use '\"' replace, not use '"', must be correct json:
	
	{
		"health":"health_level",
		"summary": ##your summary##,
		"metrics":"raw_mestics"
	}
	`
)

type CheckResult struct {
	Result struct {
		Pass   bool   `json:"pass" yaml:"pass"`
		Reason string `json:"reason" yaml:"reason"`
	} `yaml:"result"`
}

// Claude3Request represents the request payload for Claude 3 model
type Claude3Request struct {
	Messages      []Message `json:"messages"`
	MaxTokens     int       `json:"max_tokens"`
	Temperature   float64   `json:"temperature,omitempty"`
	TopP          float64   `json:"top_p,omitempty"`
	TopK          int       `json:"top_k,omitempty"`
	StopSequences []string  `json:"stop_sequences,omitempty"`
	Version       string    `json:"anthropic_version"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// InvokeModelResponse represents the response from the Claude 3 model
type InvokeModelResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func RenderTemplate(data map[string]interface{}, templateName string, templateStr string) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RunItemByID(id int64) (*MonitorMetric, error) {

	exitsItem, err := GetMonitorMetric(id)
	if err != nil {
		return nil, err
	}

	exitsConnection, err := GetMonitorConnectionByName(exitsItem.ConnectionName)
	if err != nil {
		return nil, err
	}

	fmt.Println(exitsConnection)
	result, err := RunCondition(&exitsItem,&exitsConnection)
	if err != nil {

		return nil, err
	}

	var checkResult CheckResult

	err = yaml.Unmarshal([]byte(result[1]), &checkResult)
	if err != nil {
		return nil, err
	}

	exitsItem.Status = checkResult.Result.Pass
	exitsItem.StatusDesc = checkResult.Result.Reason
	exitsItem.CheckDate = time.Now()
	exitsItem.Screen = result[0]

	err = UpdateMonitorMetric(exitsItem)
	if err != nil {
		return nil, err
	}

	return &exitsItem, nil
}

func BatchRunItems(ids []int64) ([]MonitorMetric, error) {
	var results []MonitorMetric

	for _, id := range ids {
		item, err := RunItemByID(id)
		if err != nil {
			return nil, err
		}
		results = append(results, *item)
	}

	return results, nil
}

func RunProjectSummary(items []MonitorMetric) (string, error) {

	itemsString, err := json.Marshal(items)
	if err != nil {
		return "", nil
	}

	data := map[string]interface{}{
		"metrics": string(itemsString),
	}

	//step1 extract information from image
	prompt, err := RenderTemplate(data, "txtExtrac", projectSummaryPromptTemplate)
	if err != nil {
		return "", err
	}

	output := Claude3Invok(prompt)
	if output == nil {
		return "", err
	}

	logger.Println(output)

	text := output.Content[0].Text
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")

	return text, nil

}

func RunCondition(item *MonitorMetric,connection *MonitorConnection) ([]string, error) {

	result, err := ScreenCaptureTasks([]string{
		item.DashboardURL,
	},connection)

	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"context": item.ItemDesc,
	}

	//step1 extract information from image
	prompt, err := RenderTemplate(data, "txtExtrac", imageTextExtracPromptTemplate)
	if err != nil {
		return nil, err
	}

	image := fmt.Sprintf("./%s", result[item.DashboardURL])
	base64Content, err := getBase64EncodedImage(image)
	if err != nil {
		return nil, err
	}

	output := Claude3InvokWithImage(prompt, base64Content, true)
	if output == nil {
		return nil, err
	}

	context := output.Content[0].Text

	//step2 run condition
	data = map[string]interface{}{
		"context":   context,
		"condition": item.ItemCondition,
	}

	prompt, err = RenderTemplate(data, "runCondition", runConditionPromptTemplate)
	if err != nil {
		return nil, err
	}

	output = Claude3Invok(prompt)
	if output == nil {
		return nil, err
	}

	fmt.Println(output)

	return []string{result[item.DashboardURL], output.Content[0].Text}, nil
}

// initClient initializes the Bedrock Runtime client
func initClient() *bedrockruntime.Client {
	if brc != nil {
		return brc
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	brc = bedrockruntime.NewFromConfig(cfg)
	return brc
}

// getBase64EncodedImage reads an image file and returns its base64 encoded string
func getBase64EncodedImage(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(imageData)
	return base64String, nil
}

func PrintPrettyJSON(input interface{}) {
	prettyJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Println(string(prettyJSON))
}

// Claude3InvokWithImage invokes the Claude 3 model with a prompt and an image
// Only support local file or base64
func Claude3InvokWithImage(prompt string, image string, isBase64 bool) *InvokeModelResponse {
	var base64Image string
	var err error

	client := initClient()

	if isBase64 {
		base64Image = image
	} else {
		base64Image, err = getBase64EncodedImage(image)
		if err != nil {
			logger.Fatal(err)
			return nil
		}
	}

	imageContent := []interface{}{
		map[string]interface{}{
			"type": "image",
			"source": map[string]interface{}{
				"type":       "base64",
				"media_type": "image/jpeg",
				"data":       base64Image,
			},
		},
		map[string]interface{}{
			"type": "text",
			"text": prompt,
		},
	}

	payload := Claude3Request{Messages: []Message{{Role: "user", Content: imageContent}}, MaxTokens: 2048, Version: "bedrock-2023-05-31", Temperature: 0.0}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	output, err := client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(defaultModelID),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		logger.Println("Error invoking model:", err)
		return nil
	}

	PrintPrettyJSON(output.Body)

	var response InvokeModelResponse
	err = json.Unmarshal(output.Body, &response)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	return &response
}

// Claude3Invok invokes the Claude 3 model with a text prompt
func Claude3Invok(prompt string) *InvokeModelResponse {
	client := initClient()

	payload := Claude3Request{Messages: []Message{{Role: "user", Content: prompt}}, MaxTokens: 2048, Version: "bedrock-2023-05-31"}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	output, err := client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String(defaultModelID),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		logger.Fatal(err)
		return nil
	}

	var response InvokeModelResponse
	err = json.Unmarshal(output.Body, &response)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	return &response
}
