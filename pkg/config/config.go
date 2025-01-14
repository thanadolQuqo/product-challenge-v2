package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type AWSConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Region          string `json:"region"`
	BucketName      string `json:"bucket_name"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	//Protocol int    `json:"protocol"`
}

type Config struct {
	Database  DBConfig    `json:"database"`
	Aws       AWSConfig   `json:"aws"`
	Redis     RedisConfig `json:"redis"`
	JwtSecret string      `json:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
	// load config from .env file
	godotenv.Load()

	redisDB, err := strconv.ParseInt(getEnv("REDIS_DB", "0"), 10, 64)
	if err != nil {
		fmt.Println("Unable to load redis config: ", err)
		return nil, err
	}
	config := &Config{
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "26257"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "defaultdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Aws: AWSConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_REGION", ""),
			BucketName:      getEnv("AWS_BUCKET_NAME", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Username: getEnv("REDIS_USER", "default"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       int(redisDB),
			//Protocol: 2,
		},
		JwtSecret: getEnv("JWT_SECRET", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
