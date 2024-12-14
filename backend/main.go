package main

import (
	"log"
	"time"
	"we-dear/config"
	"we-dear/handlers"
	"we-dear/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库连接
	config.InitDB()

	// 等待一下确保数据库连接完全建立
	time.Sleep(time.Second)

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
		api.POST("/patients", handlers.CreatePatient)

		// 医生相关
		api.GET("/doctors", handlers.GetAllDoctors)
		api.POST("/doctors", handlers.CreateDoctor)
		api.PUT("/doctors/:id", handlers.UpdateDoctor)
		api.DELETE("/doctors/:id", handlers.DeleteDoctor)

		// 科室相关
		api.GET("/departments", handlers.GetAllDepartments)
		api.POST("/departments", handlers.CreateDepartment)
		api.PUT("/departments/:id", handlers.UpdateDepartment)
		api.DELETE("/departments/:id", handlers.DeleteDepartment)

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
