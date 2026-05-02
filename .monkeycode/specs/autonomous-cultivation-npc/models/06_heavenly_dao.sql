-- ==========================================
-- 天道系统与事件模块
-- ==========================================
-- 包含：天道运行日志、天劫记录、世界危机、
--       因果业力流水、灵气潮汐、世界纪元、资源刷新
-- ==========================================

-- 1. 天道运行日志表（所有天道算法执行记录）
CREATE TABLE heavenly_dao_logs (
    id UUID PRIMARY KEY,
    algorithm_name VARCHAR(50),        -- 算法模块名称（karma/tribulation/balance等）
    target_entity_id UUID REFERENCES entities(id),
    input_params JSONB,                -- 输入参数快照
    output_result JSONB,               -- 算法输出结果
    execution_time TIMESTAMP DEFAULT NOW(),
    world_day INT                      -- 世界第几天
);

-- 2. 天劫记录表
CREATE TABLE heavenly_tribulations (
    id UUID PRIMARY KEY,
    entity_id UUID REFERENCES entities(id),
    target_realm VARCHAR(30),          -- 目标境界
    trigger_reason VARCHAR(50),        -- 触发原因（breakthrough/karma等）
    probability REAL,                  -- 触发时计算的概率
    strength INT,                      -- 天劫强度
    result VARCHAR(20),                -- success/failed/dead
    damage_taken REAL,                 -- 承受伤害
    reward JSONB,                      -- 渡劫成功奖励
    triggered_at TIMESTAMP DEFAULT NOW()
);

-- 3. 世界危机事件表
CREATE TABLE world_crises (
    id UUID PRIMARY KEY,
    crisis_type VARCHAR(30),           -- demon_beast_tide/heavenly_demon/void_rift等
    danger_level INT,                  -- 危险等级
    affected_regions UUID[],           -- 受影响区域
    spawn_count INT,                   -- 生成怪物/事件数量
    status VARCHAR(20),                -- active/resolved/failed
    start_day INT,
    end_day INT,
    participants JSONB,                -- 参与实体及贡献度
    resolution_rewards JSONB           -- 平息奖励分配
);

-- 4. 因果业力流水表
CREATE TABLE karma_transactions (
    id UUID PRIMARY KEY,
    entity_id UUID REFERENCES entities(id),
    action_type VARCHAR(50),           -- 触发业力的行为
    target_entity_id UUID,             -- 行为目标（可为空）
    karma_change INT,                  -- 业力变化值
    merit_change INT,                  -- 功德变化值
    context JSONB,                     -- 行为上下文
    created_at TIMESTAMP DEFAULT NOW()
);

-- 5. 灵气潮汐记录表
CREATE TABLE spiritual_tides (
    id UUID PRIMARY KEY,
    world_day INT,
    tide_level VARCHAR(20),            -- rising/ebbing/stable
    global_multiplier REAL,            -- 全局灵气倍率
    special_effects JSONB              -- 特殊效果（如某属性灵气暴涨）
);

-- 6. 世界纪元表
CREATE TABLE world_epochs (
    id UUID PRIMARY KEY,
    epoch_name VARCHAR(50),
    start_day INT,
    end_day INT,
    spiritual_density REAL,            -- 本纪元基础灵气浓度
    realm_cap VARCHAR(30),             -- 境界上限
    special_rules JSONB,               -- 特殊天道规则
    transition_cause TEXT              -- 纪元更替原因
);

-- 7. 资源刷新记录表
CREATE TABLE resource_spawns (
    id UUID PRIMARY KEY,
    region_id UUID REFERENCES world_regions(id),
    resource_type VARCHAR(30),         -- herb/ore/beast/etc
    quantity INT,
    quality_tier INT,                  -- 品质等级
    position JSONB,                    -- 坐标位置
    respawn_time TIMESTAMP,            -- 下次刷新时间
    harvested_by UUID REFERENCES entities(id), -- 最后采集者
    created_at TIMESTAMP DEFAULT NOW()
);
