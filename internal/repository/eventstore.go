package repository

import (
	"context"
	"time"
)

// Event represents a stored event in the event store
type Event struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Version       int
	Sequence      int64
	Data          []byte
	Metadata      map[string]interface{}
	CreatedAt     time.Time
}

// Snapshot represents a stored aggregate snapshot
type Snapshot struct {
	ID            string
	AggregateType string
	AggregateID   string
	Version       int
	Data          []byte
	Metadata      map[string]interface{}
	CreatedAt     time.Time
}

// Projection represents a stored projection state
type Projection struct {
	ID               string
	Type             string
	LastProcessedSeq int64
	Status           string
	UpdatedAt        time.Time
}

// EventStore defines the interface for event storage operations
type EventStore interface {
	// AppendEvents appends new events to the event store
	AppendEvents(ctx context.Context, events []Event) error

	// GetEventsByAggregateID retrieves all events for a specific aggregate
	GetEventsByAggregateID(ctx context.Context, aggregateType, aggregateID string) ([]Event, error)

	// GetEventsByType retrieves all events of a specific type
	GetEventsByType(ctx context.Context, eventType string) ([]Event, error)

	// GetEventsAfterSequence retrieves events after a specific sequence number
	GetEventsAfterSequence(ctx context.Context, sequence int64) ([]Event, error)
}

// SnapshotStore defines the interface for snapshot storage operations
type SnapshotStore interface {
	// SaveSnapshot saves a new snapshot
	SaveSnapshot(ctx context.Context, snapshot Snapshot) error

	// GetLatestSnapshot retrieves the most recent snapshot for an aggregate
	GetLatestSnapshot(ctx context.Context, aggregateType, aggregateID string) (*Snapshot, error)
}

// ProjectionStore defines the interface for projection storage operations
type ProjectionStore interface {
	// SaveProjectionState saves the current state of a projection
	SaveProjectionState(ctx context.Context, projection Projection) error

	// GetProjectionState retrieves the current state of a projection
	GetProjectionState(ctx context.Context, projectionType string) (*Projection, error)

	// UpdateProjectionStatus updates the status of a projection
	UpdateProjectionStatus(ctx context.Context, projectionID, status string) error
}
