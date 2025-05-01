package repository

import "errors"

var (
	// ErrSnapshotNotFound is returned when a snapshot is not found
	ErrSnapshotNotFound = errors.New("snapshot not found")
	// ErrProjectionNotFound is returned when a projection is not found
	ErrProjectionNotFound = errors.New("projection not found")
)
