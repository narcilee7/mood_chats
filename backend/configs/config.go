package configs

import (
	"fmt"
	"os"
)

var requiredEnvVars = []string{
	"SPARK_APP_ID",
	"SPARK_WS_API_SECRET",
	"SPARK_WS_API_KEY",
	"SPARK_HOST",
	"SPARK_WS_BASE_URL",
	"SPARK_HTTP_BASE_URL",
	"SPARK_HTTP_API_Password",
	"SPARK_MODEL",
	"MONGO_URI",
	"DB_NAME",
}

func LoadEnv() error {
	// MongoDB配置
	setEnv("MONGO_URI", "mongodb://localhost:27017")
	setEnv("DB_NAME", "mood_chatbot")

	// 讯飞星火配置
	setEnv("SPARK_APP_ID", "0b45810b")
	setEnv("SPARK_WS_API_SECRET", "ZTlmZTNiMDkwYzVkMzk4YTM0ZWRmMTI1")
	setEnv("SPARK_WS_API_KEY", "3c0127462241e39ed723627e4509ba53")
	setEnv("SPARK_HOST", "spark-api.xf-yun.com")
	setEnv("SPARK_WS_BASE_URL", "wss://spark-api.xf-yun.com/v1/x1")
	setEnv("SPARK_HTTP_BASE_URL", "https://spark-api-open.xf-yun.com/v2/chat/completions")
	setEnv("SPARK_HTTP_API_Password", "GHCvBwIdiDKsfauuRjmQ:QCCGZRqvWQdBMIhAeeNx")
	setEnv("SPARK_MODEL", "spark-x1")

	// 验证所有必需的环境变量
	for _, envVar := range requiredEnvVars {
		if value := os.Getenv(envVar); value == "" {
			return fmt.Errorf("环境变量 %s 未设置", envVar)
		}
	}

	return nil
}

func setEnv(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
