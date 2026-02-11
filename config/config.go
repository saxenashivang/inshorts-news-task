package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	LLMProvider string
	OpenAIKey   string
	GeminiKey   string
	ClaudeKey   string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "user"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "newsdb"),
		LLMProvider: getEnv("LLM_PROVIDER", "gemini"), // openai, gemini, claude
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		GeminiKey:   getEnv("GEMINI_API_KEY", ""),
		ClaudeKey:   getEnv("CLAUDE_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
