-- Drop indexes
DROP INDEX IF EXISTS idx_events_aggregate;
DROP INDEX IF EXISTS idx_events_sequence;
DROP INDEX IF EXISTS idx_snapshots_aggregate;
DROP INDEX IF EXISTS idx_projections_type;

-- Drop tables
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS projections; 