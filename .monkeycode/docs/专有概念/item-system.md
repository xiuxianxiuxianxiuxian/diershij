# 物品系统

## 概述

物品系统定义了修仙世界中所有物品的类型体系，包括物品模板、背包管理、装备系统、丹药、法宝、符箓、配方和原材料。类型定义在 `shared/types/item.go` 中。

## 物品层级

```
ItemTemplate (物品蓝图)
    |
    v
InventoryItem (背包中的实例)
    |
    +-- 装备 -> EntityEquipment -> EquipmentItem
    +-- 丹药 -> Pill
    +-- 法宝 -> Artifact
    +-- 符箓 -> Talisman
    +-- 原材料 -> Material
```

## ItemTemplate (物品模板)

物品的蓝图定义：

| 属性 | 说明 |
|------|------|
| `ID` | 物品唯一标识 |
| `Name` | 物品名称 |
| `Type` | 类型: pill/artifact/talisman/material/treasure 等 |
| `SubType` | 子类型 |
| `Rank` | 品级: 天/地/玄/黄 x 上/中/下/极品 |
| `Description` | 描述 |
| `BaseValue` | 基础价值 (下品灵石为单位) |
| `Stackable` | 是否可堆叠 |
| `MaxStack` | 最大堆叠数 |
| `Usable` | 是否可使用 |
| `Consumable` | 是否可消耗 |
| `Tradeable` | 是否可交易 |
| `Droppable` | 是否可丢弃 |

## EntityInventory (实体背包)

| 属性 | 说明 |
|------|------|
| `EntityID` | 实体 ID |
| `Items` | 背包物品列表 |
| `Capacity` | 最大物品种类数 |

### InventoryItem (背包物品)

| 属性 | 说明 |
|------|------|
| `TemplateID` | 物品模板 ID |
| `InstanceID` | 唯一实例 ID |
| `Name` | 物品名称 |
| `Quantity` | 数量 |
| `Slot` | 背包槽位 |
| `Quality` | 品质 (1-100, 随机生成) |
| `Bound` | 是否绑定 |
| `ExpiryTime` | 过期时间 (0=永久) |

## EntityEquipment (装备系统)

### 装备槽位 (10 个)

| 槽位 | 属性名 | 说明 |
|------|--------|------|
| 武器 | `Weapon` | 主武器 |
| 护甲 | `Armor` | 身体护甲 |
| 头盔 | `Helmet` | 头部装备 |
| 靴子 | `Boots` | 脚部装备 |
| 项链 | `Necklace` | 颈部装备 |
| 戒指1 | `Ring1` | 第一戒指槽 |
| 戒指2 | `Ring2` | 第二戒指槽 |
| 内甲 | `InnerArmor` | 内部护甲 |
| 腰带 | `Waist` | 腰部装备 |
| 手镯 | `Bracelet` | 手部装备 |

### EquipmentItem (装备物品)

| 属性 | 说明 |
|------|------|
| `TemplateID` | 装备模板 ID |
| `InstanceID` | 唯一实例 ID |
| `Name` | 装备名称 |
| `Slot` | 装备槽位 |
| `Rank` | 品级 |
| `Quality` | 品质 (1-100) |
| `Level` | 装备等级 |
| `Durability` | 当前耐久度 |
| `MaxDurability` | 最大耐久度 |
| `Stats` | 属性加成 (map[string]float64) |
| `Enchantments` | 附魔列表 |
| `Soulbound` | 是否灵魂绑定 |

### Enchantment (附魔)

| 属性 | 说明 |
|------|------|
| `Name` | 附魔名称 |
| `Stat` | 影响的属性 |
| `Value` | 加成值 |
| `Tier` | 附魔等级 (1-10) |

## Pill (丹药)

| 属性 | 说明 |
|------|------|
| `ID` | 丹药 ID |
| `Name` | 丹药名称 |
| `Type` | 类型: cultivation/healing/breakthrough/combat 等 |
| `Rank` | 品级 |
| `Quality` | 品质: 下品/中品/上品/极品 |
| `Effect` | 效果描述 |
| `EffectValue` | 效果值 |
| `Duration` | 持续时间 (秒) |
| `SuccessRate` | 服用成功率 |
| `FailureEffect` | 失败效果 (毒性/爆炸) |
| `Toxicity` | 丹毒积累值 |
| `Cooldown` | 冷却时间 |

## Artifact (法宝)

| 属性 | 说明 |
|------|------|
| `ID` | 法宝 ID |
| `Name` | 法宝名称 |
| `Type` | 类型: flying_sword/mirror/bell/pagoda/banner 等 |
| `Rank` | 品级 |
| `Grade` | 等级: 凡器/地器/天器/古宝 |
| `AttackBonus` | 攻击加成 (map[string]float64) |
| `DefenseBonus` | 防御加成 (map[string]float64) |
| `SpecialAbility` | 特殊能力 |
| `Energy` | 法宝能量 (用于释放能力) |
| `MaxEnergy` | 最大能量 |
| `Level` | 法宝等级 |
| `Experience` | 法宝经验值 (用于升级) |
| `Refined` | 是否已被当前主人炼化 |

## Talisman (符箓)

| 属性 | 说明 |
|------|------|
| `ID` | 符箓 ID |
| `Name` | 符箓名称 |
| `Type` | 类型: attack/defense/utility/movement 等 |
| `Rank` | 品级 |
| `Effect` | 效果 |
| `EffectValue` | 效果值 |
| `Duration` | 持续时间 |
| `Cooldown` | 冷却时间 |
| `Charges` | 剩余使用次数 |
| `MaxCharges` | 最大使用次数 |

## Recipe (配方)

| 属性 | 说明 |
|------|------|
| `ID` | 配方 ID |
| `Name` | 配方名称 |
| `Type` | 类型: pill/artifact/talisman/array |
| `ResultItemID` | 产出物品 ID |
| `ResultQuantity` | 产出数量 |
| `RequiredLevel` | 所需等级 |
| `RequiredSkill` | 所需技能: alchemy/artificing 等 |
| `RequiredSkillLevel` | 所需技能等级 |
| `Materials` | 所需材料列表 |
| `BaseSuccessRate` | 基础成功率 |
| `CreatorID` | 发现/创建者 |
| `IsSecret` | 是否为秘密配方 |

### RecipeMaterial (配方材料)

| 属性 | 说明 |
|------|------|
| `ItemID` | 材料物品 ID |
| `Name` | 材料名称 |
| `Quantity` | 所需数量 |
| `Quality` | 最低品质要求 |

## Material (原材料)

| 属性 | 说明 |
|------|------|
| `ID` | 材料 ID |
| `Name` | 材料名称 |
| `Type` | 类型: herb/ore/beast_part/essence 等 |
| `Rank` | 品级 |
| `Purity` | 纯度 (1-100) |
| `Attributes` | 材料属性 (map[string]float64) |
| `Source` | 来源地 |
| `DropRate` | 基础掉落率 |

## 待确认

- 物品的获取途径（掉落、购买、制作等）
- 物品交易系统的具体实现
- 装备耐久度消耗和修复机制
- 法宝炼化的具体流程
- 丹毒积累和清除机制
