package test

import (
	"testing"
	"time"
	"we-dear/config"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"
)

func TestDepartments(t *testing.T) {
	// 初始化配置和数据库连接
	config.Init()
	config.InitDB()

	// 获取部门存储实例
	deptStorage := storage.GetDepartmentStorage()

	// 测试部门数据
	testDepts := []models.Department{
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
		{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "儿科",
			Description: "专注于儿童疾病的诊断和治疗",
			Code:        "ER",
		},
		{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "妇产科",
			Description: "专注于妇女疾病和产科服务",
			Code:        "FU",
		},
		{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:        "骨科",
			Description: "负责骨骼和关节疾病的诊断治疗",
			Code:        "GU",
		},
	}

	// 插入测试数据
	for _, dept := range testDepts {
		err := deptStorage.CreateDepartment(&dept)
		if err != nil {
			t.Errorf("Failed to create department %s: %v", dept.Name, err)
		}
	}

	// 验证数据插入
	departments, err := deptStorage.GetAllDepartments()
	if err != nil {
		t.Errorf("Failed to get departments: %v", err)
	}

	if len(departments) < len(testDepts) {
		t.Errorf("Expected at least %d departments, got %d", len(testDepts), len(departments))
	}
}
