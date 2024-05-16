package service

import (
	"testing"
)

func TestClaude3Invok(t *testing.T) {
	// Define the expected input and output
	prompt := `
	Instruct: 你是AWS IT专家,你会根据我提供的context和condition进行判断,并且根据condition标注具体的指标,比如EC2实例id, S3 bucketname
	以中文和JSON格式返回, 返回的例子:
	result:
	   pass : true 
	   reason: |
		   your answer

	context: 
	这张图片显示了三个EC2实例在一段时间内的CPU利用率变化趋势。图中的X轴代表时间,从凌晨4:30到上午10:00,Y轴代表CPU利用率的百分比,最大为100%。有三条不同颜色的折线分别代表三个不同实例的CPU利用率变化,其中一条实例(i-0222ef85c03b5bb55)的平均CPU利用率约为4.04%,另外两个实例的平均利用率分别为0.594%和0.578%。整体来看,这三个实例的CPU利用率都比较低,大部分时间在10%以下,只有在凌晨6:30左右有一个短暂的利用率峰值。
	
	condition:
	如果EC2实例CPU最大用量超过1%则不通过
	`

	// Call the function under test
	output := Claude3Invok(prompt)

	// Assert the output
	if output == nil {
		t.Errorf("Expected non-nil output, got nil")
	}
	if output != nil {
		t.Logf("%+v", *output)
	}

}

func TestClaude3InvokWithImage(t *testing.T) {
	// Define the expected input and output
	prompt := "分析一下CPU,用中文解释."
	image := "../test.png"

	// Call the function under test
	output := Claude3InvokWithImage(prompt, image, false)

	// Assert the output
	if output == nil {
		t.Errorf("Expected non-nil output, got nil")
	}
	if output != nil {
		t.Logf("%+v", *output)
	}

}

func TestClaude3InvokWithImageBase64(t *testing.T) {
	// Define the expected input and output
	prompt := `
	根据"分析当前图片中EC2的CPU使用量"提取图片信息.
	使用中文和JSON格式返回,返回的例子
	{"text":###text exctract from image"}
	`
	image := "../test.png"
	base64Content, err := getBase64EncodedImage(image)
	if err != nil {
		t.Fatal(err)
	}

	// Call the function under test
	output := Claude3InvokWithImage(prompt, base64Content, true)

	// Assert the output
	if output == nil {
		t.Errorf("Expected non-nil output, got nil")
	}
	if output != nil {
		t.Logf("%+v", *output)
	}

}

func TestRunProjectSummary(t *testing.T) {

	InitDBConnection()

	// Retrieve the item
	items, err := GetMetricsByProject("Project1")
	if err != nil {
		t.Errorf("Failed to retrieve monitor item: %v", err)
	}

	output, err := RunProjectSummary(items)
	if err != nil {
		t.Errorf("Failed to retrieve monitor item: %v", err)
	}

	t.Log(output)

}
func TestRenderTemplate(t *testing.T) {
	data := map[string]interface{}{
		"context":   "监控每个实例的CPU",
		"condition": "如果CPU大于1%则不通过",
	}
	templateName := "myTemplate"

	prompt, err := RenderTemplate(data, templateName, runConditionPromptTemplate)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(prompt)

	prompt, err = RenderTemplate(data, templateName, imageTextExtracPromptTemplate)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(prompt)

}

