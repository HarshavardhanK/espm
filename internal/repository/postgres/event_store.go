package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/HarshavardhanK/espm/internal/events"
	"github.com/google/uuid"
)

// PostgresEventStore implements the EventStore interface using PostgreSQL
type PostgresEventStore struct {
	db *sql.DB
}

// NewPostgresEventStore creates a new PostgresEventStore
func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

// AppendEvents implements the EventStore interface
func (s *PostgresEventStore) AppendEvents(ctx context.Context, events []events.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO events (
			event_id, aggregate_type, aggregate_id, event_type,
			event_version, sequence_number, data, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, event := range events {
		_, err = stmt.ExecContext(
			ctx,
			event.EventID,
			event.AggregateType,
			event.AggregateID,
			event.EventType,
			event.EventVersion,
			event.Sequence,
			event.Data,
			event.Metadata,
			event.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetEventsByAggregateID implements the EventStore interface
func (s *PostgresEventStore) GetEventsByAggregateID(
	ctx context.Context,
	aggregateType string,
	aggregateID uuid.UUID,
) ([]events.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT event_id, aggregate_type, aggregate_id, event_type,
		       event_version, sequence_number, data, metadata, created_at
		FROM events
		WHERE aggregate_type = $1 AND aggregate_id = $2
		ORDER BY sequence_number ASC
	`, aggregateType, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []events.Event
	for rows.Next() {
		var event events.Event
		var metadataJSON []byte
		err := rows.Scan(
			&event.EventID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.EventVersion,
			&event.Sequence,
			&event.Data,
			&metadataJSON,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &event.Metadata); err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, rows.Err()
}

// GetEventsByType implements the EventStore interface
func (s *PostgresEventStore) GetEventsByType(
	ctx context.Context,
	eventType events.EventType,
) ([]events.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT event_id, aggregate_type, aggregate_id, event_type,
		       event_version, sequence_number, data, metadata, created_at
		FROM events
		WHERE event_type = $1
		ORDER BY sequence_number ASC
	`, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []events.Event
	for rows.Next() {
		var event events.Event
		var metadataJSON []byte
		err := rows.Scan(
			&event.EventID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.EventVersion,
			&event.Sequence,
			&event.Data,
			&metadataJSON,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &event.Metadata); err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, rows.Err()
}

// GetEventsAfterSequence implements the EventStore interface
func (s *PostgresEventStore) GetEventsAfterSequence(
	ctx context.Context,
	sequence int64,
) ([]events.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT event_id, aggregate_type, aggregate_id, event_type,
		       event_version, sequence_number, data, metadata, created_at
		FROM events
		WHERE sequence_number > $1
		ORDER BY sequence_number ASC
	`, sequence)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []events.Event
	for rows.Next() {
		var event events.Event
		var metadataJSON []byte
		err := rows.Scan(
			&event.EventID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.EventVersion,
			&event.Sequence,
			&event.Data,
			&metadataJSON,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(metadataJSON, &event.Metadata); err != nil {
			return nil, err
		}

		result = append(result, event)
	}

	return result, rows.Err()
}
