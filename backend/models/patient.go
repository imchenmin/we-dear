package models

import "time"

type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	Diagnosis string    `json:"diagnosis"`
	Doctor    string    `json:"doctor"`
	Avatar    string    `json:"avatar"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Type        string    `json:"type"` // "patient", "doctor", "ai_suggestion"
	ContentType string    `json:"contentType"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	ReplyTo     string    `json:"replyTo,omitempty"`
}
