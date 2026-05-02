-- ==========================================
-- 物品、装备与道具模块
-- ==========================================
-- 包含：物品模板、实体背包、装备栏、丹药库、
--       法宝库、符箓库、材料库
-- ==========================================

-- 1. 物品模板表（所有物品的基础定义）
CREATE TABLE item_templates (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    item_type VARCHAR(30),             -- weapon/armor/pill/material/talisman/etc
    rank VARCHAR(20),                  -- 品阶
    element_affinity VARCHAR(20),      -- 属性倾向
    description TEXT,
    base_stats JSONB,                  -- 基础属性（攻击/防御/回复等）
    special_effects JSONB,             -- 特殊效果
    market_reference_price REAL,       -- 市场参考价（由交易历史计算）
    is_tradable BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. 实体背包表（每个实体持有的物品实例）
CREATE TABLE entity_inventory (
    id UUID PRIMARY KEY,
    entity_id UUID REFERENCES entities(id),
    item_template_id UUID REFERENCES item_templates(id),
    quantity INT DEFAULT 1,
    quality INT,                       -- 具体品质（同模板不同品质）
    durability REAL,                   -- 耐久度（装备类）
    expiry_time TIMESTAMP,             -- 过期时间（丹药/临时道具）
    acquired_at TIMESTAMP DEFAULT NOW(),
    source VARCHAR(50)                 -- 获取来源（craft/drop/trade等）
);

-- 3. 装备栏表
CREATE TABLE entity_equipment (
    entity_id UUID REFERENCES entities(id),
    slot VARCHAR(20),                  -- weapon/armor/helmet/boots/ring/necklace等
    item_instance_id UUID REFERENCES entity_inventory(id),
    PRIMARY KEY (entity_id, slot)
);

-- 4. 丹药表（丹药详细属性）
CREATE TABLE pills (
    id UUID PRIMARY KEY,
    recipe_id UUID REFERENCES recipes(id),
    name VARCHAR(100) NOT NULL,
    rank VARCHAR(20),
    quality_tier INT,                  -- 下品/中品/上品/极品
    effects JSONB,                     -- 效果（如 cultivation_boost: 100）
    duration INT,                      -- 持续时间（秒）
    side_effects JSONB,                -- 副作用（如 toxicity: 10）
    alchemist_id UUID REFERENCES entities(id),
    crafted_at TIMESTAMP DEFAULT NOW()
);

-- 5. 法宝表
CREATE TABLE artifacts (
    id UUID PRIMARY KEY,
    blueprint_id UUID REFERENCES blueprints(id),
    name VARCHAR(100) NOT NULL,
    grade VARCHAR(20),                 -- mortal/earth/heaven/ancient
    type VARCHAR(20),                  -- offensive/defensive/utility
    attack_power REAL,
    defense_power REAL,
    special_ability TEXT,              -- 特殊能力描述
    ability_cooldown INT,              -- 技能冷却（秒）
    owner_id UUID REFERENCES entities(id),
    bound BOOLEAN DEFAULT FALSE,       -- 是否认主
    refined_at TIMESTAMP DEFAULT NOW()
);

-- 6. 符箓表
CREATE TABLE talismans (
    id UUID PRIMARY KEY,
    recipe_id UUID REFERENCES recipes(id),
    name VARCHAR(100) NOT NULL,
    rank VARCHAR(20),
    effect_type VARCHAR(30),           -- attack/defense/utility
    power INT,
    charges INT DEFAULT 1,             -- 可使用次数
    caster_id UUID REFERENCES entities(id),
    inscribed_at TIMESTAMP DEFAULT NOW()
);

-- 7. 配方表（丹药/法宝/符箓/阵法配方）
CREATE TABLE recipes (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    recipe_type VARCHAR(20),           -- pill/artifact/talisman/formation
    creator_id UUID REFERENCES entities(id),
    required_level INT,                -- 所需技能等级
    materials JSONB,                   -- 所需材料清单
    output_item_type VARCHAR(30),
    base_success_rate REAL,            -- 基础成功率
    is_secret BOOLEAN DEFAULT FALSE,   -- 是否为秘方
    discovered_at TIMESTAMP DEFAULT NOW()
);

-- 8. 材料表（炼丹/炼器/阵法材料）
CREATE TABLE materials (
    id UUID PRIMARY KEY,
    item_template_id UUID REFERENCES item_templates(id),
    material_type VARCHAR(30),         -- herb/ore/monster_part/etc
    spiritual_affinity VARCHAR(20),    -- 灵气属性
    purity REAL,                       -- 纯度
    age INT,                           -- 年份（灵草等）
    origin_region UUID REFERENCES world_regions(id)
);
