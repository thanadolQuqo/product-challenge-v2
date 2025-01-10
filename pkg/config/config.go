package config

import (
	"github.com/joho/godotenv"
	"os"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type Config struct {
	Database DBConfig `json:"database"`
}

func LoadConfig() (*Config, error) {
	// load config from .env file
	godotenv.Load()

	config := &Config{
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "26257"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "defaultdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
