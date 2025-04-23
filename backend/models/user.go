package models

// 用户
type UserProfile struct {
	ID           string         `json:"id" bson:"_id"`
	GithubID     string         `json:"githubId" bson:"githubId"`
	Username     string         `json:"username" bson:"username"`
	Email        string         `json:"email" bson:"email"`
	AvatarURL    string         `json:"avatarUrl" bson:"avatarUrl"`
	EmotionStats []EmotionStats `json:"emotionStats" bson:"emotionStats"` // 情绪统计
	EmotionTags  []EmotionTag   `json:"emotionTags" bson:"emotionTags"`   // 情绪标签
	SessionCount int            `json:"sessionCount" bson:"sessionCount"` // 会话数量
	LastActive   int64          `json:"lastActive" bson:"lastActive"`
}

type EmotionStats struct {
	Emotion   string `json:"emotion"`
	TimeStamp int64  `json:"timeStamp"`
}

type EmotionTag struct {
	Emotion string `json:"emotion"`
	Count   int    `json:"count"`
}
