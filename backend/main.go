package main

import (
	"we-dear/handlers"
	"we-dear/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
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
		api.GET("/chat/:patientId", handlers.GetChatHistory)              // 获取聊天历史
		api.POST("/chat/:patientId/doctor", handlers.SendDoctorMessage)   // 医生发送消息
		api.POST("/chat/:patientId/patient", handlers.SendPatientMessage) // 患者发送消息

		// 文件上传
		api.POST("/upload", handlers.HandleUpload)
	}

	router.Run(":8080")
}
