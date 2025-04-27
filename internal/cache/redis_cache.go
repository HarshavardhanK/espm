package cache

import (
	"context"
	"errors"
)

// Common cache errors
var (
	ErrCacheMiss  = errors.New("cache miss")
	ErrInvalidKey = errors.New("invalid cache key")
)

// Redis cache operations
type RedisCache interface {

	// Basic cache operations
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error

	// Batch operations
	BatchGet(ctx context.Context, keys []string) (map[string][]byte, error)
	BatchSet(ctx context.Context, pairs map[string][]byte) error
	BatchDelete(ctx context.Context, keys []string) error

	// Event stream specific operations
	GetEventStream(ctx context.Context, aggregateType, aggregateID string) ([]byte, error)
	SetEventStream(ctx context.Context, aggregateType, aggregateID string, value []byte) error
	BatchSetEventStreams(ctx context.Context, streams map[string]map[string][]byte) error

	// Health and maintenance
	HealthCheck(ctx context.Context) error
	Close() error
}
