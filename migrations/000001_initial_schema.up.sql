-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";

-- Create custom types
CREATE TYPE organization_type AS ENUM ('school', 'enterprise', 'event');
CREATE TYPE subscription_plan AS ENUM ('free', 'basic', 'professional', 'enterprise');
CREATE TYPE event_type AS ENUM ('demand', 'periodic');
CREATE TYPE event_status AS ENUM ('draft', 'scheduled', 'active', 'completed', 'cancelled');
CREATE TYPE participant_status AS ENUM ('pending', 'confirmed', 'denied', 'checked_in', 'no_show');
CREATE TYPE scheduler_action AS ENUM ('confirmation', 'reminder', 'closure', 'location');
CREATE TYPE scheduler_status AS ENUM ('pending', 'processed', 'failed', 'skipped');
CREATE TYPE message_template_type AS ENUM ('confirmation', 'reminder', 'location_request', 'closure', 'welcome');

-- Organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    type organization_type NOT NULL,
    subscription_plan subscription_plan NOT NULL DEFAULT 'free',
    max_events INTEGER NOT NULL DEFAULT 10,
    max_participants INTEGER NOT NULL DEFAULT 100,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT organizations_name_unique UNIQUE (name)
);

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    type event_type NOT NULL,
    status event_status NOT NULL DEFAULT 'draft',
    location_lat DOUBLE PRECISION NOT NULL,
    location_lng DOUBLE PRECISION NOT NULL,
    location_address TEXT,
    location_geometry GEOMETRY(POINT, 4326),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    rrule_string VARCHAR(500),
    confirmation_deadline TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Event instances table (for recurring events)
CREATE TABLE event_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    instance_date DATE NOT NULL,
    status event_status NOT NULL DEFAULT 'scheduled',
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT event_instances_unique UNIQUE (event_id, instance_date)
);

-- Participants table
CREATE TABLE participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    instance_id UUID REFERENCES event_instances(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    status participant_status NOT NULL DEFAULT 'pending',
    confirmed_at TIMESTAMP WITH TIME ZONE,
    checked_in_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Schedulers table
CREATE TABLE schedulers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    instance_id UUID REFERENCES event_instances(id) ON DELETE CASCADE,
    action scheduler_action NOT NULL,
    status scheduler_status NOT NULL DEFAULT 'pending',
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE,
    retries INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Message templates table
CREATE TABLE message_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type message_template_type NOT NULL,
    template_name VARCHAR(100) NOT NULL,
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    parameters JSONB,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT message_templates_unique UNIQUE (organization_id, type, language)
);

-- Message logs table
CREATE TABLE message_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    participant_id UUID REFERENCES participants(id) ON DELETE SET NULL,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    instance_id UUID REFERENCES event_instances(id) ON DELETE SET NULL,
    template_id UUID REFERENCES message_templates(id) ON DELETE SET NULL,
    phone_number VARCHAR(20) NOT NULL,
    message_type message_template_type NOT NULL,
    message_content TEXT,
    status VARCHAR(20) NOT NULL,
    external_id VARCHAR(100),
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    read_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_events_organization ON events(organization_id);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_events_created_by ON events(created_by);
CREATE INDEX idx_events_location_geometry ON events USING GIST(location_geometry);

CREATE INDEX idx_event_instances_event ON event_instances(event_id);
CREATE INDEX idx_event_instances_organization ON event_instances(organization_id);
CREATE INDEX idx_event_instances_date ON event_instances(instance_date);
CREATE INDEX idx_event_instances_status ON event_instances(status);

CREATE INDEX idx_participants_event ON participants(event_id);
CREATE INDEX idx_participants_instance ON participants(instance_id);
CREATE INDEX idx_participants_organization ON participants(organization_id);
CREATE INDEX idx_participants_phone ON participants(phone_number);
CREATE INDEX idx_participants_status ON participants(status);

CREATE INDEX idx_schedulers_organization ON schedulers(organization_id);
CREATE INDEX idx_schedulers_event ON schedulers(event_id);
CREATE INDEX idx_schedulers_status_scheduled ON schedulers(status, scheduled_at) WHERE status = 'pending';
CREATE INDEX idx_schedulers_action ON schedulers(action);

CREATE INDEX idx_message_logs_organization ON message_logs(organization_id);
CREATE INDEX idx_message_logs_participant ON message_logs(participant_id);
CREATE INDEX idx_message_logs_event ON message_logs(event_id);
CREATE INDEX idx_message_logs_created_at ON message_logs(created_at);
CREATE INDEX idx_message_logs_status ON message_logs(status);

-- Create trigger for location_geometry
CREATE OR REPLACE FUNCTION update_event_location_geometry()
RETURNS TRIGGER AS $$
BEGIN
    NEW.location_geometry := ST_SetSRID(ST_MakePoint(NEW.location_lng, NEW.location_lat), 4326);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER events_location_geometry_trigger
    BEFORE INSERT OR UPDATE OF location_lat, location_lng ON events
    FOR EACH ROW
    EXECUTE FUNCTION update_event_location_geometry();

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER event_instances_updated_at BEFORE UPDATE ON event_instances
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER participants_updated_at BEFORE UPDATE ON participants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER schedulers_updated_at BEFORE UPDATE ON schedulers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER message_templates_updated_at BEFORE UPDATE ON message_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
