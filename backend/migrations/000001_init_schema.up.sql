-- SubKeep Initial Schema
-- Version: 000001
-- Description: Create users, subscriptions, and categories tables

BEGIN;

-- ============ Enable Extensions ============
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============ Users Table ============
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$')
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- ============ Categories Table ============
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(7),
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT name_not_empty CHECK (char_length(trim(name)) > 0)
);

CREATE UNIQUE INDEX idx_categories_name ON categories(LOWER(name));

-- Insert default categories
INSERT INTO categories (name, icon, color, is_system) VALUES
('Streaming', '🎬', '#FF6B6B', TRUE),
('Music', '🎵', '#4ECDC4', TRUE),
('Cloud Storage', '☁️', '#45B7D1', TRUE),
('Software', '💻', '#96CEB4', TRUE),
('Gaming', '🎮', '#FFEAA7', TRUE),
('News', '📰', '#DFE6E9', TRUE),
('Fitness', '💪', '#FD79A8', TRUE),
('Education', '📚', '#A29BFE', TRUE),
('Other', '📦', '#B2BEC3', TRUE)
ON CONFLICT DO NOTHING;

-- ============ Subscriptions Table ============
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    
    -- Service information
    service_name VARCHAR(255) NOT NULL,
    icon_url TEXT,
    website_url TEXT,
    
    -- Billing information
    billing_amount DECIMAL(10, 2) NOT NULL,
    billing_currency VARCHAR(3) DEFAULT 'USD',
    billing_cycle VARCHAR(20) NOT NULL, -- 'MONTHLY', 'YEARLY', 'WEEKLY', 'ONE_TIME'
    
    -- Calculated monthly normalized amount (for analytics)
    monthly_cost DECIMAL(10, 2) GENERATED ALWAYS AS (
        CASE
            WHEN billing_cycle = 'MONTHLY' THEN billing_amount
            WHEN billing_cycle = 'YEARLY' THEN billing_amount / 12
            WHEN billing_cycle = 'WEEKLY' THEN billing_amount * 52 / 12
            ELSE billing_amount
        END
    ) STORED,
    
    -- Sharing information
    is_shared BOOLEAN DEFAULT FALSE,
    total_shares INTEGER DEFAULT 1,
    my_share INTEGER DEFAULT 1,
    
    -- Out-of-pocket cost (after sharing)
    out_of_pocket_cost DECIMAL(10, 2) GENERATED ALWAYS AS (
        CASE
            WHEN billing_cycle = 'MONTHLY' THEN (billing_amount * my_share / NULLIF(total_shares, 0))
            WHEN billing_cycle = 'YEARLY' THEN (billing_amount / 12 * my_share / NULLIF(total_shares, 0))
            WHEN billing_cycle = 'WEEKLY' THEN (billing_amount * 52 / 12 * my_share / NULLIF(total_shares, 0))
            ELSE billing_amount
        END
    ) STORED,
    
    -- Dates
    start_date DATE NOT NULL,
    next_billing_date DATE NOT NULL,
    end_date DATE,
    
    -- Usage tracking
    satisfaction_score INTEGER CHECK (satisfaction_score BETWEEN 1 AND 5),
    usage_frequency VARCHAR(20), -- 'DAILY', 'WEEKLY', 'MONTHLY', 'RARELY'
    
    -- Status
    status VARCHAR(20) DEFAULT 'ACTIVE', -- 'ACTIVE', 'PAUSED', 'CANCELLED'
    
    -- Metadata
    notes TEXT,
    tags TEXT[], -- Array of tags
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT billing_amount_positive CHECK (billing_amount > 0),
    CONSTRAINT total_shares_positive CHECK (total_shares > 0),
    CONSTRAINT my_share_valid CHECK (my_share > 0 AND my_share <= total_shares),
    CONSTRAINT valid_billing_cycle CHECK (billing_cycle IN ('MONTHLY', 'YEARLY', 'WEEKLY', 'ONE_TIME')),
    CONSTRAINT valid_status CHECK (status IN ('ACTIVE', 'PAUSED', 'CANCELLED')),
    CONSTRAINT valid_dates CHECK (next_billing_date >= start_date)
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_category_id ON subscriptions(category_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_next_billing ON subscriptions(next_billing_date);
CREATE INDEX idx_subscriptions_monthly_cost ON subscriptions(monthly_cost DESC);

-- ============ Triggers ============
-- Update updated_at timestamp automatically
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
