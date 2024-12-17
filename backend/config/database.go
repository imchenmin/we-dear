package config

import (
	"fmt"
	"log"
	"time"
	"we-dear/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	fmt.Println("Initializing database connection...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		GlobalConfig.DB.Host,
		GlobalConfig.DB.User,
		GlobalConfig.DB.Password,
		GlobalConfig.DB.Name,
		GlobalConfig.DB.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 验证数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// 测试连接
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Database connection established successfully")

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库结构
	err = DB.AutoMigrate(
		&models.Patient{},
		&models.Message{},
		&models.AISuggestion{},
		&models.MedicalRecord{},
		&models.Doctor{},
		&models.Department{},
		&models.Attachment{},
		&models.FollowUpRecord{},
		&models.AISuggestionFeedback{},
		&models.AIAgentTemplate{},
		&models.FollowUpTemplate{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed")
}
