package config

import (
	"time"
)

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxRetries   int
	TTL          time.Duration
}

// DefaultRedisConfig returns default Redis configuration
func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Host:         "localhost",
		Port:         6379,
		Password:     "",
		DB:           0,
		PoolSize:     10, // Increased from default
		MinIdleConns: 5,  // Keep some connections ready
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		MaxRetries:   3,
		TTL:          time.Hour * 24,
	}
}
