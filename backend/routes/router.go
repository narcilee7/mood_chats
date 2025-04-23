package routes

import (
	"chatbot-server/handlers"
	"chatbot-server/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(chatController *handlers.ChatController) *gin.Engine {
	r := gin.Default()

	// 聊天相关路由
	chat := r.Group("/api")
	{
		chat.POST("/sessions", chatController.CreateSession)
		chat.GET("/sessions/history", chatController.GetSessionsByUserID)
		chat.POST("/messages", chatController.ChatAndAnalyze)
	}

	// 登录相关路由
	auth := r.Group("/api")
	{
		auth.GET("/login", services.LoginHandler)
		auth.GET("/callback", services.CallBackHandler)
	}

	return r
}