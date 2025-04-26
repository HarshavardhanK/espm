package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/HarshavardhanK/espm/internal/config"

	"github.com/redis/go-redis/v9"
)

// redisCache implements the RedisCache interface
type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(cfg config.RedisConfig) (RedisCache, error) {

	client := redis.NewClient(&redis.Options{

		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		MaxRetries: cfg.MaxRetries,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisCache{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {

	data, err := r.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		return nil, ErrCacheMiss
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	return data, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte) error {

	if err := r.client.Set(ctx, key, value, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete implements RedisCache.Delete
func (r *redisCache) Delete(ctx context.Context, key string) error {

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}

// GetEventStream implements RedisCache.GetEventStream
func (r *redisCache) GetEventStream(ctx context.Context, aggregateType, aggregateID string) ([]byte, error) {

	key := fmt.Sprintf("events:%s:%s", aggregateType, aggregateID)

	return r.Get(ctx, key)
}

func (r *redisCache) SetEventStream(ctx context.Context, aggregateType, aggregateID string, value []byte) error {

	key := fmt.Sprintf("events:%s:%s", aggregateType, aggregateID)

	return r.Set(ctx, key, value)
}

func (r *redisCache) HealthCheck(ctx context.Context) error {

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
