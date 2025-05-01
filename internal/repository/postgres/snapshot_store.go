package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/HarshavardhanK/espm/internal/repository"
	"github.com/google/uuid"
)

// PostgresSnapshotStore implements the SnapshotStore interface using PostgreSQL
type PostgresSnapshotStore struct {
	db *sql.DB
}

// NewPostgresSnapshotStore creates a new PostgresSnapshotStore
func NewPostgresSnapshotStore(db *sql.DB) *PostgresSnapshotStore {
	return &PostgresSnapshotStore{db: db}
}

// SaveSnapshot implements the SnapshotStore interface
func (s *PostgresSnapshotStore) SaveSnapshot(
	ctx context.Context,
	aggregateType string,
	aggregateID uuid.UUID,
	version int64,
	data interface{},
) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO snapshots (
			aggregate_type, aggregate_id, version, data, created_at
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (aggregate_type, aggregate_id) 
		DO UPDATE SET version = $3, data = $4, created_at = $5
	`, aggregateType, aggregateID, version, jsonData, time.Now())

	return err
}

// GetSnapshot implements the SnapshotStore interface
func (s *PostgresSnapshotStore) GetSnapshot(
	ctx context.Context,
	aggregateType string,
	aggregateID uuid.UUID,
	data interface{},
) (int64, error) {
	var version int64
	var jsonData []byte

	err := s.db.QueryRowContext(ctx, `
		SELECT version, data
		FROM snapshots
		WHERE aggregate_type = $1 AND aggregate_id = $2
	`, aggregateType, aggregateID).Scan(&version, &jsonData)

	if err == sql.ErrNoRows {
		return 0, repository.ErrSnapshotNotFound
	}
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(jsonData, data); err != nil {
		return 0, err
	}

	return version, nil
}
