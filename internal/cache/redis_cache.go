package cache

import (
	"context"
	"errors"
)

// Cache errors
var (
	ErrCacheMiss  = errors.New("cache miss")
	ErrInvalidKey = errors.New("invalid cache key")
)

// RedisCache defines the interface for Redis caching operations
type RedisCache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache
	Set(ctx context.Context, key string, value []byte) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// GetEventStream retrieves an event stream from cache
	GetEventStream(ctx context.Context, aggregateType, aggregateID string) ([]byte, error)

	// SetEventStream stores an event stream in cache
	SetEventStream(ctx context.Context, aggregateType, aggregateID string, value []byte) error

	// HealthCheck verifies Redis connection
	HealthCheck(ctx context.Context) error

	// Close closes the Redis connection
	Close() error
}
