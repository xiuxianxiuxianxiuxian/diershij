-- ==========================================
-- 洞府、地产与个人空间模块
-- ==========================================
-- 包含：洞府表、洞府建筑、灵田、灵脉、
--       宠物/妖兽栏、道侣关系、师徒传承
-- ==========================================

-- 1. 洞府表
CREATE TABLE residences (
    id UUID PRIMARY KEY,
    owner_id UUID REFERENCES entities(id),
    region_id UUID REFERENCES world_regions(id),
    name VARCHAR(100),
    rank VARCHAR(20),                  -- 洞府等级
    position JSONB,                    -- 坐标
    spiritual_density REAL,            -- 洞府内灵气浓度
    defense_level INT,                 -- 防御等级
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. 洞府建筑表
CREATE TABLE residence_buildings (
    id UUID PRIMARY KEY,
    residence_id UUID REFERENCES residences(id),
    building_type VARCHAR(30),         -- meditation_room/alchemy_lab/forging_hall/library/etc
    level INT DEFAULT 1,
    effects JSONB,                     -- 建筑效果（如 cultivation_speed_bonus: 0.2）
    upgrade_cost JSONB,                -- 升级所需材料
    built_at TIMESTAMP DEFAULT NOW()
);

-- 3. 灵田表（种植灵草）
CREATE TABLE spirit_fields (
    id UUID PRIMARY KEY,
    residence_id UUID REFERENCES residences(id),
    soil_quality INT,                  -- 土壤品质
    planted_seed_id UUID REFERENCES materials(id),
    plant_date TIMESTAMP,
    harvest_date TIMESTAMP,
    growth_progress REAL,              -- 生长进度 0-100%
    water_level REAL,                  -- 灌溉程度
    fertilizer JSONB                   -- 施肥记录
);

-- 4. 灵脉表（洞府灵脉）
CREATE TABLE spirit_veins (
    id UUID PRIMARY KEY,
    residence_id UUID REFERENCES residences(id),
    vein_type VARCHAR(20),             -- 属性灵脉（金木水火土等）
    tier INT,                          -- 灵脉品阶
    output_rate REAL,                  -- 灵气产出速率
    last_harvested TIMESTAMP,
    next_harvest TIMESTAMP
);

-- 5. 宠物/妖兽栏表
CREATE TABLE entity_beasts (
    id UUID PRIMARY KEY,
    owner_id UUID REFERENCES entities(id),
    beast_type VARCHAR(30),            -- 妖兽种类
    name VARCHAR(50),
    level INT,
    loyalty INT,                       -- 忠诚度
    combat_power REAL,
    skills JSONB,                      -- 妖兽技能
    captured_at TIMESTAMP DEFAULT NOW()
);

-- 6. 道侣关系表
CREATE TABLE dao_partners (
    entity_a_id UUID REFERENCES entities(id),
    entity_b_id UUID REFERENCES entities(id),
    bond_strength INT,                 -- 羁绊强度
    married_at TIMESTAMP DEFAULT NOW(),
    shared_benefits JSONB,             -- 双修增益等
    PRIMARY KEY (entity_a_id, entity_b_id)
);

-- 7. 师徒传承表
CREATE TABLE master_disciple (
    master_id UUID REFERENCES entities(id),
    disciple_id UUID REFERENCES entities(id),
    relationship_type VARCHAR(20),     -- direct/inherited（嫡传/再传）
    started_at TIMESTAMP DEFAULT NOW(),
    teachings JSONB,                   -- 传授内容记录
    break_reason TEXT,                 -- 若断绝关系，原因
    PRIMARY KEY (master_id, disciple_id)
);

-- 8. 结义兄弟表
CREATE TABLE sworn_siblings (
    group_id UUID,
    entity_id UUID REFERENCES entities(id),
    rank INT,                          -- 排行
    oath_text TEXT,                    -- 结义誓词
    joined_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (group_id, entity_id)
);
