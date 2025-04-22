package models

type Message struct {
	Role      string `bson: "role"` // "user" or "ai"
	Content   string `bson: "content"`
	TimeStamp int64  `bson: "timestamp"`
}

type Session struct {
	UserId    string    `bson: "userId"`
	SessionId string    `bson: "sessionId"`
	History   []Message `bson: "history"`
}