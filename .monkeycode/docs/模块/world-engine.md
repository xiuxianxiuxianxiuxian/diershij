# World Engine 模块

## 概述

World Engine 管理修仙世界的区域、资源和事件，维护世界状态和平衡指标。

- **端口**: 50054 (gRPC)
- **框架**: gRPC
- **入口**: `server/world-engine/cmd/main.go`

> 注：项目中存在两个 World 相关服务文件 (`world_engine.go` 和 `world_service.go`)，它们定义了不同的 WorldService 实现。启动入口 (`cmd/main.go`) 使用的是 `NewWorldService()` 来自 `world_service.go`。

## 核心组件

### WorldEngineService (`world_engine.go`)

更完整的世界引擎实现，包含初始化和 Epoch 推进机制。

**初始区域** (`initializeWorld()`):

| 区域 ID | 名称 | 灵气密度 | 灵气等级 | 危险等级 | 描述 |
|---------|------|----------|----------|----------|------|
| `east_wilderness` | 东荒域 | 50 | 3 | 2 | 东荒大地，灵气充沛，是新手修士的聚集地 |
| `qingyun_town` | 青云镇 | 30 | 1 | 0 | 凡人城镇，修士的起点 |
| `spirit_mist_mountain` | 灵雾山脉 | 60 | 4 | 3 | 灵气浓郁的山脉，适合修炼 |
| `south_ridge` | 南岭域 | 70 | 5 | 4 | 南岭大地，火属性灵气浓郁 |
| `central_state` | 中州域 | 90 | 8 | 6 | 世界中心，灵气最浓郁之地 |

**区域资源示例**:
- 东荒域: 灵草 (herb, 稀有度 1, 数量 100), 灵石矿 (ore, 稀有度 2, 数量 50)
- 灵雾山脉: 千年灵草 (herb, 稀有度 3, 数量 20)
- 南岭域: 火灵石 (ore, 稀有度 4, 数量 30)
- 中州域: 天材地宝 (treasure, 稀有度 5, 数量 10)

**Epoch 推进** (`AdvanceEpoch()`):
- Epoch 计数器 +1
- 区域资源自动刷新: `Quantity += RespawnRate * 100` (上限 100)
- 更新世界状态时间戳

### WorldService (`world_service.go`)

简化的世界服务实现，当前在 `cmd/main.go` 中使用。

**功能**:
- `GetWorldState()`: 获取世界状态 (WorldTime, Cycle, Regions[])
- `GetRegion()`: 查询单个区域
- `StreamWorldEvents()`: 每 5 秒推送一次世界事件 tick
- `ModifyWorld()`: 修改世界状态

> 注：`world_service.go` 中的 Region 定义使用了额外的 `Type` 和 `SpiritualRichness` 字段，与 `types.Region` 定义不完全一致，待确认是否为代码不一致。

## gRPC 接口 (WorldService - world.proto)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `GetRegion` | `RegionRequest` | `RegionResponse` | 查询区域信息 |
| `SpawnResources` | `SpawnRequest` | `SpawnResponse` | 在指定区域生成资源 |
| `TriggerEvent` | `EventRequest` | `EventResponse` | 触发世界事件 |
| `GetWorldState` | `WorldStateRequest` | `WorldStateResponse` | 获取世界完整状态 |

## 世界状态

```go
type WorldState struct {
    Epoch          int64                     // 纪元计数
    Regions        map[RegionID]Region        // 区域映射
    ActiveEvents   []WorldEvent               // 活跃事件
    BalanceMetrics BalanceMetrics            // 平衡指标
    LastUpdated    time.Time                  // 最后更新时间
}

type BalanceMetrics struct {
    PowerDistribution   float64  // 力量分布 (0-1)
    ResourceCirculation float64  // 资源流通 (0-1)
    SectDiversity       float64  // 宗门多样性 (0-1)
    KarmaDistribution   float64  // 因果分布 (0-1)
}
```

## 依赖

- `google.golang.org/grpc` - gRPC 框架
- `github.com/google/uuid` - UUID 生成
- `github.com/cultivation-world/shared/types` - 共享类型
- `github.com/cultivation-world/shared/proto/pb` - Proto 生成代码
