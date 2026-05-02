-- ==========================================
-- 修仙体系与功法模块
-- ==========================================
-- 包含：灵根系统、功法库、功法关联
-- ==========================================

-- 1. 灵根系统表
CREATE TABLE entity_spiritual_roots (
    entity_id UUID REFERENCES entities(id),
    root_type VARCHAR(20),  -- 金木水火风雷冰光暗等
    purity INT DEFAULT 50,
    is_awakened BOOLEAN DEFAULT FALSE,
    is_mutated BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (entity_id, root_type)
);

-- 2. 功法表
CREATE TABLE cultivation_methods (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    creator_id UUID REFERENCES entities(id),
    origin_sect VARCHAR(100),
    rank VARCHAR(20),              -- 天地玄黄 x 上中下极品
    category VARCHAR(20),          -- 主修/秘术/身法/神识/辅助
    element_affinity VARCHAR(20),
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    version INTEGER DEFAULT 1,
    
    -- 修炼加成
    cultivation_speed_mult REAL DEFAULT 1.0,
    spiritual_power_cap_mult REAL DEFAULT 1.0,
    qi_cap_mult REAL DEFAULT 1.0,
    divine_sense_cap_mult REAL DEFAULT 1.0,
    lifespan_bonus INT DEFAULT 0,
    recovery_speed_mult REAL DEFAULT 1.0,
    
    -- 战斗加成
    attack_bonuses JSONB DEFAULT '{}',
    defense_bonuses JSONB DEFAULT '{}',
    utility_bonuses JSONB DEFAULT '{}',
    passive_effects TEXT[],
    active_skills JSONB DEFAULT '[]',
    ultimate_skill JSONB,
    
    -- 法则亲和
    law_affinities TEXT[],
    law_comprehension_bonus REAL DEFAULT 1.0,
    dao_compatibility TEXT[],
    
    -- 限制条件
    required_roots TEXT[],
    required_physique TEXT[],
    realm_requirement VARCHAR(30),
    alignment_restriction VARCHAR(20),
    karma_threshold INT DEFAULT 0,
    gender_restriction VARCHAR(10) DEFAULT '无',
    
    -- 传承演化
    parent_method_id UUID REFERENCES cultivation_methods(id),
    evolution_path UUID[],
    transmission_mode VARCHAR(20) DEFAULT '玉简',
    can_modify BOOLEAN DEFAULT FALSE,
    complexity INT DEFAULT 1
);

-- 3. 功法-实体关联表
CREATE TABLE entity_methods (
    entity_id UUID REFERENCES entities(id),
    method_id UUID REFERENCES cultivation_methods(id),
    proficiency REAL DEFAULT 0,     -- 熟练度
    is_active BOOLEAN DEFAULT FALSE, -- 是否为主修功法
    PRIMARY KEY (entity_id, method_id)
);
