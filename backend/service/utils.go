package service

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"net/url"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	Headless        bool
	AdminPassword   string
	Front string
	ChromeDP string
}

var logger = log.New(os.Stdout, "service|", log.LstdFlags)
var globalConfig *Config

func GetGlobalConfig() *Config {
	return globalConfig
}

func GetEnv(key, defalut string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defalut
}

func LoadConfig() (*Config, error) {
	var cfg Config

	// Load configuration from command-line flags
	flag.StringVar(&cfg.DBHost, "dbHost", "localhost", "Database host")
	flag.StringVar(&cfg.DBPort, "dbPort", "3306", "Database port")
	flag.StringVar(&cfg.DBName, "dbName", "", "Database name")
	flag.StringVar(&cfg.DBUser, "dbUser", "", "Database user")
	flag.StringVar(&cfg.DBPassword, "dbPassword", "", "Database password")
	flag.StringVar(&cfg.AdminPassword, "AdminPassword", "admin", "Admin password")
	flag.StringVar(&cfg.Front, "front", "http://localhost:3000", "front nextjs url")
	flag.StringVar(&cfg.ChromeDP, "chromeDP", "", "chromium headless port")

	


	flag.Parse()

	// Check if required configuration values are set from command-line flags
	if cfg.DBHost != "" && cfg.DBPort != "" && cfg.DBName != "" && cfg.DBUser != "" {
		globalConfig = &cfg
		return &cfg, nil
	}

	// Load configuration from environment variables
	viper.SetDefault("DBHost", "localhost")
	viper.SetDefault("DBPort", "3306")
	viper.SetDefault("DBName", "mydb")
	viper.SetDefault("DBUser", "admin")
	viper.SetDefault("DBPassword", "admin")
	viper.SetDefault("Front", "http://localhost:3000")
	viper.SetDefault("ChromeDP", "ws://localhost:9222")

	
	viper.SetDefault("AdminPassword", "admin")

	viper.AutomaticEnv()

	// Bind configuration values from viper to struct
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Check if required configuration values are set from environment variables
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.DBUser == "" {
		return nil, fmt.Errorf("required configuration values not set")
	}
	globalConfig = &cfg
	return &cfg, nil
}


func parseDateTime(str string) (time.Time, error) {
	// Adjust the layout string based on your MySQL date/time format
	layout := "2006-01-02 15:04:05"
	return time.Parse(layout, str)
}


func ReplaceHost(connectionURL, urlStr string) string {
	origin, err := url.Parse(connectionURL)
	
	if err != nil {
		return urlStr // return the original URL if it's invalid
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr // return the original URL if it's invalid
	}

	u.Scheme = origin.Scheme
	u.Host = origin.Host

	return u.String()
}