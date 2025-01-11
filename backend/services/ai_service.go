package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"we-dear/config"
	"we-dear/models"
	"we-dear/storage"
	"we-dear/utils"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
	deepseek "github.com/cohesion-org/deepseek-go"
)

type AIService struct {
	openaiClient   *openai.Client
	deepseekClient *deepseek.Client
	provider       string
}

func NewAIService() *AIService {
	provider := config.GlobalConfig.AI.Provider
	var openaiClient *openai.Client
	var deepseekClient *deepseek.Client

	// 添加更详细的日志
	log.Printf("初始化 AIService:")
	log.Printf("GlobalConfig: %+v", config.GlobalConfig)
	log.Printf("AI提供商: %s", provider)
	log.Printf("Deepseek Key长度: %d", len(config.GlobalConfig.AI.DeepseekKey))

	if provider == "openai" {
		openaiClient = openai.NewClient(config.GlobalConfig.AI.OpenAIKey)
		log.Printf("已初始化 OpenAI 客户端")
	} else if provider == "deepseek" {
		if config.GlobalConfig.AI.DeepseekKey == "" {
			log.Fatal("Deepseek API key未设置")
		}
		deepseekClient = deepseek.NewClient(config.GlobalConfig.AI.DeepseekKey)
		log.Printf("已初始化 Deepseek 客户端")
	} else {
		log.Printf("警告：未知的 AI 提供商: %s", provider)
	}

	service := &AIService{
		openaiClient:   openaiClient,
		deepseekClient: deepseekClient,
		provider:       provider,
	}
	
	// 添加服务初始化完成的日志
	log.Printf("AIService初始化完成: %+v", service)
	
	return service
}

func (s *AIService) ParseFollowUpRecords(patientID string, maxRecords int) (string, error) {
	// 通过medical_storage获取随访记录
	medicalStorage := storage.GetMedicalStorage()
	followUpRecords, err := medicalStorage.GetFollowUpRecords(patientID)
	if err != nil {
		return "", err
	}
	// 提取最多maxRecords条随访记录，如果数量不足，只返回真实的数量
	if len(followUpRecords) < maxRecords {
		maxRecords = len(followUpRecords)
	}
	followUpRecords = followUpRecords[:maxRecords]
	followUpRecordsStr := ""
	for _, record := range followUpRecords {
		followUpRecordsStr += fmt.Sprintf("```随访记录: %s\n随访日期: %s\n", record.Title, record.FollowUpDate.Format("2006-01-02 15:04:05"))
		followUpRecordsStr += fmt.Sprintf("随访内容: %s\n", record.Content)
		followUpRecordsStr += "```\n"
	}
	return followUpRecordsStr, nil
}

func (s *AIService) ParseMedicalRecords(patientID string, maxRecords int) (string, error) {
	// 通过medical_storage获取诊疗记录和随访记录
	medicalStorage := storage.GetMedicalStorage()
	medicalRecords, err := medicalStorage.GetMedicalRecords(patientID)
	if err != nil {
		return "", err
	}
	// 提取最多maxRecords条诊疗记录，如果数量不足，只返回真实的数量
	if len(medicalRecords) < maxRecords {
		maxRecords = len(medicalRecords)
	}
	medicalRecords = medicalRecords[:maxRecords]
	medicalRecordsStr := ""
	for _, record := range medicalRecords {
		medicalRecordsStr += fmt.Sprintf("```诊疗记录\n诊断日期: %s\n", record.DiagnosisDate.Format("2006-01-02 15:04:05"))
		medicalRecordsStr += fmt.Sprintf("诊断结果: %s\n", record.Diagnosis)
		medicalRecordsStr += fmt.Sprintf("治疗方案: %s\n", record.Treatment)
		medicalRecordsStr += fmt.Sprintf("处方: %s\n", record.Prescription)
		medicalRecordsStr += fmt.Sprintf("备注: %s\n", record.Notes)
		medicalRecordsStr += "```\n"
	}
	return medicalRecordsStr, nil
}

// GenerateResponse 生成回复建议
func (s *AIService) GenerateResponse(patient *models.Patient, messageID string, currentMessage string, messageHistory []models.Message) (*models.AISuggestion, error) {
	medicalRecordsStr, err := s.ParseMedicalRecords(patient.ID, 5)
	if err != nil {
		return nil, err
	}
	followUpRecordsStr, err := s.ParseFollowUpRecords(patient.ID, 5)
	if err != nil {
		return nil, err
	}
	// 构建系统提示信息
	systemPrompt := fmt.Sprintf(config.MedicalAssistantSystemPrompt,
		patient.Name,
		patient.Gender,
		patient.Age,
		patient.BloodType,
		strings.Join(patient.Allergies, "、"),
		strings.Join(patient.ChronicDiseases, "、"),
		// 诊疗记录
		medicalRecordsStr,
		// 随访记录
		followUpRecordsStr,
	)

	// 构建对话历史上下文
	contextStr := buildContext(messageHistory)

	// 记录完整的提示信息
	log.Printf("\n=== AI 请求信息 ===\n")
	log.Printf("患者ID: %s\n", patient.ID)
	log.Printf("消息ID: %s\n", messageID)
	log.Printf("系统提示:\n%s\n", systemPrompt)
	log.Printf("历史对话:\n%s\n", contextStr)
	log.Printf("当前问题: %s\n", currentMessage)
	log.Printf("================\n")

	// 构建消息历史
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
	}

	// 添加历史对话
	recentMessages := messageHistory
	if len(messageHistory) > 5 {
		recentMessages = messageHistory[len(messageHistory)-5:]
	}

	for _, msg := range recentMessages {
		role := openai.ChatMessageRoleUser
		if msg.Role == models.MessageRoleDoctor {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 添加当前问题
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: currentMessage,
	})

	// 根据不同的供应商调用不同的 API
	var aiContent string
	var genErr error

	if s.provider == "openai" {
		aiContent, genErr = s.generateOpenAIResponse(messages)
	} else if s.provider == "deepseek" {
		aiContent, genErr = s.generateDeepseekResponse(messages)
	}
	log.Printf("s.provider: %s", s.provider)

	if genErr != nil {
		return nil, genErr
	}

	log.Printf("\n=== AI 响应 ===\n%s\n=============\n", aiContent)

	now := time.Now()
	// 创建 AI 建议记录
	suggestion := &models.AISuggestion{
		BaseModel: models.BaseModel{
			ID:        fmt.Sprintf("ai_%d", now.UnixNano()),
			CreatedAt: now,
			UpdatedAt: now,
		},
		MessageID:  messageID,
		PatientID:  patient.ID,
		Content:    aiContent,
		ModelUsed:  config.DefaultAIConfig.Model,
		Confidence: 0.95, // 默认置信度
		Category:   models.AISuggestionCategoryMedication,
		Priority:   3, // 默认优先级
		Status:     models.AISuggestionStatusPending,
	}

	return suggestion, nil
}

// 添加 OpenAI 响应生成方法
func (s *AIService) generateOpenAIResponse(messages []openai.ChatCompletionMessage) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       config.GlobalConfig.AI.Model,
			Messages:    messages,
			Temperature: config.DefaultAIConfig.Temperature,
			MaxTokens:   config.DefaultAIConfig.MaxTokens,
			TopP:        config.DefaultAIConfig.TopP,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// 添加 Deepseek 响应生成方法
func (s *AIService) generateDeepseekResponse(messages []openai.ChatCompletionMessage) (string, error) {
	// 添加日志，查看转换后的消息
	log.Printf("Deepseek请求消息: %+v", messages)
	
	deepseekMessages := make([]deepseek.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		role := deepseek.ChatMessageRoleUser
		switch msg.Role {
		case openai.ChatMessageRoleSystem:
			role = "system"
		case openai.ChatMessageRoleAssistant:
			role = "assistant"
		case openai.ChatMessageRoleUser:
			role = "user"
		}
		
		deepseekMessages[i] = deepseek.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 添加更详细的错误处理
	resp, err := s.deepseekClient.CreateChatCompletion(
		ctx,
		&deepseek.ChatCompletionRequest{
			Model:    deepseek.DeepSeekChat,
				Messages: deepseekMessages,
				Temperature: 0.7,  // 添加温度参数
				MaxTokens: 2000,   // 添加最大token限制
		},
	)

	if err != nil {
		log.Printf("Deepseek API错误: %v", err)
		return "", fmt.Errorf("deepseek API调用失败: %w", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("deepseek返回了空响应")
	}

	return resp.Choices[0].Message.Content, nil
}

// buildContext 构建上下文信息
func buildContext(history []models.Message) string {
	var contextBuilder strings.Builder

	// 只取最近的5条消息作为上下文
	recentMessages := history
	if len(history) > 5 {
		recentMessages = history[len(history)-5:]
	}

	for _, msg := range recentMessages {
		contextBuilder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
			msg.CreatedAt.Format("2006-01-02 15:04:05"),
			msg.Role,
			msg.Content,
		))
	}

	return contextBuilder.String()
}

// 从聊天记录中提取生理数据（血压，血糖）
// ExtractPhysiologicalData 从聊天记录中提取生理数据
func (s *AIService) ExtractPhysiologicalData(patientID string, messageContent string) error {
	// 获取当前时间
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 构建提示信息
	systemPrompt := fmt.Sprintf(`你是一个医疗数据分析助手。请从患者的消息中提取血压、血糖数据和测量时间
（请注意患者聊天记录中说明的时间，如果是相对时间需要和当前时间进行比较）。
如果存在多个数据，请提取最新的一组。请按以下JSON格式返回:
{
    "bloodPressure": {
        "systolic": 收缩压数值(int),
        "diastolic": 舒张压数值(int),
        "measuredAt": "测量时间(YYYY-MM-DD HH:mm:ss格式)",
        "hasData": true/false
    },
    "bloodSugar": {
        "value": 血糖值(float),
        "type": "空腹/餐后/随机",
        "measuredAt": "测量时间(YYYY-MM-DD HH:mm:ss格式)", 
        "hasData": true/false
    }
}
如果消息中未明确提到测量时间，则使用当前时间: %s。`, currentTime)

	messages := []deepseek.ChatCompletionMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: messageContent,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 调用AI提取数据
	resp, err := s.deepseekClient.CreateChatCompletion(
		ctx,
		&deepseek.ChatCompletionRequest{
			Model:       deepseek.DeepSeekChat,
			Messages:    messages,
			Temperature: 0.1, // 降低温度以获得更确定的结果
		},
	)

	if err != nil {
		return fmt.Errorf("AI提取生理数据失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("AI返回空响应")
	}

	// 解析AI返回的JSON
	var result struct {
		BloodPressure struct {
			Systolic   int     `json:"systolic"`
			Diastolic  int     `json:"diastolic"`
			MeasuredAt string  `json:"measuredAt"`
			HasData    bool    `json:"hasData"`
		} `json:"bloodPressure"`
		BloodSugar struct {
			Value      float64 `json:"value"`
			Type       string  `json:"type"`
			MeasuredAt string  `json:"measuredAt"`
			HasData    bool    `json:"hasData"`
		} `json:"bloodSugar"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return fmt.Errorf("解析AI响应失败: %w", err)
	}

	// 保存有效的血压数据
	if result.BloodPressure.HasData {
		measuredAt, err := time.Parse("2006-01-02 15:04:05", result.BloodPressure.MeasuredAt)
		if err != nil {
			measuredAt = time.Now() // 如果解析失败，使用当前时间
		}

		bloodPressure := &models.PhysiologicalData{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			PatientID:  patientID,
			Type:      "blood_pressure",
			Value:     fmt.Sprintf("%d/%d", result.BloodPressure.Systolic, result.BloodPressure.Diastolic),
			MeasuredAt: measuredAt,
			Source:    "ai_extract",
			Notes:     "从聊天记录中AI提取的血压数据",
		}
		if err := storage.GetPhysiologicalDataStorage().Create(bloodPressure); err != nil {
			return fmt.Errorf("保存血压数据失败: %w", err)
		}
	}

	// 保存有效的血糖数据
	if result.BloodSugar.HasData {
		measuredAt, err := time.Parse("2006-01-02 15:04:05", result.BloodSugar.MeasuredAt)
		if err != nil {
			measuredAt = time.Now() // 如果解析失败，使用当前时间
		}

		bloodSugar := &models.PhysiologicalData{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			PatientID:  patientID,
			Type:      "blood_sugar",
			Value:     fmt.Sprintf("%.1f-%s", result.BloodSugar.Value, result.BloodSugar.Type),
			MeasuredAt: measuredAt,
			Source:    "ai_extract",
			Notes:     "从聊天记录中AI提取的血糖数据",
		}
		if err := storage.GetPhysiologicalDataStorage().Create(bloodSugar); err != nil {
			return fmt.Errorf("保存血糖数据失败: %w", err)
		}
	}

	return nil
}
