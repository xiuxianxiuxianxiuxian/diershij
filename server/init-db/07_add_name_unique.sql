-- Migration 07: Add UNIQUE constraint on entities.name to prevent duplicate usernames
ALTER TABLE entities ADD CONSTRAINT entities_name_key UNIQUE (name);
