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

// GetAllAITemplates 获取所有AI代理模板
func GetAllAITemplates(c *gin.Context) {
	var templates []models.AIAgentTemplate
	if err := config.DB.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模板列表失败"})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// GetAITemplateByID 获取指定AI代理模板
func GetAITemplateByID(c *gin.Context) {
	id := c.Param("id")
	var template models.AIAgentTemplate
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// CreateAITemplate 创建AI代理模板
func CreateAITemplate(c *gin.Context) {
	var template models.AIAgentTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证JSON内容格式
	var js json.RawMessage
	if err := json.Unmarshal([]byte(template.Content), &js); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON格式"})
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
	template.Status = models.AIAgentStatusDraft
	template.AuditStatus = models.AIAgentAuditStatusPending

	if err := config.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建模板失败"})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateAITemplate 更新AI代理模板
func UpdateAITemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.AIAgentTemplate

	// 检查模板是否存在
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 检查模板状态
	if template.Status == models.AIAgentStatusArchived {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已归档的模板不能修改"})
		return
	}

	// 绑定更新数据
	var updateData models.AIAgentTemplate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证JSON内容格式
	if updateData.Content != "" {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(updateData.Content), &js); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON格式"})
			return
		}
	}

	// 更新字段
	userID, _ := c.Get("userId")
	template.Name = updateData.Name
	template.Description = updateData.Description
	template.Content = updateData.Content
	template.Version = updateData.Version
	template.Categories = updateData.Categories
	template.Status = updateData.Status
	template.UpdatedBy = userID.(string)
	template.UpdatedAt = time.Now()

	// 如果内容有修改，重置审核状态
	if updateData.Content != "" {
		template.AuditStatus = models.AIAgentAuditStatusPending
		template.LastAuditBy = ""
		template.LastAuditAt = time.Time{}
		template.AuditNotes = ""
	}

	if err := config.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新模板失败"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteAITemplate 删除AI代理模板
func DeleteAITemplate(c *gin.Context) {
	id := c.Param("id")

	// 检查模板是否存在
	var template models.AIAgentTemplate
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 检查模板状态
	if template.Status == models.AIAgentStatusEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "启用状态的模板不能删除，请先禁用"})
		return
	}

	if err := config.DB.Delete(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除模板失败"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AuditAITemplate 审核AI代理模板
func AuditAITemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.AIAgentTemplate

	// 检查模板是否存在
	if err := config.DB.First(&template, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}

	// 绑定审核数据
	var auditData struct {
		AuditStatus string `json:"auditStatus" binding:"required"`
		AuditNotes  string `json:"auditNotes"`
	}
	if err := c.ShouldBindJSON(&auditData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新审核信息
	userID, _ := c.Get("userId")
	template.AuditStatus = auditData.AuditStatus
	template.AuditNotes = auditData.AuditNotes
	template.LastAuditBy = userID.(string)
	template.LastAuditAt = time.Now()

	// 如果审核通过，自动启用模板
	if auditData.AuditStatus == models.AIAgentAuditStatusApproved {
		template.Status = models.AIAgentStatusEnabled
	}

	if err := config.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新审核状态失败"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// GetAITemplatesByCategory 根据分类获取AI代理模板
func GetAITemplatesByCategory(c *gin.Context) {
	category := c.Query("category")
	var templates []models.AIAgentTemplate

	query := config.DB.Where("status = ? AND audit_status = ?",
		models.AIAgentStatusEnabled,
		models.AIAgentAuditStatusApproved)

	if category != "" {
		query = query.Where("? = ANY(categories)", category)
	}

	if err := query.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模板列表失败"})
		return
	}

	c.JSON(http.StatusOK, templates)
}
