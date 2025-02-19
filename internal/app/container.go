package app

import (
	indicator "github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
)

var App *Container

// Container struct holds all dependencies
type Container struct {
	RedisService      redis.RedisServiceInterface
	MarketDataService *marketdata.MarketDataService
	IndicatorService  *indicator.IndicatorService
}

// NewContainer initializes and returns all services
func NewContainer() *Container {
	redisService := redis.NewRedisService(redis.Redis)
	marketDataService := marketdata.NewMarketDataService(redisService)
	indicatorService := indicator.NewIndicatorService()
	App = &Container{
		RedisService:      redisService,
		MarketDataService: marketDataService,
		IndicatorService:  indicatorService,
	}
	return App
}
