-- Migration 06: Add password_hash column to entities table for secure authentication
ALTER TABLE entities ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255) NOT NULL DEFAULT '';
