# 修仙世界 — 微服务 MUD 修仙游戏

完全自主演化的多人在线文字修仙世界（MUD）。五服微服务架构，支持 AI 驱动的 NPC 行为、动态世界事件、天道因果系统。

## 技术栈

- **后端**: Go 微服务（Go 1.21 / 1.25），gRPC 服务间通信，Gin HTTP 框架
- **数据库**: PostgreSQL（pgxpool）+ Redis（go-redis）
- **客户端**: 终端 CLI（主）+ Bubble Tea TUI（新）+ Gio 桌面 GUI（保留）
- **AI**: DeepSeek LLM 驱动 NPC 决策 + 行为树 + 记忆系统
- **协议**: WebSocket JSON + Protocol Buffers（protoc v4.25.3）
- **基础设施**: Docker Compose，GitHub Actions CI

## 服务架构

```
┌─────────────────┐    HTTP/WS     ┌────────────┐
│  CLI 客户端      │ ◄────────────►│  Gateway   │ :8080 / :8081
│  Gio 桌面客户端  │               └─────┬──────┘
└─────────────────┘                      │ gRPC
                  ┌──────────────────────┼──────────────────────┐
                  ▼                      ▼                      ▼
            ┌──────────┐          ┌──────────┐          ┌──────────────┐
            │  Game    │          │ Heavenly │          │     AI       │
            │  Server  │◄────────►│   Dao    │          │  Scheduler   │
            │  :50051  │  gRPC    │  :50053  │          │   :50052     │
            └────┬─────┘          └──────────┘          └──────────────┘
                 │                                         │ LLM
                 ▼                                         ▼
            ┌──────────┐                            ┌────────────┐
            │  World   │                            │  DeepSeek  │
            │  Engine  │ :50054                     │    API     │
            └──────────┘                            └────────────┘
```

### 服务职责

| 服务 | 端口 | 职责 |
|------|------|------|
| **Gateway** | 8080 (HTTP/WS), 8081 | JWT 认证、WebSocket 连接管理、消息路由转发 |
| **Game Server** | 50051 | 实体管理、31+ 种操作调度、状态同步、装备/物品/功法/好友/邮件/商店/排行榜 |
| **Heavenly Dao** | 50053 | 天道引擎：修炼效率公式、突破概率、天劫判定、因果业力 |
| **AI Scheduler** | 50052 | NPC 决策：行为树 + LLM 双循环、记忆系统、NPC 个性、自主行为 |
| **World Engine** | 50054 | 区域管理、资源刷新、世界事件、世界状态持久化 |

## 快速开始

### 环境要求

- Go 1.21+（部分模块需要 1.25）
- Docker & Docker Compose
- PostgreSQL 15 / Redis 7（通过 Docker 启动）

### 一键启动（推荐）

```bash
cd server
docker-compose up -d
```

这会启动 PostgreSQL、Redis 和全部 5 个微服务。

### 本地开发模式

```bash
# 1. 启动数据库
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=cultivation postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 2. 初始化数据库（按顺序执行）
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/01_init.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/02_game_operations.sql
# ... 依次执行 03-99 迁移文件（含邮件、商店、NPC、世界状态等）

# 3. 分别启动各个服务
cd server/game-server && go run ./cmd
cd server/heavenly-dao && go run ./cmd
cd server/ai-scheduler && go run ./cmd
cd server/world-engine && go run ./cmd
cd server/gateway && go run ./cmd
```

### Windows 快速启动

```powershell
# 编译并启动全部服务
.\start-all.ps1 -rebuild

# 停止全部服务
.\stop-all.ps1
```

### 启动客户端

```bash
# CLI 客户端（主）
cd cultivation-client-cli && go run ./cmd

# Bubble Tea TUI 客户端（新）
cd cultivation-bubbletea && go run ./cmd

# Gio 桌面客户端（保留）
cd cultivation-client-go && go run ./cmd
```

CLI 客户端支持 21+ 个命令：
- `cult` / `cultivate` — 修炼
- `bt` / `breakthrough` — 突破境界
- `exp` / `explore` — 探索区域
- `combat` / `skill` / `flee` — 战斗系统
- `spell` / `cast_spell` — 法术系统
- `move` / `gather` / `craft` — 生存
- `chat` / `msg` — 社交
- `friend` / `unfriend` / `accept` — 好友系统
- `status` / `st` — 角色状态（灵根/装备/功法）
- `create_sect` / `join_sect` / `leave_sect` — 宗门

## 项目布局

```
diershij/
├── server/                        # 服务端
│   ├── gateway/                   # API 网关
│   ├── game-server/               # 游戏核心
│   │   └── internal/
│   │       ├── service/           # 操作调度（operation.go）
│   │       └── repository/        # 持久化层
│   ├── heavenly-dao/              # 天道引擎
│   ├── ai-scheduler/              # AI NPC 调度
│   ├── world-engine/              # 世界引擎
│   ├── shared/                    # 共享库
│   │   ├── types/                 # Go 类型定义
│   │   ├── proto/                 # Protobuf 定义 + 生成代码
│   │   ├── config/                # 配置管理
│   │   └── errors/                # 错误定义
│   ├── init-db/                   # SQL 迁移脚本（01-99）
│   ├── test_workflows.py          # 集成测试脚本
│   ├── docker-compose.yml
│   └── config.json                # 默认配置
├── cultivation-client-cli/        # 终端 CLI 客户端
├── cultivation-bubbletea/         # Bubble Tea TUI 客户端
├── cultivation-client-go/         # Gio 桌面 GUI 客户端
└── .github/workflows/             # CI 配置
```

## WebSocket 协议

消息采用 JSON 信封格式：

```json
{"type": "operation", "payload": {"action_type": "cultivate", "params": {}}, "timestamp": 123}
```

**客户端 → 服务端**: `operation`（action_type + params）
**服务端 → 客户端**: `op_result` / `state_sync` / `entity_update` / `world_event` / `chat` / `error`

## 游戏系统

### 境界体系

凡人 → 练气 → 筑基 → 金丹 → 元婴 → 化神 → 炼虚 → 合体 → 大乘 → 渡劫

### 灵根系统

注册时随机生成 1-3 条灵根，5% 概率变异：
- 基础元素: 金、木、水、火、土
- 变异灵根: 风、雷、冰
- 主灵根纯度 60-90，副灵根 20-50

### 修炼公式

修炼效率 = 基础速率 × 灵根加成 × 功法匹配 × 境界衰减 × 心境系数 × (1 - 衰老惩罚)

### 突破公式

突破概率 = 基准成功率 × 积累度 × 功法品质 × 资源加成 × 心境系数 × 运气，取值 [5%, 80%]

### 邮件系统

玩家间离线消息、系统通知。支持发送、收取、邮件列表，自动清理过期邮件。

### 商店系统

NPC 商店交易，支持物品买卖、货币结算。通过 `trade` 命令交互。

### 排行榜

全服玩家境界、战力等排行，定期更新。

### NPC 系统

AI Scheduler 驱动的 NPC 行为引擎：
- 行为树 + LLM 双循环决策
- NPC 个性系统（personality traits）
- 记忆系统（长期/短期记忆）
- 自主行为（探索、修炼、社交、战斗）

## 配置

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| DB_HOST | localhost | PostgreSQL 地址 |
| DB_PORT | 5432 | PostgreSQL 端口 |
| DB_PASSWORD | postgres | 数据库密码 |
| REDIS_HOST | localhost | Redis 地址 |
| REDIS_PORT | 6379 | Redis 端口 |
| JWT_SECRET | cultivation-jwt-secret-key-2024 | JWT 签名密钥 |
| LLM_API_KEY | - | DeepSeek API 密钥（AI NPC） |

服务端运行配置见 `server/config.json`，也支持通过 `config.LoadConfigFromEnv()` 从环境变量加载。

## 开发状态

- [x] 五服微服务架构搭建
- [x] Gateway — JWT 认证 + WebSocket 路由
- [x] Game Server — 31 种操作调度 + 实体管理
- [x] Heavenly Dao — 修炼/突破/天劫/业力规则引擎
- [x] AI Scheduler — 行为树 + LLM 决策流水线
- [x] World Engine — 区域/资源/事件管理
- [x] CLI 客户端 — 21+ 命令，完整交互
- [x] 灵根系统 — 随机生成 + 纯度 + 变异
- [x] 功法系统 — 学习/主修/品质影响突破
- [x] 装备系统 — 13 项属性加成 + 耐久度
- [x] 战斗系统 — NPC 掉落 + 法术 + 逃跑
- [x] 商店/交易系统 — NPC 交易、货币结算
- [x] 邮件系统 — 离线消息、系统通知
- [x] 排行榜 — 境界/战力全服排行
- [ ] 世界事件系统 — 天材地宝出世、妖兽潮
- [ ] 跨服组队副本
- [ ] 宗门战/领地争夺
