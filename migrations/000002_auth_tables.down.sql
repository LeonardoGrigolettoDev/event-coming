-- Drop constraints
ALTER TABLE events DROP CONSTRAINT IF EXISTS events_created_by_fkey;

-- Drop triggers
DROP TRIGGER IF EXISTS user_organizations_updated_at ON user_organizations;
DROP TRIGGER IF EXISTS users_updated_at ON users;

-- Drop tables
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS user_organizations CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop types
DROP TYPE IF EXISTS user_role;
