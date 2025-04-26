package router

import (
	"chatbot-server/internal/controller"
	"chatbot-server/internal/middleware"
	"chatbot-server/internal/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(cc *controller.ChatController) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	api := r.Group("/api")

	{
		// 聊天相关
		api.POST("/create-empty-session", cc.CreateSession)
		api.POST("/create-session-with-message", cc.CreateSessionWithMessage)
		api.GET("/get-session-list", cc.GetSessionsByUserID)
		api.POST("/messages", cc.ChatAndAnalyze)

		// 登录相关
		api.GET("/login", services.LoginHandler)
		api.GET("/callback", services.CallBackHandler)
		api.GET("/me", services.NewHandler)
	}

	return r
} 