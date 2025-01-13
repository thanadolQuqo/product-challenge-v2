package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"product-challenge/pkg/config"
)

func NewRedisCache(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		Username: cfg.Redis.Username,
		DB:       cfg.Redis.DB,
		//Protocol: cfg.Redis.Protocol,
	})

	status, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("Redis connection was refused")
	}

	fmt.Println("Redis status: ", status)
	return client, nil
}
