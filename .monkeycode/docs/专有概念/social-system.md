# 社交系统

## 概述

社交系统定义了修仙世界中实体之间的关系网络，包括宗门体系、人际关系和 NPC 性格系统。类型定义在 `shared/types/social.go` 中。

## 宗门系统

### Sect (宗门)

| 属性 | 说明 |
|------|------|
| `ID` | 宗门唯一标识 |
| `Name` | 宗门名称 |
| `FounderID` | 创建者 ID |
| `Philosophy` | 宗门理念 |
| `EntryRequirements` | 入门条件 (map[string]any) |
| `Territory` | 势力范围 (区域 ID 列表) |
| `Rules` | 宗门规则 (map[string]any) |
| `Alignment` | 立场: 正道/魔道/中立 |
| `CreatedAt` | 创建时间 |
| `MemberCount` | 成员数量 |
| `Prestige` | 宗门声望 |
| `Wealth` | 宗门财富 (灵石) |
| `FacilityScore` | 设施评分 |
| `CultivationResources` | 宗门修炼资源列表 |

### SectMember (宗门成员)

| 属性 | 说明 |
|------|------|
| `SectID` | 所属宗门 ID |
| `EntityID` | 实体 ID |
| `Rank` | 职位 |
| `Contribution` | 贡献值 |
| `JoinedAt` | 加入时间 |
| `Privileges` | 特权列表 |

## 人际关系

### Relationship (人际关系)

| 属性 | 说明 |
|------|------|
| `ID` | 关系唯一标识 |
| `EntityAID` | 实体 A ID |
| `EntityBID` | 实体 B ID |
| `RelationshipType` | 关系类型: 师徒/仇敌/盟友/恋人/结义等 |
| `Strength` | 关系强度 (0-100) |
| `History` | 关系历史 |
| `CreatedAt` | 创建时间 |

### 实体属性中的关系字段

在 `types.Attributes` 中也包含了关系相关的快捷字段：

| 属性 | 说明 |
|------|------|
| `RelationshipCount` | 人际关系总数 |
| `MentorID` | 师尊 ID |
| `DiscipleIDs` | 弟子 ID 列表 |
| `SwornSiblings` | 结拜兄弟/姐妹列表 |
| `Enemies` | 仇敌列表 |
| `Lovers` | 道侣列表 |

## NPC 性格系统

### NPCPersonality (NPC 性格配置)

| 属性 | 说明 |
|------|------|
| `NPCID` | NPC ID |
| `PersonalityType` | 性格类型 |
| `MoralAlignment` | 道德倾向 (lawful_good/neutral/chaotic_evil 等) |
| `AmbitionLevel` | 野心程度 (1-100) |
| `RiskTolerance` | 风险承受度 (0-1) |
| `SocialPreference` | 社交偏好 (extrovert/introvert/balanced) |
| `BackgroundStory` | 背景故事 |
| `CurrentGoal` | 当前目标 |
| `HiddenSecrets` | 隐藏秘密列表 |
| `LLMSystemPrompt` | LLM 系统提示词 |
| `BehaviorTreeConfig` | 行为树配置 (map[string]any) |
| `InitialActions` | 初始行为模式列表 |

### NPCDecisionLog (NPC 决策日志)

| 属性 | 说明 |
|------|------|
| `ID` | 决策日志 ID |
| `NPCID` | NPC ID |
| `DecisionType` | 决策类型 |
| `Context` | 决策上下文 (map[string]any) |
| `ActionTaken` | 采取的行动 (map[string]any) |
| `Reasoning` | 决策推理 |
| `ModelUsed` | 使用的模型 (deepseek-chat/deepseek-reasoner) |
| `Source` | 决策来源 (behavior_tree/llm) |
| `TokenCost` | API 调用成本 |
| `Timestamp` | 时间戳 |

## 实体社交属性

实体属性中包含以下社交相关字段：

| 属性 | 说明 |
|------|------|
| `Reputation` | 声望值 |
| `SectContribution` | 宗门贡献值 |
| `FactionStandings` | 各势力立场 (map[string]int，可正可负) |

## AI 驱动的 NPC 行为

NPC 的行为决策由 AI Scheduler 服务驱动，结合以下方式：

1. **行为树模板**: 预定义的行为模式匹配
2. **LLM 决策**: 通过 DeepSeek API 进行推理决策
3. **性格影响**: 性格配置影响决策倾向
4. **目标导向**: 当前目标驱动行为选择

### NPC 注册到 AI Scheduler

注册时需要提供：
- NPC ID
- 性格类型
- 道德倾向
- 野心程度 (1-100)
- 风险承受度 (0-1)
- 背景故事
- 当前目标

## 待确认

- 宗门的具体创建流程和管理机制
- 宗门之间的互动（战争、联盟等）
- 关系强度的增减规则
- NPC 性格对行为决策的具体影响权重
- 社交互动（交易、组队、交谈等）的具体实现
