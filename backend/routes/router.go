package routes

import (
	"chatbot-server/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(chatController *handlers.ChatController) *gin.Engine {
	r := gin.Default()

	// 添加中间件
	//r.Use(middleware.CORS())
	//r.Use(middleware.Auth())

	api := r.Group("/api")
	{
		/*会话相关*/
		// 构建新会话
		api.POST("/sessions", chatController.CreateSession)

		// 根据userID获取会话列表
		api.GET("/sessions/history", chatController.GetSessionsByUserID)

		// /*消息相关*/
		// api.POST("/messages", chatController.SendMessage)
		// 发消息
		api.POST("/messages", chatController.ChatAndAnalyze)
	}

	return r
}
