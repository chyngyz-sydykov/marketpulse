package app

import (
	"github.com/chyngyz-sydykov/marketpulse/internal/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/redis"
)

var App *Container

// Container struct holds all dependencies
type Container struct {
	RedisService      redis.RedisServiceInterface
	MarketDataService *marketdata.MarketDataService
}

// NewContainer initializes and returns all services
func NewContainer() *Container {
	redisService := redis.NewRedisService(redis.Redis)
	marketDataService := marketdata.NewMarketDataService(redisService)
	App = &Container{
		RedisService:      redisService,
		MarketDataService: marketDataService,
	}
	return App
}
