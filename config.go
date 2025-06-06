package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type LarkApp struct {
	LarkID      string
	LarkSecret  string
	LarkBaseURL string
	LarkWebURL  string
}

// GitLabConfig GitLab相关配置
type GitLabConfig struct {
	WebhookSecret string // GitLab webhook secret token
	BaseURL       string // GitLab实例的基础URL
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *LarkApp {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	return &LarkApp{
		LarkID:      os.Getenv("LARK_APP_ID"),
		LarkSecret:  os.Getenv("LARK_APP_SECRET"),
		LarkBaseURL: getEnvOrDefault("LARK_BASE_URL", "https://open.feishu.cn"),
		LarkWebURL:  getEnvOrDefault("LARK_WEB_URL", ""),
	}
}

// LoadGitLabConfig 从环境变量加载GitLab配置
func LoadGitLabConfig() *GitLabConfig {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	return &GitLabConfig{
		WebhookSecret: os.Getenv("GITLAB_WEBHOOK_SECRET"),
		BaseURL:       getEnvOrDefault("GITLAB_BASE_URL", "https://gitlab.com"),
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
