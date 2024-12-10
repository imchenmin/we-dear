package main

import (
	"log"
	"we-dear/config"
	"we-dear/handlers"
	"we-dear/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	router := gin.Default()

	// 中间件
	router.Use(middleware.Cors())

	// 静态文件服务
	router.Static("/uploads", "./uploads")

	// API 路由
	api := router.Group("/api")
	{
		// 患者相关
		api.GET("/patients", handlers.GetAllPatients)
		api.GET("/patients/:id", handlers.GetPatientById)

		// 消息相关
		api.GET("/chat/:patientId", handlers.GetChatHistory)               // 获取聊天历史
		api.POST("/chat/:patientId/doctor", handlers.SendDoctorMessage)    // 医生发送消息
		api.POST("/chat/:patientId/patient", handlers.SendPatientMessage)  // 患者发送消息
		api.GET("/chat/:patientId/suggestions", handlers.GetAISuggestions) // 获取 AI 建议
	}

	log.Printf("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
