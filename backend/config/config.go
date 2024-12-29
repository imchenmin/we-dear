package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DatabaseConfig
	AI AIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var GlobalConfig Config

func Init() {
	// 添加更详细的日志
	log.Printf("开始加载配置...")
	
	// 尝试加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
		// 尝试从其他位置加载
		err = godotenv.Load("../backend/.env")
		if err != nil {
			log.Printf("Warning: ../backend/.env not found: %v", err)
		}
	}

	// 打印所有环境变量
	log.Printf("AI_PROVIDER=%s", os.Getenv("AI_PROVIDER"))
	log.Printf("DEEPSEEK_API_KEY length=%d", len(os.Getenv("DEEPSEEK_API_KEY")))

	GlobalConfig = Config{
		DB: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", ""),
			Name:     getEnvOrDefault("DB_NAME", "wedear"),
		},
		AI: AIConfig{
			Provider:    getEnvOrDefault("AI_PROVIDER", DefaultAIConfig.Provider),
			OpenAIKey:   getEnvOrDefault("OPENAI_API_KEY", ""),
			DeepseekKey: getEnvOrDefault("DEEPSEEK_API_KEY", ""),
			Model:       getEnvOrDefault("AI_MODEL", DefaultAIConfig.Model),
			Temperature: DefaultAIConfig.Temperature,
			MaxTokens:   DefaultAIConfig.MaxTokens,
			TopP:        DefaultAIConfig.TopP,
		},
	}

	// 打印加载后的配置
	log.Printf("配置加载完成:")
	log.Printf("AI Provider: %s", GlobalConfig.AI.Provider)
	log.Printf("AI Model: %s", GlobalConfig.AI.Model)
	log.Printf("Deepseek Key length: %d", len(GlobalConfig.AI.DeepseekKey))

	// 验证配置
	if GlobalConfig.AI.Provider == "openai" && GlobalConfig.AI.OpenAIKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}
	if GlobalConfig.AI.Provider == "deepseek" && GlobalConfig.AI.DeepseekKey == "" {
		panic("DEEPSEEK_API_KEY environment variable is not set")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
