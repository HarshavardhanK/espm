package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/HarshavardhanK/espm/internal/cache"
	"github.com/HarshavardhanK/espm/internal/repository"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventStore implements repository.EventStore interface for testing
type MockEventStore struct {
	mock.Mock
}

func (m *MockEventStore) AppendEvents(ctx context.Context, events []repository.Event) error {
	args := m.Called(ctx, events)
	return args.Error(0)
}

func (m *MockEventStore) GetEventsByAggregateID(ctx context.Context, aggregateType, aggregateID string) ([]repository.Event, error) {
	args := m.Called(ctx, aggregateType, aggregateID)
	return args.Get(0).([]repository.Event), args.Error(1)
}

func (m *MockEventStore) GetEventsByType(ctx context.Context, eventType string) ([]repository.Event, error) {
	args := m.Called(ctx, eventType)
	return args.Get(0).([]repository.Event), args.Error(1)
}

func (m *MockEventStore) GetEventsAfterSequence(ctx context.Context, sequence int64) ([]repository.Event, error) {
	args := m.Called(ctx, sequence)
	return args.Get(0).([]repository.Event), args.Error(1)
}

// MockRedisCache implements cache.RedisCache interface for testing
type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockRedisCache) Set(ctx context.Context, key string, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockRedisCache) GetEventStream(ctx context.Context, aggregateType, aggregateID string) ([]byte, error) {
	args := m.Called(ctx, aggregateType, aggregateID)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockRedisCache) SetEventStream(ctx context.Context, aggregateType, aggregateID string, data []byte) error {
	args := m.Called(ctx, aggregateType, aggregateID, data)
	return args.Error(0)
}

func (m *MockRedisCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRedisCache) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisCache) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestCachedEventStore_AppendEvents(t *testing.T) {

	ctx := context.Background()

	mockStore := new(MockEventStore)
	mockCache := new(MockRedisCache)

	cachedStore := repository.NewCachedEventStore(mockStore, mockCache, time.Hour)

	// Create test events
	events := []repository.Event{
		{
			ID:            uuid.New().String(),
			AggregateType: "Order",
			AggregateID:   uuid.New().String(),
			EventType:     "OrderCreated",
			Version:       1,
			Sequence:      1,
			Data:          []byte(`{"orderId": "123"}`),
			Metadata:      map[string]interface{}{"source": "test"},
			CreatedAt:     time.Now(),
		},
	}

	// Set up expectations
	mockStore.On("AppendEvents", ctx, events).Return(nil)
	mockCache.On("Delete", ctx, mock.Anything).Return(nil)

	// Test append events
	err := cachedStore.AppendEvents(ctx, events)
	assert.NoError(t, err)

	// Verify expectations
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedEventStore_GetEventsByAggregateID_CacheHit(t *testing.T) {

	ctx := context.Background()

	mockStore := new(MockEventStore)
	mockCache := new(MockRedisCache)

	cachedStore := repository.NewCachedEventStore(mockStore, mockCache, time.Hour)

	// Create test data
	aggregateType := "Order"
	aggregateID := uuid.New().String()

	now := time.Now()
	cachedEvents := []repository.Event{

		{
			ID:            uuid.New().String(),
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			EventType:     "OrderCreated",
			Version:       1,
			Sequence:      1,
			Data:          []byte(`{"orderId": "123"}`),
			Metadata:      map[string]interface{}{"source": "test"},
			CreatedAt:     now,
		},
	}

	// Set up expectations for cache hit
	cachedData, _ := json.Marshal(cachedEvents)
	mockCache.On("GetEventStream", ctx, aggregateType, aggregateID).Return(cachedData, nil)

	// Test get events with cache hit
	events, err := cachedStore.GetEventsByAggregateID(ctx, aggregateType, aggregateID)
	assert.NoError(t, err)

	// Compare all fields except CreatedAt
	assert.Equal(t, len(cachedEvents), len(events))

	for i := range events {

		assert.Equal(t, cachedEvents[i].ID, events[i].ID)
		assert.Equal(t, cachedEvents[i].AggregateType, events[i].AggregateType)
		assert.Equal(t, cachedEvents[i].AggregateID, events[i].AggregateID)
		assert.Equal(t, cachedEvents[i].EventType, events[i].EventType)
		assert.Equal(t, cachedEvents[i].Version, events[i].Version)
		assert.Equal(t, cachedEvents[i].Sequence, events[i].Sequence)
		assert.Equal(t, cachedEvents[i].Data, events[i].Data)
		assert.Equal(t, cachedEvents[i].Metadata, events[i].Metadata)
	}

	// Verify expectations
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedEventStore_GetEventsByAggregateID_CacheMiss(t *testing.T) {

	ctx := context.Background()

	mockStore := new(MockEventStore)
	mockCache := new(MockRedisCache)

	cachedStore := repository.NewCachedEventStore(mockStore, mockCache, time.Hour)

	// Create test data
	aggregateType := "Order"
	aggregateID := uuid.New().String()

	now := time.Now()
	storeEvents := []repository.Event{

		{
			ID:            uuid.New().String(),
			AggregateType: aggregateType,
			AggregateID:   aggregateID,
			EventType:     "OrderCreated",
			Version:       1,
			Sequence:      1,
			Data:          []byte(`{"orderId": "123"}`),
			Metadata:      map[string]interface{}{"source": "test"},
			CreatedAt:     now,
		},
	}

	// Set up expectations for cache miss
	mockCache.On("GetEventStream", ctx, aggregateType, aggregateID).Return([]byte{}, cache.ErrCacheMiss)
	mockStore.On("GetEventsByAggregateID", ctx, aggregateType, aggregateID).Return(storeEvents, nil)
	mockCache.On("SetEventStream", ctx, aggregateType, aggregateID, mock.Anything).Return(nil)

	// Test get events with cache miss
	events, err := cachedStore.GetEventsByAggregateID(ctx, aggregateType, aggregateID)
	assert.NoError(t, err)

	// Compare all fields except CreatedAt
	assert.Equal(t, len(storeEvents), len(events))

	for i := range events {

		assert.Equal(t, storeEvents[i].ID, events[i].ID)
		assert.Equal(t, storeEvents[i].AggregateType, events[i].AggregateType)
		assert.Equal(t, storeEvents[i].AggregateID, events[i].AggregateID)
		assert.Equal(t, storeEvents[i].EventType, events[i].EventType)
		assert.Equal(t, storeEvents[i].Version, events[i].Version)
		assert.Equal(t, storeEvents[i].Sequence, events[i].Sequence)
		assert.Equal(t, storeEvents[i].Data, events[i].Data)
		assert.Equal(t, storeEvents[i].Metadata, events[i].Metadata)
	}

	// Verify expectations
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedEventStore_GetEventsByType(t *testing.T) {

	ctx := context.Background()

	mockStore := new(MockEventStore)
	mockCache := new(MockRedisCache)

	cachedStore := repository.NewCachedEventStore(mockStore, mockCache, time.Hour)

	// Create test data
	eventType := "OrderCreated"
	now := time.Now()
	events := []repository.Event{
		{
			ID:            uuid.New().String(),
			AggregateType: "Order",
			AggregateID:   uuid.New().String(),
			EventType:     eventType,
			Version:       1,
			Sequence:      1,
			Data:          []byte(`{"orderId": "123"}`),
			Metadata:      map[string]interface{}{"source": "test"},
			CreatedAt:     now,
		},
	}

	// Set up expectations
	mockStore.On("GetEventsByType", ctx, eventType).Return(events, nil)

	// Test get events by type
	result, err := cachedStore.GetEventsByType(ctx, eventType)
	assert.NoError(t, err)

	// Compare all fields except CreatedAt
	assert.Equal(t, len(events), len(result))

	for i := range result {
		assert.Equal(t, events[i].ID, result[i].ID)
		assert.Equal(t, events[i].AggregateType, result[i].AggregateType)
		assert.Equal(t, events[i].AggregateID, result[i].AggregateID)
		assert.Equal(t, events[i].EventType, result[i].EventType)
		assert.Equal(t, events[i].Version, result[i].Version)
		assert.Equal(t, events[i].Sequence, result[i].Sequence)
		assert.Equal(t, events[i].Data, result[i].Data)
		assert.Equal(t, events[i].Metadata, result[i].Metadata)
	}

	// Verify expectations
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCachedEventStore_GetEventsAfterSequence(t *testing.T) {
	ctx := context.Background()
	mockStore := new(MockEventStore)
	mockCache := new(MockRedisCache)
	cachedStore := repository.NewCachedEventStore(mockStore, mockCache, time.Hour)

	// Create test data
	sequence := int64(1)
	now := time.Now()
	events := []repository.Event{
		{
			ID:            uuid.New().String(),
			AggregateType: "Order",
			AggregateID:   uuid.New().String(),
			EventType:     "OrderCreated",
			Version:       1,
			Sequence:      2,
			Data:          []byte(`{"orderId": "123"}`),
			Metadata:      map[string]interface{}{"source": "test"},
			CreatedAt:     now,
		},
	}

	// Set up expectations
	mockStore.On("GetEventsAfterSequence", ctx, sequence).Return(events, nil)

	// Test get events after sequence
	result, err := cachedStore.GetEventsAfterSequence(ctx, sequence)
	assert.NoError(t, err)

	// Compare all fields except CreatedAt
	assert.Equal(t, len(events), len(result))
	for i := range result {
		assert.Equal(t, events[i].ID, result[i].ID)
		assert.Equal(t, events[i].AggregateType, result[i].AggregateType)
		assert.Equal(t, events[i].AggregateID, result[i].AggregateID)
		assert.Equal(t, events[i].EventType, result[i].EventType)
		assert.Equal(t, events[i].Version, result[i].Version)
		assert.Equal(t, events[i].Sequence, result[i].Sequence)
		assert.Equal(t, events[i].Data, result[i].Data)
		assert.Equal(t, events[i].Metadata, result[i].Metadata)
	}

	// Verify expectations
	mockStore.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
