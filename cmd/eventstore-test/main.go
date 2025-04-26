package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/espm/internal/repository"
	"github.com/yourusername/espm/internal/repository/postgres"
)

func main() {
	// Database connection string
	connStr := "postgres://espm:espm123@localhost:5432/espm?sslmode=disable"

	// Create a new PostgreSQL store
	store, err := postgres.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create a test event
	event := repository.Event{
		ID:            uuid.New().String(),
		AggregateType: "Order",
		AggregateID:   uuid.New().String(),
		EventType:     "OrderCreated",
		Version:       1,
		Sequence:      1,
		Data:          []byte(`{"orderId": "123", "customerId": "456"}`),
		Metadata: map[string]interface{}{
			"source":    "test",
			"timestamp": time.Now().Unix(),
		},
		CreatedAt: time.Now(),
	}

	// Append the event
	ctx := context.Background()
	err = store.AppendEvents(ctx, []repository.Event{event})
	if err != nil {
		log.Fatalf("Failed to append event: %v", err)
	}

	fmt.Println("Successfully appended event")

	// Retrieve events by aggregate ID
	events, err := store.GetEventsByAggregateID(ctx, "Order", event.AggregateID)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	fmt.Println("\nRetrieved events:")
	for _, e := range events {
		fmt.Printf("Event ID: %s\n", e.ID)
		fmt.Printf("Aggregate Type: %s\n", e.AggregateType)
		fmt.Printf("Aggregate ID: %s\n", e.AggregateID)
		fmt.Printf("Event Type: %s\n", e.EventType)
		fmt.Printf("Version: %d\n", e.Version)
		fmt.Printf("Sequence: %d\n", e.Sequence)
		fmt.Printf("Data: %s\n", string(e.Data))
		metadata, _ := json.MarshalIndent(e.Metadata, "", "  ")
		fmt.Printf("Metadata: %s\n", string(metadata))
		fmt.Printf("Created At: %s\n", e.CreatedAt)
		fmt.Println("---")
	}
}
