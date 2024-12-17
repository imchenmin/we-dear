package handlers

import (
	"net/http"
	"time"

	"we-dear/config"
	"we-dear/models"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

// CreateAISuggestionFeedback 创建AI建议评价
func CreateAISuggestionFeedback(c *gin.Context) {
	var feedback models.AISuggestionFeedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查AI建议是否存在
	var suggestion models.AISuggestion
	if err := config.DB.First(&suggestion, "id = ?", feedback.SuggestionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AI建议不存在"})
		return
	}

	// 设置基础字段
	userID, _ := c.Get("userId")
	now := time.Now()
	feedback.BaseModel = models.BaseModel{
		ID:        utils.GenerateID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
	feedback.CreatedBy = userID.(string)
	feedback.UpdatedBy = userID.(string)
	feedback.Status = models.AISuggestionFeedbackStatusPending
	feedback.PatientID = suggestion.PatientID

	if err := config.DB.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评价失败"})
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// UpdateAISuggestionFeedback 更新AI建议评价
func UpdateAISuggestionFeedback(c *gin.Context) {
	id := c.Param("id")
	var feedback models.AISuggestionFeedback

	// 检查评价是否存在
	if err := config.DB.First(&feedback, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评价不存在"})
		return
	}

	// 绑定更新数据
	var updateData models.AISuggestionFeedback
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	userID, _ := c.Get("userId")
	feedback.Rating = updateData.Rating
	feedback.Comment = updateData.Comment
	feedback.Tags = updateData.Tags
	feedback.UpdatedBy = userID.(string)
	feedback.UpdatedAt = time.Now()
	feedback.Status = models.AISuggestionFeedbackStatusPending

	if err := config.DB.Save(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评价失败"})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// GetAISuggestionFeedbacks 获取AI建议评价列表
func GetAISuggestionFeedbacks(c *gin.Context) {
	suggestionID := c.Query("suggestionId")
	patientID := c.Query("patientId")
	status := c.Query("status")

	query := config.DB.Model(&models.AISuggestionFeedback{})

	if suggestionID != "" {
		query = query.Where("id = ?", suggestionID)
	}
	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var feedbacks []models.AISuggestionFeedback
	if err := query.Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价列表失败"})
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

// ReviewAISuggestionFeedback 审核AI建议评价
func ReviewAISuggestionFeedback(c *gin.Context) {
	id := c.Param("id")
	var feedback models.AISuggestionFeedback

	// 检查评价是否存在
	if err := config.DB.First(&feedback, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评价不存在"})
		return
	}

	// 绑定审核数据
	var reviewData struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reviewData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新审核信息
	userID, _ := c.Get("userId")
	feedback.Status = reviewData.Status
	feedback.ReviewedBy = userID.(string)
	feedback.ReviewedAt = time.Now()

	if err := config.DB.Save(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新审核状态失败"})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// GetFeedbackStats 获取AI建议评价统计信息
func GetFeedbackStats(c *gin.Context) {
	suggestionID := c.Query("suggestionId")
	patientID := c.Query("patientId")

	query := config.DB.Model(&models.AISuggestionFeedback{}).
		Where("status = ?", models.AISuggestionFeedbackStatusApproved)

	if suggestionID != "" {
		query = query.Where("id = ?", suggestionID)
	}
	if patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// 统计点赞和踩的数量
	var stats struct {
		Likes    int64 `json:"likes"`
		Dislikes int64 `json:"dislikes"`
	}

	if err := query.Where("rating = ?", models.AISuggestionFeedbackRatingLike).Count(&stats.Likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "统计点赞数失败"})
		return
	}

	if err := query.Where("rating = ?", models.AISuggestionFeedbackRatingDislike).Count(&stats.Dislikes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "统计踩数失败"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
