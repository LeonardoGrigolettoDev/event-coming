-- Drop continuous aggregate
DROP MATERIALIZED VIEW IF EXISTS locations_hourly CASCADE;

-- Drop trigger and function
DROP TRIGGER IF EXISTS locations_geometry_trigger ON locations;
DROP FUNCTION IF EXISTS update_location_geometry();

-- Drop hypertable (this will also drop the table)
DROP TABLE IF EXISTS locations CASCADE;
