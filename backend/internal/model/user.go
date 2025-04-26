package models

import "time"

// 用户
type UserProfile struct {
	ID           string    `json:"id" bson:"_id"`
	EmotionTags  []Emotion `json:"emotionTags" bson:"emotionTags"`   // 情绪标签
	SessionCount int       `json:"sessionCount" bson:"sessionCount"` // 会话数量
	LastActive   int64     `json:"lastActive" bson:"lastActive"`
	GithubConfig GithubConf  `json:"githubConfig" bson:"githubConfig"`
}

type GithubConf struct {
	ClientId string
	ClientSecret string
	RedirectUrl string 
	RawData           map[string]interface{}
  Provider          string
  Email             string
  Name              string
  FirstName         string
  LastName          string
  NickName          string
  Description       string
  UserID            string
  AvatarURL         string
  Location          string
  AccessToken       string
  AccessTokenSecret string
  RefreshToken      string
  ExpiresAt         time.Time
  IDToken           string
}

