package configs

import (
	"os"
)

type SparkConfig struct {
	AppID           string
	SecretKey       string
	APIKey          string
	Host            string
	WSBaseURL       string
	HTTPBaseURL     string
	HTTPAPIPassword string
	Model           string
}

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GlobalConfig struct {
	MongoURI     string
	DBName       string
	SparkConfig  *SparkConfig
	GithubConfig *GithubConfig
	JWTSecret    string
}

var Config *GlobalConfig

func LoadEnv() error {
	// 默认配置设置
	setEnv("MONGO_URI", "mongodb://localhost:27017")
	setEnv("DB_NAME", "mood_chatbot")
	setEnv("SPARK_APP_ID", "0b45810b")
	setEnv("SPARK_WS_API_SECRET", "ZTlmZTNiMDkwYzVkMzk4YTM0ZWRmMTI1")
	setEnv("SPARK_WS_API_KEY", "3c0127462241e39ed723627e4509ba53")
	setEnv("SPARK_HOST", "spark-api.xf-yun.com")
	setEnv("SPARK_WS_BASE_URL", "wss://spark-api.xf-yun.com/v1/x1")
	setEnv("SPARK_HTTP_BASE_URL", "https://spark-api-open.xf-yun.com/v2/chat/completions")
	setEnv("SPARK_HTTP_API_Password", "GHCvBwIdiDKsfauuRjmQ:QCCGZRqvWQdBMIhAeeNx")
	setEnv("SPARK_MODEL", "spark-x1")
	setEnv("GITHUB_CLIENT_ID", "Ov23ct6h8TUBSbpkSsLu")
	setEnv("GITHUB_CLIENT_SECRET", "1d0c04cae758cee3db0cd47134fcdbdff411294e")
	setEnv("GITHUB_REDIRECT_URL", "http://localhost:8081/api/callback")
	setEnv("JWT_SECRET", "mood_chatbot_jwt_secret")

	Config = &GlobalConfig{
		MongoURI: os.Getenv("MONGO_URI"),
		DBName:   os.Getenv("DB_NAME"),
		SparkConfig: &SparkConfig{
			AppID:           os.Getenv("SPARK_APP_ID"),
			SecretKey:       os.Getenv("SPARK_WS_API_SECRET"),
			APIKey:          os.Getenv("SPARK_WS_API_KEY"),
			Host:            os.Getenv("SPARK_HOST"),
			WSBaseURL:       os.Getenv("SPARK_WS_BASE_URL"),
			HTTPBaseURL:     os.Getenv("SPARK_HTTP_BASE_URL"),
			HTTPAPIPassword: os.Getenv("SPARK_HTTP_API_Password"),
			Model:           os.Getenv("SPARK_MODEL"),
		},
		GithubConfig: &GithubConfig{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		},
		JWTSecret: os.Getenv("JWT_SECRET"),
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
