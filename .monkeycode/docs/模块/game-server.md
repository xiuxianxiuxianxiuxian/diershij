# Game Server 模块

## 概述

Game Server 是修仙世界的核心游戏逻辑服务，负责实体管理、操作执行、数据持久化。

- **端口**: 50051 (gRPC)
- **框架**: gRPC
- **数据库**: PostgreSQL (pgx)
- **缓存**: Redis
- **入口**: `server/game-server/cmd/main.go`

## 架构

```
Gateway (gRPC client) -> GameService (gRPC server)
                           |
                           v
                    OperationService
                           |
                           v
                    EntityRepository
                      /        \
                PostgreSQL    Redis (cache)
```

## 核心组件

### GameService (`service/game_service.go`)

gRPC 服务实现，提供以下 RPC 方法：

| RPC 方法 | 说明 |
|----------|------|
| `CreateEntity` | 创建新实体，初始化默认属性 |
| `AuthenticateEntity` | 按名称查找实体进行认证 |
| `ExecuteOperation` | 执行操作，委托给 OperationService |
| `GetEntity` | 按 ID 查询实体 |
| `SyncState` | 同步实体状态 |
| `StreamEntityUpdates` | 实体更新流 (当前未实现) |

**实体创建默认值** (`CreateEntity`):
- 初始境界: `RealmMortal` (凡人)
- 初始位置: `qingyun_town` (青云镇), (0, 0)
- 初始 Qi: 100, MaxQi: 100
- 初始灵力: 100, MaxSpiritualPower: 100
- 神识: 10, 悟性: 50, 体质: 50, 运气: 50
- 攻击力: 10, 防御力: 10, 速度: 10
- 心智稳定度: 50
- 寿命: 80 年
- 天道印记: "clear" (清白)

### OperationService (`service/operation.go`)

操作执行引擎，支持以下操作类型：

| 操作类型 | 方法 | 说明 |
|----------|------|------|
| `cultivate` | `executeCultivate()` | 修炼，增加修为进度 |
| `move` | `executeMove()` | 移动到指定区域和坐标 |
| `meditate` | `executeMeditate()` | 打坐恢复 Qi 和灵力 |
| `sleep` | `executeSleep()` | 休息，完全恢复 Qi 和灵力 |
| `breakthrough` | `executeBreakthrough()` | 突破境界 |

**修炼逻辑** (`executeCultivate`):
```
cultivationGain = 0.1 * (Comprehension / 100)
CultivationProgress += cultivationGain
if CultivationProgress > 100: CultivationProgress = 100
```

**突破逻辑** (`executeBreakthrough`):
```
前置条件: CultivationProgress >= 100
成功率 = 0.5 + (Luck / 200)，上限 0.8
突破后:
  - 境界提升一级
  - CultivationProgress = 0
  - MaxQi *= 1.5
  - MaxSpiritualPower *= 1.5
  - MaxLifespan = 新境界寿命
```

**境界寿命表**:

| 境界 | 寿命 (年) |
|------|-----------|
| mortal | 80 |
| qi_condensation | 120 |
| foundation | 200 |
| golden_core | 500 |
| nascent_soul | 1000 |
| soul_transformation | 3000 |
| void_refinement | 5000 |
| integration | 8000 |
| mahayana | 10000 |
| tribulation | 15000 |

### EntityRepository (`repository/database.go`)

数据持久化层，使用 PostgreSQL + Redis 双层存储：

**PostgreSQL 操作**:
- `Create()`: 插入实体记录
- `GetByID()`: 按 ID 查询实体
- `GetByName()`: 按名称查询实体 (同时加载属性)
- `Update()`: 更新实体基本信息
- `GetAttributes()`: 查询基础属性 (15 个字段)
- `UpdateAttributes()`: 更新基础属性 (UPSERT)

**Redis 操作**:
- `CacheEntity()`: 缓存实体，5 分钟过期
- `GetCachedEntity()`: 获取缓存实体

**数据库表**:
- `entities`: 存储实体基本信息 (10 个字段)
- `base_attributes`: 存储基础修炼属性 (16 个字段，含 entity_id)

> 注：`base_attributes` 表仅存储 15 个基础属性。`types.Attributes` 中定义的 83+ 属性中，其余属性（如灵根、战斗属性、生活技能、社交属性等）尚未在数据库层实现持久化。

## 启动流程

```
1. 加载配置 (环境变量)
2. 连接 PostgreSQL
3. 连接 Redis
4. 创建 EntityRepository
5. 创建 OperationService
6. 创建 GameService
7. 注册 gRPC 服务
8. 开始监听
```

## 依赖

- `google.golang.org/grpc` - gRPC 框架
- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL 连接池
- `github.com/redis/go-redis/v9` - Redis 客户端
- `github.com/cultivation-world/shared/types` - 共享类型
- `github.com/cultivation-world/shared/config` - 共享配置
- `github.com/cultivation-world/shared/errors` - 共享错误
- `github.com/cultivation-world/shared/proto/pb` - Proto 生成代码
