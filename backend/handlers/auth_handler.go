package handlers

import (
	"net/http"
	"time"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := storage.GetDoctorStorage().GetDoctorByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	hashedPassword := utils.HashPassword(req.Password, doctor.Salt)
	if hashedPassword != doctor.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 更新最后登录时间
	doctor.LastLoginAt = time.Now()
	storage.GetDoctorStorage().UpdateDoctor(doctor)

	// 生成token
	token, err := utils.GenerateToken(doctor.ID, doctor.Username, doctor.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       doctor.ID,
			"username": doctor.Username,
			"name":     doctor.Name,
			"role":     doctor.Role,
			"avatar":   doctor.Avatar,
		},
	})
}

type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Title        string `json:"title"`
	DepartmentID string `json:"departmentId"`
	License      string `json:"license"`
	Specialty    string `json:"specialty"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	_, err := storage.GetDoctorStorage().GetDoctorByUsername(req.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 生成密码盐和哈希
	salt, err := utils.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成密码盐失败"})
		return
	}
	hashedPassword := utils.HashPassword(req.Password, salt)

	// 创建新医生
	doctor := &models.Doctor{
		BaseModel: models.BaseModel{
			ID:        utils.GenerateID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     req.Username,
		Password:     hashedPassword,
		Salt:         salt,
		Name:         req.Name,
		Title:        req.Title,
		DepartmentID: req.DepartmentID,
		License:      req.License,
		Specialty:    req.Specialty,
		Status:       "active",
		Role:         "doctor", // 默认角色为普通医生
	}

	if err := storage.GetDoctorStorage().CreateDoctor(doctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"id":      doctor.ID,
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取当前用户ID
	userID, _ := c.Get("userId")

	// 获取用户信息
	doctor, err := storage.GetDoctorStorage().GetDoctorByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 验证旧密码
	hashedOldPassword := utils.HashPassword(req.OldPassword, doctor.Salt)
	if hashedOldPassword != doctor.Password {
		c.JSON(http.StatusBadRequest, gin.H{"error": "旧密码错误"})
		return
	}

	// 生成新的密码盐和哈希
	newSalt, err := utils.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成密码盐失败"})
		return
	}
	hashedNewPassword := utils.HashPassword(req.NewPassword, newSalt)

	// 更新密码
	doctor.Password = hashedNewPassword
	doctor.Salt = newSalt
	if err := storage.GetDoctorStorage().UpdateDoctor(doctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
