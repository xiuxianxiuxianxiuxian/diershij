# 修仙世界项目文档索引

> Cultivation World - 基于 AI 驱动的修仙世界模拟器

## 项目概述

修仙世界是一个多微服务架构的在线修仙世界模拟器，采用 gRPC + WebSocket 技术栈。玩家和 NPC 可在修仙世界中进行修炼、突破、探索、社交等活动，NPC 由 AI 驱动自主决策。

## 文档目录

### 核心文档

| 文档 | 说明 |
|------|------|
| [系统架构](../ARCHITECTURE.md) | 整体架构设计、微服务说明、数据流 |
| [接口文档](../INTERFACES.md) | gRPC 接口定义、Proto 消息类型、数据库 Schema |
| [开发者指南](../DEVELOPER_GUIDE.md) | 环境搭建、构建、运行、测试 |

### 模块文档

| 模块 | 说明 |
|------|------|
| [Gateway](模块/gateway.md) | 网关服务 - WebSocket 连接、HTTP API、认证 |
| [Game Server](模块/game-server.md) | 游戏服务器 - 实体管理、操作执行、状态同步 |
| [World Engine](模块/world-engine.md) | 世界引擎 - 区域管理、资源生成、世界事件 |
| [Heavenly Dao](模块/heavenly-dao.md) | 天道系统 - 因果评判、天劫管理、世界平衡 |
| [AI Scheduler](模块/ai-scheduler.md) | AI 调度器 - NPC 行为决策、行为树、LLM 集成 |

### 专有概念文档

| 概念 | 说明 |
|------|------|
| [境界系统](专有概念/realm-system.md) | 10 个修仙境界及突破机制 |
| [属性系统](专有概念/attributes-system.md) | 83+ 核心属性分类详解 |
| [功法系统](专有概念/method-system.md) | 功法体系、技能系统、功法传承 |
| [物品系统](专有概念/item-system.md) | 物品模板、背包、装备、丹药、法宝、符箓 |
| [因果系统](专有概念/karma-system.md) | 业力、功德、天道印记机制 |
| [社交系统](专有概念/social-system.md) | 宗门、人际关系、NPC 性格 |
| [世界系统](专有概念/world-system.md) | 区域、资源、事件、世界平衡 |

## 项目结构

```
server/
├── gateway/                  # 网关服务 (HTTP + WebSocket, 端口 8080)
│   ├── cmd/main.go
│   └── internal/
│       ├── handler/          # WebSocket Hub, HTTP 路由
│       └── service/          # 认证、gRPC 客户端
├── game-server/              # 游戏服务器 (gRPC, 端口 50051)
│   ├── cmd/main.go
│   └── internal/
│       ├── service/          # 游戏逻辑、操作执行
│       └── repository/       # PostgreSQL + Redis 数据层
├── world-engine/             # 世界引擎 (gRPC, 端口 50054)
│   ├── cmd/main.go
│   └── internal/service/     # 区域、资源、事件管理
├── heavenly-dao/             # 天道服务 (gRPC, 端口 50053)
│   ├── cmd/main.go
│   └── internal/service/     # 因果、天劫、平衡检查
├── ai-scheduler/             # AI 调度器 (gRPC, 端口 50052)
│   ├── cmd/main.go
│   └── internal/service/     # NPC 决策、行为树、LLM
└── shared/                   # 共享代码
    ├── config/               # 配置管理
    ├── errors/               # 错误定义
    ├── proto/                # Proto 定义
    │   ├── game.proto
    │   ├── entity.proto
    │   ├── world.proto
    │   ├── heavenly_dao.proto
    │   └── ai_scheduler.proto
    └── types/                # 核心类型定义
        ├── entity.go         # 实体、属性、境界
        ├── method.go         # 功法、技能
        ├── item.go           # 物品、装备、丹药、法宝
        ├── social.go         # 宗门、关系、NPC 性格
        ├── world.go          # 区域、资源、事件
        ├── operation.go      # 操作类型
        ├── message.go        # WebSocket 消息
        └── id.go             # ID 生成
```

## 技术栈

- **语言**: Go
- **RPC**: gRPC + Protocol Buffers
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **Web 框架**: Gin (Gateway)
- **WebSocket**: Gorilla WebSocket
- **认证**: JWT (HS256)
- **AI**: DeepSeek API (deepseek-chat / deepseek-reasoner)
- **部署**: Docker Compose

## 快速链接

- [游戏服务器类型定义 - entity.go](../../server/shared/types/entity.go)
- [游戏服务器类型定义 - method.go](../../server/shared/types/method.go)
- [游戏服务器类型定义 - item.go](../../server/shared/types/item.go)
- [游戏服务器类型定义 - social.go](../../server/shared/types/social.go)
- [Docker Compose 配置](../../server/docker-compose.yml)
