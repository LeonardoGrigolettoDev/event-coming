-- Drop triggers
DROP TRIGGER IF EXISTS message_templates_updated_at ON message_templates;
DROP TRIGGER IF EXISTS schedulers_updated_at ON schedulers;
DROP TRIGGER IF EXISTS participants_updated_at ON participants;
DROP TRIGGER IF EXISTS event_instances_updated_at ON event_instances;
DROP TRIGGER IF EXISTS events_updated_at ON events;
DROP TRIGGER IF EXISTS events_location_geometry_trigger ON events;
DROP TRIGGER IF EXISTS organizations_updated_at ON organizations;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS update_event_location_geometry();

-- Drop tables
DROP TABLE IF EXISTS message_logs CASCADE;
DROP TABLE IF EXISTS message_templates CASCADE;
DROP TABLE IF EXISTS schedulers CASCADE;
DROP TABLE IF EXISTS participants CASCADE;
DROP TABLE IF EXISTS event_instances CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;

-- Drop types
DROP TYPE IF EXISTS message_template_type;
DROP TYPE IF EXISTS scheduler_status;
DROP TYPE IF EXISTS scheduler_action;
DROP TYPE IF EXISTS participant_status;
DROP TYPE IF EXISTS event_status;
DROP TYPE IF EXISTS event_type;
DROP TYPE IF EXISTS subscription_plan;
DROP TYPE IF EXISTS organization_type;
