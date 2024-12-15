package handlers

import (
	"net/http"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

func GetFollowUpRecords(c *gin.Context) {
	patientID := c.Param("id")
	records, err := storage.GetMedicalStorage().GetFollowUpRecords(patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func CreateFollowUpRecord(c *gin.Context) {
	var record models.FollowUpRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置基础字段
	record.ID = utils.GenerateID()
	userID, _ := c.Get("userId")
	record.DoctorID = userID.(string)

	if err := storage.GetMedicalStorage().CreateFollowUpRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func UpdateFollowUpRecord(c *gin.Context) {
	var record models.FollowUpRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := storage.GetMedicalStorage().UpdateFollowUpRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func DeleteFollowUpRecord(c *gin.Context) {
	id := c.Param("id")
	if err := storage.GetMedicalStorage().DeleteFollowUpRecord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// 医疗记录的处理函数
func GetMedicalRecords(c *gin.Context) {
	patientID := c.Param("id")
	records, err := storage.GetMedicalStorage().GetMedicalRecords(patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func CreateMedicalRecord(c *gin.Context) {
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.ID = utils.GenerateID()
	userID, _ := c.Get("userId")
	record.DoctorID = userID.(string)

	if err := storage.GetMedicalStorage().CreateMedicalRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func UpdateMedicalRecord(c *gin.Context) {
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := storage.GetMedicalStorage().UpdateMedicalRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func DeleteMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	if err := storage.GetMedicalStorage().DeleteMedicalRecord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}
