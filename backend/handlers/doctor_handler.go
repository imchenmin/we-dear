package handlers

import (
	"net/http"
	"time"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

func initDoctorStorage() *storage.DoctorStorage {
	return storage.GetDoctorStorage()
}

func CreateDoctor(c *gin.Context) {
	var doctor models.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置创建时间等基础字段
	now := time.Now()
	doctor.BaseModel = models.BaseModel{
		ID:        utils.GenerateID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// if err := doctor.Validate(); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	if err := initDoctorStorage().CreateDoctor(&doctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doctor)
}

func GetAllDoctors(c *gin.Context) {
	doctors, err := initDoctorStorage().GetAllDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

func UpdateDoctor(c *gin.Context) {
	id := c.Param("id")
	var doctor models.Doctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor.ID = id
	if err := initDoctorStorage().UpdateDoctor(&doctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func DeleteDoctor(c *gin.Context) {
	id := c.Param("id")
	if err := initDoctorStorage().DeleteDoctor(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
