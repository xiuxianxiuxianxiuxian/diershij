# AI Scheduler 模块

## 概述

AI Scheduler (AI 调度器) 负责 NPC 的自主行为决策，结合行为树和 LLM 实现智能 NPC 系统。

- **端口**: 50052 (gRPC)
- **框架**: gRPC
- **AI 提供商**: DeepSeek (deepseek-chat / deepseek-reasoner)
- **入口**: `server/ai-scheduler/cmd/main.go`

> 注：项目中存在两个 AI 调度器实现。`ai_scheduler.go` 包含完整的行为树和 LLM 逻辑，`ai_scheduler_service.go` 提供了 gRPC 服务接口和简化的调度器。

## 核心组件

### AISchedulerService (`ai_scheduler.go`)

完整的 AI 调度器实现，包含 NPC 注册、行为模板、限流和 LLM 集成。

**NPCProfile (NPC 档案)**:
- NPCID, PersonalityType, MoralAlignment
- AmbitionLevel (1-100), RiskTolerance (0-1)
- BackgroundStory, CurrentGoal

**BehaviorTemplateLibrary (行为模板库)**:

| 行为类型 | 模板 | 条件 | 动作 | 权重 |
|----------|------|------|------|------|
| cultivate | daily_cultivation | qi<50% | cultivate | 0.8 |
| cultivate | breakthrough_attempt | progress>=100% | breakthrough | 0.6 |
| explore | resource_gathering | low_resources | gather | 0.7 |
| explore | region_exploration | curious | explore | 0.5 |
| social | seek_alliance | weak | form_alliance | 0.4 |
| social | trade | surplus_resources | trade | 0.6 |

**RateLimiter (限流器)**:
- 令牌桶算法
- `maxTokens`: 由配置 `LLM.RateLimit` 决定 (默认 600)
- `refillRate`: 10 tokens/秒

**决策流程** (`ScheduleDecision`):
```
1. 获取 NPC 档案 (不存在则使用默认配置)
2. 尝试行为模板匹配
3. 如果模板匹配且随机数 < 0.7: 返回模板决策 (source=behavior_tree)
4. 否则检查 LLM 限流
5. 如果限流通过: 调用 LLM (source=llm, token_cost=100)
6. 如果限流未通过: 返回默认决策 (source=fallback)
```

**行为树逻辑** (`executeBehaviorTreeLogic`):

| 行为树名称 | 输出动作 | 参数 |
|------------|----------|------|
| daily_routine | cultivate | duration=1h |
| combat | attack | target=enemy_id |
| exploration | explore | direction=random |
| 默认 | meditate | (空) |

### AISchedulerService (`ai_scheduler_service.go`)

gRPC 服务实现：

| RPC 方法 | 说明 |
|----------|------|
| `GetAIAction` | 获取 AI 生成的动作 |
| `SetEntityBehavior` | 设置实体行为模式 |
| `StreamAIActions` | 每 30 秒推送 AI 动作 |

**行为模式**:
- `idle`: 空闲
- `cultivate`: 修炼
- `explore`: 探索

**Scheduler**:
- `GenerateAction()`: 随机生成 cultivate/move/explore 三种动作之一
- `ScheduleDecision()`: 基于世界状态生成动作列表 (当前返回空列表)

## gRPC 接口 (AISchedulerService - ai_scheduler.proto)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `ScheduleDecision` | `DecisionRequest` | `DecisionResponse` | NPC 行为决策 |
| `ExecuteBehaviorTree` | `BehaviorTreeRequest` | `BehaviorTreeResponse` | 执行行为树 |
| `RegisterNPC` | `NPCRegisterRequest` | `NPCRegisterResponse` | 注册 NPC |
| `UnregisterNPC` | `NPCUnregisterRequest` | `NPCUnregisterResponse` | 注销 NPC |

## NPC 注册

注册 NPC 时需提供：
- `npc_id`: NPC 唯一标识
- `personality_type`: 性格类型
- `moral_alignment`: 道德倾向 (lawful_good, neutral, chaotic_evil 等)
- `ambition_level`: 野心程度 (1-100)
- `risk_tolerance`: 风险承受度 (0-1)
- `background_story`: 背景故事
- `current_goal`: 当前目标

## LLM 集成

**配置** (`shared/config/config.go`):
```go
type LLMConfig struct {
    Provider    string  // "deepseek"
    APIKey      string  // API 密钥
    DailyModel  string  // "deepseek-chat"
    ReasonModel string  // "deepseek-reasoner"
    RateLimit   int     // 600 (tokens/秒)
    Timeout     int     // 10 (秒)
}
```

> 注：当前代码中 `callLLM()` 方法为占位实现，随机选择可用动作返回。完整的 LLM 调用逻辑需要后续实现。

## 依赖

- `google.golang.org/grpc` - gRPC 框架
- `github.com/cultivation-world/shared/config` - 共享配置
- `github.com/cultivation-world/shared/types` - 共享类型
- `github.com/cultivation-world/shared/proto/pb` - Proto 生成代码
