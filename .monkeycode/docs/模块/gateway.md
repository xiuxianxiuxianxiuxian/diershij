# Gateway 模块

## 概述

Gateway 是修仙世界的入口服务，负责处理客户端的 HTTP 请求和 WebSocket 连接，提供认证、消息路由和客户端管理功能。

- **端口**: 8080
- **协议**: HTTP + WebSocket
- **框架**: Gin + Gorilla WebSocket
- **入口**: `server/gateway/cmd/main.go`

## 架构

```
Client <-> Gateway (HTTP/WS) <-> Game Server (gRPC :50051)
                         <-> Heavenly Dao (gRPC :50053)
```

## 核心组件

### Server (`handler/server.go`)

HTTP 服务器，使用 Gin 框架，提供以下路由：

| 路由 | 方法 | 说明 |
|------|------|------|
| `/auth/register` | POST | 注册新用户 |
| `/auth/login` | POST | 用户登录 |
| `/ws` | GET | WebSocket 连接 (需 token 参数) |
| `/health` | GET | 健康检查 |

**启动流程**:
1. 加载配置 (环境变量)
2. 创建 Game Service gRPC 客户端
3. 创建 WebSocket Hub 并启动
4. 创建 AuthService
5. 启动 Gin HTTP 服务器

### WebSocketHub (`handler/websocket.go`)

WebSocket 连接管理中心，使用 goroutine 通道模式管理客户端生命周期：

```
register   chan -> 新客户端注册
unregister chan -> 客户端断开
broadcast  chan -> 消息广播
```

**核心功能**:
- `Run()`: 主循环，处理注册/注销/广播
- `Register()`: 注册客户端到 Hub
- `Unregister()`: 移除断开的客户端
- `BroadcastToEntity()`: 向指定实体发送消息
- `BroadcastToAll()`: 向所有客户端广播消息

### WebSocketClient (`handler/websocket.go`)

单个 WebSocket 连接封装：

**连接参数**:
- `writeWait`: 10 秒
- `pongWait`: 60 秒
- `pingPeriod`: 54 秒 (pongWait * 9/10)
- `maxMessageSize`: 512 字节

**消息处理**:
- `ReadPump()`: 从 WebSocket 读取消息，JSON 反序列化为 `types.Message`
- `WritePump()`: 向 WebSocket 写入消息，定时发送 Ping
- `handleMessage()`: 根据消息类型分发处理
  - `operation` -> `handleOperation()`: 转发到 Game Server 执行操作
  - `chat` -> `handleChat()`: 广播聊天消息
  - 其他 -> 返回错误

### AuthService (`service/auth.go`)

JWT 认证服务：

- `Register(username, password)`: 调用 Game Server 创建实体，生成 JWT token
- `Login(username, password)`: 调用 Game Server 认证，生成 JWT token
- `ValidateToken(tokenString)`: 验证 JWT token，返回 EntityID
- `generateToken(entityID, username)`: 生成 HS256 JWT，24 小时过期

**JWT Claims**:
```go
type Claims struct {
    EntityID types.EntityID `json:"entity_id"`
    Username string         `json:"username"`
    jwt.RegisteredClaims
}
```

### GameServiceClient (`service/game_client.go`)

gRPC 客户端封装，连接 Game Server：

- `ExecuteOperation()`: 执行操作
- `CreateEntity()`: 创建实体
- `AuthenticateEntity()`: 认证实体
- `GetEntity()`: 获取实体信息

## WebSocket 消息协议

客户端通过 WebSocket 发送 JSON 格式消息：

```json
{
  "type": "operation",
  "payload": {
    "action_type": "cultivate",
    "params": {}
  },
  "request_id": "optional-request-id"
}
```

支持的客户端消息类型：
- `operation`: 执行游戏操作
- `chat`: 发送聊天消息

服务端返回消息类型：
- `op_result`: 操作结果
- `entity_update`: 实体状态更新
- `world_event`: 世界事件
- `error`: 错误信息
- `chat`: 聊天消息

## 依赖

- `github.com/gin-gonic/gin` - HTTP 框架
- `github.com/gorilla/websocket` - WebSocket
- `github.com/golang-jwt/jwt/v5` - JWT 认证
- `google.golang.org/grpc` - gRPC 客户端
