package redis

import (
	"context"
	"fmt"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func ConnectRedis() error {

	cfg := config.LoadConfig()
	dsn := fmt.Sprintf(
		"%s:%s",
		cfg.RedisHost, cfg.RedisPort,
	)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     dsn, // Service name as per docker-compose
		Password: "",  // No password by default
		DB:       0,   // Default DB
	})

	// Test Redis connection
	ctx := context.Background()
	err := redisClient.Set(ctx, "greeting", "Hello, Redis!", 0).Err()
	if err != nil {
		return err
	}

	_, err = redisClient.Get(ctx, "greeting").Result()
	if err != nil {
		return err
	}
	Redis = redisClient
	return nil
}
