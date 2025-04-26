package app

import (
	"chatbot-server/configs"
	"chatbot-server/internal/controller"
	"chatbot-server/internal/database"
	"chatbot-server/internal/router"
	"chatbot-server/internal/services"
	"chatbot-server/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Application struct {
	engine *gin.Engine	
}

func InitializeApp() (*Application, error) {
	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		return nil, fmt.Errorf("初始化日志失败: %w", err)
	}

	// 加载环境配置
	if err := configs.LoadEnv(); err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	zap.L().Info("配置加载成功")

	// 连接数据库
	db, err := database.ConnectDB()
	
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	zap.L().Info("数据库连接成功")

	// 初始化第三方服务
	sparkProvider := services.NewSparkProvider(
		configs.GetEnv("SPARK_APP_ID"),
		configs.GetEnv("SPARK_WS_API_SECRET"),
		configs.GetEnv("SPARK_WS_API_KEY"),
		configs.GetEnv("SPARK_HOST"),
		configs.GetEnv("SPARK_WS_BASE_URL"),
		configs.GetEnv("SPARK_HTTP_BASE_URL"),
		configs.GetEnv("SPARK_HTTP_API_PASSWORD"), // 注意你那边一个拼错了 Password
	)

	// 初始化业务服务
	chatService := services.NewChatService(db, sparkProvider)
	chatController := controller.NewChatController(chatService)

	// 初始化路由
	r := router.NewRouter(chatController)

	return &Application{
		engine: r,
	}, nil
}

func (a *Application) Run() error {
	zap.L().Info("服务启动成功", zap.String("port", ":8081"))
	return a.engine.Run(":8081")
}
