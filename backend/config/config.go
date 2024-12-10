package config

import (
	"os"
)

type Config struct {
	OpenAIAPIKey string
	OpenAIModel  string
}

var GlobalConfig Config

func Init() {
	GlobalConfig = Config{
		OpenAIAPIKey: getEnvOrDefault("OPENAI_API_KEY", "sk-proj-efqvozACpCMXtMeEBWm9T3BlbkFJcQ0YiNNK1GSk5Iil7Dyg"),
		OpenAIModel:  getEnvOrDefault("OPENAI_MODEL", "gpt-3.5-turbo"),
	}

	if GlobalConfig.OpenAIAPIKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
