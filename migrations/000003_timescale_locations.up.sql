-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Locations table optimized for time-series data
CREATE TABLE locations (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    participant_id UUID NOT NULL,
    event_id UUID NOT NULL,
    instance_id UUID,
    organization_id UUID NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    geometry GEOMETRY(POINT, 4326),
    accuracy DOUBLE PRECISION,
    altitude DOUBLE PRECISION,
    speed DOUBLE PRECISION,
    heading DOUBLE PRECISION,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT locations_participant_fkey FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT locations_event_fkey FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    CONSTRAINT locations_instance_fkey FOREIGN KEY (instance_id) REFERENCES event_instances(id) ON DELETE CASCADE,
    CONSTRAINT locations_organization_fkey FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- Convert to hypertable (partitioned by time)
SELECT create_hypertable('locations', 'time', if_not_exists => TRUE);

-- Create indexes for better query performance
CREATE INDEX idx_locations_participant ON locations(participant_id, time DESC);
CREATE INDEX idx_locations_event ON locations(event_id, time DESC);
CREATE INDEX idx_locations_organization ON locations(organization_id, time DESC);
CREATE INDEX idx_locations_geometry ON locations USING GIST(geometry);
CREATE INDEX idx_locations_time ON locations(time DESC);

-- Create trigger for geometry column
CREATE OR REPLACE FUNCTION update_location_geometry()
RETURNS TRIGGER AS $$
BEGIN
    NEW.geometry := ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER locations_geometry_trigger
    BEFORE INSERT OR UPDATE OF latitude, longitude ON locations
    FOR EACH ROW
    EXECUTE FUNCTION update_location_geometry();

-- Add compression policy (compress data older than 7 days)
ALTER TABLE locations SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'participant_id, event_id, organization_id',
    timescaledb.compress_orderby = 'time DESC'
);

SELECT add_compression_policy('locations', INTERVAL '7 days', if_not_exists => TRUE);

-- Add retention policy (retain data for 90 days)
SELECT add_retention_policy('locations', INTERVAL '90 days', if_not_exists => TRUE);

-- Create continuous aggregate for hourly location summaries
CREATE MATERIALIZED VIEW locations_hourly
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 hour', time) AS bucket,
    participant_id,
    event_id,
    organization_id,
    COUNT(*) as location_count,
    AVG(latitude) as avg_latitude,
    AVG(longitude) as avg_longitude,
    AVG(speed) as avg_speed,
    MAX(time) as last_update
FROM locations
GROUP BY bucket, participant_id, event_id, organization_id
WITH NO DATA;

-- Add refresh policy for continuous aggregate
SELECT add_continuous_aggregate_policy('locations_hourly',
    start_offset => INTERVAL '3 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => TRUE);
