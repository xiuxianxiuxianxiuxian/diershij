-- ==========================================
-- AI 决策与 NPC 行为模块
-- ==========================================
-- 包含：NPC 人格配置、决策日志
-- ==========================================

-- 1. NPC 人格配置表
CREATE TABLE npc_personalities (
    npc_id UUID PRIMARY KEY REFERENCES entities(id),
    personality_type VARCHAR(30),      -- 性格类型
    moral_alignment VARCHAR(20),       -- 道德倾向
    ambition_level INTEGER,            -- 野心程度
    risk_tolerance REAL,               -- 风险承受度
    social_preference VARCHAR(20),     -- 社交偏好
    background_story TEXT,             -- 背景故事
    current_goal TEXT,                 -- 当前目标
    hidden_secrets JSONB,              -- 隐藏秘密
    llm_system_prompt TEXT,            -- DeepSeek 系统提示词
    behavior_tree_config JSONB,        -- 行为树配置
    initial_actions JSONB              -- 初始行为模式
);

-- 2. NPC 决策日志表
CREATE TABLE npc_decision_log (
    id UUID PRIMARY KEY,
    npc_id UUID REFERENCES entities(id),
    decision_type VARCHAR(30),
    context JSONB,                     -- 决策上下文
    action_taken JSONB,                -- 采取的行动
    reasoning TEXT,                    -- 决策推理（DeepSeek 输出）
    model_used VARCHAR(20),            -- deepseek-chat / deepseek-reasoner
    source VARCHAR(10),                -- 'behavior_tree' | 'llm'
    token_cost REAL,                   -- API 调用成本
    timestamp TIMESTAMP DEFAULT NOW()
);
