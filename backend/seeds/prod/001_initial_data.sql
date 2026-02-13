-- ====================================
-- Production Initial Data
-- ====================================
-- Essential data for production environment
-- Run ONLY ONCE during initial deployment

BEGIN;

-- Add any essential production data here
-- Examples: default settings, system users, lookup tables, etc.

-- Log
DO $$
BEGIN
    RAISE NOTICE 'Production initial data inserted at %', NOW();
END $$;

COMMIT;
