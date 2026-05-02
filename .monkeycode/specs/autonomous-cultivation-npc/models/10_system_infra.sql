-- ==========================================
-- 系统配置、服务器与模板数据模块
-- ==========================================
-- 包含：系统配置、服务器状态、世界种子、
--       预置 NPC 模板、初始世界数据、版本记录
-- ==========================================

-- 1. 系统配置表
CREATE TABLE system_config (
    key VARCHAR(50) PRIMARY KEY,
    value JSONB,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 2. 服务器状态表
CREATE TABLE server_status (
    id UUID PRIMARY KEY,
    node_name VARCHAR(50),
    entity_count INT DEFAULT 0,        -- 当前承载实体数
    max_entity_capacity INT,           -- 最大容量
    cpu_usage REAL,
    memory_usage REAL,
    status VARCHAR(20),                -- active/maintenance/offline
    last_heartbeat TIMESTAMP DEFAULT NOW()
);

-- 3. 世界种子表
CREATE TABLE world_seeds (
    id UUID PRIMARY KEY,
    seed_value BIGINT,                 -- 随机种子
    template_name VARCHAR(50),         -- 使用的世界模板
    generation_params JSONB,           -- 生成参数
    generated_at TIMESTAMP DEFAULT NOW()
);

-- 4. NPC 初始模板表
CREATE TABLE npc_templates (
    id UUID PRIMARY KEY,
    template_name VARCHAR(50),
    realm VARCHAR(30),
    personality_type VARCHAR(30),
    moral_alignment VARCHAR(20),
    background_story TEXT,
    initial_attributes JSONB,
    initial_methods UUID[],
    initial_equipment JSONB,
    spawn_region_id UUID REFERENCES world_regions(id),
    behavior_presets JSONB             -- 预设行为模式
);

-- 5. 初始世界数据表
CREATE TABLE world_initial_data (
    id UUID PRIMARY KEY,
    data_type VARCHAR(30),             -- region/sect/npc/resource
    data_content JSONB,                -- 具体内容
    load_order INT,                    -- 加载顺序
    is_loaded BOOLEAN DEFAULT FALSE
);

-- 6. 版本记录表
CREATE TABLE version_history (
    id UUID PRIMARY KEY,
    version VARCHAR(20),
    migration_script TEXT,
    applied_at TIMESTAMP DEFAULT NOW()
);

-- 7. 缓存数据表（热数据缓存）
CREATE TABLE cache_data (
    key VARCHAR(100) PRIMARY KEY,
    value JSONB,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 8. 审计日志表
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    actor_id UUID,
    action VARCHAR(50),
    target_type VARCHAR(30),
    target_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW()
);
