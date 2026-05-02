# 修仙世界 - 微服务架构

完全自主演化的多人在线文字MUD修仙世界。

## 项目结构

```
.
├── server/                    # 服务端
│   ├── shared/               # 共享库
│   │   ├── types/           # 类型定义
│   │   ├── config/          # 配置管理
│   │   ├── errors/          # 错误处理
│   │   └── proto/           # gRPC协议定义
│   ├── gateway/              # API网关服务 (端口8080)
│   ├── game-server/          # 游戏核心服务 (端口50051)
│   ├── heavenly-dao/         # 天道引擎服务 (端口50052)
│   ├── ai-scheduler/         # AI调度服务 (端口50053)
│   ├── world-engine/         # 世界引擎服务 (端口50054)
│   ├── init-db/              # 数据库初始化脚本
│   └── docker-compose.yml    # Docker编排配置
│
└── client/                    # Tauri桌面客户端
    ├── src/                   # React前端源码
    │   ├── pages/            # 页面组件
    │   ├── components/       # 通用组件
    │   ├── stores/           # 状态管理
    │   ├── services/         # API/WebSocket服务
    │   └── types/            # TypeScript类型
    └── src-tauri/            # Tauri后端
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- Rust (用于Tauri)
- Docker & Docker Compose

### 启动服务端

```bash
cd server

# 启动数据库和服务
docker-compose up -d

# 或本地开发模式
# 1. 启动PostgreSQL和Redis
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cultivation postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 2. 初始化数据库
psql -h localhost -U postgres -d cultivation -f init-db/01_init.sql

# 3. 启动各服务
cd game-server && go run ./cmd &
cd heavenly-dao && go run ./cmd &
cd ai-scheduler && go run ./cmd &
cd world-engine && go run ./cmd &
cd gateway && go run ./cmd &
```

### 启动客户端

```bash
cd client

# 安装依赖
npm install

# 开发模式
npm run dev

# 或启动Tauri桌面应用
npm run tauri dev
```

## 服务架构

```
┌─────────────┐    WebSocket    ┌──────────┐
│   Client    │ ◄──────────────►│ Gateway  │ :8080
└─────────────┘                 └────┬─────┘
                                     │ gRPC
                    ┌────────────────┼────────────────┐
                    ▼                ▼                ▼
              ┌──────────┐    ┌──────────┐    ┌──────────┐
              │  Game    │    │ Heavenly │    │    AI    │
              │  Server  │    │   Dao    │    │Scheduler │
              │  :50051  │    │  :50052  │    │  :50053  │
              └────┬─────┘    └──────────┘    └──────────┘
                   │
                   ▼
              ┌──────────┐
              │  World   │
              │  Engine  │ :50054
              └──────────┘
```

## 核心功能

### 服务端
- **Gateway**: WebSocket接入、JWT认证、消息路由
- **Game Server**: 实体管理、操作验证、状态同步
- **Heavenly Dao**: 因果业力、天劫判定、世界平衡
- **AI Scheduler**: NPC决策、行为树、LLM调用
- **World Engine**: 区域管理、资源刷新、世界事件

### 客户端
- 现代单页应用风格UI
- 多标签页切换（角色、世界、社交、战斗、设置）
- 实时WebSocket通信
- Zustand状态管理

## 配置

环境变量配置：

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cultivation

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key

# LLM (可选)
LLM_PROVIDER=deepseek
LLM_API_KEY=your-api-key
```

## 开发状态

- [x] 项目结构搭建
- [x] 共享库基础类型
- [x] Gateway服务框架
- [x] Game Server服务框架
- [x] Heavenly Dao服务框架
- [x] AI Scheduler服务框架
- [x] World Engine服务框架
- [x] Docker Compose配置
- [x] Tauri客户端框架
- [x] 客户端核心页面
- [ ] 完整业务逻辑实现
- [ ] 单元测试
- [ ] 集成测试
