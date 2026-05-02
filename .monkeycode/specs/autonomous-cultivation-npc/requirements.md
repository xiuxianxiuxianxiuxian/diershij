# Requirements Document

## Introduction

本系统是一个**完全自主演化的多人在线文字MUD修仙世界**，核心理念是：

1. **去中心化自治**：无预设任务或剧情引导，所有修炼、突破、探索、结盟、背叛均由参与者自主决定
2. **有机社会演化**：门派、散修联盟、拍卖行、功法交易市场由社群自发形成，强弱更替自然发生
3. **无上限创造性**：玩家可自创功法、炼丹配方、阵法，可能性仅受想象力与资源限制
4. **动态经济系统**：灵石、法宝、灵脉资源的价值由实际供需决定，无系统商店强行定价
5. **力量即自由**：顶级强者可制定规则（划定禁区、设立税收），弱小者可结盟或钻研奇术破局
6. **反脆弱世界生态**：无安全区或绝对平衡机制，世界在危机（妖兽潮、天魔入侵）中自我调整

**核心技术创新**：NPC与现实玩家拥有完全一致的操作权限和接口，NPC通过混合AI系统（行为树 + LLM）自主决策，使得第二世界中每个实体都是自治个体。

## Glossary

- **玩家 (Player)**: 由真实人类通过客户端操作的修仙角色
- **NPC (Non-Player Character)**: 由AI系统自主控制的修仙角色，拥有与玩家相同的操作权限和接口
- **统一操作接口 (Unified Action Interface)**: 玩家和NPC共同使用的操作抽象层，包含修炼、战斗、社交、交易等行为
- **行为树 (Behavior Tree)**: 基于规则的AI决策系统，处理NPC的日常行为模式
- **LLM驱动 (LLM-Driven)**: 使用大语言模型处理NPC的复杂决策和自然语言交互
- **境界 (Cultivation Realm)**: 修仙者的修为等级，如炼气、筑基、金丹、元婴、化神等
- **游戏服务器 (Game Server)**: 处理多人在线同步、状态管理、世界演化的后端服务
- **桌面客户端 (Desktop Client)**: 基于Tauri框架的桌面端游戏界面
- **天道系统 (Heavenly Dao System)**: 底层规则引擎，处理因果业力、天劫惩罚、世界平衡
- **功法 (Cultivation Method)**: 修炼法门，可由玩家/NPC自创或传承
- **世界状态 (World State)**: 当前世界的整体状态，包括势力分布、资源分布、规则变更等

## Requirements

### Requirement 1: 统一操作接口

**User Story:** AS 游戏架构师, I WANT 玩家和NPC共享同一套操作接口, SO THAT NPC能够执行玩家能做的所有操作，实现真正的平等自治

#### Acceptance Criteria

1. WHEN 定义游戏操作, THE 系统 SHALL 提供统一的操作抽象层, 包含修炼、战斗、社交、交易、探索、功法创建等基础行为
2. WHEN 玩家执行操作, THE 系统 SHALL 通过客户端输入调用统一操作接口
3. WHEN NPC执行操作, THE 系统 SHALL 通过AI决策模块调用统一操作接口
4. WHILE 操作执行中, THE 系统 SHALL 记录操作日志用于调试和回放
5. IF 操作参数非法, THE 系统 SHALL 返回错误信息并拒绝执行

### Requirement 2: 角色创建与属性系统

**User Story:** AS 修仙者, I WANT 创建角色并自由发展, SO THAT 我能够按照自己的道途修行

#### Acceptance Criteria

1. WHEN 新角色首次进入世界, THE 系统 SHALL 提供基础角色创建流程, 包含姓名、初始灵根属性分配
2. WHEN 角色创建完成, THE 系统 SHALL 生成唯一角色ID并初始化基础属性（气血、灵力、悟性、根骨等）
3. WHILE 角色在线, THE 系统 SHALL 实时更新并同步属性变化
4. WHEN 角色查看自身状态, THE 系统 SHALL 显示当前境界、属性、功法、背包、因果业力等完整状态
5. IF 角色名称已存在, THE 系统 SHALL 提示重新输入

### Requirement 3: 自主修仙与境界系统

**User Story:** AS 修仙者, I WANT 自主选择修炼路径并突破境界, SO THAT 我能够按照自己的道途变强

#### Acceptance Criteria

1. WHEN 角色执行修炼操作, THE 系统 SHALL 根据功法、灵根、环境、心境计算修炼效率
2. WHEN 修为达到突破阈值, THE 系统 SHALL 触发境界突破流程, 包含天劫判定
3. WHILE 境界突破中, THE 系统 SHALL 计算成功率并处理成功或失败结果
4. WHEN 境界提升, THE 系统 SHALL 更新角色属性上限并解锁新能力
5. IF 突破失败, THE 系统 SHALL 施加惩罚（修为损失、心境损伤、cooldown时间）

### Requirement 4: 功法与道法自创系统

**User Story:** AS 修仙者, I WANT 自创功法和道法, SO THAT 我能够开辟独特的修行路径

#### Acceptance Criteria

1. WHEN 角色尝试自创功法, THE 系统 SHALL 验证角色是否拥有至少10000极品灵石
2. WHEN 验证通过, THE 系统 SHALL 扣除10000极品灵石并调用DeepSeek R1进行功法验证和生成
3. WHEN 功法创建完成, THE 系统 SHALL 生成功法实体并记录创始人信息
4. WHILE 功法传承中, THE 系统 SHALL 追踪功法版本演化树
5. WHEN 其他角色学习功法, THE 系统 SHALL 根据创始人设定决定是否需要传承许可
6. IF 功法存在逻辑矛盾, THE 系统 SHALL 拒绝创建并退还极品灵石

### Requirement 5: NPC自主决策系统（混合AI）

**User Story:** AS 游戏世界, I WANT NPC拥有与玩家相同的自主决策能力, SO THAT 世界中的每个实体都是真正的自治个体

#### Acceptance Criteria

1. WHEN NPC需要执行日常行为, THE 系统 SHALL 使用行为树模块处理修炼、采集、打坐等规则化操作
2. WHEN NPC面临复杂决策（结盟、背叛、自创功法、建立宗门）, THE 系统 SHALL 调用LLM模块进行分析和决策
3. WHILE NPC在线, THE 系统 SHALL 每周期评估当前状态并选择下一步行动
4. WHEN NPC与其他角色交互, THE 系统 SHALL 生成符合角色设定和当前情境的自然语言对话
5. IF LLM响应超时, THE 系统 SHALL 降级使用行为树默认决策

### Requirement 6: 动态社会与经济系统

**User Story:** AS 世界参与者, I WANT 社会关系和经济体系完全由参与者自发形成, SO THAT 世界具有真实的社会演化体验

#### Acceptance Criteria

1. WHEN 角色创建宗门, THE 系统 SHALL 提供宗门创建接口, 包含宗门名称、理念、入门条件设定
2. WHEN 角色发起交易, THE 系统 SHALL 提供玩家/NPC间交易接口, 无系统定价干预
3. WHILE 市场运行中, THE 系统 SHALL 记录交易历史供参与者参考
4. WHEN 宗门设立规则（税收、禁区）, THE 系统 SHALL 在宗门势力范围内生效
5. IF 宗门瓦解, THE 系统 SHALL 释放相关资源并通知成员

### Requirement 7: 天道与因果系统

**User Story:** AS 世界规则, I WANT 通过底层机制维持世界秩序, SO THAT 世界能够自我调节而非依赖人为干预

#### Acceptance Criteria

1. WHEN 角色做出重大行为, THE 系统 SHALL 记录因果业力值
2. WHEN 业力值触发阈值, THE 系统 SHALL 引入天劫或机缘进行平衡
3. WHILE 世界运行中, THE 系统 SHALL 周期性评估世界生态健康度
4. WHEN 触发世界危机事件, THE 系统 SHALL 生成妖兽潮、天魔入侵等事件
5. IF 某势力过于强大, THE 系统 SHALL 通过自然机制（非人为干预）引入制衡因素

### Requirement 8: 多人在线同步

**User Story:** AS 世界参与者, I WANT 与其他实体在同一世界实时互动, SO THAT 我能够体验有机演化的修仙社会

#### Acceptance Criteria

1. WHEN 角色连接到世界, THE 系统 SHALL 建立WebSocket连接并同步当前世界状态
2. WHILE 角色在线, THE 系统 SHALL 实时广播区域内其他实体（玩家和NPC）的状态变化
3. WHEN 实体执行操作, THE 系统 SHALL 验证操作合法性并同步结果给相关参与者
4. WHEN 角色进入新区域, THE 系统 SHALL 加载该区域的实体列表和环境信息
5. IF 连接断开, THE 系统 SHALL 保存角色状态并在重连时恢复

### Requirement 9: 桌面客户端界面

**User Story:** AS 世界参与者, I WANT 通过桌面客户端进入修仙世界, SO THAT 我获得沉浸式的文字MUD体验

#### Acceptance Criteria

1. WHEN 客户端启动, THE 系统 SHALL 显示登录界面并连接到世界服务器
2. WHEN 参与者发送指令, THE 系统 SHALL 以文本形式展示操作过程和结果
3. WHILE 游戏进行中, THE 系统 SHALL 提供主界面包含：聊天区、角色状态区、操作指令区、世界日志区、势力地图
4. WHEN 收到服务器消息, THE 系统 SHALL 实时渲染到对应界面区域
5. IF 服务器响应延迟, THE 系统 SHALL 显示加载状态提示

### Requirement 10: 游戏服务器架构

**User Story:** AS 世界系统, I WANT 高并发的服务器架构, SO THAT 支持大量玩家和NPC同时在线自主演化

#### Acceptance Criteria

1. WHEN 服务器启动, THE 系统 SHALL 初始化数据库连接、WebSocket监听器、AI调度模块、天道引擎
2. WHILE 服务器运行, THE 系统 SHALL 每个游戏周期处理所有在线实体的AI决策
3. WHEN 实体数量超过阈值, THE 系统 SHALL 支持水平扩展新增服务器节点
4. WHEN 世界状态变更, THE 系统 SHALL 持久化存储到数据库
5. IF 服务器节点故障, THE 系统 SHALL 转移实体到其他节点并恢复连接

## Decisions Made

| 决策项 | 选择 |
|--------|------|
| 实体规模 | 根据节点配置自动调整，动态伸缩 |
| LLM集成方式 | DeepSeek API（v3日常决策 + R1复杂推理） |
| 世界初始化 | 丰富的预置内容（多区域、NPC群体、资源、势力、世界历史等） |

## Pending Clarifications

以下问题需要与用户进一步确认以完善需求：

1. **世界历史与传说**：是否需要预置世界背景故事和传说？（如上古大战、陨落大能、秘境起源等，影响NPC行为和世界探索）
2. **初始宗门数量**：预置几个基础宗门？（建议2-3个不同理念的宗门，如正道、魔道、中立散修联盟）
3. **NPC初始修为分布**：初始NPC的修为跨度？（如大部分炼气/筑基，少数金丹作为世界"锚点"）
