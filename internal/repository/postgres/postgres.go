package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HarshavardhanK/espm/internal/repository"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// PostgresStore implements the repository interfaces using PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store instance
func NewPostgresStore(connStr string) (*PostgresStore, error) {

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// Close closes the database connection
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// AppendEvents implements the EventStore interface
func (s *PostgresStore) AppendEvents(ctx context.Context, events []repository.Event) error {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	for _, event := range events {

		metadata, err := json.Marshal(event.Metadata)

		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		query := `
			INSERT INTO events (
				event_id, aggregate_type, aggregate_id, event_type,
				event_version, sequence_number, data, metadata, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

		_, err = tx.ExecContext(ctx, query, event.ID, event.AggregateType, event.AggregateID, event.EventType,
			event.Version, event.Sequence, event.Data, metadata, event.CreatedAt)

		if err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}
	}

	return tx.Commit()
}

// GetEventsByAggregateID implements the EventStore interface
func (s *PostgresStore) GetEventsByAggregateID(ctx context.Context, aggregateType, aggregateID string) ([]repository.Event, error) {

	query := `
		SELECT event_id, aggregate_type, aggregate_id, event_type,
			   event_version, sequence_number, data, metadata, created_at
		FROM events
		WHERE aggregate_type = $1 AND aggregate_id = $2
		ORDER BY sequence_number ASC
	`
	rows, err := s.db.QueryContext(ctx, query, aggregateType, aggregateID)

	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []repository.Event

	for rows.Next() {

		var event repository.Event
		var metadata []byte

		err := rows.Scan(
			&event.ID, &event.AggregateType, &event.AggregateID, &event.EventType,
			&event.Version, &event.Sequence, &event.Data, &metadata, &event.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if err := json.Unmarshal(metadata, &event.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

// GetEventsByType implements the EventStore interface
func (s *PostgresStore) GetEventsByType(ctx context.Context, eventType string) ([]repository.Event, error) {

	query := `
		SELECT event_id, aggregate_type, aggregate_id, event_type,
			   event_version, sequence_number, data, metadata, created_at
		FROM events
		WHERE event_type = $1
		ORDER BY sequence_number ASC
	`

	rows, err := s.db.QueryContext(ctx, query, eventType)

	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []repository.Event

	for rows.Next() {

		var event repository.Event
		var metadata []byte

		err := rows.Scan(
			&event.ID, &event.AggregateType, &event.AggregateID, &event.EventType,
			&event.Version, &event.Sequence, &event.Data, &metadata, &event.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if err := json.Unmarshal(metadata, &event.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

// GetEventsAfterSequence implements the EventStore interface
func (s *PostgresStore) GetEventsAfterSequence(ctx context.Context, sequence int64) ([]repository.Event, error) {

	query := `
		SELECT event_id, aggregate_type, aggregate_id, event_type,
			   event_version, sequence_number, data, metadata, created_at
		FROM events
		WHERE sequence_number > $1
		ORDER BY sequence_number ASC
	`
	rows, err := s.db.QueryContext(ctx, query, sequence)

	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}

	defer rows.Close()

	var events []repository.Event

	for rows.Next() {

		var event repository.Event

		var metadata []byte

		err := rows.Scan(
			&event.ID, &event.AggregateType, &event.AggregateID, &event.EventType,
			&event.Version, &event.Sequence, &event.Data, &metadata, &event.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if err := json.Unmarshal(metadata, &event.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

// SaveSnapshot implements the SnapshotStore interface
func (s *PostgresStore) SaveSnapshot(ctx context.Context, snapshot repository.Snapshot) error {

	metadata, err := json.Marshal(snapshot.Metadata)

	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO snapshots (
			snapshot_id, aggregate_type, aggregate_id, aggregate_version,
			data, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = s.db.ExecContext(ctx, query, snapshot.ID, snapshot.AggregateType, snapshot.AggregateID,
		snapshot.Version, snapshot.Data, metadata, snapshot.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	return nil
}

// GetLatestSnapshot implements the SnapshotStore interface
func (s *PostgresStore) GetLatestSnapshot(ctx context.Context, aggregateType, aggregateID string) (*repository.Snapshot, error) {

	var snapshot repository.Snapshot
	var metadata []byte

	query := `
		SELECT snapshot_id, aggregate_type, aggregate_id, aggregate_version,
			   data, metadata, created_at
		FROM snapshots
		WHERE aggregate_type = $1 AND aggregate_id = $2
		ORDER BY aggregate_version DESC
		LIMIT 1
	`

	err := s.db.QueryRowContext(ctx, query, aggregateType, aggregateID).Scan(
		&snapshot.ID, &snapshot.AggregateType, &snapshot.AggregateID,
		&snapshot.Version, &snapshot.Data, &metadata, &snapshot.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	if err := json.Unmarshal(metadata, &snapshot.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &snapshot, nil
}

// SaveProjectionState implements the ProjectionStore interface
func (s *PostgresStore) SaveProjectionState(ctx context.Context, projection repository.Projection) error {

	query := `
		INSERT INTO projections (
			projection_id, projection_type, last_processed_event,
			status, updated_at
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (projection_id) DO UPDATE SET
			last_processed_event = EXCLUDED.last_processed_event,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`
	_, err := s.db.ExecContext(ctx, query, projection.ID, projection.Type, projection.LastProcessedSeq,
		projection.Status, projection.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to save projection state: %w", err)
	}

	return nil
}

// GetProjectionState implements the ProjectionStore interface
func (s *PostgresStore) GetProjectionState(ctx context.Context, projectionType string) (*repository.Projection, error) {
	var projection repository.Projection

	query := `
		SELECT projection_id, projection_type, last_processed_event,
			   status, updated_at
		FROM projections
		WHERE projection_type = $1
	`

	err := s.db.QueryRowContext(ctx, query, projectionType).Scan(
		&projection.ID, &projection.Type, &projection.LastProcessedSeq,
		&projection.Status, &projection.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			// Create a new projection if it doesn't exist
			projection = repository.Projection{
				ID:               uuid.New().String(),
				Type:             projectionType,
				LastProcessedSeq: 0,
				Status:           "active",
				UpdatedAt:        time.Now(),
			}

			if err := s.SaveProjectionState(ctx, projection); err != nil {
				return nil, err
			}

			return &projection, nil
		}

		return nil, fmt.Errorf("failed to get projection state: %w", err)
	}

	return &projection, nil
}

// UpdateProjectionStatus implements the ProjectionStore interface
func (s *PostgresStore) UpdateProjectionStatus(ctx context.Context, projectionID, status string) error {

	query := `
		UPDATE projections
		SET status = $1, updated_at = $2
		WHERE projection_id = $3
	`
	_, err := s.db.ExecContext(ctx, query, status, time.Now(), projectionID)

	if err != nil {
		return fmt.Errorf("failed to update projection status: %w", err)
	}

	return nil
}
