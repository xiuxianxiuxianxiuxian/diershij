-- Add spiritual roots support
CREATE TABLE IF NOT EXISTS spiritual_roots (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    element VARCHAR(20) NOT NULL,  -- gold, wood, water, fire, earth, wind, thunder, ice, light, dark
    purity INTEGER NOT NULL DEFAULT 50,  -- 1-100
    PRIMARY KEY (entity_id, element)
);

-- Add entity methods support (cultivation techniques)
CREATE TABLE IF NOT EXISTS entity_methods (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    method_id UUID NOT NULL,
    mastery_level DOUBLE PRECISION DEFAULT 0,  -- 0-100%
    is_main_method BOOLEAN DEFAULT false,
    learned_at TIMESTAMP DEFAULT NOW(),
    last_practiced TIMESTAMP,
    backlash_risk DOUBLE PRECISION DEFAULT 0,  -- 0-1
    modified BOOLEAN DEFAULT false,
    modified_notes TEXT,
    PRIMARY KEY (entity_id, method_id)
);

-- Add missing entity columns for complete character system
ALTER TABLE entities ADD COLUMN IF NOT EXISTS age INTEGER DEFAULT 18;
ALTER TABLE entities ADD COLUMN IF NOT EXISTS gender VARCHAR(10) DEFAULT 'unknown';
