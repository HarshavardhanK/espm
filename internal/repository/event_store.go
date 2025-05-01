package repository

import (
	"context"

	"github.com/HarshavardhanK/espm/internal/events"
	"github.com/google/uuid"
)

// EventStore defines the interface for event storage
type EventStore interface {
	// AppendEvents appends new events to the store
	AppendEvents(ctx context.Context, events []events.Event) error

	// GetEventsByAggregateID retrieves all events for a specific aggregate
	GetEventsByAggregateID(ctx context.Context, aggregateType string, aggregateID uuid.UUID) ([]events.Event, error)

	// GetEventsByType retrieves all events of a specific type
	GetEventsByType(ctx context.Context, eventType events.EventType) ([]events.Event, error)

	// GetEventsAfterSequence retrieves all events after a specific sequence number
	GetEventsAfterSequence(ctx context.Context, sequence int64) ([]events.Event, error)
}

// SnapshotStore defines the interface for aggregate snapshots
type SnapshotStore interface {
	// SaveSnapshot saves a snapshot of an aggregate
	SaveSnapshot(ctx context.Context, aggregateType string, aggregateID uuid.UUID, version int, data []byte) error

	// GetLatestSnapshot retrieves the latest snapshot for an aggregate
	GetLatestSnapshot(ctx context.Context, aggregateType string, aggregateID uuid.UUID) ([]byte, int, error)
}

// ProjectionStore defines the interface for projection state storage
type ProjectionStore interface {
	// SaveProjectionState saves the state of a projection
	SaveProjectionState(ctx context.Context, projectionType string, data []byte) error

	// GetProjectionState retrieves the state of a projection
	GetProjectionState(ctx context.Context, projectionType string) ([]byte, error)

	// UpdateProjectionStatus updates the status of a projection
	UpdateProjectionStatus(ctx context.Context, projectionType string, status string) error
}
