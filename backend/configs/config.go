package configs

import (
	"os"
)

func LoadEnv() {
	// MongoDB配置
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("DB_NAME", "mood_chatbot")

	// 讯飞星火配置
	os.Setenv("SPARK_APP_ID", "0b45810b")
	os.Setenv("SPARK_API_SECRET", "ZTlmZTNiMDkwYzVkMzk4YTM0ZWRmMTI1")
	os.Setenv("SPARK_API_KEY", "3c0127462241e39ed723627e4509ba53")
	os.Setenv("SPARK_HOST", "spark-api.xf-yun.com")
	os.Setenv("SPARK_BASE_URL", "wss://spark-api.xf-yun.com/v1/x1")
	os.Setenv("SPARK_MODEL", "spark-x1")
}

func GetEnv(key string) string {
	return os.Getenv(key)
}