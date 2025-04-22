package main

import (
	"chatbot-server/configs"
	"chatbot-server/database"
	"chatbot-server/routes"
)

func main() {
	// 加载环境变量
	configs.LoadEnv()
	database.ConnectDB()
	r := routes.SetupRouter()

	r.Run(":8080")
}