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
		// 公开路由
		api.POST("/login", handlers.Login)
		// api.POST("/register", handlers.Register)
	}

	// 需要认证的路由
	authorized := api.Group("")
	authorized.Use(middleware.AuthRequired())
	{
		// 患者相关
		authorized.GET("/patients", handlers.GetAllPatients)
		authorized.GET("/patients/:id", handlers.GetPatientById)
		authorized.GET("/patients/:id/followup", handlers.GetFollowUpRecords)
		authorized.POST("/patients", handlers.CreatePatient)

		// 医生相关
		authorized.GET("/doctors", handlers.GetAllDoctors)
		authorized.POST("/doctors", middleware.AdminRequired(), handlers.CreateDoctor)
		authorized.PUT("/doctors/:id", handlers.UpdateDoctor)
		authorized.DELETE("/doctors/:id", middleware.AdminRequired(), handlers.DeleteDoctor)

		// 科室相关
		authorized.GET("/departments", handlers.GetAllDepartments)
		authorized.POST("/departments", middleware.AdminRequired(), handlers.CreateDepartment)
		authorized.PUT("/departments/:id", middleware.AdminRequired(), handlers.UpdateDepartment)
		authorized.DELETE("/departments/:id", middleware.AdminRequired(), handlers.DeleteDepartment)

		// 消息相关
		authorized.GET("/chat/list", handlers.GetChatList)                        // 获取聊天列表
		authorized.GET("/chat/:patientId", handlers.GetChatHistory)               // 获取聊天历史
		authorized.POST("/chat/:patientId/doctor", handlers.SendDoctorMessage)    // 医生发送消息
		authorized.POST("/chat/:patientId/patient", handlers.SendPatientMessage)  // 患者发送消息
		authorized.GET("/chat/:patientId/suggestions", handlers.GetAISuggestions) // 获取 AI 建议

		// 用户认证相关
		authorized.POST("/change-password", handlers.ChangePassword)

		// 随访记录相关路由
		authorized.POST("/followup", handlers.CreateFollowUpRecord)
		authorized.PUT("/followup/:id", handlers.UpdateFollowUpRecord)
		authorized.DELETE("/followup/:id", handlers.DeleteFollowUpRecord)

		// 医疗记录相关路由
		authorized.GET("/patients/:id/medical", handlers.GetMedicalRecords)
		authorized.POST("/medical", handlers.CreateMedicalRecord)
		authorized.PUT("/medical/:id", handlers.UpdateMedicalRecord)
		authorized.DELETE("/medical/:id", handlers.DeleteMedicalRecord)
	}

	log.Printf("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
