package main

import (
	"fmt"
	"log"
	"we-dear/models"
	"we-dear/services"
)

func main() {
	// 测试患者数据
	testPatient := &models.Patient{
		Name:            "测试患者",
		Gender:          "男",
		Age:             30,
		BloodType:       "A",
		Allergies:       []string{"青霉素"},
		ChronicDiseases: []string{"高血压"},
	}

	// 测试消息
	testMessage := "我最近经常头痛，该怎么办？"

	// 测试历史消息
	testHistory := []models.Message{
		{
			Role:    models.MessageRolePatient,
			Content: "医生，我最近睡眠不好",
		},
	}

	// 初始化AI服务
	aiService := services.NewAIService()

	// 生成回复
	suggestion, err := aiService.GenerateResponse(testPatient, "test_msg_001", testMessage, testHistory)
	if err != nil {
		log.Fatalf("生成AI回复失败: %v", err)
	}

	// 打印结果
	fmt.Printf("\n=== AI回复测试结果 ===\n")
	fmt.Printf("建议内容：\n%s\n", suggestion.Content)
	fmt.Printf("置信度：%.2f\n", suggestion.Confidence)
	fmt.Printf("类别：%s\n", suggestion.Category)
	fmt.Printf("优先级：%d\n", suggestion.Priority)
}
