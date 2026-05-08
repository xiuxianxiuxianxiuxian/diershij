# 修仙世界 — Cultivation World 🏯

**完全自主演化的多人在线文字修仙世界（MUD）** | **A fully self-evolving multiplayer online text-based cultivation world (MUD)**

五服微服务架构，支持 AI 驱动的 NPC 行为、动态世界事件、天道因果系统。  
Five-service microservice architecture with AI-driven NPC behavior, dynamic world events, and a heavenly dao karma system.

---

## 技术栈 / Tech Stack

- **后端 Backend**: Go 微服务（Go 1.21 / 1.25），gRPC 服务间通信，Gin HTTP 框架
- **数据库 Database**: PostgreSQL（pgxpool）+ Redis（go-redis）
- **客户端 Clients**: 终端 CLI（主）+ Bubble Tea TUI（新）+ Gio 桌面 GUI（保留）
- **AI**: DeepSeek LLM 驱动 NPC 决策 + 行为树 + 记忆系统
- **协议 Protocols**: WebSocket JSON + Protocol Buffers（protoc v4.25.3）
- **基础设施 Infrastructure**: Docker Compose，GitHub Actions CI

---

## 服务架构 / Service Architecture

```
┌─────────────────┐    HTTP/WS     ┌────────────┐
│  CLI / TUI /    │ ◄────────────►│  Gateway   │ :8080 / :8081
│  Gio Clients    │               └─────┬──────┘
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

### 服务职责 / Service Responsibilities

| 服务 Service | 端口 Port | 职责 Responsibilities |
|-------------|-----------|----------------------|
| **Gateway** | 8080 (HTTP/WS), 8081 | JWT 认证、WebSocket 连接管理、消息路由转发 / JWT auth, WebSocket conn management, message routing |
| **Game Server** | 50051 | 实体管理、31+ 种操作调度、状态同步、装备/物品/功法/好友/邮件/商店/排行榜 / Entity management, 31+ operations, state sync, equipment/items/spells/friends/mail/shop/leaderboard |
| **Heavenly Dao** | 50053 | 天道引擎：修炼效率公式、突破概率、天劫判定、因果业力 / Heavenly engine: cultivation formulas, breakthrough probability, tribulation, karma |
| **AI Scheduler** | 50052 | NPC 决策：行为树 + LLM 双循环、记忆系统、NPC 个性、自主行为 / NPC decisions: behavior tree + LLM dual loop, memory, personality, autonomous behavior |
| **World Engine** | 50054 | 区域管理、资源刷新、世界事件、世界状态持久化 / Region management, resource respawn, world events, state persistence |

---

## 快速开始 / Quick Start

### 环境要求 / Prerequisites

- Go 1.21+（some modules require 1.25）
- Docker & Docker Compose
- PostgreSQL 15 / Redis 7（via Docker）

### 一键启动（推荐）/ One-Click Start (Recommended)

```bash
cd server
docker-compose up -d
```

This starts PostgreSQL, Redis, and all 5 microservices.

### 本地开发模式 / Local Development

```bash
# 1. 启动数据库 / Start databases
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=cultivation postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 2. 初始化数据库（按顺序执行）/ Initialize database (run in order)
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/01_init.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/02_game_operations.sql
# ... run 03-99 migration files in sequence (includes mail, shop, NPC, world state, etc.)

# 3. 分别启动各个服务 / Start each service individually
cd server/game-server && go run ./cmd
cd server/heavenly-dao && go run ./cmd
cd server/ai-scheduler && go run ./cmd
cd server/world-engine && go run ./cmd
cd server/gateway && go run ./cmd
```

### Windows 快速启动 / Windows Quick Start

```powershell
# 编译并启动全部服务 / Build & start all services
.\start-all.ps1 -rebuild

# 停止全部服务 / Stop all services
.\stop-all.ps1
```

### 启动客户端 / Launching Clients

```bash
# CLI 客户端（主）/ CLI client (primary)
cd cultivation-client-cli && go run ./cmd

# Bubble Tea TUI 客户端（新）/ Bubble Tea TUI client (new)
cd cultivation-bubbletea && go run ./cmd

# Gio 桌面客户端（保留）/ Gio desktop client (retained)
cd cultivation-client-go && go run ./cmd
```

CLI 客户端支持 21+ 个命令 / The CLI client supports 21+ commands:

| 命令 Command | 别名 Alias | 说明 Description |
|-------------|-----------|-----------------|
| `cultivate` | `cult` | 修炼 / Cultivate |
| `breakthrough` | `bt` | 突破境界 / Breakthrough |
| `explore` | `exp` | 探索区域 / Explore |
| `combat` | — | 战斗 / Combat |
| `use_skill` | `skill` | 使用技能 / Use skill |
| `flee` | — | 逃跑 / Flee |
| `cast_spell` | `spell` | 法术系统 / Cast spell |
| `move` | — | 移动 / Move |
| `gather` | — | 采集 / Gather |
| `craft` | — | 制作 / Craft |
| `create_method` | `cm` | 创造功法 / Create method |
| `chat` | — | 聊天 / Chat |
| `send_message` | `msg` | 发送消息 / Send message |
| `add_friend` | `friend` | 添加好友 / Add friend |
| `remove_friend` | `unfriend` | 删除好友 / Remove friend |
| `accept_friend` | `accept` | 接受好友 / Accept friend |
| `form_sect` | `create_sect` | 创建宗门 / Create sect |
| `join_sect` | — | 加入宗门 / Join sect |
| `leave_sect` | — | 离开宗门 / Leave sect |
| `trade` | — | 交易 / Trade |
| `status` | `st` | 角色状态 / Character status |
| `help` | — | 帮助 / Help |
| `clear` | `cls` | 清屏 / Clear screen |
| `exit` | — | 退出 / Exit |

---

## 项目布局 / Project Layout

```
diershij/
├── server/                        # 服务端 / Server
│   ├── gateway/                   # API 网关 / API Gateway
│   ├── game-server/               # 游戏核心 / Game Core
│   │   └── internal/
│   │       ├── service/           # 操作调度（operation.go）/ Operation dispatcher
│   │       └── repository/        # 持久化层 / Persistence layer
│   ├── heavenly-dao/              # 天道引擎 / Heavenly Dao Engine
│   ├── ai-scheduler/              # AI NPC 调度 / AI NPC Scheduler
│   ├── world-engine/              # 世界引擎 / World Engine
│   ├── shared/                    # 共享库 / Shared Library
│   │   ├── types/                 # Go 类型定义 / Type definitions
│   │   ├── proto/                 # Protobuf 定义 + 生成代码 / Proto definitions + generated code
│   │   ├── config/                # 配置管理 / Configuration
│   │   └── errors/                # 错误定义 / Error definitions
│   ├── init-db/                   # SQL 迁移脚本（01-99）/ SQL migrations
│   ├── test_workflows.py          # 集成测试脚本 / Integration test script
│   ├── docker-compose.yml
│   └── config.json                # 默认配置 / Default config
├── cultivation-client-cli/        # 终端 CLI 客户端 / Terminal CLI Client
├── cultivation-bubbletea/         # Bubble Tea TUI 客户端 / Bubble Tea TUI Client
├── cultivation-client-go/         # Gio 桌面 GUI 客户端 / Gio Desktop Client
└── .github/workflows/             # CI 配置 / CI configuration
```

---

## WebSocket 协议 / WebSocket Protocol

消息采用 JSON 信封格式 / Messages use JSON envelope format:

```json
{"type": "operation", "payload": {"action_type": "cultivate", "params": {}}, "timestamp": 123}
```

**Client → Server**: `operation`（action_type + params）
**Server → Client**: `op_result` / `state_sync` / `entity_update` / `world_event` / `chat` / `error`

---

## 游戏系统 / Game Systems

### 境界体系 / Realm System

凡人 Mortal → 练气 Qi Refining → 筑基 Foundation → 金丹 Golden Core → 元婴 Nascent Soul → 化神 Spirit Transformation → 炼虚 Void Refining → 合体 Unity → 大乘 Mahayana → 渡劫 Tribulation Transcendence

### 灵根系统 / Spiritual Root System

注册时随机生成 1-3 条灵根，5% 概率变异 / Randomly generates 1-3 spiritual roots on registration, 5% chance of mutation:
- 基础元素 Basic: 金 Metal、木 Wood、水 Water、火 Fire、土 Earth
- 变异灵根 Mutated: 风 Wind、雷 Lightning、冰 Ice
- 主灵根纯度 Primary root purity: 60-90，副灵根 Secondary roots: 20-50

### 修炼公式 / Cultivation Formula

```
Cultivation Efficiency = Base Rate × Root Bonus × Spell Match × Realm Decay × Mental State × (1 - Aging Penalty)
```

### 突破公式 / Breakthrough Formula

```
Breakthrough Probability = Base Success × Accumulation × Spell Quality × Resource Bonus × Mental State × Luck
Clamped to [5%, 80%]
```

### 邮件系统 / Mail System

玩家间离线消息、系统通知。支持发送、收取、邮件列表，自动清理过期邮件。  
Offline messaging between players and system notifications. Supports send, receive, inbox listing, auto-cleanup of expired mail.

### 商店系统 / Shop System

NPC 商店交易，支持物品买卖、货币结算。通过 `trade` 命令交互。  
NPC shop trading with item buy/sell and currency settlement. Interact via the `trade` command.

### 排行榜 / Leaderboard

全服玩家境界、战力等排行，定期更新。  
Global player rankings by realm, combat power, etc., updated periodically.

### NPC 系统 / NPC System

AI Scheduler 驱动的 NPC 行为引擎 / AI Scheduler-driven NPC behavior engine:
- 行为树 + LLM 双循环决策 / Behavior tree + LLM dual-loop decision making
- NPC 个性系统（personality traits）/ NPC personality system
- 记忆系统（长期/短期记忆）/ Memory system (long-term/short-term)
- 自主行为（探索、修炼、社交、战斗）/ Autonomous behavior (explore, cultivate, socialize, combat)

---

## 配置 / Configuration

| 环境变量 Env Var | 默认值 Default | 说明 Description |
|-----------------|---------------|-----------------|
| DB_HOST | localhost | PostgreSQL 地址 / PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL 端口 / PostgreSQL port |
| DB_PASSWORD | postgres | 数据库密码 / Database password |
| REDIS_HOST | localhost | Redis 地址 / Redis host |
| REDIS_PORT | 6379 | Redis 端口 / Redis port |
| JWT_SECRET | cultivation-jwt-secret-key-2024 | JWT 签名密钥 / JWT signing key |
| LLM_API_KEY | — | DeepSeek API 密钥（AI NPC）/ DeepSeek API key for AI NPC |

服务端运行配置见 `server/config.json`，也支持通过 `config.LoadConfigFromEnv()` 从环境变量加载。  
See `server/config.json` for server runtime config, also supports loading from environment variables via `config.LoadConfigFromEnv()`.

---

## 开发状态 / Development Status

- [x] 五服微服务架构搭建 / Five-service microservice architecture
- [x] Gateway — JWT 认证 + WebSocket 路由 / JWT auth + WebSocket routing
- [x] Game Server — 31 种操作调度 + 实体管理 / 31 operation types + entity management
- [x] Heavenly Dao — 修炼/突破/天劫/业力规则引擎 / Cultivation/breakthrough/tribulation/karma rules engine
- [x] AI Scheduler — 行为树 + LLM 决策流水线 / Behavior tree + LLM decision pipeline
- [x] World Engine — 区域/资源/事件管理 / Region/resource/event management
- [x] CLI 客户端 — 21+ 命令，完整交互 / 21+ commands, full interaction
- [x] 灵根系统 — 随机生成 + 纯度 + 变异 / Spiritual root system — random generation + purity + mutation
- [x] 功法系统 — 学习/主修/品质影响突破 / Spell system — learn/major/quality affects breakthrough
- [x] 装备系统 — 13 项属性加成 + 耐久度 / Equipment system — 13 stat bonuses + durability
- [x] 战斗系统 — NPC 掉落 + 法术 + 逃跑 / Combat system — NPC drops + spells + flee
- [x] 商店/交易系统 — NPC 交易、货币结算 / Shop/trade system — NPC trading, currency settlement
- [x] 邮件系统 — 离线消息、系统通知 / Mail system — offline messages, system notifications
- [x] 排行榜 — 境界/战力全服排行 / Leaderboard — realm/combat power rankings
- [x] 世界事件系统 — 妖兽潮、天道异象、灵潮、秘境开启、宗门战 / World events — beast tides, heavenly anomalies, spirit tides, secret realms, sect wars
- [x] Bubble Tea TUI 客户端 — 终端交互界面 / Terminal UI client
- [ ] 跨服组队副本 / Cross-server party dungeons
- [ ] 宗门战/领地争夺 / Sect wars / territory control
