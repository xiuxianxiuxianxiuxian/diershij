-- 99_seed_data.sql - 补充种子数据

-- 基础区域
INSERT INTO world_regions (id, name, description, spiritual_density, spiritual_tier, danger_level, resources, rules, lore)
VALUES
  ('qingyun_town', 'Qingyun Town', 'Starting town for cultivators', 0.3, 1, 1, '{}', '{}', 'A peaceful town where new cultivators begin their journey.'),
  ('misty_mountains', 'Misty Mountains', 'Spirit-rich mountains for cultivation', 0.6, 2, 2, '{"herbs": "abundant", "ores": "moderate"}', '{}', 'Mountains shrouded in spiritual mist.'),
  ('spirit_forest', 'Spirit Forest', 'Ancient forest with abundant resources', 0.5, 2, 2, '{"wood": "abundant", "herbs": "abundant"}', '{}', 'An ancient forest teeming with spiritual energy.'),
  ('zhongzhou_city', 'Zhongzhou City', 'Major cultivator city', 0.8, 3, 3, '{}', '{}', 'The largest cultivator city in the land.')
ON CONFLICT (id) DO NOTHING;

-- 商店
INSERT INTO shops (id, name, description, shop_type, region_id, npc_owner, markup_rate, buy_rate)
VALUES
  ('qingyun_trading', 'Qingyun Shop', 'Basic cultivation supplies', 'general', 'qingyun_town', 'Shopkeeper Wang', 1.2, 0.5),
  ('misty_herb_shop', 'Misty Herb Shop', 'Herbs and medicines', 'herbs', 'misty_mountains', 'Herbalist Li', 1.3, 0.5),
  ('zhongzhou_auction', 'Zhongzhou Auction', 'Premium item auctions', 'auction', 'zhongzhou_city', 'Auctioneer', 1.5, 0.4)
ON CONFLICT (id) DO NOTHING;

-- 商店库存
INSERT INTO shop_inventory (shop_id, item_name, item_type, rarity, price, quantity, min_realm)
VALUES
  ('qingyun_trading', 'Low-grade Spirit Stone', 'material', 1, 10, 1000, 'mortal'),
  ('qingyun_trading', 'Healing Pill', 'potion', 1, 50, 100, 'mortal'),
  ('qingyun_trading', 'Qi Gathering Powder', 'potion', 1, 100, 50, 'mortal'),
  ('misty_herb_shop', 'Ganoderma', 'material', 2, 80, 30, 'mortal'),
  ('misty_herb_shop', 'Millennial Ginseng', 'material', 3, 500, 10, 'qi_refining')
ON CONFLICT (shop_id, item_name) DO NOTHING;

-- 默认功法种子数据
INSERT INTO methods (id, name, description, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus)
SELECT gen_random_uuid(), 'Basic Breathing', 'The most basic cultivation method.', '黄级', 'mortal', 'neutral', 1.0, 0.0
WHERE NOT EXISTS (SELECT 1 FROM methods WHERE name = 'Basic Breathing');

INSERT INTO methods (id, name, description, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus)
SELECT gen_random_uuid(), 'Qingyun Heart Method', 'Basic inner method from Qingyun Town.', '黄级', 'mortal', 'neutral', 1.5, 0.1
WHERE NOT EXISTS (SELECT 1 FROM methods WHERE name = 'Qingyun Heart Method');

INSERT INTO methods (id, name, description, quality, realm_requirement, element_affinity, cultivation_multiplier, breakthrough_bonus)
SELECT gen_random_uuid(), 'Misty Art', 'Misty mountain inheritance for qi refining cultivators.', '玄级', 'qi_refining', 'water', 2.0, 0.2
WHERE NOT EXISTS (SELECT 1 FROM methods WHERE name = 'Misty Art');
