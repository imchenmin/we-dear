package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"we-dear/models"
	"we-dear/services"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

var (
	aiService = services.NewAIService()
)

func initPatientStorage() *storage.PatientStorage {
	return storage.GetPatientStorage()
}

// 通用的消息请求结构
type MessageRequest struct {
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Sender    string `json:"sender"`
	Avatar    string `json:"avatar,omitempty"`
}

func GetAllPatients(c *gin.Context) {
	patientStorage := initPatientStorage()
	patients := patientStorage.GetAllPatients()
	c.JSON(http.StatusOK, patients)
}

func GetPatientById(c *gin.Context) {
	id := c.Param("id")
	patient, err := initPatientStorage().GetPatientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func GetChatHistory(c *gin.Context) {
	patientId := c.Param("patientId")
	messages, err := initPatientStorage().GetChatHistory(patientId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

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
		Status:    models.MessageStatusUnread,
	}

	err := initPatientStorage().AddMessage(patientId, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func SendPatientMessage(c *gin.Context) {
	patientId := c.Param("patientId")
	var req MessageRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取患者信息和历史消息
	patient, err := initPatientStorage().GetPatientByID(patientId)
	if err != nil {
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
		Status:    models.MessageStatusUnread,
	}

	// 保存患者消息
	err = initPatientStorage().AddMessage(patientId, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成 AI 建议（异步）
	go func() {
		// 获取历史消息
		messages, err := initPatientStorage().GetChatHistory(patientId)
		if err != nil {
			log.Printf("获取聊天历史失败: %v", err)
			return
		}

		suggestion, err := aiService.GenerateResponse(patient, messageID, req.Content, messages)
		if err != nil {
			log.Printf("生成 AI 建议失败: %v", err)
			return
		}

		// 保存 AI 建议
		err = initPatientStorage().SaveAISuggestion(suggestion)
		if err != nil {
			log.Printf("保存 AI 建议失败: %v", err)
		}
	}()

	c.JSON(http.StatusOK, message)
}

// GetAISuggestions 获取医生视图的 AI 建议
func GetAISuggestions(c *gin.Context) {
	patientId := c.Param("patientId")
	messageId := c.Query("messageId")

	suggestions, err := initPatientStorage().GetAISuggestions(patientId, messageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

func CreatePatient(c *gin.Context) {
	var patient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置创建时间等基础字段
	now := time.Now()
	patient.BaseModel = models.BaseModel{
		ID:        utils.GenerateID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := patient.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initPatientStorage().CreatePatient(&patient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}
