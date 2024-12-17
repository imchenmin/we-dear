package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"we-dear/config"
	"we-dear/models"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

// GetAllTemplates 获取所有随访模板
func GetAllTemplates(c *gin.Context) {
	var templates []models.FollowUpTemplate
	if err := config.DB.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模板列表失败"})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// GetTemplateByID 获取指定模板
func GetTemplateByID(c *gin.Context) {
	id := c.Param("id")
	var template models.FollowUpTemplate
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// CreateTemplate 创建随访模板
func CreateTemplate(c *gin.Context) {
	var template models.FollowUpTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证JSON Schema格式
	var js json.RawMessage
	if err := json.Unmarshal([]byte(template.Schema), &js); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON Schema格式"})
		return
	}

	// 设置基础字段
	userID, _ := c.Get("userId")
	now := time.Now()
	template.BaseModel = models.BaseModel{
		ID:        utils.GenerateID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	template.CreatedBy = userID.(string)
	template.UpdatedBy = userID.(string)
	template.Status = models.TemplateStatusEnabled

	if err := config.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建模板失败"})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateTemplate 更新随访模板
func UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.FollowUpTemplate

	// 检查模板是否存在
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 绑定更新数据
	var updateData models.FollowUpTemplate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证JSON Schema格式
	if updateData.Schema != "" {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(updateData.Schema), &js); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON Schema格式"})
			return
		}
	}

	// 更新字段
	userID, _ := c.Get("userId")
	template.Name = updateData.Name
	template.Description = updateData.Description
	template.Schema = updateData.Schema
	template.Version = updateData.Version
	template.Categories = updateData.Categories
	template.Status = updateData.Status
	template.UpdatedBy = userID.(string)
	template.UpdatedAt = time.Now()

	if err := config.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新模板失败"})
		return
	}

	c.JSON(http.StatusOK, template)
}

func GetDefaultSchema(c *gin.Context) {
	c.JSON(http.StatusOK, utils.GetDefaultSchema())
}

// DeleteTemplate 删除随访模板
func DeleteTemplate(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有随访记录使用此模板
	var count int64
	if err := config.DB.Model(&models.FollowUpRecord{}).Where("template_id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查模板使用情况失败"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该模板已被使用，无法删除"})
		return
	}

	if err := config.DB.Delete(&models.FollowUpTemplate{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除模板失败"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetTemplatesByCategory 根据分类获取模板
func GetTemplatesByCategory(c *gin.Context) {
	category := c.Query("category")
	var templates []models.FollowUpTemplate

	query := config.DB.Where("status = ?", models.TemplateStatusEnabled)
	if category != "" {
		query = query.Where("? = ANY(categories)", category)
	}

	if err := query.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模板列表失败"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// ValidateTemplateData 验证随访记录数据是否符合模板要求
func ValidateTemplateData(c *gin.Context) {
	templateID := c.Query("templateId")
	var template models.FollowUpTemplate

	// 获取模板
	if err := config.DB.First(&template, "id = ?", templateID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 获取要验证的数据
	var data interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON数据"})
		return
	}

	// 将数据转换为JSON字符串
	dataBytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据序列化失败"})
		return
	}

	// 验证数据是否符合模板schema
	valid, err := utils.ValidateJSONSchema(template.Schema, string(dataBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证失败: " + err.Error()})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": "数据不符合模板要求",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "验证通过",
	})
}
