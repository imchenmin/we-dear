package handlers

import (
	"net/http"
	"time"

	"we-dear/models"
	"we-dear/storage"

	"github.com/gin-gonic/gin"
)

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

func SendMessage(c *gin.Context) {
	id := c.Param("id")
	var message models.Message
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.Timestamp = time.Now()

	err := storage.GetPatientStorage().AddMessage(id, message)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}
