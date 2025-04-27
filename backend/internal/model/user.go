package models

// 用户
type UserProfile struct {
	ID           string    `json:"id" bson:"_id"`
	EmotionTags  []Emotion `json:"emotionTags" bson:"emotionTags"`   // 情绪标签
	SessionCount int       `json:"sessionCount" bson:"sessionCount"` // 会话数量
	LastActive   int64     `json:"lastActive" bson:"lastActive"`
	GithubConfig GithubConf  `json:"githubConfig" bson:"githubConfig"`
}

type GithubConf struct {
  Username string
  Email             string
  Name              string
  NickName          string
  Description       string
  AvatarURL         string
  Location          string
  Bio              string
  Company string
  Blog string
}