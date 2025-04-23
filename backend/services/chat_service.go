package services

import (
	"chatbot-server/models"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// ChatService 定义聊天服务接口
type ChatService interface {
	// 创建新会话
	CreateSession(userID string) (*models.Session, error)
	// 获取会话列表
	GetSessionsByUserID(userID string) ([]*models.Session, error)
	// 发送消息并收集情绪
	ChatAndAnalyze(ctx context.Context, userID string, sessionID string, message string) (string, error)
	// 获取会话历史
	GetSessionHistory(sessionID string) ([]models.Message, error)
	// 更新用户画像
	UpdateUserProfile(userID string, emotion *models.Emotion) error
}

// ChatServiceImpl 聊天服务实现
type ChatServiceImpl struct {
	db *mongo.Database
	ai SparkProvider
}

func NewChatService(db *mongo.Database, ai *SparkProvider) ChatService {
	return &ChatServiceImpl{
		db: db,
		ai: *ai,
	}
}

/*
创建新会话
*/
func (cs *ChatServiceImpl) CreateSession(userID string) (*models.Session, error) {
	session := &models.Session{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		Title:     "新会话",
		Messages:  []models.Message{},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	_, err := cs.db.Collection("sessions").InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *ChatServiceImpl) GetSessionHistory(sessionID string) ([]models.Message, error) {
	var session models.Session
	err := s.db.Collection("sessions").FindOne(
		context.Background(),
		bson.M{"_id": sessionID},
	).Decode(&session)

	if err != nil {
		return nil, err
	}

	return session.Messages, nil
}

func (s *ChatServiceImpl) UpdateUserProfile(userID string, emotion *models.Emotion) error {
	// 更新用户情绪统计
	update := bson.M{
		"$inc": bson.M{
			"emotionStats." + emotion.Type: 1,
		},
		"$set": bson.M{
			"lastActive": time.Now().Unix(),
		},
	}

	_, err := s.db.Collection("user_profiles").UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		update,
		options.Update().SetUpsert(true),
	)

	return err
}

func (c *ChatServiceImpl) GetSessionsByUserID(UserID string) ([]*models.Session, error) {
	var sessions []*models.Session
	cursor, err := c.db.Collection("sessions").Find(context.Background(), bson.M{"userId": UserID})

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &sessions)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (cs *ChatServiceImpl) ChatAndAnalyze(ctx context.Context, userID string, sessionID string, message string) (string, error) {
	// 获取历史消息
	history, err := cs.GetSessionHistory(sessionID)
	if err != nil {
		zap.L().Error("获取会话历史失败", zap.Error(err))
		return "", err
	}

	// 分析用户情绪
	userEmotion, err := cs.ai.AnalyzeEmotion(userID, message)
	if err != nil {
		zap.L().Error("用户情绪分析失败", zap.Error(err))
		// 即使情绪分析失败，也继续处理对话
	}

	// 调用大模型回复
	answer, err := cs.ai.ChatWithHttp(message, userID, history)
	if err != nil {
		zap.L().Error("大模型回复失败", zap.Error(err))
		return "", err
	}

	// 创建消息对象
	userMessage := models.DBMessage{
		ID:        primitive.NewObjectID().Hex(),
		Role:      models.User,
		Content:   message,
		Timestamp: time.Now().Unix(),
		SessionID: sessionID,
	}

	assistantMessage := models.DBMessage{
		ID:        primitive.NewObjectID().Hex(),
		Role:      models.Assistant,
		Content:   answer,
		Timestamp: time.Now().Unix(),
		SessionID: sessionID,
	}

	// 存储消息
	_, err = cs.db.Collection("messages").InsertOne(ctx, userMessage)
	if err != nil {
		zap.L().Error("存储用户消息失败", zap.Error(err))
		return "", err
	}

	_, err = cs.db.Collection("messages").InsertOne(ctx, assistantMessage)
	if err != nil {
		zap.L().Error("存储助手消息失败", zap.Error(err))
		return "", err
	}

	// 更新会话
	update := bson.M{
		"$push": bson.M{
			"messages": bson.M{
				"$each": []models.Message{
					{Role: models.User, Content: message, Timestamp: userMessage.Timestamp},
					{Role: models.Assistant, Content: answer, Timestamp: assistantMessage.Timestamp},
				},
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now().Unix(),
		},
	}

	_, err = cs.db.Collection("sessions").UpdateOne(
		ctx,
		bson.M{"_id": sessionID},
		update,
	)
	if err != nil {
		zap.L().Error("更新会话失败", zap.Error(err))
		return "", err
	}

	// 更新用户情绪标签（仅在情绪分析成功时）
	if userEmotion != nil {
		fmt.Println("userEmotion:", userEmotion)
		err = cs.UpdateUserProfile(userID, userEmotion)
		if err != nil {
			zap.L().Error("更新用户情绪标签失败", zap.Error(err))
			// 情绪标签更新失败不影响对话流程
		}
	}

	return answer, nil
}
