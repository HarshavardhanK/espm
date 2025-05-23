package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HarshavardhanK/espm/internal/cache"
)

// CachedEventStore wraps an EventStore with Redis caching
type CachedEventStore struct {
	store EventStore
	cache cache.RedisCache
	ttl   time.Duration
}

// NewCachedEventStore creates a new cached event store
func NewCachedEventStore(store EventStore, redisCache cache.RedisCache, ttl time.Duration) *CachedEventStore {

	return &CachedEventStore{

		store: store,
		cache: redisCache,
		ttl:   ttl,
	}
}

// AppendEvents implements EventStore.AppendEvents with caching
func (c *CachedEventStore) AppendEvents(ctx context.Context, events []Event) error {

	// First append to the store
	if err := c.store.AppendEvents(ctx, events); err != nil {
		return fmt.Errorf("failed to append events: %w", err)
	}

	// Group events by aggregate for batch invalidation
	aggregateKeys := make([]string, 0, len(events))

	for _, event := range events {
		key := fmt.Sprintf("events:%s:%s", event.AggregateType, event.AggregateID)
		aggregateKeys = append(aggregateKeys, key)
	}

	// Batch delete cache entries
	if err := c.cache.BatchDelete(ctx, aggregateKeys); err != nil {
		fmt.Printf("Warning: failed to invalidate cache for %d aggregates: %v\n", len(aggregateKeys), err)
	}

	return nil
}

// GetEventsByAggregateID implements EventStore.GetEventsByAggregateID with caching
func (c *CachedEventStore) GetEventsByAggregateID(ctx context.Context, aggregateType, aggregateID string) ([]Event, error) {
	// Try to get from cache first
	cached, err := c.cache.GetEventStream(ctx, aggregateType, aggregateID)

	if err == nil {

		var events []Event

		if err := json.Unmarshal(cached, &events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached events: %w", err)
		}

		return events, nil
	}

	// If not in cache, get from store
	events, err := c.store.GetEventsByAggregateID(ctx, aggregateType, aggregateID)

	if err != nil {
		return nil, fmt.Errorf("failed to get events from store: %w", err)
	}

	// Cache the result
	data, err := json.Marshal(events)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal events: %w", err)
	}

	// Use batch operation for potential future optimization
	streams := map[string]map[string][]byte{

		aggregateType: {
			aggregateID: data,
		},
	}

	if err := c.cache.BatchSetEventStreams(ctx, streams); err != nil {

		fmt.Printf("Warning: failed to cache events: %v\n", err)
	}

	return events, nil
}

// GetEventsByType implements EventStore.GetEventsByType
func (c *CachedEventStore) GetEventsByType(ctx context.Context, eventType string) ([]Event, error) {

	// This operation is not cached as it's not frequently used
	return c.store.GetEventsByType(ctx, eventType)
}

// GetEventsAfterSequence implements EventStore.GetEventsAfterSequence
func (c *CachedEventStore) GetEventsAfterSequence(ctx context.Context, sequence int64) ([]Event, error) {

	return c.store.GetEventsAfterSequence(ctx, sequence)
}
