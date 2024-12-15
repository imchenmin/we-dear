package handlers

import (
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"we-dear/config"
	"we-dear/models"
	"we-dear/services"
)

var (
	aiService = services.NewAIService()
)

// 通用的消息请求结构
type MessageRequest struct {
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Sender    string `json:"sender"`
	Avatar    string `json:"avatar,omitempty"`
}

// GetChatList 获取医生的聊天列表
func GetChatList(c *gin.Context) {
	// 从上下文获取当前医生信息
	userID, _ := c.Get("userId")
	role, _ := c.Get("role")

	var patients []models.Patient
	db := config.DB

	if role == "admin" {
		// 管理员可以看到所有患者
		err := db.Preload("Doctor").Find(&patients).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取患者列表失败"})
			return
		}
	} else {
		// 普通医生只能看到自己的患者
		err := db.Preload("Doctor").Where("doctor_id = ?", userID).Find(&patients).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取患者列表失败"})
			return
		}
	}

	// 获取每个患者的最后一条消息
	type ChatItem struct {
		PatientID     string    `json:"patientId"`
		PatientName   string    `json:"patientName"`
		PatientAvatar string    `json:"patientAvatar"`
		DoctorID      string    `json:"doctorId"`
		DoctorName    string    `json:"doctorName"`
		LastMessage   string    `json:"lastMessage"`
		LastMessageAt time.Time `json:"lastMessageAt"`
		UnreadCount   int       `json:"unreadCount"`
	}

	var chatList []ChatItem
	for _, patient := range patients {
		var lastMessage models.Message
		var unreadCount int64

		// 获取最后一条消息
		err := db.Where("patient_id = ?", patient.ID).
			Order("created_at desc").
			First(&lastMessage).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}

		// 获取未读消息数量
		db.Model(&models.Message{}).
			Where("patient_id = ? AND role = ? AND read = ?", patient.ID, "patient", false).
			Count(&unreadCount)

		chatList = append(chatList, ChatItem{
			PatientID:     patient.ID,
			PatientName:   patient.Name,
			PatientAvatar: patient.Avatar,
			DoctorID:      patient.Doctor.ID,
			DoctorName:    patient.Doctor.Name,
			LastMessage:   lastMessage.Content,
			LastMessageAt: lastMessage.CreatedAt,
			UnreadCount:   int(unreadCount),
		})
	}

	// 按最后消息时间排序
	sort.Slice(chatList, func(i, j int) bool {
		return chatList[i].LastMessageAt.After(chatList[j].LastMessageAt)
	})

	c.JSON(http.StatusOK, chatList)
}

// GetChatHistory 获取具体聊天记录
func GetChatHistory(c *gin.Context) {
	patientID := c.Param("patientId")
	userID, _ := c.Get("userId")
	role, _ := c.Get("role")

	// 检查权限
	if role != "admin" {
		// 非管理员只能查看自己的患者聊天记录
		var patient models.Patient
		if err := config.DB.First(&patient, "id = ?", patientID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "患者不存在"})
			return
		}

		if patient.DoctorID != userID.(string) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此患者的聊天记录"})
			return
		}
	}

	// 获取聊天记录
	var messages []models.Message
	if err := config.DB.Where("patient_id = ?", patientID).
		Order("created_at asc").
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取聊天记录失败"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// SendDoctorMessage 医生发送消息
func SendDoctorMessage(c *gin.Context) {
	patientId := c.Param("patientId")
	var req MessageRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := models.Message{
		BaseModel: models.BaseModel{
			ID:        strconv.FormatInt(time.Now().UnixNano(), 10),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PatientID: patientId,
		DoctorID:  req.Sender,
		Content:   req.Content,
		Type:      models.MessageTypeText,
		Role:      models.MessageRoleDoctor,
		Read:      false,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

// SendPatientMessage 患者发送消息
func SendPatientMessage(c *gin.Context) {
	patientId := c.Param("patientId")
	var req MessageRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取患者信息
	var patient models.Patient
	if err := config.DB.First(&patient, "id = ?", patientId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// 生成消息
	messageID := strconv.FormatInt(time.Now().UnixNano(), 10)
	message := models.Message{
		BaseModel: models.BaseModel{
			ID:        messageID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		PatientID: patientId,
		Content:   req.Content,
		Type:      models.MessageTypeText,
		Role:      models.MessageRolePatient,
		Read:      false,
	}

	// 保存患者消息
	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成 AI 建议（异步）
	go func() {
		// 获取历史消息
		var messages []models.Message
		if err := config.DB.Where("patient_id = ?", patientId).
			Order("created_at asc").
			Find(&messages).Error; err != nil {
			log.Printf("获取聊天历史失败: %v", err)
			return
		}

		suggestion, err := aiService.GenerateResponse(&patient, messageID, req.Content, messages)
		if err != nil {
			log.Printf("生成 AI 建议失败: %v", err)
			return
		}

		// 保存 AI 建议
		if err := config.DB.Create(&suggestion).Error; err != nil {
			log.Printf("保存 AI 建议失败: %v", err)
		}
	}()

	c.JSON(http.StatusOK, message)
}

// GetAISuggestions 获取医生视图的 AI 建议
func GetAISuggestions(c *gin.Context) {
	patientId := c.Param("patientId")
	messageId := c.Query("messageId")

	var suggestions []models.AISuggestion
	if err := config.DB.Where("patient_id = ? AND message_id = ?", patientId, messageId).
		Find(&suggestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}
