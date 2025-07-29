-- Database initialization script
-- This script runs when PostgreSQL container starts for the first time

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Set timezone
SET timezone = 'UTC';

-- Create initial database if it doesn't exist
-- (This is handled by POSTGRES_DB environment variable)