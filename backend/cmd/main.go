package main

import (
	"chatbot-server/configs"
	"chatbot-server/database"
	"chatbot-server/handlers"
	"chatbot-server/routes"
	"chatbot-server/services"
	"go.uber.org/zap"
	"log"
)

func main() {
	// 加载配置
	configs.LoadEnv()
	// 初始化日志
	zap.L().Info("Logger initialized")

	// 连接数据库
	database.ConnectDB()
	zap.L().Info("Database initialized")

	// 初始化服务
	sparkProvider := services.NewSparkProvider(
		configs.GetEnv("SPARK_APP_ID"),
		configs.GetEnv("SPARK_API_SECRET"),
		configs.GetEnv("SPARK_API_KEY"),
		configs.GetEnv("SPARK_HOST"),
		configs.GetEnv("SPARK_BASE_URL"),
		configs.GetEnv("SPARK_MODEL"),
	)

	chatService := services.NewChatService(database.DB, sparkProvider)
	chatController := controllers.NewChatController(chatService)

	// 设置路由
	r := routes.SetupRouter(chatController)

	// 启动服务器
	log.Println("Server starting on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
