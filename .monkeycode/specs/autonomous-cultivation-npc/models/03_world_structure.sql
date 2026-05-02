-- ==========================================
-- 世界结构与势力模块
-- ==========================================
-- 包含：区域表、宗门表、宗门成员、世界传说
-- ==========================================

-- 1. 区域表
CREATE TABLE world_regions (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_region_id UUID REFERENCES world_regions(id),
    spiritual_density REAL DEFAULT 0,  -- 灵气浓度
    spiritual_tier INTEGER DEFAULT 1,  -- 灵气品阶（1-9）
    danger_level INTEGER DEFAULT 0,    -- 危险等级
    resources JSONB,                    -- 资源分布
    rules JSONB,                        -- 区域规则（禁区等）
    description TEXT,                   -- 区域描述
    lore TEXT                           -- 区域传说/历史
);

-- 2. 宗门表
CREATE TABLE sects (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    founder_id UUID REFERENCES entities(id),
    philosophy TEXT,                    -- 宗门理念
    entry_requirements JSONB,           -- 入门条件
    territory JSONB,                    -- 势力范围
    rules JSONB,                        -- 宗门规则
    alignment VARCHAR(20),              -- 正道/魔道/中立
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3. 宗门成员表
CREATE TABLE sect_members (
    sect_id UUID REFERENCES sects(id),
    entity_id UUID REFERENCES entities(id),
    rank VARCHAR(30),                   -- 职位
    contribution REAL DEFAULT 0,        -- 贡献值
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (sect_id, entity_id)
);

-- 4. 世界历史/传说表
CREATE TABLE world_lore (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    category VARCHAR(30),               -- 大战/传说/遗迹/未解之谜
    content TEXT,                       -- 详细内容
    related_regions JSONB,              -- 关联区域
    related_entities JSONB,             -- 关联实体
    hints JSONB,                        -- 探索线索
    is_discovered BOOLEAN DEFAULT FALSE -- 是否已被发现
);
