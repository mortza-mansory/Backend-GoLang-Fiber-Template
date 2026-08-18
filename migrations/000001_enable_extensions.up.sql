-- Sample migration: enable commonly useful extensions.
-- This is infrastructure, not business data.
-- Add your own feature migrations here with the pattern:
--   000002_<feature>.up.sql / 000002_<feature>.down.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";