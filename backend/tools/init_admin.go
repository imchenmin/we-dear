package main

import (
	"fmt"
	"log"
	"time"
	"we-dear/config"
	"we-dear/models"
	"we-dear/utils"
)

func main() {
	// 初始化配置和数据库连接
	config.Init()
	config.InitDB()

	// 创建管理员科室
	adminDept := &models.Department{
		BaseModel: models.BaseModel{
			ID:        utils.GenerateID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        "系统管理",
		Code:        "ADMIN",
		Description: "系统管理部门",
	}

	// 检查管理员科室是否已存在
	var deptCount int64
	if err := config.DB.Model(&models.Department{}).Where("code = ?", adminDept.Code).Count(&deptCount).Error; err != nil {
		log.Fatalf("检查管理员科室失败: %v", err)
	}

	if deptCount == 0 {
		if err := config.DB.Create(adminDept).Error; err != nil {
			log.Fatalf("创建管理员科室失败: %v", err)
		}
		log.Printf("管理员科室创建成功，ID: %s", adminDept.ID)
	} else {
		// 如果科室已存在，获取其ID
		if err := config.DB.Where("code = ?", adminDept.Code).First(adminDept).Error; err != nil {
			log.Fatalf("获取管理员科室失败: %v", err)
		}
	}

	// 管理员信息
	adminInfo := struct {
		Username string
		Password string
		Name     string
	}{
		Username: "admin",
		Password: "admin123", // 在实际使用时应该使用更强的密码
		Name:     "系统管理员",
	}

	// 检查管理员是否已存在
	var count int64
	if err := config.DB.Model(&models.Doctor{}).Where("username = ?", adminInfo.Username).Count(&count).Error; err != nil {
		log.Fatalf("检查管理员失败: %v", err)
	}

	if count > 0 {
		log.Printf("管理员账户 %s 已存在", adminInfo.Username)
		return
	}

	// 生成密码盐和哈希
	salt, err := utils.GenerateSalt()
	if err != nil {
		log.Fatalf("生成密码盐失败: %v", err)
	}
	hashedPassword := utils.HashPassword(adminInfo.Password, salt)

	// 创建管理员用户
	admin := &models.Doctor{
		BaseModel: models.BaseModel{
			ID:        utils.GenerateID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:     adminInfo.Username,
		Password:     hashedPassword,
		Salt:         salt,
		Name:         adminInfo.Name,
		Status:       "active",
		Role:         "admin",
		DepartmentID: adminDept.ID, // 关联到管理员科室
		Title:        "系统管理员",
	}

	if err := config.DB.Create(admin).Error; err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}

	fmt.Printf("初始化完成:\n")
	fmt.Printf("管理员科室:\n")
	fmt.Printf("  - ID: %s\n", adminDept.ID)
	fmt.Printf("  - 名称: %s\n", adminDept.Name)
	fmt.Printf("  - 代码: %s\n", adminDept.Code)
	fmt.Printf("\n管理员账户:\n")
	fmt.Printf("  - ID: %s\n", admin.ID)
	fmt.Printf("  - 用户名: %s\n", adminInfo.Username)
	fmt.Printf("  - 密码: %s\n", adminInfo.Password)
}
