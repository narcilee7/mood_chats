package controller

import (
	services "chatbot-server/internal/service"
	"fmt"
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
	var req struct {
		UserID string `json:"userId"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := req.UserID

	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	session, err := c.chatService.CreateSession(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

func (c *ChatController) CreateSessionWithMessage(ctx *gin.Context) {
	var req struct {
		UserID string `json:"userId"`
		NewMessage string `json:"newMessage"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := req.UserID
	newMessage := req.NewMessage

	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	if newMessage == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "newMessage is required"})
		return
	}

	session, err := c.chatService.CreateSessionWithMessage(userId, newMessage)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}

// // SendMessage 发送消息
// func (c *ChatController) SendMessage(ctx *gin.Context) {
// 	var req struct {
// 		SessionID string `json:"sessionId" binding:"required"`
// 		Content   string `json:"content" binding:"required"`
// 	}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	message, err := c.chatService.SendMessage(req.SessionID, req.Content)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, message)
// }

// 根据userId获取sessions历史
func (c *ChatController) GetSessionsByUserID(ctx *gin.Context) {
	fmt.Println(ctx)
	userId := ctx.Query("userId")

	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	sessions, err := c.chatService.GetSessionsByUserID(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, sessions)
}

func (c *ChatController) ChatAndAnalyze(ctx *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
		Content   string `json:"content" binding:"required"`
		UserID    string `json:"userId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := c.chatService.ChatAndAnalyze(ctx, req.UserID, req.SessionID, req.Content)

	if err != nil {
		fmt.Println("Err: from ChatAndAnalyze")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, message)
}