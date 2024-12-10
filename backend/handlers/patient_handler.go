package handlers

import (
	"net/http"
	"strconv"
	"time"

	"we-dear/models"
	"we-dear/storage"

	"github.com/gin-gonic/gin"
)

// 通用的消息请求结构
type MessageRequest struct {
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Sender    string `json:"sender"`
	Avatar    string `json:"avatar,omitempty"`
}

func GetAllPatients(c *gin.Context) {
	patients := storage.GetPatientStorage().GetAllPatients()
	c.JSON(http.StatusOK, patients)
}

func GetPatientById(c *gin.Context) {
	id := c.Param("id")
	patient, err := storage.GetPatientStorage().GetPatientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func GetChatHistory(c *gin.Context) {
	patientId := c.Param("patientId")
	messages, err := storage.GetPatientStorage().GetChatHistory(patientId)
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
		ID:        strconv.FormatInt(time.Now().UnixNano(), 10),
		Content:   req.Content,
		Timestamp: time.Now(),
		Role:      "doctor",
		Sender:    req.Sender,
		Avatar:    req.Avatar,
	}

	err := storage.GetPatientStorage().AddMessage(patientId, message)
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

	message := models.Message{
		ID:        strconv.FormatInt(time.Now().UnixNano(), 10),
		Content:   req.Content,
		Timestamp: time.Now(),
		Role:      "patient",
		Sender:    req.Sender,
		Avatar:    req.Avatar,
	}

	err := storage.GetPatientStorage().AddMessage(patientId, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}
