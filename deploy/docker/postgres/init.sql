-- Create events table
CREATE TABLE IF NOT EXISTS events (
    event_id UUID PRIMARY KEY,
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    event_version INTEGER NOT NULL,
    sequence_number BIGINT NOT NULL,
    data JSONB NOT NULL,
    metadata JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE (aggregate_type, aggregate_id, sequence_number)
);

-- Create snapshots table
CREATE TABLE IF NOT EXISTS snapshots (
    aggregate_type VARCHAR(255) NOT NULL,
    aggregate_id UUID NOT NULL,
    version BIGINT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (aggregate_type, aggregate_id)
);

-- Create projections table
CREATE TABLE IF NOT EXISTS projections (
    projection_name VARCHAR(255) PRIMARY KEY,
    state JSONB NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_events_aggregate ON events (aggregate_type, aggregate_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON events (event_type);
CREATE INDEX IF NOT EXISTS idx_events_sequence ON events (sequence_number); 