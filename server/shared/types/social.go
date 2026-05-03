package types

// Sect represents a cultivation sect/faction.
type Sect struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	FounderID        string            `json:"founder_id"`
	Philosophy       string            `json:"philosophy"`         // 宗门理念
	EntryRequirements map[string]any   `json:"entry_requirements"`  // 入门条件
	Territory        []string          `json:"territory"`          // 势力范围 (region IDs)
	Rules            map[string]any    `json:"rules"`              // 宗门规则
	Alignment        string            `json:"alignment"`          // 正道/魔道/中立
	CreatedAt        int64             `json:"created_at"`
	MemberCount      int               `json:"member_count"`
	Prestige         int               `json:"prestige"`           // 宗门声望
	Wealth           int64             `json:"wealth"`             // 宗门财富 (spirit stones)
	FacilityScore    int               `json:"facility_score"`     // 设施评分
	CultivationResources []string      `json:"cultivation_resources"` // 宗门修炼资源
}

// SectMember represents a member's information within a sect.
type SectMember struct {
	SectID       string  `json:"sect_id"`
	EntityID     string  `json:"entity_id"`
	Rank         string  `json:"rank"`            // 职位
	Contribution float64 `json:"contribution"`     // 贡献值
	JoinedAt     int64   `json:"joined_at"`
	Privileges   []string `json:"privileges"`      // 特权列表
}

// Relationship represents a relationship between two entities.
type Relationship struct {
	ID               string  `json:"id"`
	EntityAID        string  `json:"entity_a_id"`
	EntityBID        string  `json:"entity_b_id"`
	RelationshipType string  `json:"relationship_type"` // 师徒/仇敌/盟友/恋人/结义等
	Strength         float64 `json:"strength"`          // 关系强度 (0-100)
	History          string  `json:"history"`           // 关系历史
	CreatedAt        int64   `json:"created_at"`
}

// NPCPersonality represents an NPC's personality configuration.
type NPCPersonality struct {
	NPCID            string  `json:"npc_id"`
	PersonalityType  string  `json:"personality_type"`   // 性格类型
	MoralAlignment   string  `json:"moral_alignment"`    // 道德倾向 (lawful_good, neutral, chaotic_evil, etc.)
	AmbitionLevel    int     `json:"ambition_level"`     // 野心程度 (1-100)
	RiskTolerance    float64 `json:"risk_tolerance"`     // 风险承受度 (0-1)
	SocialPreference string  `json:"social_preference"`  // 社交偏好 (extrovert, introvert, balanced)
	BackgroundStory  string  `json:"background_story"`   // 背景故事
	CurrentGoal      string  `json:"current_goal"`       // 当前目标
	HiddenSecrets    []string `json:"hidden_secrets"`    // 隐藏秘密
	LLMSystemPrompt  string  `json:"llm_system_prompt"`  // LLM 系统提示词
	BehaviorTreeConfig map[string]any `json:"behavior_tree_config"` // 行为树配置
	InitialActions   []string `json:"initial_actions"`   // 初始行为模式
}

// NPCDecisionLog represents a single NPC decision record.
type NPCDecisionLog struct {
	ID            string `json:"id"`
	NPCID         string `json:"npc_id"`
	DecisionType  string `json:"decision_type"`
	Context       map[string]any `json:"context"`       // 决策上下文
	ActionTaken   map[string]any `json:"action_taken"`  // 采取的行动
	Reasoning     string `json:"reasoning"`             // 决策推理
	ModelUsed     string `json:"model_used"`            // deepseek-chat / deepseek-reasoner
	Source        string `json:"source"`                // behavior_tree / llm
	TokenCost     float64 `json:"token_cost"`           // API调用成本
	Timestamp     int64  `json:"timestamp"`
}
