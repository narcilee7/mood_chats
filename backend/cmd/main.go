package main

import (
	"chatbot-server/configs"
	"chatbot-server/internal/controller"
	"chatbot-server/internal/database"
	"chatbot-server/internal/router"
	services "chatbot-server/internal/service"
	"chatbot-server/pkg/logger"
	"log"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化配置
	if err := configs.LoadEnv(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	zap.L().Info("配置加载成功")

	// 连接数据库
	database.ConnectDB()
	zap.L().Info("数据库连接成功")

	// 初始化服务
	sparkProvider := services.NewSparkProvider(
		configs.GetEnv("SPARK_APP_ID"),
		configs.GetEnv("SPARK_WS_API_SECRET"),
		configs.GetEnv("SPARK_WS_API_KEY"),
		configs.GetEnv("SPARK_HOST"),
		configs.GetEnv("SPARK_WS_BASE_URL"),
		configs.GetEnv("SPARK_HTTP_BASE_URL"),
		configs.GetEnv("SPARK_HTTP_API_Password"),
	)

	chatService := services.NewChatService(database.DB, sparkProvider)
	chatController := controller.NewChatController(chatService)

	// 设置路由
	r := router.SetupRouter(chatController)

	// 启动服务器
	zap.L().Info("服务启动成功")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
