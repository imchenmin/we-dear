package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"we-dear/config"
	"we-dear/models"

	openai "github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
}

func NewAIService() *AIService {
	client := openai.NewClient("sk-proj-efqvozACpCMXtMeEBWm9T3BlbkFJcQ0YiNNK1GSk5Iil7Dyg")
	fmt.Println(config.GlobalConfig.OpenAIKey)
	return &AIService{
		client: client,
	}
}

// GenerateResponse 使用 OpenAI 生成回复建议
func (s *AIService) GenerateResponse(patient *models.Patient, messageID string, currentMessage string, messageHistory []models.Message) (*models.AISuggestion, error) {
	// 构建系统提示信息
	systemPrompt := fmt.Sprintf(config.MedicalAssistantSystemPrompt,
		patient.Name,
		patient.Gender,
		patient.Age,
		patient.BloodType,
		strings.Join(patient.Allergies, "、"),
		strings.Join(patient.ChronicDiseases, "、"),
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

	// 创建请求上下文（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 发送请求
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       config.DefaultAIConfig.Model,
			Messages:    messages,
			Temperature: config.DefaultAIConfig.Temperature,
			MaxTokens:   config.DefaultAIConfig.MaxTokens,
			TopP:        config.DefaultAIConfig.TopP,
		},
	)

	if err != nil {
		log.Printf("OpenAI API 错误: %v\n", err)
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Printf("AI 未返回响应\n")
		return nil, fmt.Errorf("no response from AI")
	}

	aiContent := resp.Choices[0].Message.Content
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
