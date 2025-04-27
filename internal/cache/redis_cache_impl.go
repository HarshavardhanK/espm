package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/HarshavardhanK/espm/internal/config"

	"github.com/redis/go-redis/v9"
)

// Simple wrapper around Redis client
type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// Creates a new Redis client with the given config
func NewRedisCache(cfg config.RedisConfig) (RedisCache, error) {

	// Enhanced connection pooling configuration

	client := redis.NewClient(&redis.Options{

		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),

		Password: cfg.Password,
		DB:       cfg.DB,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,

		MaxIdleConns: cfg.PoolSize, // Match pool size for optimal performance
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		MaxRetries:  cfg.MaxRetries,
		PoolTimeout: time.Second * 30, // Increased pool timeout for high concurrency
	})

	// Quick ping to verify connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisCache{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

// Get a value from Redis
func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, ErrInvalidKey
	}

	data, err := r.client.Get(ctx, key).Bytes()

	if err == redis.Nil {
		return nil, ErrCacheMiss
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	return data, nil
}

// Store a value in Redis
func (r *redisCache) Set(ctx context.Context, key string, value []byte) error {
	if key == "" {
		return ErrInvalidKey
	}

	if err := r.client.Set(ctx, key, value, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// BatchSet stores multiple key-value pairs in Redis
func (r *redisCache) BatchSet(ctx context.Context, pairs map[string][]byte) error {
	if len(pairs) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()
	for key, value := range pairs {
		if key == "" {
			continue
		}
		pipe.Set(ctx, key, value, r.ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to batch set cache: %w", err)
	}

	return nil
}

// BatchGet retrieves multiple values from Redis
func (r *redisCache) BatchGet(ctx context.Context, keys []string) (map[string][]byte, error) {
	if len(keys) == 0 {
		return make(map[string][]byte), nil
	}

	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, key := range keys {
		if key == "" {
			continue
		}
		cmds[key] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to batch get from cache: %w", err)
	}

	result := make(map[string][]byte)

	for key, cmd := range cmds {

		val, err := cmd.Bytes()

		if err == redis.Nil {
			continue
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		result[key] = val
	}

	return result, nil
}

func (r *redisCache) Delete(ctx context.Context, key string) error {

	if key == "" {
		return ErrInvalidKey
	}

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}

// BatchDelete removes multiple keys from Redis
func (r *redisCache) BatchDelete(ctx context.Context, keys []string) error {

	if len(keys) == 0 {
		return nil
	}

	pipe := r.client.Pipeline()

	for _, key := range keys {

		if key == "" {
			continue
		}

		pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to batch delete from cache: %w", err)
	}

	return nil
}

// Get event stream for an aggregate
func (r *redisCache) GetEventStream(ctx context.Context, aggregateType, aggregateID string) ([]byte, error) {

	if aggregateType == "" || aggregateID == "" {
		return nil, ErrInvalidKey
	}

	key := fmt.Sprintf("events:%s:%s", aggregateType, aggregateID)

	return r.Get(ctx, key)
}

// Store event stream for an aggregate
func (r *redisCache) SetEventStream(ctx context.Context, aggregateType, aggregateID string, value []byte) error {

	if aggregateType == "" || aggregateID == "" {
		return ErrInvalidKey
	}

	key := fmt.Sprintf("events:%s:%s", aggregateType, aggregateID)

	return r.Set(ctx, key, value)
}

// BatchSetEventStreams stores multiple event streams
func (r *redisCache) BatchSetEventStreams(ctx context.Context, streams map[string]map[string][]byte) error {

	if len(streams) == 0 {
		return nil
	}

	pairs := make(map[string][]byte)

	for aggregateType, typeStreams := range streams {

		for aggregateID, value := range typeStreams {
			key := fmt.Sprintf("events:%s:%s", aggregateType, aggregateID)
			pairs[key] = value
		}
	}

	return r.BatchSet(ctx, pairs)
}

// Check if Redis is responding
func (r *redisCache) HealthCheck(ctx context.Context) error {

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	return nil
}

// Clean up Redis connection
func (r *redisCache) Close() error {
	return r.client.Close()
}
