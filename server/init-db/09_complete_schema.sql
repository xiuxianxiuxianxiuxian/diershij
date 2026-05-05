-- Migration 09: Complete schema — fill gaps between Go types and database
-- Addresses: missing columns, missing tables, JSONB complex types, indexes

-- ============================================================================
-- PART 1: Extend base_attributes with missing scalar fields + JSONB complex types
-- ============================================================================

-- Missing scalar fields from Attributes Go struct
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS appearance INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS charisma INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS obsession_count INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS inner_demon_resistance INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS aging_penalty DOUBLE PRECISION DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS root_awakened BOOLEAN DEFAULT false;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS mutated_root VARCHAR(30) DEFAULT '';
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS inner_demon_resistance INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS aging_penalty DOUBLE PRECISION DEFAULT 0;

-- JSONB complex type columns (no separate tables needed — queried by ID)
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS injuries JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS buffs JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS debuffs JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS laws JSONB DEFAULT '{}'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS faction_standings JSONB DEFAULT '{}'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS real_estate JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS disciple_ids JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS sworn_siblings JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS enemies JSONB DEFAULT '[]'::jsonb;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS lovers JSONB DEFAULT '[]'::jsonb;

-- Special attributes
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS bloodline VARCHAR(30) DEFAULT '';
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS bloodline_purity INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS physique VARCHAR(30) DEFAULT '';
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS physique_awakened BOOLEAN DEFAULT false;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS destiny INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS world_favor INTEGER DEFAULT 0;

-- Law attributes
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS law_resonance INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS domain_power DOUBLE PRECISION DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS domain_range DOUBLE PRECISION DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS law_suppression DOUBLE PRECISION DEFAULT 0;

-- Dao attributes
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS dao_seed_type VARCHAR(30) DEFAULT '';
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS dao_seed_level INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS dao_seed_growth DOUBLE PRECISION DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS dao_marks INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS dao_heart_comprehension INTEGER DEFAULT 0;
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS destiny_path VARCHAR(50) DEFAULT '';

-- Social relationship IDs (mentor stored as simple FK)
ALTER TABLE base_attributes ADD COLUMN IF NOT EXISTS mentor_id UUID REFERENCES entities(id) ON DELETE SET NULL;

-- ============================================================================
-- PART 3: Create cultivation_methods table (60+ fields from types.CultivationMethod)
-- ============================================================================
CREATE TABLE IF NOT EXISTS cultivation_methods (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    creator_id UUID REFERENCES entities(id) ON DELETE SET NULL,
    origin_sect VARCHAR(100),
    rank VARCHAR(20) NOT NULL,
    category VARCHAR(30) NOT NULL,
    element_affinity VARCHAR(20),
    description TEXT,
    version INTEGER DEFAULT 1,

    -- Core cultivation bonuses
    cultivation_speed_mult DOUBLE PRECISION DEFAULT 1.0,
    spiritual_power_cap_mult DOUBLE PRECISION DEFAULT 1.0,
    qi_cap_mult DOUBLE PRECISION DEFAULT 1.0,
    divine_sense_cap_mult DOUBLE PRECISION DEFAULT 1.0,
    lifespan_bonus INTEGER DEFAULT 0,
    recovery_speed_mult DOUBLE PRECISION DEFAULT 1.0,

    -- Combat bonuses (JSONB maps)
    attack_bonuses JSONB DEFAULT '{}'::jsonb,
    defense_bonuses JSONB DEFAULT '{}'::jsonb,
    utility_bonuses JSONB DEFAULT '{}'::jsonb,

    -- Special effects
    passive_effects JSONB DEFAULT '[]'::jsonb,
    active_skills JSONB DEFAULT '[]'::jsonb,
    ultimate_skill JSONB,

    -- Law and Dao affinity
    law_affinities JSONB DEFAULT '[]'::jsonb,
    law_comprehension_bonus DOUBLE PRECISION DEFAULT 0,
    dao_compatibility JSONB DEFAULT '[]'::jsonb,

    -- Restrictions
    required_roots JSONB DEFAULT '[]'::jsonb,
    required_physique JSONB DEFAULT '[]'::jsonb,
    realm_requirement VARCHAR(30) DEFAULT 'mortal',
    alignment_restriction VARCHAR(20) DEFAULT 'none',
    karma_threshold INTEGER DEFAULT 0,
    gender_restriction VARCHAR(10) DEFAULT 'none',

    -- Inheritance and evolution
    parent_method_id UUID REFERENCES cultivation_methods(id) ON DELETE SET NULL,
    evolution_path JSONB DEFAULT '[]'::jsonb,
    transmission_mode VARCHAR(20) DEFAULT 'jade_slip',
    can_modify BOOLEAN DEFAULT false,
    complexity INTEGER DEFAULT 1,

    -- Evaluation
    power_score INTEGER DEFAULT 0,
    potential INTEGER DEFAULT 0,
    popularity INTEGER DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- PART 4: Create world_events table
-- ============================================================================
CREATE TABLE IF NOT EXISTS world_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,
    description TEXT,
    region_id VARCHAR(50) REFERENCES world_regions(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    participants JSONB DEFAULT '[]'::jsonb,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- PART 5: Create npc_decision_logs table
-- ============================================================================
CREATE TABLE IF NOT EXISTS npc_decision_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    npc_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    decision_type VARCHAR(50) NOT NULL,
    context JSONB,
    action_taken JSONB,
    reasoning TEXT,
    model_used VARCHAR(50),
    source VARCHAR(20) DEFAULT 'behavior_tree',
    token_cost DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- PART 6: Extend npc_personalities
-- ============================================================================
ALTER TABLE npc_personalities ADD COLUMN IF NOT EXISTS llm_system_prompt TEXT;
ALTER TABLE npc_personalities ADD COLUMN IF NOT EXISTS initial_actions JSONB DEFAULT '[]'::jsonb;

-- ============================================================================
-- PART 7: Extend sects
-- ============================================================================
ALTER TABLE sects ADD COLUMN IF NOT EXISTS prestige INTEGER DEFAULT 0;
ALTER TABLE sects ADD COLUMN IF NOT EXISTS wealth BIGINT DEFAULT 0;
ALTER TABLE sects ADD COLUMN IF NOT EXISTS facility_score INTEGER DEFAULT 0;
ALTER TABLE sects ADD COLUMN IF NOT EXISTS cultivation_resources JSONB DEFAULT '[]'::jsonb;

-- ============================================================================
-- PART 8: Extend sect_members
-- ============================================================================
ALTER TABLE sect_members ADD COLUMN IF NOT EXISTS privileges JSONB DEFAULT '[]'::jsonb;

-- ============================================================================
-- PART 9: Create relationships table (general entity-to-entity relationships)
-- ============================================================================
CREATE TABLE IF NOT EXISTS relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_a_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    entity_b_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    relationship_type VARCHAR(30) NOT NULL,
    strength DOUBLE PRECISION DEFAULT 50,
    history TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(entity_a_id, entity_b_id, relationship_type)
);

-- ============================================================================
-- PART 10: Add missing indexes for query performance
-- ============================================================================

-- Safety net: ensure spiritual_roots and entity_methods exist (created in migration 08)
-- before creating indexes on them
CREATE TABLE IF NOT EXISTS spiritual_roots (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    element VARCHAR(20) NOT NULL,
    purity INTEGER NOT NULL DEFAULT 50,
    PRIMARY KEY (entity_id, element)
);
CREATE TABLE IF NOT EXISTS entity_methods (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    method_id UUID NOT NULL,
    mastery_level DOUBLE PRECISION DEFAULT 0,
    is_main_method BOOLEAN DEFAULT false,
    learned_at TIMESTAMP DEFAULT NOW(),
    last_practiced TIMESTAMP,
    backlash_risk DOUBLE PRECISION DEFAULT 0,
    modified BOOLEAN DEFAULT false,
    modified_notes TEXT,
    PRIMARY KEY (entity_id, method_id)
);

CREATE INDEX IF NOT EXISTS idx_spiritual_roots_entity ON spiritual_roots(entity_id);
CREATE INDEX IF NOT EXISTS idx_entity_methods_entity ON entity_methods(entity_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_unread ON messages(receiver_id, is_read);
CREATE INDEX IF NOT EXISTS idx_world_events_region ON world_events(region_id);
CREATE INDEX IF NOT EXISTS idx_world_events_status ON world_events(status);
CREATE INDEX IF NOT EXISTS idx_world_events_type ON world_events(type);
CREATE INDEX IF NOT EXISTS idx_npc_decision_logs_npc ON npc_decision_logs(npc_id);
CREATE INDEX IF NOT EXISTS idx_npc_decision_logs_type ON npc_decision_logs(decision_type);
CREATE INDEX IF NOT EXISTS idx_relationships_entity_a ON relationships(entity_a_id);
CREATE INDEX IF NOT EXISTS idx_relationships_entity_b ON relationships(entity_b_id);
CREATE INDEX IF NOT EXISTS idx_relationships_type ON relationships(relationship_type);
CREATE INDEX IF NOT EXISTS idx_cultivation_methods_rank ON cultivation_methods(rank);
CREATE INDEX IF NOT EXISTS idx_cultivation_methods_category ON cultivation_methods(category);
CREATE INDEX IF NOT EXISTS idx_cultivation_methods_element ON cultivation_methods(element_affinity);
CREATE INDEX IF NOT EXISTS idx_cultivation_methods_realm ON cultivation_methods(realm_requirement);
CREATE INDEX IF NOT EXISTS idx_cultivation_methods_creator ON cultivation_methods(creator_id);
