-- Migration 10: Shops, shop inventory, and auction house system
-- ============================================================================

-- ============================================================================
-- PART 1: Shops table
-- ============================================================================
CREATE TABLE IF NOT EXISTS shops (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    region_id VARCHAR(50) REFERENCES world_regions(id) ON DELETE CASCADE,
    shop_type VARCHAR(30) DEFAULT 'general',  -- general, herb, weapon, auction
    npc_owner VARCHAR(100),
    markup_rate DOUBLE PRECISION DEFAULT 1.0,  -- price multiplier
    buy_rate DOUBLE PRECISION DEFAULT 0.5,     -- percent of price when buying from player
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- PART 2: Shop inventory table — each row is a stocked item template
-- ============================================================================
CREATE TABLE IF NOT EXISTS shop_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id VARCHAR(50) NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    item_type VARCHAR(30) NOT NULL,
    rarity INTEGER DEFAULT 1,
    price BIGINT NOT NULL,                -- buy price in low-grade spirit stones
    quantity INTEGER DEFAULT -1,          -- -1 = unlimited
    refresh_hours INTEGER DEFAULT 24,     -- restock interval
    min_realm VARCHAR(30) DEFAULT 'mortal',
    UNIQUE(shop_id, item_name)
);

-- ============================================================================
-- PART 3: Auctions table — player-to-player auction listings
-- ============================================================================
CREATE TABLE IF NOT EXISTS auctions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seller_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    item_name VARCHAR(100) NOT NULL,
    quantity INTEGER DEFAULT 1,
    price BIGINT NOT NULL,              -- buyout price in low-grade spirit stones
    deposit BIGINT NOT NULL DEFAULT 0,  -- listing fee
    status VARCHAR(20) DEFAULT 'active', -- active, sold, cancelled
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    buyer_id UUID REFERENCES entities(id) ON DELETE SET NULL,
    sold_at TIMESTAMP
);

-- ============================================================================
-- PART 4: Indexes
-- ============================================================================
CREATE INDEX IF NOT EXISTS idx_shops_region ON shops(region_id);
CREATE INDEX IF NOT EXISTS idx_shops_type ON shops(shop_type);
CREATE INDEX IF NOT EXISTS idx_shop_inventory_shop ON shop_inventory(shop_id);
CREATE INDEX IF NOT EXISTS idx_auctions_seller ON auctions(seller_id);
CREATE INDEX IF NOT EXISTS idx_auctions_status ON auctions(status);
CREATE INDEX IF NOT EXISTS idx_auctions_expires ON auctions(expires_at);

-- ============================================================================
-- PART 5: Seed data — regional shops
-- ============================================================================

-- 青云镇杂货铺 (basic supplies for new cultivators)
INSERT INTO shops (id, name, description, region_id, shop_type, npc_owner, markup_rate, buy_rate)
VALUES ('qingyun_shop', '青云镇杂货铺', '青云镇唯一的杂货铺，售卖基础修炼物资', 'qingyun_town', 'general', '王掌柜', 1.0, 0.4)
ON CONFLICT (id) DO NOTHING;

INSERT INTO shop_inventory (shop_id, item_name, item_type, rarity, price, quantity, refresh_hours, min_realm)
VALUES
    ('qingyun_shop', '回气丹', 'pill', 1, 50, 100, 24, 'mortal'),
    ('qingyun_shop', '聚气散', 'pill', 1, 100, 50, 24, 'mortal'),
    ('qingyun_shop', '疗伤膏', 'pill', 1, 30, 80, 12, 'mortal'),
    ('qingyun_shop', '粗制飞剑', 'weapon', 1, 200, 10, 48, 'mortal'),
    ('qingyun_shop', '布衣', 'armor', 1, 80, 20, 48, 'mortal'),
    ('qingyun_shop', '草鞋', 'boots', 1, 40, 20, 48, 'mortal'),
    ('qingyun_shop', '铁护腕', 'armor', 1, 120, 15, 48, 'mortal')
ON CONFLICT (shop_id, item_name) DO NOTHING;

-- 灵雾山脉药材铺 (herbs and alchemy materials)
INSERT INTO shops (id, name, description, region_id, shop_type, npc_owner, markup_rate, buy_rate)
VALUES ('lingwu_herb_shop', '灵雾药材铺', '灵雾山脉脚下的药材铺，专卖珍稀草药', 'lingwu_mountains', 'herb', '药老', 1.2, 0.5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO shop_inventory (shop_id, item_name, item_type, rarity, price, quantity, refresh_hours, min_realm)
VALUES
    ('lingwu_herb_shop', '十年灵参', 'material', 2, 150, 30, 24, 'qi_condensation'),
    ('lingwu_herb_shop', '百年灵芝', 'material', 3, 500, 15, 48, 'foundation'),
    ('lingwu_herb_shop', '灵雾草', 'material', 1, 80, 50, 12, 'mortal'),
    ('lingwu_herb_shop', '凝血花', 'material', 2, 200, 25, 24, 'qi_condensation'),
    ('lingwu_herb_shop', '培元丹', 'pill', 2, 300, 20, 24, 'qi_condensation'),
    ('lingwu_herb_shop', '筑基丹', 'pill', 3, 2000, 5, 72, 'foundation')
ON CONFLICT (shop_id, item_name) DO NOTHING;

-- 中州城拍卖行 (high-end equipment and treasures)
INSERT INTO shops (id, name, description, region_id, shop_type, npc_owner, markup_rate, buy_rate)
VALUES ('zhongzhou_auction', '中州城拍卖行', '中州城最大的拍卖行，常有珍品流出', 'zhongzhou_city', 'auction', '金算子', 1.5, 0.6)
ON CONFLICT (id) DO NOTHING;

INSERT INTO shop_inventory (shop_id, item_name, item_type, rarity, price, quantity, refresh_hours, min_realm)
VALUES
    ('zhongzhou_auction', '青霜剑', 'weapon', 3, 5000, 3, 72, 'foundation'),
    ('zhongzhou_auction', '金蚕丝甲', 'armor', 3, 8000, 2, 72, 'foundation'),
    ('zhongzhou_auction', '玄铁重剑', 'weapon', 4, 20000, 1, 168, 'golden_core'),
    ('zhongzhou_auction', '天灵丹', 'pill', 4, 15000, 3, 168, 'golden_core'),
    ('zhongzhou_auction', '紫金冠', 'helmet', 3, 6000, 3, 72, 'foundation'),
    ('zhongzhou_auction', '追风靴', 'boots', 3, 4000, 4, 72, 'foundation'),
    ('zhongzhou_auction', '聚灵玉佩', 'necklace', 3, 10000, 2, 168, 'foundation')
ON CONFLICT (shop_id, item_name) DO NOTHING;
