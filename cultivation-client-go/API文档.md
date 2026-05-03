# 修仙世界客户端 API 对接文档

## 概述

本文档描述了 Go 客户端与修仙世界服务端 API 的对接方式。客户端通过 HTTP REST API 进行认证，通过 WebSocket 进行实时通信。

## 基础配置

### 服务端地址

```go
// 默认配置
baseURL: "http://localhost:8080"
wsURL:   "ws://localhost:8080/ws"
```

### 客户端初始化

```go
import (
    "cultivation-client/internal/network"
    "cultivation-client/internal/store"
)

// 获取 API 客户端
apiClient := network.GetAPIClient()

// 获取 WebSocket 客户端
wsClient := network.GetWebSocketClient()

// 注册默认消息处理器
network.RegisterDefaultHandlers()
```

---

## HTTP API

### 1. 认证接口

#### 1.1 用户注册

**请求**
```http
POST /auth/register
Content-Type: application/json

{
    "username": "string",  // 用户名，必填
    "password": "string"   // 密码，必填
}
```

**响应**
```json
{
    "success": true,
    "token": "jwt_token_string",
    "entity": {
        "id": "entity_uuid",
        "name": "username",
        "realm": "mortal",
        "position": {
            "region_id": "qingyun_town",
            "x": 0,
            "y": 0
        }
    }
}
```

**Go 代码示例**
```go
resp, err := apiClient.Register("username", "password")
if err != nil {
    // 处理错误
    return
}
if resp.Success {
    // 注册成功，token 已自动保存
    fmt.Println("注册成功，实体ID:", resp.Entity.ID)
}
```

#### 1.2 用户登录

**请求**
```http
POST /auth/login
Content-Type: application/json

{
    "username": "string",
    "password": "string"
}
```

**响应**
```json
{
    "success": true,
    "token": "jwt_token_string",
    "entity": {
        "id": "entity_uuid",
        "name": "username",
        "realm": "mortal"
    }
}
```

**Go 代码示例**
```go
resp, err := apiClient.Login("username", "password")
if err != nil {
    // 处理错误
    return
}
if resp.Success {
    // 登录成功，token 已自动保存到 store
    fmt.Println("登录成功")
}
```

#### 1.3 用户登出

**Go 代码示例**
```go
err := apiClient.Logout()
// 清除本地存储的 token 和实体信息
```

---

## WebSocket API

### 连接方式

```go
// 必须先登录获取 token
wsClient := network.GetWebSocketClient()

// 设置自定义地址（可选）
wsClient.url = "ws://localhost:8080/ws"

// 连接 WebSocket
err := wsClient.Connect()
if err != nil {
    fmt.Printf("连接失败: %v\n", err)
    return
}

// 断开连接
wsClient.Disconnect()
```

### 消息格式

所有 WebSocket 消息使用 JSON 格式：

```json
{
    "type": "message_type",
    "payload": {
        // 消息内容
    }
}
```

### 发送消息

```go
// 发送操作消息
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "cultivate",
    "params": {},
})

// 发送聊天消息
err := wsClient.Send("chat", map[string]interface{}{
    "content": "大家好！",
    "channel": "world",
})
```

### 注册消息处理器

```go
wsClient.RegisterHandler("message_type", func(payload []byte) {
    // 处理消息
    var data YourType
    json.Unmarshal(payload, &data)
    // ...
})
```

---

## 游戏操作

### 操作类型列表

| 操作类型 | 说明 | 必需参数 |
|---------|------|---------|
| cultivate | 修炼 | 无 |
| move | 移动 | region_id, x, y |
| meditate | 打坐 | 无 |
| sleep | 休息 | 无 |
| breakthrough | 突破 | 无 |
| combat | 战斗 | target_id |
| explore | 探索 | 无 |
| gather | 采集 | resource_type, quantity |
| craft | 炼器/炼丹 | recipe_id |
| create_method | 自创功法 | method_name, method_type |
| trade | 交易 | target_id, item_id, price |
| form_sect | 创建宗门 | sect_name |
| join_sect | 加入宗门 | sect_id |
| send_message | 发送消息 | content, message_type |
| cast_spell | 施法 | spell_id |

### 操作示例

#### 修炼

```go
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "cultivate",
    "params": {},
})
```

#### 移动

```go
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "move",
    "params": map[string]interface{}{
        "region_id": "east_wilderness",
        "x": 10.5,
        "y": 20.3,
    },
})
```

#### 战斗

```go
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "combat",
    "params": map[string]interface{}{
        "target_id": "target_entity_uuid",
    },
})
```

#### 采集

```go
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "gather",
    "params": map[string]interface{}{
        "resource_type": "herb",
        "quantity": 5,
    },
})
```

#### 发送消息

```go
// 私聊
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "send_message",
    "params": map[string]interface{}{
        "content": "你好！",
        "message_type": "private",
        "receiver_id": "target_entity_uuid",
    },
})

// 世界频道
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "send_message",
    "params": map[string]interface{}{
        "content": "大家好！",
        "message_type": "world",
    },
})
```

#### 施法

```go
err := wsClient.Send("operation", map[string]interface{}{
    "action_type": "cast_spell",
    "params": map[string]interface{}{
        "spell_id": "spell_uuid",
        "target_id": "target_entity_uuid", // 可选
    },
})
```

---

## 消息推送

### 默认处理器

客户端已注册以下默认消息处理器：

```go
// 战斗更新
ws.RegisterHandler("combat_update", func(payload []byte) {
    var update types.CombatState
    json.Unmarshal(payload, &update)
    store.GetGameStore().SetCombat(&update)
})

// 世界更新
ws.RegisterHandler("world_update", func(payload []byte) {
    var update types.WorldState
    json.Unmarshal(payload, &update)
    store.GetGameStore().SetWorld(&update)
})

// 社交更新
ws.RegisterHandler("social_update", func(payload []byte) {
    var update types.SocialInfo
    json.Unmarshal(payload, &update)
    store.GetGameStore().SetSocial(&update)
})

// 新消息
ws.RegisterHandler("new_message", func(payload []byte) {
    var update types.Message
    json.Unmarshal(payload, &update)
    store.GetGameStore().AddMessage(update)
})
```

### 自定义处理器

```go
// 注册自定义处理器
wsClient.RegisterHandler("entity_update", func(payload []byte) {
    var entity types.Entity
    if err := json.Unmarshal(payload, &entity); err != nil {
        return
    }
    // 更新 UI 或处理数据
    fmt.Printf("实体更新: %s\n", entity.Name)
})
```

---

## 数据存储

### AuthStore - 认证信息

```go
authStore := store.GetAuthStore()

// 获取 token
token := authStore.GetToken()

// 获取实体信息
entity := authStore.GetEntity()
entityID := authStore.GetEntityID()
entityName := authStore.GetEntityName()

// 设置 token
authStore.SetToken("jwt_token")

// 登出
authStore.Logout()
```

### GameStore - 游戏数据

```go
gameStore := store.GetGameStore()

// 角色信息
character := gameStore.GetCharacter()
gameStore.SetCharacter(character)

// 世界状态
world := gameStore.GetWorld()
gameStore.SetWorld(world)

// 战斗状态
combat := gameStore.GetCombat()
gameStore.SetCombat(combat)

// 社交信息
social := gameStore.GetSocial()
gameStore.SetSocial(social)

// 设置
settings := gameStore.GetSettings()
gameStore.SetSettings(settings)

// 添加消息
gameStore.AddMessage(message)

// 清除数据
gameStore.Clear()
```

---

## 类型定义

### 核心类型

```go
// 实体
type Entity struct {
    ID       string
    Name     string
    Realm    string
    Position WorldPosition
}

// 角色
type Character struct {
    ID               string
    Name             string
    Level            int
    Health           int
    MaxHealth        int
    Energy           int
    MaxEnergy        int
    Attack           int
    Defense          int
    Speed            int
    CultivationRealm string
}

// 世界状态
type WorldState struct {
    CurrentMap    string
    PlayersOnline int
    Events        []WorldEvent
    Announcements []Announcement
}

// 战斗状态
type CombatState struct {
    InCombat   bool
    BattleLog  []CombatLog
    TurnNumber int
}

// 社交信息
type SocialInfo struct {
    Friends  []Friend
    Messages []Message
    Requests []FriendRequest
}

// WebSocket 消息
type WSMessage struct {
    Type    WSMessageType
    Payload map[string]interface{}
}
```

---

## 错误处理

### HTTP 错误

```go
resp, err := apiClient.Login("user", "pass")
if err != nil {
    // 错误类型：
    // - 网络错误：connection refused, timeout
    // - API 错误：invalid credentials, entity not found
    fmt.Printf("登录失败: %v\n", err)
    return
}
```

### WebSocket 错误

```go
err := wsClient.Connect()
if err != nil {
    // 常见错误：
    // - 未认证：not authenticated
    // - 连接失败：failed to connect websocket
    fmt.Printf("连接失败: %v\n", err)
    return
}
```

### 操作结果错误

操作结果通过 WebSocket 返回：

```json
{
    "type": "op_result",
    "payload": {
        "success": false,
        "message": "not enough qi",
        "error_code": 5
    }
}
```

---

## 完整使用示例

```go
package main

import (
    "fmt"
    "cultivation-client/internal/network"
    "cultivation-client/internal/store"
)

func main() {
    // 1. 登录
    apiClient := network.GetAPIClient()
    resp, err := apiClient.Login("username", "password")
    if err != nil {
        fmt.Printf("登录失败: %v\n", err)
        return
    }
    fmt.Printf("登录成功: %s\n", resp.Entity.Name)

    // 2. 连接 WebSocket
    wsClient := network.GetWebSocketClient()
    
    // 注册自定义处理器
    wsClient.RegisterHandler("entity_update", func(payload []byte) {
        fmt.Println("收到实体更新")
    })
    
    // 连接
    err = wsClient.Connect()
    if err != nil {
        fmt.Printf("WebSocket 连接失败: %v\n", err)
        return
    }
    defer wsClient.Disconnect()

    // 3. 发送操作
    // 修炼
    err = wsClient.Send("operation", map[string]interface{}{
        "action_type": "cultivate",
        "params": {},
    })
    if err != nil {
        fmt.Printf("发送失败: %v\n", err)
    }

    // 4. 获取游戏数据
    character := store.GetGameStore().GetCharacter()
    fmt.Printf("角色: %s, 境界: %s\n", character.Name, character.CultivationRealm)

    // 5. 保持运行
    select {}
}
```

---

## 注意事项

1. **认证顺序**：必须先调用 HTTP API 登录/注册，获取 token 后才能连接 WebSocket
2. **自动重连**：WebSocket 客户端支持自动重连，断线后会每 5 秒尝试重连
3. **并发安全**：所有 store 操作都是线程安全的
4. **资源释放**：程序退出前调用 `wsClient.Disconnect()` 释放资源
5. **端口配置**：确保服务端端口（8080）与客户端配置一致

---

## 更新日志

- **2024-01-XX**: 初始版本，支持基础操作和 WebSocket 通信
