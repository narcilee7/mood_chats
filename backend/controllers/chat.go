package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{ "message": "Hello World" })
}