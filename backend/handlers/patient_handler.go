package handlers

import (
	"net/http"
	"time"

	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

func initPatientStorage() *storage.PatientStorage {
	return storage.GetPatientStorage()
}

func GetAllPatients(c *gin.Context) {
	// 从上下文获取当前医生信息
	userID, _ := c.Get("userId")
	role, _ := c.Get("role")

	patientStorage := initPatientStorage()
	var patients []models.Patient

	if role == "admin" {
		// 管理员可以看到所有患者
		patients = patientStorage.GetAllPatients()
	} else {
		// 普通医生只能看到自己的患者
		patients = patientStorage.GetPatientsByDoctorID(userID.(string))
	}

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
