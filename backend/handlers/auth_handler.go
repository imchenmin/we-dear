package handlers

import (
	"net/http"
	"time"
	"we-dear/storage"
	"we-dear/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	UserID      string `json:"userId"`      // 要修改密码的用户ID（管理员使用）
	OldPassword string `json:"oldPassword"` // 旧密码（普通用户必填）
	NewPassword string `json:"newPassword" binding:"required"`
}

// Login 处理登录请求
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
