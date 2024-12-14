package config

import (
	"os"
)

type Config struct {
	OpenAIKey   string
	OpenAIModel string
	DB          DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var GlobalConfig = Config{
	OpenAIKey:   "your-api-key",
	OpenAIModel: "gpt-3.5-turbo",
	DB: DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "your-password",
		Name:     "wedear",
	},
}

func Init() {
	GlobalConfig = Config{
		OpenAIKey:   "sk-proj-efqvozACpCMXtMeEBWm9T3BlbkFJcQ0YiNNK1GSk5Iil7Dyg",
		OpenAIModel: "gpt-3.5-turbo",
		DB: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "pccmdxy1998",
			Name:     "wedear",
		},
	}

	if GlobalConfig.OpenAIKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
