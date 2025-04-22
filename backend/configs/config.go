package configs

import "os"

func LoadEnv() {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	os.Setenv("DB_NAME", "mood_chatbot")
}