-- initial schema for PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
    IF NOT EXISTS contacts (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        name text NOT NULL,
        phone_number text UNIQUE NOT NULL,
        created_at timestamp
        WITH
            TIME ZONE DEFAULT now (),
            updated_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );

CREATE TABLE
    IF NOT EXISTS users (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        email text UNIQUE NOT NULL,
        hashed_password text NOT NULL,
        contact_id uuid NOT NULL REFERENCES contacts (id) ON DELETE CASCADE,
        created_at timestamp
        WITH
            TIME ZONE DEFAULT now (),
            updated_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );

CREATE TABLE
    IF NOT EXISTS schedulers (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        name text NOT NULL,
        description text,
        schedule_type text NOT NULL,
        cron_expression text,
        start_date timestamp
        WITH
            TIME ZONE,
            end_date timestamp
        WITH
            TIME ZONE,
            created_at timestamp
        WITH
            TIME ZONE DEFAULT now (),
            updated_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );

CREATE TABLE
    IF NOT EXISTS scheduler_contacts (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        scheduler_id uuid NOT NULL REFERENCES schedulers (id) ON DELETE CASCADE,
        contact_id uuid NOT NULL REFERENCES contacts (id) ON DELETE CASCADE,
        created_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );

CREATE TABLE
    IF NOT EXISTS places (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        location_description text NOT NULL,
        longitude float4 NOT NULL,
        latitude float4 NOT NULL,
        altitude float4 NOT NULL,
        updated_at timestamp
        WITH
            TIME ZONE DEFAULT now () created_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );

CREATE TYPE IF NOT EXISTS event_type AS ENUM (
    'notification',
    'confirmation',
    'location',
    'checked'
);

CREATE TYPE IF NOT EXISTS event_status AS ENUM (
    'pending',
    'in_progress',
    'success',
    'skipped',
    'error'
);

CREATE TABLE
    IF NOT EXISTS events (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        scheduler_id uuid NOT NULL REFERENCES schedulers (id) ON DELETE CASCADE,
        contact_id uuid NOT NULL REFERENCES contacts (id) ON DELETE CASCADE,
        event_type text NOT NULL,
        status text NOT NULL,
        payload jsonb,
        created_at timestamp
        WITH
            TIME ZONE DEFAULT now (),
            updated_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );

CREATE TYPE IF NOT EXISTS final_status AS ENUM (
    'success',
    'skipped',
    'error',
    'timeout',
    'no_reply'
);

CREATE TABLE
    IF NOT EXISTS consolidated (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        scheduler_id uuid NOT NULL REFERENCES schedulers (id) ON DELETE CASCADE,
        contact_id uuid NOT NULL REFERENCES contacts (id) ON DELETE CASCADE,
        final_status text NOT NULL,
        timeline jsonb,
        payload_final jsonb,
        started_at timestamp
        WITH
            TIME ZONE,
            finished_at timestamp
        WITH
            TIME ZONE,
            updated_at timestamp
        WITH
            TIME ZONE DEFAULT now ()
    );