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
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
