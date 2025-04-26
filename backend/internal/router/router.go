package router

import (
	"chatbot-server/internal/controller"
	"chatbot-server/internal/middleware"
	services "chatbot-server/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(chatController *controller.ChatController) *gin.Engine {
	r := gin.Default()

	// 跨域设置
	r.Use(middleware.CORS())

	// 聊天相关路由
	chat := r.Group("/api")
	{
		chat.POST("/create-empty-session", chatController.CreateSession)
		chat.POST("/create-session-with-message", chatController.CreateSessionWithMessage)
		chat.GET("/get-session-list", chatController.GetSessionsByUserID)
		chat.POST("/messages", chatController.ChatAndAnalyze)
	}

	// 登录相关路由
	auth := r.Group("/api")
	{
		auth.GET("/login", services.LoginHandler)
		auth.GET("/callback", services.CallBackHandler)
		auth.GET("/me", services.NewHandler)
	}

	return r
}