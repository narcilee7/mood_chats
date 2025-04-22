package routes

import (
	"chatbot-server/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(chatController *controllers.ChatController) *gin.Engine {
	r := gin.Default()

	// 添加中间件
	//r.Use(middleware.CORS())
	//r.Use(middleware.Auth())

	api := r.Group("/api")
	{
		// 会话相关
		api.POST("/sessions", chatController.CreateSession)
		api.GET("/sessions/:sessionId/history", chatController.GetHistory)

		// 消息相关
		api.POST("/messages", chatController.SendMessage)
	}

	return r
}
