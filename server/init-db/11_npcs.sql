-- 11_npcs.sql - NPC profiles, memory and relationship storage
-- Run after 10_mail.sql / 10_shops.sql

CREATE TABLE IF NOT EXISTS npc_profiles (
    npc_id VARCHAR(255) PRIMARY KEY,
    entity_id VARCHAR(255) NOT NULL,
    personality_type VARCHAR(50) NOT NULL DEFAULT 'balanced',
    moral_alignment VARCHAR(50) NOT NULL DEFAULT 'neutral',
    ambition_level INT NOT NULL DEFAULT 50,
    risk_tolerance DOUBLE PRECISION NOT NULL DEFAULT 0.5,
    background_story TEXT DEFAULT '',
    current_goal VARCHAR(500) DEFAULT '',
    current_region VARCHAR(255) DEFAULT '',
    realm VARCHAR(50) NOT NULL DEFAULT 'mortal',
    status VARCHAR(20) NOT NULL DEFAULT 'idle',
    last_active_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS npc_memory (
    id SERIAL PRIMARY KEY,
    npc_id VARCHAR(255) NOT NULL REFERENCES npc_profiles(npc_id) ON DELETE CASCADE,
    memory_type VARCHAR(20) NOT NULL DEFAULT 'short_term',
    memory_key VARCHAR(255) NOT NULL DEFAULT '',
    content TEXT NOT NULL,
    importance DOUBLE PRECISION NOT NULL DEFAULT 0.5,
    related_entity_id VARCHAR(255) DEFAULT '',
    related_entity_name VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS npc_relationships (
    npc_id VARCHAR(255) NOT NULL REFERENCES npc_profiles(npc_id) ON DELETE CASCADE,
    target_id VARCHAR(255) NOT NULL,
    target_name VARCHAR(255) NOT NULL DEFAULT '',
    relationship_type VARCHAR(20) NOT NULL DEFAULT 'player',
    affinity INT NOT NULL DEFAULT 0,
    familiarity INT NOT NULL DEFAULT 0,
    last_interaction_at TIMESTAMP DEFAULT NOW(),
    interaction_count INT NOT NULL DEFAULT 0,
    notes TEXT DEFAULT '',
    PRIMARY KEY (npc_id, target_id)
);

CREATE INDEX IF NOT EXISTS idx_npc_memory_npc_id ON npc_memory(npc_id);
CREATE INDEX IF NOT EXISTS idx_npc_memory_type ON npc_memory(npc_id, memory_type);
CREATE INDEX IF NOT EXISTS idx_npc_relationships_affinity ON npc_relationships(npc_id, affinity);
CREATE INDEX IF NOT EXISTS idx_npc_profiles_region ON npc_profiles(current_region);
CREATE INDEX IF NOT EXISTS idx_npc_profiles_status ON npc_profiles(status);
