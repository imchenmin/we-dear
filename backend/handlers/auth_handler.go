package handlers

import (
	"net/http"
	"time"
	"we-dear/config"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=doctor patient"`
}

type ChangePasswordRequest struct {
	UserID      string `json:"userId"`      // 要修改密码的用户ID（管理员使用）
	OldPassword string `json:"oldPassword"` // 旧密码（普通用户必填）
	NewPassword string `json:"newPassword" binding:"required"`
}

type PatientRegistrationRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	DoctorID string `json:"doctorId" binding:"required"`
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "doctor" {
		handleDoctorLogin(c, req)
	} else {
		handlePatientLogin(c, req)
	}
}

func handleDoctorLogin(c *gin.Context, req LoginRequest) {
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
	token, err := utils.GenerateToken(doctor.ID, doctor.Name, "doctor")
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
			"role":     "doctor",
			"avatar":   doctor.Avatar,
		},
	})
}

func handlePatientLogin(c *gin.Context, req LoginRequest) {
	var patient models.Patient
	if err := config.DB.Where("username = ?", req.Username).First(&patient).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	hashedPassword := utils.HashPassword(req.Password, patient.Salt)
	if hashedPassword != patient.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成token - 使用username而不是name
	token, err := utils.GenerateToken(patient.ID, patient.Username, "patient")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       patient.ID,
			"username": patient.Username,
			"name":     patient.Name,
			"role":     "patient",
			"doctorId": patient.DoctorID,
		},
	})
}

// ChangePassword 处理修改密码请求
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	currentUserID, _ := c.Get("userId")
	currentRole, _ := c.Get("role")

	// 确定要修改密码的用户ID
	targetUserID := currentUserID.(string)
	if req.UserID != "" && currentRole == "admin" {
		targetUserID = req.UserID
	}

	// 获取目标用户信息
	doctor, err := storage.GetDoctorStorage().GetDoctorByID(targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 如果不是管理员，需要验证旧密码
	if currentRole != "admin" {
		if req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请输入原密码"})
			return
		}

		hashedOldPassword := utils.HashPassword(req.OldPassword, doctor.Salt)
		if hashedOldPassword != doctor.Password {
			c.JSON(http.StatusBadRequest, gin.H{"error": "原密码错误"})
			return
		}
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

func Register(c *gin.Context) {
	if !config.GlobalConfig.App.AllowPatientRegistration {
		c.JSON(http.StatusForbidden, gin.H{"error": "Patient registration is currently disabled"})
		return
	}

	var req PatientRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if doctor exists
	var doctor models.Doctor
	if err := config.DB.First(&doctor, "id = ?", req.DoctorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Selected doctor not found"})
		return
	}

	// Check if username already exists
	var existingPatient models.Patient
	if err := config.DB.Where("username = ?", req.Username).First(&existingPatient).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already registered"})
		return
	}

	// Generate salt and hash password
	salt, err := utils.GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
		return
	}
	hashedPassword := utils.HashPassword(req.Password, salt)

	patient := models.Patient{
		Name:     req.Name,
		Username: req.Username,
		Password: hashedPassword,
		Salt:     salt,
		DoctorID: req.DoctorID,
	}

	if err := config.DB.Create(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"patient": gin.H{
			"id":       patient.ID,
			"username": patient.Username,
			"name":     patient.Name,
			"doctorId": patient.DoctorID,
		},
	})
}
