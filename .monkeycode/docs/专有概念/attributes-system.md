# 属性系统

## 概述

实体属性系统包含 83+ 个核心属性，覆盖修仙世界的各个维度。属性定义在 `shared/types/entity.go` 的 `Attributes` 结构体中。

> 注：当前数据库 `base_attributes` 表仅存储 15 个基础属性。其余 68+ 属性已定义在类型系统中，但尚未在数据库层实现持久化。

## 属性分类

### 1. 基础属性 (4 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `Age` | int | 年龄 |
| `Gender` | string | 性别 |
| `Appearance` | int | 外貌值 |
| `Charisma` | int | 魅力值 |

### 2. 修炼属性 (9 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `Qi` | float64 | 当前灵气值 |
| `MaxQi` | float64 | 最大灵气值 |
| `SpiritualPower` | float64 | 当前灵力值 |
| `MaxSpiritualPower` | float64 | 最大灵力值 |
| `DivineSense` | float64 | 神识强度 |
| `Comprehension` | int | 悟性 (影响修炼速度) |
| `Constitution` | int | 体质 |
| `Luck` | int | 运气 (影响突破成功率) |
| `CultivationProgress` | float64 | 修炼进度 (0-100) |

**已持久化到数据库**: 以上 9 个属性中，除 DivineSense 外的属性已存储在 `base_attributes` 表中。

### 3. 战斗属性 (9 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `AttackPower` | float64 | 攻击力 |
| `Defense` | float64 | 防御力 |
| `Speed` | float64 | 速度 |
| `CritRate` | float64 | 暴击率 (%) |
| `CritDamage` | float64 | 暴击伤害 (%) |
| `DodgeRate` | float64 | 闪避率 (%) |
| `HitRate` | float64 | 命中率 (%) |
| `Penetration` | float64 | 穿透力 |
| `DamageReduction` | float64 | 伤害减免 (%) |

### 4. 灵根系统 (4 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `SpiritualRoots` | []SpiritualRoot | 灵根列表 (金木水火土风雷冰光暗等) |
| `RootPurity` | int | 灵根纯度 (1-100) |
| `RootAwakened` | bool | 灵根是否觉醒 |
| `MutatedRoot` | string | 变异灵根类型 |

**SpiritualRoot 结构**:
```go
type SpiritualRoot struct {
    Element string  // gold, wood, water, fire, earth, wind, thunder, ice, light, dark
    Purity  int     // 1-100
}
```

### 5. 心境系统 (5 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `MentalStability` | int | 心智稳定度 |
| `ObsessionCount` | int | 执念数量 |
| `DaoHeart` | int | 道心 |
| `InnerDemonResistance` | int | 心魔抗性 |
| `Enlightenment` | int | 领悟度 |

### 6. 生活技能 (8 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `AlchemyLevel` | int | 炼丹等级 |
| `ArtificingLevel` | int | 炼器等级 |
| `FormationLevel` | int | 阵法等级 |
| `FireControl` | int | 控火术 |
| `HerbKnowledge` | int | 灵草辨识 |
| `MiningSkill` | int | 采矿技能 |
| `TalismanSkill` | int | 符箓技能 |
| `BeastTaming` | int | 御兽术 |

### 7. 社交属性 (9 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `Reputation` | int | 声望 |
| `SectContribution` | int | 宗门贡献 |
| `FactionStandings` | map[string]int | 各势力立场 |
| `RelationshipCount` | int | 人际关系数量 |
| `MentorID` | string | 师尊 ID |
| `DiscipleIDs` | []string | 弟子 ID 列表 |
| `SwornSiblings` | []string | 结拜兄弟/姐妹 |
| `Enemies` | []string | 仇敌列表 |
| `Lovers` | []string | 道侣列表 |

### 8. 财富属性 (4 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `SpiritStones` | SpiritStones | 灵石 (四级: 下品/中品/上品/极品) |
| `PropertyValue` | int | 财产总值 |
| `RealEstate` | []string | 房产列表 |
| `BusinessIncome` | int | 商业收入 |

**SpiritStones 结构**:
```go
type SpiritStones struct {
    LowGrade     int64  // 下品灵石
    MediumGrade  int64  // 中品灵石
    HighGrade    int64  // 上品灵石
    PremiumGrade int64  // 极品灵石
}
```

### 9. 特殊属性 (6 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `Bloodline` | string | 血脉类型 |
| `BloodlinePurity` | int | 血脉纯度 |
| `Physique` | string | 体质类型 |
| `PhysiqueAwakened` | bool | 体质是否觉醒 |
| `Destiny` | int | 天命值 |
| `WorldFavor` | int | 世界眷顾度 |

### 10. 法则属性 (5 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `Laws` | map[string]float64 | 法则领悟度 (key: 法则名称) |
| `LawResonance` | int | 法则共鸣度 |
| `DomainPower` | float64 | 领域力量 |
| `DomainRange` | float64 | 领域范围 |
| `LawSuppression` | float64 | 法则压制力 |

### 11. 道属性 (6 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `DaoSeedType` | string | 道种类型 |
| `DaoSeedLevel` | int | 道种等级 |
| `DaoSeedGrowth` | float64 | 道种成长度 |
| `DaoMarks` | int | 道痕数量 |
| `DaoHeartComprehension` | int | 道心领悟 |
| `DestinyPath` | string | 命运之路 |

### 12. 寿命属性 (3 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `RemainingLifespan` | int | 剩余寿命 (年) |
| `MaxLifespan` | int | 最大寿命 (年) |
| `AgingPenalty` | float64 | 衰老惩罚系数 |

### 13. 状态效果 (5 个)

| 属性 | 类型 | 说明 |
|------|------|------|
| `Injuries` | []Injury | 伤势列表 |
| `Buffs` | []Buff | 增益效果列表 |
| `Debuffs` | []Debuff | 减益效果列表 |
| `PoisonLevel` | int | 中毒等级 |
| `CurseLevel` | int | 诅咒等级 |

**Injury 结构**:
```go
type Injury struct {
    Type        string  // 伤势类型 (internal/external)
    Severity    int     // 严重程度
    Cause       string  // 原因
    HealTime    int64   // 愈合时间 (unix timestamp)
    Description string  // 描述
}
```

**Buff 结构**:
```go
type Buff struct {
    Name       string  // 增益名称
    Effect     string  // 效果类型
    Value      float64 // 效果值
    Source     string  // 来源 (pill/technique/etc)
    ExpiryTime int64   // 过期时间
}
```

**Debuff 结构**: 同 Buff。

## 属性持久化状态

| 类别 | 数据库 | Proto | 说明 |
|------|--------|-------|------|
| 基础属性 | 部分 | 完整 | base_attributes 表存储大部分修炼属性 |
| 战斗属性 | 未实现 | 未定义 | 类型已定义 |
| 灵根 | 未实现 | 未定义 | 类型已定义 |
| 心境 | 部分 | 未定义 | MentalStability 已存储 |
| 生活技能 | 未实现 | 未定义 | 类型已定义 |
| 社交属性 | 未实现 | 未定义 | 类型已定义 |
| 财富属性 | 未实现 | SpiritStones | SpiritStones Proto 已定义 |
| 特殊属性 | 未实现 | 未定义 | 类型已定义 |
| 法则属性 | 未实现 | 未定义 | 类型已定义 |
| 道属性 | 未实现 | 未定义 | 类型已定义 |
| 寿命属性 | 部分 | 未定义 | RemainingLifespan/MaxLifespan 已存储 |
| 状态效果 | 未实现 | 未定义 | 类型已定义 |
