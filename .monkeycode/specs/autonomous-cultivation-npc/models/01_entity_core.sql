-- ==========================================
-- 实体核心与属性模块
-- ==========================================
-- 包含：实体基表、基础属性、战斗属性、心境属性、
--       寿命状态、财富资产、特殊属性、法则属性、大道属性
-- ==========================================

-- 1. 实体基表
CREATE TABLE entities (
    id UUID PRIMARY KEY,
    entity_type VARCHAR(10) NOT NULL CHECK (entity_type IN ('player', 'npc')),
    name VARCHAR(50) UNIQUE NOT NULL,
    realm VARCHAR(30) NOT NULL DEFAULT 'mortal',
    created_at TIMESTAMP DEFAULT NOW(),
    last_active_at TIMESTAMP,
    is_online BOOLEAN DEFAULT FALSE
);

-- 2. 基础属性表
CREATE TABLE entity_base_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    age INT DEFAULT 16,
    gender VARCHAR(10),
    appearance INT DEFAULT 50,
    charisma INT DEFAULT 50,
    qi REAL DEFAULT 100,
    spiritual_power REAL DEFAULT 50,
    divine_sense REAL DEFAULT 10,
    comprehension INT DEFAULT 50,
    constitution INT DEFAULT 50,
    luck INT DEFAULT 50,
    cultivation_progress REAL DEFAULT 0,
    cultivation_method_proficiency REAL DEFAULT 0
);

-- 3. 战斗属性表
CREATE TABLE entity_combat_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    attack_power REAL DEFAULT 10,
    defense REAL DEFAULT 5,
    speed REAL DEFAULT 10,
    crit_rate REAL DEFAULT 5,
    crit_damage REAL DEFAULT 150,
    dodge_rate REAL DEFAULT 5,
    hit_rate REAL DEFAULT 95,
    penetration REAL DEFAULT 0,
    damage_reduction REAL DEFAULT 0
);

-- 4. 心境属性表
CREATE TABLE entity_mental_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    mental_stability INT DEFAULT 80,
    obsession_count INT DEFAULT 0,
    dao_heart INT DEFAULT 50,
    inner_demon_resistance INT DEFAULT 50,
    enlightenment INT DEFAULT 0
);

-- 5. 寿命与状态表
CREATE TABLE entity_lifespan_status (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    remaining_lifespan INT DEFAULT 80,
    max_lifespan INT DEFAULT 80,
    aging_penalty REAL DEFAULT 0,
    injuries JSONB DEFAULT '[]',
    buffs JSONB DEFAULT '[]',
    debuffs JSONB DEFAULT '[]',
    poison_level INT DEFAULT 0,
    curse_level INT DEFAULT 0
);

-- 6. 财富与资产表
CREATE TABLE entity_wealth_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    property_value BIGINT DEFAULT 0,
    real_estate JSONB DEFAULT '[]',
    business_income BIGINT DEFAULT 0
);

-- 7. 特殊属性表
CREATE TABLE entity_special_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    bloodline VARCHAR(30) DEFAULT '凡人',
    bloodline_purity INT DEFAULT 0,
    physique VARCHAR(50) DEFAULT '凡体',
    physique_awakened BOOLEAN DEFAULT FALSE,
    destiny INT DEFAULT 50,
    world_favor INT DEFAULT 0
);

-- 8. 法则属性表
CREATE TABLE entity_law_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    laws JSONB DEFAULT '{}',  -- {"metal": 20.5, "wood": 10.0}
    law_resonance INT DEFAULT 0,
    domain_power REAL DEFAULT 0,
    domain_range REAL DEFAULT 0,
    law_suppression REAL DEFAULT 0
);

-- 9. 大道属性表
CREATE TABLE entity_dao_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    dao_seed_type VARCHAR(30) DEFAULT '无',
    dao_seed_level INT DEFAULT 0,
    dao_seed_growth REAL DEFAULT 0,
    dao_marks INT DEFAULT 0,
    dao_heart_comprehension INT DEFAULT 0,
    destiny_path VARCHAR(50) DEFAULT '凡途'
);
