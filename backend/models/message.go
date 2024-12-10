package models

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"-"`         // 内部使用 time.Time
	UnixTime  int64     `json:"timestamp"` // 对外使用 Unix 时间戳
	Role      string    `json:"role"`      // "doctor" 或 "patient"
	Sender    string    `json:"sender"`    // 发送者名称
	Avatar    string    `json:"avatar,omitempty"`
}

// 在序列化之前设置 UnixTime
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias:     (*Alias)(m),
		Timestamp: m.Timestamp.UnixMilli(),
	})
}

// 在反序列化时处理时间戳
func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		*Alias
		Timestamp int64 `json:"timestamp"`
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	m.Timestamp = time.Unix(0, aux.Timestamp*int64(time.Millisecond))
	m.UnixTime = aux.Timestamp
	return nil
}
