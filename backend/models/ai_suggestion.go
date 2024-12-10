package models

import "time"

type AISuggestion struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patientId"`
	MessageID   string    `json:"messageId"` // 关联的患者消息ID
	Content     string    `json:"content"`   // AI 建议内容
	Timestamp   time.Time `json:"-"`
	UnixTime    int64     `json:"timestamp"`
	PromptUsed  string    `json:"-"` // 用于生成建议的提示词
	ContextUsed string    `json:"-"` // 使用的上下文
	ModelUsed   string    `json:"-"` // 使用的模型
}
