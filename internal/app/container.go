package app

import (
	"github.com/chyngyz-sydykov/marketpulse/internal/app/event"
	indicator "github.com/chyngyz-sydykov/marketpulse/internal/core/indicator"
	"github.com/chyngyz-sydykov/marketpulse/internal/core/marketdata"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/grpc"
	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
)

var App *Container

// Container struct holds all dependencies
type Container struct {
	RedisService      redis.RedisServiceInterface
	MarketDataService *marketdata.MarketDataService
	IndicatorService  *indicator.IndicatorService
	EventListener     *event.EventListener
	GrpcServer        *grpc.GrpcServer
}

// NewContainer initializes and returns all services
func NewContainer() *Container {
	redisService := redis.NewRedisService(redis.Redis)
	marketDataService := marketdata.NewMarketDataService(redisService)
	indicatorService := indicator.NewIndicatorService(redisService)
	GrpcServcer := grpc.NewGrpcService(marketDataService, indicatorService)

	EventListener := event.NewEventListener(
		marketDataService,
		indicatorService,
		redisService,
	)

	App = &Container{
		RedisService:      redisService,
		MarketDataService: marketDataService,
		IndicatorService:  indicatorService,
		EventListener:     EventListener,
		GrpcServer:        GrpcServcer,
	}
	return App
}
