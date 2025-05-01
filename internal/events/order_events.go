package events

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of an event
type EventType string

const (
	OrderCreatedEventType     EventType = "OrderCreated"
	OrderItemAddedEventType   EventType = "OrderItemAdded"
	OrderItemRemovedEventType EventType = "OrderItemRemoved"
	OrderSubmittedEventType   EventType = "OrderSubmitted"
	OrderCancelledEventType   EventType = "OrderCancelled"
)

// Event represents the base event structure
type Event struct {
	EventID       uuid.UUID
	AggregateType string
	AggregateID   uuid.UUID
	EventType     EventType
	EventVersion  int
	Sequence      int64
	Data          []byte
	Metadata      map[string]interface{}
	CreatedAt     time.Time
}

// OrderCreatedEvent represents the event when an order is created
type OrderCreatedEvent struct {
	CustomerID uuid.UUID
	CreatedAt  time.Time
}

// OrderItemAddedEvent represents the event when an item is added to an order
type OrderItemAddedEvent struct {
	ProductID uuid.UUID
	Quantity  int
	UnitPrice float64
}

// OrderItemRemovedEvent represents the event when an item is removed from an order
type OrderItemRemovedEvent struct {
	ProductID uuid.UUID
}

// OrderSubmittedEvent represents the event when an order is submitted
type OrderSubmittedEvent struct {
	SubmittedAt time.Time
}

// OrderCancelledEvent represents the event when an order is cancelled
type OrderCancelledEvent struct {
	CancelledAt time.Time
	Reason      string
}

// NewEvent creates a new event with the given parameters
func NewEvent(
	aggregateType string,
	aggregateID uuid.UUID,
	eventType EventType,
	eventVersion int,
	sequence int64,
	data []byte,
	metadata map[string]interface{},
) Event {
	return Event{
		EventID:       uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		EventVersion:  eventVersion,
		Sequence:      sequence,
		Data:          data,
		Metadata:      metadata,
		CreatedAt:     time.Now(),
	}
}
