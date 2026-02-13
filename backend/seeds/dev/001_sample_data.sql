-- ====================================
-- Development Seed Data
-- ====================================
-- Sample data for local development and testing

BEGIN;

-- Sample Users
INSERT INTO users (id, email, username, password_hash, full_name, email_verified, is_active) VALUES
    ('11111111-1111-1111-1111-111111111111', 'test@example.com', 'testuser', '$2a$10$dummy.hash.for.testing', 'Test User', TRUE, TRUE),
    ('22222222-2222-2222-2222-222222222222', 'demo@example.com', 'demouser', '$2a$10$dummy.hash.for.testing', 'Demo User', TRUE, TRUE)
ON CONFLICT (email) DO NOTHING;

-- Log
DO $$
BEGIN
    RAISE NOTICE 'Development seed data inserted at %', NOW();
END $$;

COMMIT;
