# 开发者指南

## 环境要求

- **Go**: 1.21+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Protocol Buffers**: protoc 编译器 (用于生成 proto 代码)

## 项目结构

```
server/
├── gateway/              # 网关服务 (HTTP + WebSocket)
├── game-server/          # 游戏服务器 (gRPC)
├── world-engine/         # 世界引擎 (gRPC)
├── heavenly-dao/         # 天道服务 (gRPC)
├── ai-scheduler/         # AI 调度器 (gRPC)
└── shared/               # 共享代码
    ├── config/           # 配置
    ├── errors/           # 错误定义
    ├── proto/            # Proto 文件
    │   └── pb/           # 生成的 Go 代码
    └── types/            # 类型定义
```

## 构建与运行

### 方式一：Docker Compose (推荐)

```bash
# 启动所有服务 (PostgreSQL, Redis, Gateway, Game Server, Heavenly Dao, AI Scheduler, World Engine)
cd server
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f gateway
docker compose logs -f game-server

# 停止所有服务
docker compose down

# 停止并清除数据
docker compose down -v
```

### 方式二：本地开发

#### 1. 启动基础设施

```bash
# 仅启动 PostgreSQL 和 Redis
cd server
docker compose up -d postgres redis
```

#### 2. 启动各服务

每个服务可独立启动，通过环境变量配置：

```bash
# Gateway
cd server/gateway
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=cultivation \
REDIS_HOST=localhost REDIS_PORT=6379 \
GAME_SERVER_HOST=localhost GAME_SERVER_PORT=50051 \
JWT_SECRET=cultivation-secret-key \
go run cmd/main.go

# Game Server
cd server/game-server
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=cultivation \
REDIS_HOST=localhost REDIS_PORT=6379 \
GRPC_PORT=50051 \
go run cmd/main.go

# World Engine
cd server/world-engine
GRPC_PORT=50054 \
go run cmd/main.go

# Heavenly Dao
cd server/heavenly-dao
GRPC_PORT=50053 \
go run cmd/main.go

# AI Scheduler
cd server/ai-scheduler
GRPC_PORT=50052 \
LLM_PROVIDER=deepseek \
LLM_API_KEY=your-api-key \
LLM_RATE_LIMIT=600 \
go run cmd/main.go
```

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| Gateway | 8080 | HTTP/HTTPS + WebSocket |
| Game Server | 50051 | gRPC |
| AI Scheduler | 50052 | gRPC |
| Heavenly Dao | 50053 | gRPC |
| World Engine | 50054 | gRPC |
| PostgreSQL | 5432 | TCP |
| Redis | 6379 | TCP |

## Proto 代码生成

```bash
cd server/shared/proto

# 生成所有 proto 文件的 Go 代码
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    *.proto
```

需要安装的 protoc 插件:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 测试

### 运行所有测试

```bash
cd server
go test ./... -v
```

### 运行特定包的测试

```bash
# 类型定义测试
go test ./shared/types/... -v

# 游戏服务器测试
go test ./game-server/... -v
```

### 当前测试覆盖

项目中已包含以下测试文件：

| 测试文件 | 测试内容 |
|----------|----------|
| `shared/types/entity_test.go` | Entity ID 生成、实体初始化、Attributes (83+ 字段)、Karma、SpiritualRoot、Injury、Buff、Debuff、SpiritStones、完整实体 |
| `shared/types/method_test.go` | CultivationMethod 初始化、Skill、SkillEffect、EntityMethod、默认值、AttackBonuses、多技能 |
| `shared/types/item_test.go` | 待确认 (文件存在) |
| `shared/types/social_test.go` | 待确认 (文件存在) |
| `shared/types/world_test.go` | 待确认 (文件存在) |

## 配置

### 环境变量

所有服务支持通过环境变量配置，详见 `shared/config/config.go` 中的 `LoadConfigFromEnv()` 函数。

核心环境变量:

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DB_HOST` | localhost | PostgreSQL 主机 |
| `DB_PORT` | 5432 | PostgreSQL 端口 |
| `DB_USER` | postgres | 数据库用户 |
| `DB_PASSWORD` | postgres | 数据库密码 |
| `DB_NAME` | cultivation | 数据库名 |
| `REDIS_HOST` | localhost | Redis 主机 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | (空) | Redis 密码 |
| `REDIS_DB` | 0 | Redis 数据库编号 |
| `SERVER_HOST` | 0.0.0.0 | HTTP 服务器主机 |
| `SERVER_PORT` | 8080 | HTTP 服务器端口 |
| `GRPC_HOST` | 0.0.0.0 | gRPC 服务器主机 |
| `GRPC_PORT` | 50051 | gRPC 服务器端口 |
| `JWT_SECRET` | cultivation-secret-key | JWT 签名密钥 |
| `LLM_PROVIDER` | deepseek | LLM 提供商 |
| `LLM_API_KEY` | (空) | LLM API 密钥 |
| `LLM_DAILY_MODEL` | deepseek-chat | 日常模型 |
| `LLM_REASON_MODEL` | deepseek-reasoner | 推理模型 |
| `LLM_RATE_LIMIT` | 600 | LLM 限流 (tokens) |
| `LLM_TIMEOUT` | 10 | LLM 超时 (秒) |

### JSON 配置文件

也可通过 JSON 文件配置：

```bash
go run cmd/main.go --config config.json
```

配置结构参考 `shared/config/config.go` 中的 `Config` 类型。

## 开发规范

### 添加新操作类型

1. 在 `shared/types/operation.go` 中添加新的 `ActionType` 常量
2. 在 `game-server/internal/service/operation.go` 的 `Execute()` 方法中添加 case 分支
3. 实现对应的 `executeXxx()` 方法

### 添加新属性

1. 在 `shared/types/entity.go` 的 `Attributes` 结构体中添加字段
2. 在 Proto 文件中添加对应字段（如果需要跨服务传输）
3. 在数据库 repository 中添加对应的查询/更新逻辑
4. 在 `entityToProto()` 和 `protoToEntity()` 中添加字段映射

### 添加新 gRPC 服务

1. 在 `shared/proto/` 下创建新的 `.proto` 文件
2. 运行 `protoc` 生成代码
3. 在 `server/` 下创建服务目录
4. 实现 gRPC 接口
5. 在 `docker-compose.yml` 中添加服务配置

## 常见问题

### 数据库连接失败

确保 PostgreSQL 已启动且配置正确：
```bash
docker compose ps postgres
docker compose logs postgres
```

### Proto 代码未更新

修改 `.proto` 文件后需要重新生成：
```bash
cd server/shared/proto
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    *.proto
```

### LLM API 调用失败

检查 `LLM_API_KEY` 环境变量是否正确设置，以及 `LLM_RATE_LIMIT` 是否合理。
