package services

import (
	"chatbot-server/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// ChatService 定义聊天服务接口
type ChatService interface {
	// 创建新会话
	CreateSession(userID string) (*models.Session, error)
	// 发送消息
	SendMessage(sessionID string, content string) (*models.Message, error)
	// 获取会话历史
	GetSessionHistory(sessionID string) ([]models.Message, error)
	// 分析情绪
	AnalyzeEmotion(content string) (*models.Emotion, error)
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
func (s *ChatServiceImpl) CreateSession(userID string) (*models.Session, error) {
	session := &models.Session{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		Title:     "新会话",
		Messages:  []models.Message{},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	_, err := s.db.Collection("sessions").InsertOne(context.Background(), session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

/*
发送消息
*/
func (s *ChatServiceImpl) SendMessage(sessionID string, content string) (*models.Message, error) {
	// 1. 分析情绪
	emotion, err := s.AnalyzeEmotion(content)
	if err != nil {
		return nil, err
	}

	// 2. 生成AI回复
	history := []models.Message{}
	_, err = s.ai.Chat(content, history)

	if err != nil {
		return nil, err
	}

	// 3. 保存消息
	message := &models.Message{
		Role:      "user",
		Content:   content,
		Timestamp: time.Now().Unix(),
		Emotion:   *emotion,
	}

	// 4. 更新会话
	update := bson.M{
		"$push": bson.M{"messages": message},
		"$set":  bson.M{"updatedAt": time.Now().Unix()},
	}

	_, err = s.db.Collection("sessions").UpdateOne(
		context.Background(),
		bson.M{"_id": sessionID},
		update,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
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

/*
实现情感分析
*/
func (s *ChatServiceImpl) AnalyzeEmotion(content string) (*models.Emotion, error) {
	// TODO: 实现情绪分析
	// 这里可以集成第三方情绪分析API
	return &models.Emotion{
		Type:     "neutral",
		Score:    0.5,
		Keywords: []string{},
	}, nil
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
