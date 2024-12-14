package main

import (
	"log"
	"time"
	"we-dear/config"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"
)

func InitDepartments() {
	// 初始化配置和数据库连接
	config.Init()
	config.InitDB()

	// 获取部门存储实例
	deptStorage := storage.GetDepartmentStorage()

	// 初始化部门数据
	departments := []models.Department{
		{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "内科",
			Description: "负责内科疾病的诊断和治疗",
			Code:        "NEI",
		},
		{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "外科",
			Description: "负责外科手术和治疗",
			Code:        "WAI",
		},
		// ... 其他部门数据
	}

	// 插入部门数据
	for _, dept := range departments {
		err := deptStorage.CreateDepartment(&dept)
		if err != nil {
			log.Printf("Failed to create department %s: %v", dept.Name, err)
		} else {
			log.Printf("Successfully created department: %s", dept.Name)
		}
	}
	// 打印所有部门
	departments, err := deptStorage.GetAllDepartments()
	if err != nil {
		log.Printf("Failed to get all departments: %v", err)
	} else {
		log.Printf("All departments: %v", departments)
	}
}

func main() {
	InitDepartments()
}
