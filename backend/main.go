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
		api.GET("/patients", handlers.GetAllPatients)
		api.GET("/patients/:id", handlers.GetPatientById)
		api.POST("/patients/:id/messages", handlers.SendMessage)
		api.POST("/upload", handlers.HandleUpload)
	}

	router.Run(":8080")
}
