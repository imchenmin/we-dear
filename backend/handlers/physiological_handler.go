package handlers

import (
	"net/http"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"
	"we-dear/websocket"

	"github.com/gin-gonic/gin"
)

// GetPhysiologicalData 获取患者的生理数据
func GetPhysiologicalData(c *gin.Context) {
	patientID := c.Param("id")
	if patientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "患者ID不能为空"})
		return
	}

	// 获取数据类型参数（可选）
	dataType := c.Query("type")
	
	var data []models.PhysiologicalData
	var err error
	
	store := storage.GetPhysiologicalDataStorage()
	if dataType != "" {
		data, err = store.GetByPatientIDAndType(patientID, dataType)
	} else {
		data, err = store.GetByPatientID(patientID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取生理数据失败"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// CreatePhysiologicalData 创建生理数据记录
func CreatePhysiologicalData(c *gin.Context) {
	var data models.PhysiologicalData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 生成ID
	data.ID = utils.GenerateID()

	store := storage.GetPhysiologicalDataStorage()
	if err := store.Create(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建生理数据失败"})
		return
	}

	// 发送WebSocket通知
	websocket.GetService().NotifyNewPhysiologicalData(&data)

	c.JSON(http.StatusOK, data)
}

// UpdatePhysiologicalData 更新生理数据记录
func UpdatePhysiologicalData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	var data models.PhysiologicalData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	data.ID = id
	store := storage.GetPhysiologicalDataStorage()
	if err := store.Update(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新生理数据失败"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeletePhysiologicalData 删除生理数据记录
func DeletePhysiologicalData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	store := storage.GetPhysiologicalDataStorage()
	if err := store.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除生理数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
} 