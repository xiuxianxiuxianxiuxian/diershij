# Heavenly Dao 模块

## 概述

Heavenly Dao (天道) 服务管理修仙世界的因果法则、天劫系统和世界平衡机制。

- **端口**: 50053 (gRPC)
- **框架**: gRPC
- **入口**: `server/heavenly-dao/cmd/main.go`

> 注：项目中存在两个 HeavenlyDao 服务实现文件。`heavenly_dao_service.go` 定义了 gRPC 服务，`heavenly_dao.go` 定义了天道核心逻辑。

## 核心组件

### HeavenlyDaoService (`heavenly_dao_service.go`)

gRPC 服务实现，包含三个子引擎：

**KarmaEngine (因果引擎)**:
- `UpdateKarma()`: 更新实体因果 (当前为占位实现)

**TribulationManager (天劫管理器)**:
- `TriggerTribulation()`: 触发天劫
- `CheckTribulation()`: 检查是否需要天劫

**LawManager (法则管理器)**:
- `EnforceLaws()`: 执行世界法则 (当前为占位实现)

### HeavenlyDaoService 核心逻辑 (`heavenly_dao.go`)

另一个 HeavenlyDaoService 实现，包含完整的天道逻辑：

#### 因果评估 (`EvaluateKarma`)

根据行动类型评估因果变化：

| 行动类型 | 业力变化 | 说明 |
|----------|----------|------|
| `kill_innocent` | +500 | 杀害无辜，业力缠身 |
| `kill_cultivator` | +200 | 同修相残，有损天道 |
| `kill_demon` | -50 | 斩妖除魔，功德无量 |
| `save_life` | -100 | 救人性命，功德加身 |
| `teach_method` | -200 | 传道授业，功德无量 |
| `betray_master` | +1000 | 欺师灭祖，业力滔天 |
| `break_oath` | +300 | 背信弃义，业力加身 |
| `destroy_sect` | +800 | 毁人道统，业力深重 |
| `create_method` | -500 | 开宗立派，功德无量 |

> 注：负值代表功德（善行），正值代表业力（恶行）。

#### 天道印记

根据业力值计算天道印记等级：

| 业力阈值 | 印记等级 |
|----------|----------|
| < 100 | clear (清白) |
| < 500 | slight (轻微) |
| < 1000 | heavy (沉重) |
| < 5000 | notorious (臭名昭著) |
| >= 5000 | heaven_fury (天怒) |

#### 天劫系统 (`CheckTribulation`)

**天劫基础概率**:

| 境界 | 基础概率 |
|------|----------|
| qi_condensation | 0.1 |
| foundation | 0.2 |
| golden_core | 0.3 |
| nascent_soul | 0.5 |
| soul_transformation | 0.6 |
| void_refinement | 0.7 |
| integration | 0.8 |
| mahayana | 0.9 |
| tribulation | 0.95 |

**天劫强度计算**:
```
baseStrength = 100.0
multiplier = realmMultiplier[realm] (1.0 ~ 500.0)
karmaFactor = 1.0 + karma/500.0
strength = baseStrength * multiplier * karmaFactor
```

**天劫类型**:

| 境界 | 天劫类型 |
|------|----------|
| qi_condensation | thunder (雷劫) |
| foundation | thunder (雷劫) |
| golden_core | thunder_fire (雷火劫) |
| nascent_soul | thunder_fire_wind (雷火风劫) |
| soul_transformation | five_element (五行劫) |
| void_refinement | heart_demon (心魔劫) |
| integration | dao_tribulation (道劫) |
| mahayana | extinction (寂灭劫) |
| tribulation | ascension (飞升劫) |

#### 因果衰减 (`ApplyKarmaDecay`)

```
decayAmount = oldKarma * karmaDecayRate (0.01)
newKarma = oldKarma - decayAmount
```

#### 世界平衡检查 (`BalanceCheck`)

监控以下指标并提出调整建议：

| 指标 | 阈值 | 调整建议 |
|------|------|----------|
| PowerDistribution > 0.7 | spawn_opportunity_for_weak (为弱者生成机缘) |
| ResourceCirculation < 0.1 | increase_resource_spawn (增加资源生成) |
| KarmaDistribution > 500 | trigger_heavenly_cleansing (触发天道净化) |

## gRPC 接口 (HeavenlyDaoService - heavenly_dao.proto)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `EvaluateKarma` | `KarmaRequest` | `KarmaResponse` | 评估行动的因果变化 |
| `CheckTribulation` | `TribulationRequest` | `TribulationResponse` | 检查天劫触发概率和强度 |
| `BalanceCheck` | `BalanceCheckRequest` | `BalanceCheckResponse` | 世界平衡检查 |
| `ApplyKarmaDecay` | `DecayRequest` | `DecayResponse` | 应用因果衰减 |
| `WatchHeavenlyEvents` | `HeavenlyEventRequest` | `stream HeavenlyEvent` | 天道事件流 (每 10 秒) |

## 配置

天道相关配置在 `shared/config/config.go` 的 `HeavenlyDaoConfig` 中：

```go
type HeavenlyDaoConfig struct {
    KarmaDecayRate     float64              // 因果衰减率 (默认 0.01)
    TribulationBase    map[string]float64   // 天劫基础概率
    RealmLifespan      map[string]int       // 境界寿命
    KarmaThresholds    map[string]int       // 因果阈值
}
```

## 依赖

- `google.golang.org/grpc` - gRPC 框架
- `github.com/cultivation-world/shared/types` - 共享类型
- `github.com/cultivation-world/shared/proto/pb` - Proto 生成代码
