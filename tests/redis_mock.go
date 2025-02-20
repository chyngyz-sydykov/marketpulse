package main

import (
	"context"

	"github.com/chyngyz-sydykov/marketpulse/internal/infrastructure/redis"
	"github.com/stretchr/testify/mock"
)

// MockRedisService implements RedisServiceInterface for testing
type MockRedisService struct {
	mock.Mock
}

// PublishEvent is a mock implementation that does nothing
func (m *MockRedisService) PublishEvent(ctx context.Context, eventName, source string) error {
	args := m.Called(ctx, eventName, source)
	return args.Error(0)
}

// SubscribeToEvent is a mock implementation that does nothing
func (m *MockRedisService) SubscribeToEvent(ctx context.Context, eventName string, handler func(event redis.Event)) {
	m.Called(ctx, eventName, handler)
}
