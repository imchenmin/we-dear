package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"we-dear/models"
)

var (
	wsService *Service
	wsOnce    sync.Once
)

// Service WebSocket服务
type Service struct {
	manager *Manager
}

// GetService 获取WebSocket服务单例
func GetService() *Service {
	wsOnce.Do(func() {
		wsService = &Service{
			manager: NewManager(),
		}
		// 启动WebSocket管理器
		go wsService.manager.Start()
	})
	return wsService
}

// GetManager 获取WebSocket管理器
func (s *Service) GetManager() *Manager {
	return s.manager
}

// NotifyNewMessage 通知新消息
func (s *Service) NotifyNewMessage(message *models.Message) {
	wsMessage := WSMessage{
		Type:   "chat",
		Action: "created",
		Payload: map[string]interface{}{
			"id":        message.ID,
			"content":   message.Content,
			"type":      message.Type,
			"role":      message.Role,
			"patientId": message.PatientID,
			"doctorId":  message.DoctorID,
			"createdAt": message.CreatedAt,
		},
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Failed to marshal message notification: %v", err)
		return
	}

	// 通知相关用户
	if message.Role == "doctor" {
		// 医生发送的消息，通知患者
		s.manager.SendToUser(message.PatientID, data)
	} else {
		// 患者发送的消息，通知医生
		s.manager.SendToUser(message.DoctorID, data)
	}
}

// NotifyNewAISuggestion 通知新的AI建议
func (s *Service) NotifyNewAISuggestion(suggestion *models.AISuggestion) {
	wsMessage := WSMessage{
		Type:   "ai_suggestion",
		Action: "created",
		Payload: map[string]interface{}{
			"id":        suggestion.ID,
			"messageId": suggestion.MessageID,
			"content":   suggestion.Content,
			"category":  suggestion.Category,
			"priority":  suggestion.Priority,
			"createdAt": suggestion.CreatedAt,
		},
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Failed to marshal AI suggestion notification: %v", err)
		return
	}

	// 通知医生有新的AI建议
	s.manager.SendToRole("doctor", data)
}

// NotifyNewPhysiologicalData 通知新的生理数据
func (s *Service) NotifyNewPhysiologicalData(data *models.PhysiologicalData) {
	wsMessage := WSMessage{
		Type:   "physiological",
		Action: "created",
		Payload: map[string]interface{}{
			"id":         data.ID,
			"patientId":  data.PatientID,
			"type":       data.Type,
			"value":      data.Value,
			"measuredAt": data.MeasuredAt,
			"source":     data.Source,
			"createdAt":  data.CreatedAt,
		},
	}

	messageData, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Failed to marshal physiological data notification: %v", err)
		return
	}

	// 通知医生有新的生理数据
	s.manager.SendToRole("doctor", messageData)
} 