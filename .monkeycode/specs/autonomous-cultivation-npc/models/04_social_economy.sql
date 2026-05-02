-- ==========================================
-- 社交、因果与经济模块
-- ==========================================
-- 包含：生活技能、声望社交、因果业力、
--       灵石资产、交易记录、NPC 关系网络
-- ==========================================

-- 1. 生活技能表
CREATE TABLE entity_life_skills (
    entity_id UUID REFERENCES entities(id),
    alchemy_level INT DEFAULT 0,
    artificing_level INT DEFAULT 0,
    formation_level INT DEFAULT 0,
    fire_control INT DEFAULT 0,
    herb_knowledge INT DEFAULT 0,
    mining_skill INT DEFAULT 0,
    talisman_skill INT DEFAULT 0,
    beast_taming INT DEFAULT 0,
    PRIMARY KEY (entity_id)
);

-- 2. 声望与社交表
CREATE TABLE entity_social_attributes (
    entity_id UUID REFERENCES entities(id),
    reputation INT DEFAULT 0,
    sect_contribution INT DEFAULT 0,
    relationship_count INT DEFAULT 0,
    mentor_id UUID REFERENCES entities(id),
    disciple_ids UUID[],
    sworn_siblings UUID[],
    enemies UUID[],
    lovers UUID[],
    faction_standings JSONB DEFAULT '{}',
    PRIMARY KEY (entity_id)
);

-- 3. 因果业力表
CREATE TABLE entity_karma_attributes (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    karma INT DEFAULT 0,
    merit INT DEFAULT 0,
    karmic_debt INT DEFAULT 0,
    heavenly_mark VARCHAR(20) DEFAULT '清白'
);

-- 4. 灵石资产表
CREATE TABLE entity_spirit_stones (
    entity_id UUID PRIMARY KEY REFERENCES entities(id),
    low_grade BIGINT DEFAULT 0,     -- 下品灵石
    medium_grade BIGINT DEFAULT 0,  -- 中品灵石
    high_grade BIGINT DEFAULT 0,    -- 上品灵石
    premium_grade BIGINT DEFAULT 0  -- 极品灵石
);

-- 5. 交易记录表
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    seller_id UUID REFERENCES entities(id),
    buyer_id UUID REFERENCES entities(id),
    item_id UUID,
    item_type VARCHAR(30),
    price REAL,
    currency VARCHAR(20) DEFAULT 'spirit_stone',
    created_at TIMESTAMP DEFAULT NOW()
);

-- 6. NPC 关系网络表
CREATE TABLE npc_relationships (
    id UUID PRIMARY KEY,
    entity_a_id UUID REFERENCES entities(id),
    entity_b_id UUID REFERENCES entities(id),
    relationship_type VARCHAR(30),      -- 师徒/仇敌/盟友/恋人等
    strength REAL,                      -- 关系强度
    history TEXT,                       -- 关系历史
    created_at TIMESTAMP DEFAULT NOW()
);
