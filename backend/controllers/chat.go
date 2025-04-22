package controllers

import (
	"chatbot-server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService services.ChatService
}

func NewChatController(chatService services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

// CreateSession 创建新会话
func (c *ChatController) CreateSession(ctx *gin.Context) {
	userID := ctx.GetString("userID") // 从中间件获取用户ID
	
	session, err := c.chatService.CreateSession(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

// SendMessage 发送消息
func (c *ChatController) SendMessage(ctx *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := c.chatService.SendMessage(req.SessionID, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, message)
}

// GetHistory 获取会话历史
func (c *ChatController) GetHistory(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")
	
	messages, err := c.chatService.GetSessionHistory(sessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}