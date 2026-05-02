-- ==========================================
-- 战斗、技能与行为记录模块
-- ==========================================
-- 包含：战斗记录、技能库、实体技能、
--       操作日志、行为模板、世界消息
-- ==========================================

-- 1. 战斗记录表
CREATE TABLE combat_logs (
    id UUID PRIMARY KEY,
    attacker_id UUID REFERENCES entities(id),
    defender_id UUID REFERENCES entities(id),
    location_id UUID REFERENCES world_regions(id),
    attacker_skills JSONB,             -- 使用的技能
    defender_skills JSONB,
    damage_dealt REAL,
    damage_taken REAL,
    result VARCHAR(20),                -- attacker_win/defender_win/draw/flee
    loot JSONB,                        -- 战利品
    karma_change INT,                  -- 战斗导致的业力变化
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP
);

-- 2. 技能库表（所有可学习技能）
CREATE TABLE skills (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    skill_type VARCHAR(30),            -- attack/defense/utility/formation/alchemy
    element VARCHAR(20),
    description TEXT,
    base_power REAL,
    cooldown INT,                      -- 冷却时间（秒）
    resource_cost JSONB,               -- 消耗（灵力/气血等）
    required_realm VARCHAR(30),
    required_methods UUID[],           -- 需要掌握的功法
    learned_from UUID REFERENCES cultivation_methods(id)
);

-- 3. 实体技能表（实体已掌握的技能）
CREATE TABLE entity_skills (
    entity_id UUID REFERENCES entities(id),
    skill_id UUID REFERENCES skills(id),
    mastery REAL DEFAULT 0,            -- 熟练度 0-100
    is_unlocked BOOLEAN DEFAULT FALSE,
    last_used TIMESTAMP,
    PRIMARY KEY (entity_id, skill_id)
);

-- 4. 操作日志表（所有实体的操作记录，用于回放和调试）
CREATE TABLE operation_logs (
    id UUID PRIMARY KEY,
    entity_id UUID REFERENCES entities(id),
    action_type VARCHAR(30),           -- cultivate/combat/trade/explore等
    params JSONB,                      -- 操作参数
    result JSONB,                      -- 操作结果
    world_day INT,
    world_hour INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 5. 行为模板表（AI 使用的预生成模板）
CREATE TABLE behavior_templates (
    id UUID PRIMARY KEY,
    template_type VARCHAR(30),         -- behavior/dialogue/decision
    pattern JSONB,                     -- 匹配模式
    action_template JSONB,             -- 动作模板
    usage_count INT DEFAULT 0,         -- 使用次数
    source VARCHAR(20),                -- pre_generated/llm_generated
    created_at TIMESTAMP DEFAULT NOW()
);

-- 6. 世界消息表（全服/区域广播消息）
CREATE TABLE world_messages (
    id UUID PRIMARY KEY,
    message_type VARCHAR(30),          -- breakthrough/crisis/discovery/etc
    content TEXT,
    target_audience VARCHAR(20),       -- all/region/sect
    related_entity_id UUID,
    related_region_id UUID,
    importance INT,                    -- 重要程度（决定显示方式）
    created_at TIMESTAMP DEFAULT NOW()
);

-- 7. 探索日志表
CREATE TABLE exploration_logs (
    id UUID PRIMARY KEY,
    entity_id UUID REFERENCES entities(id),
    region_id UUID REFERENCES world_regions(id),
    exploration_type VARCHAR(30),      -- search/discover/harvest/etc
    findings JSONB,                    -- 发现的内容
    danger_encountered JSONB,          -- 遇到的危险
    rewards JSONB,                     -- 获得的奖励
    explored_at TIMESTAMP DEFAULT NOW()
);
