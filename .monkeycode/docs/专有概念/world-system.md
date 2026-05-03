# 世界系统

## 概述

世界系统管理修仙世界的地理区域、资源分布、世界事件和整体平衡。类型定义在 `shared/types/world.go` 中。

## 区域系统

### Region (区域)

| 属性 | 说明 |
|------|------|
| `ID` | 区域唯一标识 |
| `Name` | 区域名称 |
| `ParentRegionID` | 父区域 ID (支持层级结构) |
| `SpiritualDensity` | 灵气密度 |
| `SpiritualTier` | 灵气等级 |
| `DangerLevel` | 危险等级 |
| `Resources` | 资源列表 |
| `Rules` | 区域规则 |
| `Description` | 描述 |
| `Lore` | 背景故事 |

### 初始区域

| 区域 ID | 名称 | 父区域 | 灵气密度 | 灵气等级 | 危险等级 | 描述 |
|---------|------|--------|----------|----------|----------|------|
| `east_wilderness` | 东荒域 | - | 50 | 3 | 2 | 东荒大地，灵气充沛，是新手修士的聚集地 |
| `qingyun_town` | 青云镇 | 东荒域 | 30 | 1 | 0 | 凡人城镇，修士的起点 |
| `spirit_mist_mountain` | 灵雾山脉 | 东荒域 | 60 | 4 | 3 | 灵气浓郁的山脉，适合修炼 |
| `south_ridge` | 南岭域 | - | 70 | 5 | 4 | 南岭大地，火属性灵气浓郁 |
| `central_state` | 中州域 | - | 90 | 8 | 6 | 世界中心，灵气最浓郁之地 |

### RegionRules (区域规则)

| 属性 | 说明 |
|------|------|
| `IsRestricted` | 是否为限制区域 |
| `RestrictedBy` | 限制方 |
| `TaxRate` | 税率 |
| `ForbiddenActions` | 禁止行为列表 |

## 资源系统

### Resource (资源)

| 属性 | 说明 |
|------|------|
| `ID` | 资源 ID |
| `Name` | 资源名称 |
| `Type` | 类型: herb/ore/treasure 等 |
| `Rarity` | 稀有度 |
| `Quantity` | 数量 |
| `RespawnRate` | 刷新率 |
| `LastHarvested` | 最后采集时间 |

### 初始资源分布

| 区域 | 资源 | 类型 | 稀有度 | 数量 | 刷新率 |
|------|------|------|--------|------|--------|
| 东荒域 | 灵草 | herb | 1 | 100 | 0.1 |
| 东荒域 | 灵石矿 | ore | 2 | 50 | 0.05 |
| 灵雾山脉 | 千年灵草 | herb | 3 | 20 | 0.02 |
| 南岭域 | 火灵石 | ore | 4 | 30 | 0.03 |
| 中州域 | 天材地宝 | treasure | 5 | 10 | 0.01 |

### 资源刷新

Epoch 推进时，区域资源自动刷新：

```
if Quantity < 100:
    Quantity += RespawnRate * 100
```

## 世界事件

### WorldEvent (世界事件)

| 属性 | 说明 |
|------|------|
| `ID` | 事件 ID |
| `Name` | 事件名称 |
| `Type` | 事件类型 |
| `Description` | 事件描述 |
| `RegionID` | 发生区域 |
| `StartTime` | 开始时间 |
| `EndTime` | 结束时间 (可为空) |
| `Participants` | 参与者列表 |
| `Status` | 状态 |

## 世界状态

### WorldState (世界状态)

| 属性 | 说明 |
|------|------|
| `Epoch` | 纪元计数 |
| `Regions` | 区域映射 |
| `ActiveEvents` | 活跃事件列表 |
| `BalanceMetrics` | 平衡指标 |
| `LastUpdated` | 最后更新时间 |

### BalanceMetrics (平衡指标)

| 属性 | 说明 | 初始值 |
|------|------|--------|
| `PowerDistribution` | 力量分布 (0-1) | 0.5 |
| `ResourceCirculation` | 资源流通 (0-1) | 0.3 |
| `SectDiversity` | 宗门多样性 (0-1) | 0.4 |
| `KarmaDistribution` | 因果分布 | 0.5 |

## 世界引擎功能

World Engine 服务提供以下能力：

| 功能 | 方法 | 说明 |
|------|------|------|
| 查询区域 | `GetRegion()` | 获取区域信息 |
| 生成资源 | `SpawnResources()` | 在指定区域生成指定数量的资源 |
| 触发事件 | `TriggerEvent()` | 触发世界事件 |
| 获取世界状态 | `GetWorldState()` | 获取完整世界状态 |
| 推进纪元 | `AdvanceEpoch()` | 推进世界纪元，刷新资源 |
| 事件流 | `StreamWorldEvents()` | 每 5 秒推送世界事件 tick |

## 世界平衡机制

Heavenly Dao 服务监控世界平衡并提出调整建议：

| 指标 | 阈值 | 调整建议 |
|------|------|----------|
| PowerDistribution > 0.7 | 力量过于集中 | spawn_opportunity_for_weak (为弱者生成机缘) |
| ResourceCirculation < 0.1 | 资源流通不足 | increase_resource_spawn (增加资源生成) |
| KarmaDistribution > 500 | 因果失衡 | trigger_heavenly_cleansing (触发天道净化) |

## 待确认

- 区域之间的移动规则（是否有进入限制）
- 世界事件的类型和触发条件
- 资源采集的具体流程和限制
- Epoch 推进的时间间隔
- 新区域的创建机制
