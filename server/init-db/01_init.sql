CREATE TABLE IF NOT EXISTS entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(20) NOT NULL DEFAULT 'player',
    name VARCHAR(100) NOT NULL,
    realm VARCHAR(30) NOT NULL DEFAULT 'mortal',
    region_id VARCHAR(50) DEFAULT 'qingyun_town',
    x DOUBLE PRECISION DEFAULT 0,
    y DOUBLE PRECISION DEFAULT 0,
    status VARCHAR(20) DEFAULT 'normal',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS base_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id) ON DELETE CASCADE,
    qi DOUBLE PRECISION DEFAULT 100,
    max_qi DOUBLE PRECISION DEFAULT 100,
    spiritual_power DOUBLE PRECISION DEFAULT 100,
    max_spiritual_power DOUBLE PRECISION DEFAULT 100,
    divine_sense DOUBLE PRECISION DEFAULT 10,
    comprehension INTEGER DEFAULT 50,
    constitution INTEGER DEFAULT 50,
    luck INTEGER DEFAULT 50,
    cultivation_progress DOUBLE PRECISION DEFAULT 0,
    attack_power DOUBLE PRECISION DEFAULT 10,
    defense DOUBLE PRECISION DEFAULT 10,
    speed DOUBLE PRECISION DEFAULT 10,
    mental_stability INTEGER DEFAULT 50,
    remaining_lifespan INTEGER DEFAULT 80,
    max_lifespan INTEGER DEFAULT 80
);

CREATE TABLE IF NOT EXISTS karma_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id) ON DELETE CASCADE,
    karma_value INTEGER DEFAULT 0,
    merit INTEGER DEFAULT 0,
    heavenly_mark VARCHAR(20) DEFAULT 'clear'
);

CREATE TABLE IF NOT EXISTS spirit_stones (
    entity_id UUID PRIMARY KEY REFERENCES entities(id) ON DELETE CASCADE,
    low_grade BIGINT DEFAULT 0,
    medium_grade BIGINT DEFAULT 0,
    high_grade BIGINT DEFAULT 0,
    premium_grade BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS operation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES entities(id),
    action_type VARCHAR(50) NOT NULL,
    params JSONB,
    result JSONB,
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS world_regions (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_region_id VARCHAR(50),
    spiritual_density DOUBLE PRECISION DEFAULT 0,
    spiritual_tier INTEGER DEFAULT 1,
    danger_level INTEGER DEFAULT 0,
    resources JSONB,
    rules JSONB,
    description TEXT,
    lore TEXT
);

CREATE TABLE IF NOT EXISTS sects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    founder_id UUID REFERENCES entities(id),
    philosophy TEXT,
    entry_requirements JSONB,
    territory JSONB,
    rules JSONB,
    alignment VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sect_members (
    sect_id UUID REFERENCES sects(id) ON DELETE CASCADE,
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    rank VARCHAR(30),
    contribution DOUBLE PRECISION DEFAULT 0,
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (sect_id, entity_id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID REFERENCES entities(id),
    buyer_id UUID REFERENCES entities(id),
    item_id UUID,
    item_type VARCHAR(30),
    price DOUBLE PRECISION,
    currency VARCHAR(20) DEFAULT 'spirit_stone',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS npc_personalities (
    npc_id UUID PRIMARY KEY REFERENCES entities(id) ON DELETE CASCADE,
    personality_type VARCHAR(30),
    moral_alignment VARCHAR(20),
    ambition_level INTEGER,
    risk_tolerance DOUBLE PRECISION,
    social_preference VARCHAR(20),
    background_story TEXT,
    current_goal TEXT,
    hidden_secrets JSONB,
    behavior_tree_config JSONB
);

CREATE INDEX idx_entities_name ON entities(name);
CREATE INDEX idx_entities_region ON entities(region_id);
CREATE INDEX idx_operation_logs_actor ON operation_logs(actor_id);
CREATE INDEX idx_operation_logs_timestamp ON operation_logs(timestamp);
