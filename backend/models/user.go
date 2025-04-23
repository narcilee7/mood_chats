package models

// 用户
type UserProfile struct {
	ID           string    `json:"id" bson:"_id"`
	GithubID     string    `json:"githubId" bson:"githubId"`
	Username     string    `json:"username" bson:"username"`
	Email        string    `json:"email" bson:"email"`
	AvatarURL    string    `json:"avatarUrl" bson:"avatarUrl"`
	EmotionTags  []Emotion `json:"emotionTags" bson:"emotionTags"`   // 情绪标签
	SessionCount int       `json:"sessionCount" bson:"sessionCount"` // 会话数量
	LastActive   int64     `json:"lastActive" bson:"lastActive"`
}
