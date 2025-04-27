package respository

import (
	models "chatbot-server/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	// 创建用户
	CreateUser(userProfile *models.UserProfile) error
}


func CreateUser(userProfile *models.UserProfile) error {
	*userProfile = models.UserProfile{
		ID: primitive.NewObjectID().Hex(),
		EmotionTags:  []models.Emotion{},
		SessionCount: 0,
		LastActive:  0,
		GithubConfig: models.GithubConf{
			Username: userProfile.GithubConfig.Username,
			Email: userProfile.GithubConfig.Email,
			Name: userProfile.GithubConfig.Name,
			NickName: userProfile.GithubConfig.NickName,
			Description: userProfile.GithubConfig.Description,
			AvatarURL: userProfile.GithubConfig.AvatarURL,
			Location: userProfile.GithubConfig.Location,
			Bio: userProfile.GithubConfig.Bio,
			Company: userProfile.GithubConfig.Company,
			Blog: userProfile.GithubConfig.Blog,
		},
	}
	return nil
}