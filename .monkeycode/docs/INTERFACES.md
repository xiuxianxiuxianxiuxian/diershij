# 接口文档

## gRPC 服务接口

### GameService (`game.proto`)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `CreateEntity` | `CreateEntityRequest` | `CreateEntityResponse` | 创建新实体 |
| `AuthenticateEntity` | `AuthRequest` | `AuthResponse` | 实体认证 |
| `ExecuteOperation` | `OperationRequest` | `OperationResponse` | 执行操作 |
| `GetEntity` | `EntityRequest` | `EntityResponse` | 查询实体 |
| `SyncState` | `SyncRequest` | `SyncResponse` | 状态同步 |
| `StreamEntityUpdates` | `EntityStreamRequest` | `stream EntityUpdate` | 实体更新流 |

### WorldService (`world.proto`)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `GetRegion` | `RegionRequest` | `RegionResponse` | 查询区域 |
| `SpawnResources` | `SpawnRequest` | `SpawnResponse` | 生成资源 |
| `TriggerEvent` | `EventRequest` | `EventResponse` | 触发世界事件 |
| `GetWorldState` | `WorldStateRequest` | `WorldStateResponse` | 获取世界状态 |

### HeavenlyDaoService (`heavenly_dao.proto`)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `EvaluateKarma` | `KarmaRequest` | `KarmaResponse` | 评估因果变化 |
| `CheckTribulation` | `TribulationRequest` | `TribulationResponse` | 检查天劫 |
| `BalanceCheck` | `BalanceCheckRequest` | `BalanceCheckResponse` | 世界平衡检查 |
| `ApplyKarmaDecay` | `DecayRequest` | `DecayResponse` | 应用因果衰减 |

### AISchedulerService (`ai_scheduler.proto`)

| RPC 方法 | 请求 | 响应 | 说明 |
|----------|------|------|------|
| `ScheduleDecision` | `DecisionRequest` | `DecisionResponse` | NPC 行为决策 |
| `ExecuteBehaviorTree` | `BehaviorTreeRequest` | `BehaviorTreeResponse` | 执行行为树 |
| `RegisterNPC` | `NPCRegisterRequest` | `NPCRegisterResponse` | 注册 NPC |
| `UnregisterNPC` | `NPCUnregisterRequest` | `NPCUnregisterResponse` | 注销 NPC |

## Proto 消息类型

### Entity 相关 (`entity.proto`)

```protobuf
message Entity {
    string id = 1;
    string entity_type = 2;    // "player" | "npc"
    string name = 3;
    string realm = 4;
    WorldPosition position = 5;
    Attributes attributes = 6;
    Karma karma = 7;
    string status = 8;
    int64 created_at = 9;
    int64 updated_at = 10;
}

message WorldPosition {
    string region_id = 1;
    double x = 2;
    double y = 3;
}

message Attributes {
    double qi = 1;
    double max_qi = 2;
    double spiritual_power = 3;
    double max_spiritual_power = 4;
    double divine_sense = 5;
    int32 comprehension = 6;
    int32 constitution = 7;
    int32 luck = 8;
    double cultivation_progress = 9;
    double attack_power = 10;
    double defense = 11;
    double speed = 12;
    int32 mental_stability = 13;
    int32 remaining_lifespan = 14;
    int32 max_lifespan = 15;
}

message Karma {
    int32 karma_value = 1;
    int32 merit = 2;
    string heavenly_mark = 3;
}

message SpiritStones {
    int64 low_grade = 1;
    int64 medium_grade = 2;
    int64 high_grade = 3;
    int64 premium_grade = 4;
}
```

> 注：Proto 中的 Attributes 仅包含 15 个字段。Go 类型 `types.Attributes` 定义了 83+ 属性，其余属性尚未在 Proto 中定义。

### World 相关 (`world.proto`)

```protobuf
message Region {
    string id = 1;
    string name = 2;
    string parent_region_id = 3;
    double spiritual_density = 4;
    int32 spiritual_tier = 5;
    int32 danger_level = 6;
    repeated Resource resources = 7;
    RegionRules rules = 8;
    string description = 9;
    string lore = 10;
}

message Resource {
    string id = 1;
    string name = 2;
    string type = 3;
    int32 rarity = 4;
    int32 quantity = 5;
    double respawn_rate = 6;
    int64 last_harvested = 7;
}

message RegionRules {
    bool is_restricted = 1;
    string restricted_by = 2;
    double tax_rate = 3;
    repeated string forbidden_actions = 4;
}

message WorldEvent {
    string id = 1;
    string name = 2;
    string type = 3;
    string description = 4;
    string region_id = 5;
    int64 start_time = 6;
    int64 end_time = 7;
    repeated string participant_ids = 8;
    string status = 9;
}

message WorldState {
    int64 epoch = 1;
    map<string, Region> regions = 2;
    repeated WorldEvent active_events = 3;
    BalanceMetrics balance_metrics = 4;
    int64 last_updated = 5;
}

message BalanceMetrics {
    double power_distribution = 1;
    double resource_circulation = 2;
    double sect_diversity = 3;
    double karma_distribution = 4;
}
```

## Go 类型定义

### 实体类型 (`types/entity.go`)

**EntityID / EntityType**:
```go
type EntityID string
type EntityType string  // "player" | "npc"
```

**CultivationRealm (10 个境界)**:
```go
RealmMortal         // mortal - 凡人
RealmQiCondensation // qi_condensation - 凝气期
RealmFoundation     // foundation - 筑基期
RealmGoldenCore     // golden_core - 金丹期
RealmNascentSoul    // nascent_soul - 元婴期
RealmSoulTransform  // soul_transformation - 化神期
RealmVoidRefinement // void_refinement - 炼虚期
RealmIntegration    // integration - 合体期
RealmMahayana       // mahayana - 大乘期
RealmTribulation    // tribulation - 渡劫期
```

**EntityStatus (8 种状态)**:
```go
StatusNormal, StatusCultivating, StatusCombat, StatusResting,
StatusDead, StatusExploring, StatusCrafting, StatusMeditating
```

**Attributes (83+ 属性)** - 按类别分组:

| 类别 | 属性数量 | 关键属性 |
|------|----------|----------|
| 基础属性 | 4 | Age, Gender, Appearance, Charisma |
| 修炼属性 | 9 | Qi, MaxQi, SpiritualPower, CultivationProgress, Comprehension 等 |
| 战斗属性 | 9 | AttackPower, Defense, Speed, CritRate, CritDamage, DodgeRate 等 |
| 灵根 | 4 | SpiritualRoots[], RootPurity, RootAwakened, MutatedRoot |
| 心境 | 5 | MentalStability, ObsessionCount, DaoHeart, InnerDemonResistance, Enlightenment |
| 生活技能 | 8 | AlchemyLevel, ArtificingLevel, FormationLevel, FireControl 等 |
| 社交属性 | 9 | Reputation, SectContribution, FactionStandings, MentorID 等 |
| 财富属性 | 4 | SpiritStones, PropertyValue, RealEstate, BusinessIncome |
| 特殊属性 | 6 | Bloodline, Physique, Destiny, WorldFavor 等 |
| 法则属性 | 5 | Laws{}, LawResonance, DomainPower, DomainRange, LawSuppression |
| 道属性 | 6 | DaoSeedType, DaoSeedLevel, DaoMarks, DestinyPath 等 |
| 寿命属性 | 3 | RemainingLifespan, MaxLifespan, AgingPenalty |
| 状态效果 | 5 | Injuries[], Buffs[], Debuffs[], PoisonLevel, CurseLevel |

### 功法类型 (`types/method.go`)

**CultivationMethod (60+ 属性)**:
- 基本信息: ID, Name, CreatorID, OriginSect, Rank, Category, ElementAffinity
- 修炼加成: CultivationSpeedMult, SpiritualPowerCapMult, QiCapMult, LifespanBonus 等
- 战斗加成: AttackBonuses{}, DefenseBonuses{}, UtilityBonuses{}
- 特殊效果: PassiveEffects[], ActiveSkills[], UltimateSkill
- 法则亲和: LawAffinities[], LawComprehensionBonus, DaoCompatibility[]
- 限制: RequiredRoots[], RequiredPhysique[], RealmRequirement, KarmaThreshold
- 传承: ParentMethodID, EvolutionPath[], TransmissionMode, CanModify, Complexity
- 评价: PowerScore, Potential, Popularity

**Rank 等级**: 天/地/玄/黄 x 上/中/下/极品 (共 12 级)

**Category 分类**: 主修功法 / 秘术 / 身法 / 神识 / 辅助 / 生活

**Skill**:
- ID, Name, Description, Category (active/passive/ultimate), Element
- DamageMult, Cooldown, ManaCost, Range, Duration
- Effects[]: Type (damage/heal/buff/debuff/shield/teleport), Value, Target, Duration, Condition

**EntityMethod** (实体已学功法):
- MethodID, EntityID, MasteryLevel (0-100%), IsMainMethod
- LearnedAt, LastPracticed, BacklashRisk (0-1), Modified, ModifiedNotes

### 物品类型 (`types/item.go`)

| 类型 | 说明 |
|------|------|
| `ItemTemplate` | 物品蓝图 (ID, Name, Type, SubType, Rank, BaseValue, Stackable, Tradeable 等) |
| `EntityInventory` | 实体背包 (EntityID, Items[], Capacity) |
| `InventoryItem` | 背包物品 (TemplateID, InstanceID, Quantity, Slot, Quality, Bound, ExpiryTime) |
| `EntityEquipment` | 装备栏 (10 个槽位: Weapon, Armor, Helmet, Boots, Necklace, Ring1, Ring2, InnerArmor, Waist, Bracelet) |
| `EquipmentItem` | 装备 (TemplateID, Rank, Quality, Level, Durability, Stats{}, Enchantments[], Soulbound) |
| `Enchantment` | 附魔 (Name, Stat, Value, Tier 1-10) |
| `Pill` | 丹药 (Type: cultivation/healing/breakthrough/combat, Effect, SuccessRate, Toxicity, Cooldown) |
| `Artifact` | 法宝 (Type: flying_sword/mirror/bell/pagoda/banner, Grade: 凡器/地器/天器/古宝, Energy, Refined) |
| `Talisman` | 符箓 (Type: attack/defense/utility/movement, Charges, MaxCharges) |
| `Recipe` | 配方 (Type, ResultItemID, Materials[], BaseSuccessRate, IsSecret) |
| `Material` | 原材料 (Type: herb/ore/beast_part/essence, Purity, Attributes{}, Source, DropRate) |

### 社交类型 (`types/social.go`)

**Sect (宗门)**:
- ID, Name, FounderID, Philosophy, EntryRequirements{}, Territory[], Rules{}
- Alignment (正道/魔道/中立), Prestige, Wealth, FacilityScore, CultivationResources[]

**SectMember (宗门成员)**:
- SectID, EntityID, Rank (职位), Contribution, JoinedAt, Privileges[]

**Relationship (人际关系)**:
- ID, EntityAID, EntityBID, RelationshipType (师徒/仇敌/盟友/恋人/结义等)
- Strength (0-100), History, CreatedAt

**NPCPersonality (NPC 性格)**:
- NPCID, PersonalityType, MoralAlignment (lawful_good/neutral/chaotic_evil 等)
- AmbitionLevel (1-100), RiskTolerance (0-1), SocialPreference
- BackgroundStory, CurrentGoal, HiddenSecrets[], LLMSystemPrompt, BehaviorTreeConfig{}, InitialActions[]

**NPCDecisionLog (NPC 决策日志)**:
- ID, NPCID, DecisionType, Context{}, ActionTaken{}, Reasoning
- ModelUsed (deepseek-chat/deepseek-reasoner), Source (behavior_tree/llm)
- TokenCost, Timestamp

### 世界类型 (`types/world.go`)

- `RegionID`: 区域 ID 类型
- `Region`: ID, Name, ParentRegionID, SpiritualDensity, SpiritualTier, DangerLevel, Resources[], Rules{}, Description, Lore
- `Resource`: ID, Name, Type, Rarity, Quantity, RespawnRate, LastHarvested
- `RegionRules`: IsRestricted, RestrictedBy, TaxRate, ForbiddenActions[]
- `WorldEvent`: ID, Name, Type, Description, RegionID, StartTime, EndTime, Participants[], Status
- `WorldState`: Epoch, Regions{}, ActiveEvents[], BalanceMetrics, LastUpdated
- `BalanceMetrics`: PowerDistribution, ResourceCirculation, SectDiversity, KarmaDistribution

### 操作类型 (`types/operation.go`)

**14 种 ActionType**:
```
cultivate, breakthrough, combat, explore, gather, craft,
create_method, trade, form_sect, join_sect, send_message,
cast_spell, meditate, sleep, move
```

**Operation**:
- ID, ActorID, ActionType, Params{}, Timestamp, Signature

**OperationResult**:
- Success, Message, Effects{}, Timestamp

**ValidationResult**:
- Valid, Errors[], Warnings[]

### WebSocket 消息类型 (`types/message.go`)

**9 种 MessageType**:
```
auth, auth_result, operation, op_result, state_sync,
entity_update, world_event, chat, system, error
```

**Message**: Type, Payload{}, Timestamp, RequestID

**载荷类型**:
- `AuthPayload`: Username, Password, Token
- `AuthResultPayload`: Success, Token, Entity, Message
- `StateSyncPayload`: Entity, Region, NearbyEntities[], WorldTime
- `ChatPayload`: SenderID, SenderName, Channel, Content
- `ErrorPayload`: Code, Message

## HTTP API (Gateway)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/auth/register` | 注册 (username, password) |
| POST | `/auth/login` | 登录 (username, password) |
| GET | `/ws?token=xxx` | WebSocket 连接 (需 JWT token) |
| GET | `/health` | 健康检查 |

## 错误码

| 错误码 | 名称 | 说明 |
|--------|------|------|
| 0 | ErrUnknown | 未知错误 |
| 1 | ErrInvalidOperation | 无效操作 |
| 2 | ErrUnauthorized | 未授权 |
| 3 | ErrEntityNotFound | 实体不存在 |
| 4 | ErrRegionNotFound | 区域不存在 |
| 5 | ErrInsufficientResources | 资源不足 |
| 6 | ErrInvalidParams | 参数无效 |
| 7 | ErrCooldownActive | 冷却中 |
| 8 | ErrOperationFailed | 操作失败 |
| 9 | ErrInternalError | 内部错误 |
| 10 | ErrServiceUnavailable | 服务不可用 |

预定义错误实例:
- `ErrInvalidOperationType`: 无效的操作类型
- `ErrUnauthorizedAccess`: 未授权访问
- `ErrEntityNotFound_`: 实体不存在
- `ErrRegionNotFound_`: 区域不存在
- `ErrInsufficientFunds`: 灵石不足
- `ErrInvalidParams_`: 参数无效
- `ErrCooldownActive_`: 冷却中
- `ErrBreakthroughFailed`: 突破失败
- `ErrInternalError_`: 内部服务器错误
