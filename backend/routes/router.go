package routes

import (
	"chatbot-server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	{
		api.POST("/chat", controllers.ChatHandler)
	}

	return r
}