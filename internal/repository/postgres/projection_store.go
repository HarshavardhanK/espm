package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/HarshavardhanK/espm/internal/repository"
)

// PostgresProjectionStore implements the ProjectionStore interface using PostgreSQL
type PostgresProjectionStore struct {
	db *sql.DB
}

// NewPostgresProjectionStore creates a new PostgresProjectionStore
func NewPostgresProjectionStore(db *sql.DB) *PostgresProjectionStore {
	return &PostgresProjectionStore{db: db}
}

// SaveProjectionState implements the ProjectionStore interface
func (s *PostgresProjectionStore) SaveProjectionState(
	ctx context.Context,
	projectionName string,
	state interface{},
) error {
	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO projections (
			projection_name, state, updated_at
		) VALUES ($1, $2, $3)
		ON CONFLICT (projection_name) 
		DO UPDATE SET state = $2, updated_at = $3
	`, projectionName, jsonState, time.Now())

	return err
}

// GetProjectionState implements the ProjectionStore interface
func (s *PostgresProjectionStore) GetProjectionState(
	ctx context.Context,
	projectionName string,
	state interface{},
) error {
	var jsonState []byte

	err := s.db.QueryRowContext(ctx, `
		SELECT state
		FROM projections
		WHERE projection_name = $1
	`, projectionName).Scan(&jsonState)

	if err == sql.ErrNoRows {
		return repository.ErrProjectionNotFound
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonState, state)
}

// UpdateProjectionState implements the ProjectionStore interface
func (s *PostgresProjectionStore) UpdateProjectionState(
	ctx context.Context,
	projectionName string,
	state interface{},
) error {
	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
	}

	result, err := s.db.ExecContext(ctx, `
		UPDATE projections
		SET state = $1, updated_at = $2
		WHERE projection_name = $3
	`, jsonState, time.Now(), projectionName)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return repository.ErrProjectionNotFound
	}

	return nil
}
