-- ====================================
-- SubKeep Initial Database Schema
-- ====================================
-- This script is automatically executed when the PostgreSQL container is first created

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create schemas (if needed for separation)
-- CREATE SCHEMA IF NOT EXISTS subkeep;
-- SET search_path TO subkeep, public;

-- Set timezone
SET timezone = 'Asia/Seoul';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE subkeep_db TO subkeep_user;

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Log initialization
DO $$
BEGIN
    RAISE NOTICE 'SubKeep database initialized successfully at %', NOW();
END $$;
