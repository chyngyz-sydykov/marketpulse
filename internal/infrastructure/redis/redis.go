package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chyngyz-sydykov/marketpulse/config"
	"github.com/redis/go-redis/v9"
)

// Global Redis client
var Redis *redis.Client

// Event structure
type Event struct {
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

type RedisServiceInterface interface {
	PublishEvent(ctx context.Context, eventName, source string) error
	SubscribeToEvent(ctx context.Context, eventName string, handler func(event Event))
}

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{client: client}
}

// PublishEvent publishes an event to Redis Pub/Sub
func (r *RedisService) PublishEvent(ctx context.Context, eventName, source string) error {
	event := Event{
		Name:      eventName,
		Source:    source,
		Timestamp: time.Now(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = r.client.Publish(ctx, eventName, eventData).Err()
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}

func (r *RedisService) SubscribeToEvent(ctx context.Context, eventName string, handler func(event Event)) {
	pubsub := r.client.Subscribe(ctx, eventName)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		var event Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Error decoding event: %v\n", err)
			continue
		}

		// Handle event
		handler(event)
	}
}

func ConnectRedis() error {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: cfg.RedisPassword,
		DB:       0, // Default DB
	})

	// Test Redis connection
	ctx := context.Background()
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		return err
	}

	Redis = redisClient
	return nil
}
