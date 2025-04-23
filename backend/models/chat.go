package models

// Message Role 枚举
type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
	System    Role = "system"
)

// 消息类型
type Message struct {
	Role      Role     `json:"role" bson:"role"`           // 消息角色
	Content   string   `json:"content" bson:"content"`     // 消息内容
	Timestamp int64    `json:"timestamp" bson:"timestamp"` // 时间戳
	Emotion   Emotion  `json:"emotion" bson:"emotion"`     // 情绪分析
	Tags      []string `json:"tags" bson:"tags"`           // 语义标签
}

// 情绪分析结果
type Emotion struct {
	Type     string   `json:"type" bson:"type"`         // 情绪类型：happy, sad, angry, etc.
	Score    float64  `json:"score" bson:"score"`       // 情绪强度
	Keywords []string `json:"keywords" bson:"keywords"` // 关键词
}

// 会话
type Session struct {
	ID        string    `json:"id" bson:"_id"`              // 会话ID
	UserID    string    `json:"userId" bson:"userId"`       // 用户ID
	Title     string    `json:"title" bson:"title"`         // 会话标题
	Messages  []Message `json:"messages" bson:"messages"`   // 消息历史
	CreatedAt int64     `json:"createdAt" bson:"createdAt"` // 创建时间
	UpdatedAt int64     `json:"updatedAt" bson:"updatedAt"` // 更新时间
}
