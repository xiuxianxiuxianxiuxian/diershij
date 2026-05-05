-- 功法参考表 + 种子数据
-- 在 entity_methods 表基础上建立 methods 参考表

CREATE TABLE IF NOT EXISTS methods (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    quality VARCHAR(10) NOT NULL CHECK (quality IN ('黄级', '玄级', '地级', '天级')),
    realm_requirement VARCHAR(30) NOT NULL,
    element_affinity VARCHAR(20) DEFAULT '',  -- 元素亲和, ''表示无属性
    cultivation_multiplier DOUBLE PRECISION DEFAULT 1.0,  -- 修炼效率倍率
    breakthrough_bonus DOUBLE PRECISION DEFAULT 0,  -- 突破成功率加成
    description TEXT DEFAULT ''
);

-- ========== 凡人期 (mortal) ==========
INSERT INTO methods (id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description)
VALUES
    ('a0000000-0000-0000-0000-000000000001', '基础吐纳术', '黄级', 'mortal', '', 1.0, 0, '最基础的呼吸吐纳之法，凡人入门的必修功课。'),
    ('a0000000-0000-0000-0000-000000000002', '养气诀', '黄级', 'mortal', '', 1.1, 0.01, '温养体内先天之气的法门，修炼速度略有提升。')
ON CONFLICT (id) DO NOTHING;

-- ========== 练气期 (qi_condensation) ==========
INSERT INTO methods (id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description)
VALUES
    ('a0000000-0000-0000-0000-000000000003', '火灵诀', '黄级', 'qi_condensation', 'fire', 1.2, 0.02, '火属性基础功法，吸纳天地火灵气锤炼自身。'),
    ('a0000000-0000-0000-0000-000000000004', '青木诀', '黄级', 'qi_condensation', 'wood', 1.2, 0.02, '木属性基础功法，借草木生机滋养经脉。'),
    ('a0000000-0000-0000-0000-000000000005', '流水诀', '黄级', 'qi_condensation', 'water', 1.2, 0.02, '水属性基础功法，以柔克刚润物无声。')
ON CONFLICT (id) DO NOTHING;

-- ========== 筑基期 (foundation) ==========
INSERT INTO methods (id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description)
VALUES
    ('a0000000-0000-0000-0000-000000000006', '烈焰诀', '玄级', 'foundation', 'fire', 1.5, 0.05, '烈火焚天，以狂暴火灵力淬炼根基。'),
    ('a0000000-0000-0000-0000-000000000007', '寒冰诀', '玄级', 'foundation', 'ice', 1.5, 0.05, '极寒冰法，凝天地寒气固本培元。'),
    ('a0000000-0000-0000-0000-000000000008', '金刚诀', '玄级', 'foundation', 'metal', 1.5, 0.05, '金行功法，铸就金刚不坏之基。')
ON CONFLICT (id) DO NOTHING;

-- ========== 金丹期 (golden_core) ==========
INSERT INTO methods (id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description)
VALUES
    ('a0000000-0000-0000-0000-000000000009', '大日焚天诀', '玄级', 'golden_core', 'fire', 1.8, 0.08, '如大日当空，焚尽万物，修炼至极致可熔炼金丹。'),
    ('a0000000-0000-0000-0000-00000000000a', '九转玄功', '地级', 'golden_core', '', 2.0, 0.12, '上古炼气士传承，九转成丹，无属性限制。'),
    ('a0000000-0000-0000-0000-00000000000b', '万剑诀', '玄级', 'golden_core', 'metal', 1.8, 0.08, '以金灵气凝万剑归宗，攻守兼备。')
ON CONFLICT (id) DO NOTHING;

-- ========== 元婴期 (nascent_soul) ==========
INSERT INTO methods (id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description)
VALUES
    ('a0000000-0000-0000-0000-00000000000c', '天地烘炉诀', '地级', 'nascent_soul', 'fire', 2.5, 0.15, '以天地为烘炉，以万物为薪炭，铸就不灭元婴。'),
    ('a0000000-0000-0000-0000-00000000000d', '太虚诀', '地级', 'nascent_soul', '', 2.3, 0.18, '参悟太虚大道，元婴遨游太虚，不受五行束缚。'),
    ('a0000000-0000-0000-0000-00000000000e', '五行轮回诀', '地级', 'nascent_soul', '', 2.8, 0.10, '五行轮转生生不息，修炼极快但突破略难。')
ON CONFLICT (id) DO NOTHING;

-- ========== 化神期 (soul_transformation) ==========
INSERT INTO methods (id, name, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus, description)
VALUES
    ('a0000000-0000-0000-0000-00000000000f', '混沌天功', '天级', 'soul_transformation', '', 3.0, 0.20, '直指大道的无上功法，引混沌之气淬炼神魂。'),
    ('a0000000-0000-0000-0000-000000000010', '阴阳玄法', '地级', 'soul_transformation', '', 2.8, 0.22, '调和阴阳，化神返虚，对突破有极大助益。')
ON CONFLICT (id) DO NOTHING;
