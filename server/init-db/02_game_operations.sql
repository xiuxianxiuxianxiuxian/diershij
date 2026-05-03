-- 游戏业务操作相关数据库表

-- 扩展 base_attributes 表，添加生活技能字段
ALTER TABLE base_attributes
ADD COLUMN IF NOT EXISTS alchemy_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS artificing_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS mining_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS herb_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS talisman_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS formation_level INTEGER DEFAULT 0;

-- 扩展 base_attributes 表，添加战斗属性、社交属性等
ALTER TABLE base_attributes
ADD COLUMN IF NOT EXISTS crit_rate DOUBLE PRECISION DEFAULT 0,
ADD COLUMN IF NOT EXISTS crit_damage DOUBLE PRECISION DEFAULT 0,
ADD COLUMN IF NOT EXISTS dodge_rate DOUBLE PRECISION DEFAULT 0,
ADD COLUMN IF NOT EXISTS hit_rate DOUBLE PRECISION DEFAULT 0,
ADD COLUMN IF NOT EXISTS penetration DOUBLE PRECISION DEFAULT 0,
ADD COLUMN IF NOT EXISTS damage_reduction DOUBLE PRECISION DEFAULT 0,
ADD COLUMN IF NOT EXISTS fire_control INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS beast_taming INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS reputation INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS sect_contribution INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS dao_heart INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS enlightenment INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS property_value INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS business_income INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS root_purity INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS poison_level INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS curse_level INTEGER DEFAULT 0;

-- 物品表
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL, -- weapon, armor, pill, material, talisman, etc
    rarity INTEGER DEFAULT 1, -- 1-5, 对应 黄玄地天神
    description TEXT,
    attributes JSONB, -- 物品属性，如攻击力、防御力等
    stackable BOOLEAN DEFAULT false,
    max_stack INTEGER DEFAULT 1,
    usable BOOLEAN DEFAULT false,
    level_requirement INTEGER DEFAULT 0,
    realm_requirement VARCHAR(30) DEFAULT 'mortal',
    created_at TIMESTAMP DEFAULT NOW()
);

-- 玩家背包表
CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 1,
    equipped BOOLEAN DEFAULT false,
    slot VARCHAR(20), -- head, body, weapon, accessory, etc
    durability INTEGER, -- 耐久度
    bound BOOLEAN DEFAULT false, -- 是否绑定
    acquired_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(entity_id, item_id, slot)
);

-- 配方表
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL, -- alchemy, artificing, talisman
    difficulty INTEGER DEFAULT 1, -- 1-10
    description TEXT,
    materials JSONB NOT NULL, -- [{"item_id": "xxx", "quantity": 5}, ...]
    result_item_id UUID REFERENCES items(id),
    result_quantity INTEGER DEFAULT 1,
    skill_level_required INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 玩家已学配方表
CREATE TABLE IF NOT EXISTS entity_recipes (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    recipe_id UUID REFERENCES recipes(id) ON DELETE CASCADE,
    learned_at TIMESTAMP DEFAULT NOW(),
    proficiency INTEGER DEFAULT 0, -- 熟练度
    PRIMARY KEY (entity_id, recipe_id)
);

-- 法术表
CREATE TABLE IF NOT EXISTS spells (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL, -- attack, heal, buff, debuff
    element VARCHAR(20), -- fire, water, earth, metal, wood, etc
    cost INTEGER DEFAULT 10, -- 灵气消耗
    base_damage INTEGER DEFAULT 0,
    base_heal INTEGER DEFAULT 0,
    duration INTEGER DEFAULT 0, -- 持续时间（秒）
    cooldown INTEGER DEFAULT 0, -- 冷却时间（秒）
    description TEXT,
    realm_requirement VARCHAR(30) DEFAULT 'mortal',
    created_at TIMESTAMP DEFAULT NOW()
);

-- 玩家已学法术表
CREATE TABLE IF NOT EXISTS entity_spells (
    entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
    spell_id UUID REFERENCES spells(id) ON DELETE CASCADE,
    learned_at TIMESTAMP DEFAULT NOW(),
    proficiency INTEGER DEFAULT 0,
    last_cast_at TIMESTAMP,
    PRIMARY KEY (entity_id, spell_id)
);

-- 消息表
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID REFERENCES entities(id),
    receiver_id UUID REFERENCES entities(id), -- NULL for broadcast
    type VARCHAR(20) NOT NULL, -- private, sect, world, system
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 战斗日志表
CREATE TABLE IF NOT EXISTS combat_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attacker_id UUID REFERENCES entities(id),
    defender_id UUID REFERENCES entities(id),
    damage_dealt INTEGER DEFAULT 0,
    damage_received INTEGER DEFAULT 0,
    is_crit BOOLEAN DEFAULT false,
    is_dodge BOOLEAN DEFAULT false,
    skill_used VARCHAR(100),
    result VARCHAR(20), -- win, lose, draw, escape
    location_region_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 探索日志表
CREATE TABLE IF NOT EXISTS explore_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES entities(id),
    region_id VARCHAR(50),
    result_type VARCHAR(30), -- resource, event, entity, nothing
    result_data JSONB,
    qi_consumed INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 采集日志表
CREATE TABLE IF NOT EXISTS gather_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES entities(id),
    region_id VARCHAR(50),
    resource_type VARCHAR(50),
    resource_name VARCHAR(100),
    quantity INTEGER DEFAULT 0,
    skill_exp_gained INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 制作日志表
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

-- 交易记录表（扩展原有 transactions 表）
ALTER TABLE transactions
ADD COLUMN IF NOT EXISTS item_quantity INTEGER DEFAULT 1,
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'completed'; -- pending, completed, cancelled

-- 功法创造记录表
CREATE TABLE IF NOT EXISTS method_creation_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID REFERENCES entities(id),
    method_name VARCHAR(100),
    method_type VARCHAR(30),
    method_rank VARCHAR(20), -- 黄玄地天
    qi_consumed INTEGER,
    divine_sense_consumed INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 施法记录表
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

-- 创建索引
CREATE INDEX idx_inventory_entity ON inventory(entity_id);
CREATE INDEX idx_inventory_item ON inventory(item_id);
CREATE INDEX idx_inventory_equipped ON inventory(equipped);
CREATE INDEX idx_entity_recipes_entity ON entity_recipes(entity_id);
CREATE INDEX idx_entity_spells_entity ON entity_spells(entity_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_receiver ON messages(receiver_id);
CREATE INDEX idx_messages_type ON messages(type);
CREATE INDEX idx_messages_created ON messages(created_at);
CREATE INDEX idx_combat_logs_attacker ON combat_logs(attacker_id);
CREATE INDEX idx_combat_logs_defender ON combat_logs(defender_id);
CREATE INDEX idx_combat_logs_created ON combat_logs(created_at);
CREATE INDEX idx_explore_logs_entity ON explore_logs(entity_id);
CREATE INDEX idx_gather_logs_entity ON gather_logs(entity_id);
CREATE INDEX idx_craft_logs_entity ON craft_logs(entity_id);
