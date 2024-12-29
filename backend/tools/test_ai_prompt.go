package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"we-dear/config"

	deepseek "github.com/cohesion-org/deepseek-go"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化 Deepseek 客户端
	client := deepseek.NewClient(config.GlobalConfig.AI.DeepseekKey)

	// 构建测试消息
	messages := []deepseek.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "你是一位专业的医生，请用专业且易懂的语言回答病人的问题。",
		},
		{
			Role:    "user",
			Content: "我最近血压一直在波动，早上测了收缩压160，舒张压95，我现在在服用络活喜，需要调整用药吗？",
		},
	}

	// 创建聊天请求
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		&deepseek.ChatCompletionRequest{
			Model:    deepseek.DeepSeekChat,
			Messages: messages,
		},
	)

	if err != nil {
		log.Fatalf("Deepseek API 调用失败: %v", err)
	}

	// 打印响应
	fmt.Printf("\n=== Deepseek 测试结果 ===\n")
	fmt.Printf("问题: %s\n", messages[1].Content)
	fmt.Printf("\n回复:\n%s\n", resp.Choices[0].Message.Content)
	fmt.Printf("===========================\n")
}
