# 功法系统

## 概述

功法系统定义了修仙世界中修炼方法的完整体系，包含功法定义、技能系统和功法传承机制。类型定义在 `shared/types/method.go` 中。

## CultivationMethod (功法定义)

功法包含 60+ 属性，分为以下类别：

### 基本信息

| 属性 | 说明 |
|------|------|
| `ID` | 功法唯一标识 |
| `Name` | 功法名称 |
| `CreatorID` | 创建者 ID |
| `OriginSect` | 起源宗门 |
| `Version` | 版本号 |
| `Description` | 功法描述 |

### 等级体系

**Rank (品级)**: 天/地/玄/黄 x 上/中/下/极品，共 12 级
- 天品 (Heaven) - 最高品级
- 地品 (Earth)
- 玄品 (Mystic)
- 黄品 (Yellow) - 最低品级

每品分上、中、下、极品。

**Category (分类)**:
- 主修功法 (main)
- 秘术 (secret)
- 身法 (movement)
- 神识 (divine_sense)
- 辅助 (support)
- 生活 (life)

**ElementAffinity (属性亲和)**: 金/木/水/火/土/风/雷/冰/光/暗/无

### 修炼加成

| 属性 | 说明 |
|------|------|
| `CultivationSpeedMult` | 修炼速度倍率 (如 1.5x) |
| `SpiritualPowerCapMult` | 灵力上限倍率 |
| `QiCapMult` | 灵气上限倍率 |
| `DivineSenseCapMult` | 神识上限倍率 |
| `LifespanBonus` | 额外寿命 (年) |
| `RecoverySpeedMult` | 恢复速度倍率 |

### 战斗加成

| 属性 | 类型 | 说明 |
|------|------|------|
| `AttackBonuses` | map[string]float64 | 攻击加成 (如 {"fire_damage": 1.2, "penetration": 0.5}) |
| `DefenseBonuses` | map[string]float64 | 防御加成 (如 {"magic_resist": 0.3, "damage_reduction": 0.1}) |
| `UtilityBonuses` | map[string]float64 | 辅助加成 (如 {"loot_rate": 0.1, "alchemy_success": 0.2}) |

### 特殊效果

| 属性 | 说明 |
|------|------|
| `PassiveEffects` | 被动效果列表 (如 "mana_shield", "life_steal", "reflect") |
| `ActiveSkills` | 主动技能列表 |
| `UltimateSkill` | 终极技能 (精通时解锁) |

### 法则亲和

| 属性 | 说明 |
|------|------|
| `LawAffinities` | 加速领悟的法则列表 |
| `LawComprehensionBonus` | 法则领悟速度倍率 |
| `DaoCompatibility` | 兼容的道 (如 "sword_dao" 适合剑修) |

### 修习限制

| 属性 | 说明 |
|------|------|
| `RequiredRoots` | 所需灵根 (如 ["fire", "metal"]) |
| `RequiredPhysique` | 所需体质 (如 ["pure_yang_body"]) |
| `RealmRequirement` | 最低境界要求 |
| `AlignmentRestriction` | 立场限制 (正/魔/中立/无) |
| `KarmaThreshold` | 业力阈值 (超过则无法修炼) |
| `GenderRestriction` | 性别限制 (男/女/无) |

### 传承与进化

| 属性 | 说明 |
|------|------|
| `ParentMethodID` | 衍生自哪个功法 |
| `EvolutionPath` | 进化方向 (升级版本列表) |
| `TransmissionMode` | 传承方式 (玉简/口授/血脉/神念) |
| `CanModify` | 是否可被后继者修改 |
| `Complexity` | 复杂度 (影响学习难度和反噬风险) |

### 评价属性

| 属性 | 说明 |
|------|------|
| `PowerScore` | 综合战力评分 (系统计算) |
| `Potential` | 潜力评分 (影响后期进化上限) |
| `Popularity` | 流行度 (学习该功法的实体数量) |

## Skill (技能)

功法提供的技能，分三种类型：

**Category**:
- `active`: 主动技能
- `passive`: 被动技能
- `ultimate`: 终极技能

**技能属性**:

| 属性 | 说明 |
|------|------|
| `ID` | 技能 ID |
| `Name` | 技能名称 |
| `Description` | 描述 |
| `Element` | 属性元素 |
| `DamageMult` | 伤害倍率 |
| `Cooldown` | 冷却时间 (秒) |
| `ManaCost` | 灵力消耗 |
| `Range` | 作用范围 (米) |
| `Duration` | 效果持续时间 (秒) |
| `Effects` | 效果列表 |

**SkillEffect (技能效果)**:

| 属性 | 说明 |
|------|------|
| `Type` | 效果类型: damage/heal/buff/debuff/shield/teleport 等 |
| `Value` | 效果值 |
| `Target` | 目标: self/enemy/area/ally |
| `Duration` | 持续时间 |
| `Condition` | 触发条件 |

## EntityMethod (实体已学功法)

实体学习功法后的实例化数据：

| 属性 | 说明 |
|------|------|
| `MethodID` | 功法 ID |
| `EntityID` | 实体 ID |
| `MasteryLevel` | 精通程度 (0-100%) |
| `IsMainMethod` | 是否为主修功法 |
| `LearnedAt` | 学习时间 |
| `LastPracticed` | 最后练习时间 |
| `BacklashRisk` | 当前反噬风险 (0-1) |
| `Modified` | 是否已修改 |
| `ModifiedNotes` | 修改备注 |

## 功法示例

测试文件中的示例功法 "Nine Sun Scripture"：

- **品级**: heaven_upper (天品上)
- **分类**: main (主修功法)
- **属性**: fire (火)
- **修炼速度**: 1.5x
- **被动效果**: fire_aura, heat_resistance
- **主动技能**: Sun Flare (伤害倍率 2.0, 冷却 10s, 灵力消耗 50)
- **终极技能**: Nine Suns Annihilation (伤害倍率 10.0, 冷却 300s, 灵力消耗 500)
- **法则亲和**: fire, light
- **道兼容**: sword_dao
- **需求**: fire 灵根, foundation 境界
- **复杂度**: 8
- **战力评分**: 90, 潜力: 95

## 待确认

- 功法的具体学习途径和方式 (当前仅定义了类型，未实现学习逻辑)
- 功法修改机制的具体实现
- 功法反噬的计算公式
- 功法进化的触发条件
