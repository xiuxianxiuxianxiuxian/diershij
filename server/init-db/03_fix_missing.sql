-- fix encoding issues: recreate tables that failed due to GBK encoding

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,
    rarity INTEGER DEFAULT 1,
    description TEXT,
    attributes JSONB,
    stackable BOOLEAN DEFAULT false,
    max_stack INTEGER DEFAULT 1,
    usable BOOLEAN DEFAULT false,
    level_requirement INTEGER DEFAULT 0,
    realm_requirement VARCHAR(30) DEFAULT 'mortal',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 1,
    equipped BOOLEAN DEFAULT false,
    slot VARCHAR(20),
    durability INTEGER,
    bound BOOLEAN DEFAULT false,
    acquired_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(entity_id, item_id, slot)
);

CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,
    difficulty INTEGER DEFAULT 1,
    description TEXT,
    materials JSONB NOT NULL,
    result_item_id UUID REFERENCES items(id),
    result_quantity INTEGER DEFAULT 1,
    skill_level_required INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entity_recipes (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    learned_at TIMESTAMP DEFAULT NOW(),
    proficiency INTEGER DEFAULT 0,
    PRIMARY KEY (entity_id, recipe_id)
);

CREATE TABLE IF NOT EXISTS spells (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL,
    element VARCHAR(20),
    cost INTEGER DEFAULT 10,
    base_damage INTEGER DEFAULT 0,
    base_heal INTEGER DEFAULT 0,
    duration INTEGER DEFAULT 0,
    cooldown INTEGER DEFAULT 0,
    description TEXT,
    realm_requirement VARCHAR(30) DEFAULT 'mortal',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entity_spells (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    spell_id UUID REFERENCES spells(id) ON DELETE CASCADE,
    learned_at TIMESTAMP DEFAULT NOW(),
    proficiency INTEGER DEFAULT 0,
    last_cast_at TIMESTAMP,
    PRIMARY KEY (entity_id, spell_id)
);

CREATE TABLE IF NOT EXISTS craft_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES entities(id),
    recipe_id UUID REFERENCES recipes(id),
    success BOOLEAN DEFAULT false,
    result_item_id UUID REFERENCES items(id),
    materials_consumed JSONB,
    skill_exp_gained INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cast_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caster_id UUID REFERENCES entities(id),
    spell_id UUID REFERENCES spells(id),
    target_id UUID REFERENCES entities(id),
    damage_dealt INTEGER DEFAULT 0,
    heal_amount INTEGER DEFAULT 0,
    qi_consumed INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_inventory_entity ON inventory(entity_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item ON inventory(item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_equipped ON inventory(equipped);
CREATE INDEX IF NOT EXISTS idx_entity_recipes_entity ON entity_recipes(entity_id);
CREATE INDEX IF NOT EXISTS idx_entity_spells_entity ON entity_spells(entity_id);
