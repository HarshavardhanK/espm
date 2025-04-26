package main

import (
	"context"
	"encoding/json"

	"fmt"
	"log"

	"github.com/HarshavardhanK/espm/internal/cache"
	"github.com/HarshavardhanK/espm/internal/config"
)

func main() {

	cfg := config.DefaultRedisConfig()

	//Create Redis cache
	redisCache, err := cache.NewRedisCache(cfg)

	if err != nil {
		log.Fatalf("Failed to create Redis cache: %v", err)
	}

	defer redisCache.Close()

	// Test basic operations
	testBasicOperations(redisCache)

	// Test event stream operations
	testEventStreamOperations(redisCache)

	// Test error handling
	testErrorHandling(redisCache)

	// Test health check
	testHealthCheck(redisCache)
}

func testBasicOperations(redisCache cache.RedisCache) {
	ctx := context.Background()
	key := "test:key"
	value := []byte("test value")

	// Test Set
	fmt.Println("\nTesting Set operation...")

	if err := redisCache.Set(ctx, key, value); err != nil {
		log.Printf("Set failed: %v", err)
		return
	}

	fmt.Println("Set successful")

	// Test Get
	fmt.Println("\nTesting Get operation.")
	retrieved, err := redisCache.Get(ctx, key)

	if err != nil {
		log.Printf("Get failed: %v", err)
		return
	}

	fmt.Printf("Get successful: %s\n", string(retrieved))

	// Test Delete
	fmt.Println("\nTesting Delete operation...")

	if err := redisCache.Delete(ctx, key); err != nil {
		log.Printf("Delete failed: %v", err)
		return
	}

	fmt.Println("Delete successful")
}

func testEventStreamOperations(redisCache cache.RedisCache) {

	ctx := context.Background()

	aggregateType := "Order"
	aggregateID := "123"

	// Create test event data
	events := []map[string]interface{}{
		{"type": "OrderCreated", "data": map[string]interface{}{"orderId": "123"}},
		{"type": "OrderUpdated", "data": map[string]interface{}{"status": "processing"}},
	}

	value, err := json.Marshal(events)

	if err != nil {
		log.Printf("Failed to marshal events: %v", err)
		return
	}

	// Test SetEventStream
	fmt.Println("\nTesting SetEventStream operation")

	if err := redisCache.SetEventStream(ctx, aggregateType, aggregateID, value); err != nil {
		log.Printf("SetEventStream failed: %v", err)
		return
	}

	fmt.Println("SetEventStream successful")

	// Test GetEventStream
	fmt.Println("\nTesting GetEventStream operation.")

	retrieved, err := redisCache.GetEventStream(ctx, aggregateType, aggregateID)

	if err != nil {
		log.Printf("GetEventStream failed: %v", err)
		return
	}

	fmt.Printf("GetEventStream successful: %s\n", string(retrieved))
}

func testErrorHandling(redisCache cache.RedisCache) {

	ctx := context.Background()

	// Test cache miss
	fmt.Println("\nTesting cache miss...")

	_, err := redisCache.Get(ctx, "nonexistent:key")

	if err == cache.ErrCacheMiss {
		fmt.Println("✓ Cache miss handled correctly")

	} else {
		log.Printf("Unexpected error: %v", err)
	}

	// Test invalid key
	fmt.Println("\nTesting invalid key.")

	err = redisCache.Set(ctx, "", []byte("value"))

	if err != nil {
		fmt.Println("✓ Invalid key handled correctly")

	} else {
		log.Println("Invalid key not handled")
	}
}

func testHealthCheck(redisCache cache.RedisCache) {
	ctx := context.Background()

	fmt.Println("\nTesting health check...")

	if err := redisCache.HealthCheck(ctx); err != nil {
		log.Printf("Health check failed: %v", err)
		return
	}

	fmt.Println("✓ Health check successful")
}
