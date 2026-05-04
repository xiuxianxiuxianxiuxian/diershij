package types

import "time"

// ===== API 响应类型（与服务端匹配） =====

// ServerEntity 服务端返回的完整实体
type ServerEntity struct {
	ID         string                 `json:"id"`
	EntityType string                 `json:"entity_type"`
	Name       string                 `json:"name"`
	Realm      string                 `json:"realm"`
	Position   map[string]interface{} `json:"position,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Karma      map[string]interface{} `json:"karma,omitempty"`
	Status     string                 `json:"status,omitempty"`
	CreatedAt  string                 `json:"created_at,omitempty"`
	UpdatedAt  string                 `json:"updated_at,omitempty"`
}

// AuthResponse 登录/注册统一响应
type AuthResponse struct {
	Success bool          `json:"success"`
	Token   string        `json:"token"`
	Entity  *ServerEntity `json:"entity"`
	Error   string        `json:"error,omitempty"`
}

// ===== 客户端本地实体类型 =====

// Entity 简化的用户/角色信息
type Entity struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Realm      string `json:"realm"`
	EntityType string `json:"entity_type"`
}

// Character 角色属性（客户端 UI 用）
type Character struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Level            int    `json:"level"`
	Health           int    `json:"health"`
	MaxHealth        int    `json:"max_health"`
	Energy           int    `json:"energy"`
	MaxEnergy        int    `json:"max_energy"`
	Attack           int    `json:"attack"`
	Defense          int    `json:"defense"`
	Speed            int    `json:"speed"`
	CultivationRealm string `json:"cultivation_realm,omitempty"`
	CultivationProgress float64 `json:"cultivation_progress"`
	Qi               float64 `json:"qi"`
	MaxQi            float64 `json:"max_qi"`
	SpiritualPower   float64 `json:"spiritual_power"`
	MaxSpiritualPower float64 `json:"max_spiritual_power"`

	// 修炼资质
	Comprehension    int     `json:"comprehension"`
	Constitution     int     `json:"constitution"`
	Luck             int     `json:"luck"`
	DivineSense      float64 `json:"divine_sense"`

	// 心境 / 寿元
	MentalStability     int `json:"mental_stability"`
	RemainingLifespan   int `json:"remaining_lifespan"`
	MaxLifespan         int `json:"max_lifespan"`

	// 位置
	RegionID  string  `json:"region_id"`
	PositionX float64 `json:"position_x"`
	PositionY float64 `json:"position_y"`

	// 状态
	Status string `json:"status"`

	// 业力
	KarmaValue   int    `json:"karma_value"`
	Merit        int    `json:"merit"`
	KarmicDebt   int    `json:"karmic_debt"`
	HeavenlyMark string `json:"heavenly_mark"`

	// 战斗属性
	CritRate        float64 `json:"crit_rate"`
	CritDamage      float64 `json:"crit_damage"`
	DodgeRate       float64 `json:"dodge_rate"`
	HitRate         float64 `json:"hit_rate"`
	Penetration     float64 `json:"penetration"`
	DamageReduction float64 `json:"damage_reduction"`

	// 生活技能
	AlchemyLevel    int `json:"alchemy_level"`
	ArtificingLevel int `json:"artificing_level"`
	FormationLevel  int `json:"formation_level"`
	FireControl     int `json:"fire_control"`
	HerbKnowledge   int `json:"herb_knowledge"`
	MiningSkill     int `json:"mining_skill"`
	TalismanSkill   int `json:"talisman_skill"`
	BeastTaming     int `json:"beast_taming"`

	// 社交
	Reputation       int `json:"reputation"`
	SectContribution int `json:"sect_contribution"`

	// 心境 / 灵性
	DaoHeart      int `json:"dao_heart"`
	Enlightenment int `json:"enlightenment"`
	RootPurity    int `json:"root_purity"`
	PoisonLevel   int `json:"poison_level"`
	CurseLevel    int `json:"curse_level"`

	// 灵石
	LowGradeStones     int64 `json:"low_grade_stones"`
	MediumGradeStones  int64 `json:"medium_grade_stones"`
	HighGradeStones    int64 `json:"high_grade_stones"`
	PremiumGradeStones int64 `json:"premium_grade_stones"`
}

// ===== 请求类型 =====

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ===== 世界 / 战斗 / 社交 =====

type WorldState struct {
	CurrentMap    string         `json:"current_map"`
	PlayersOnline int            `json:"players_online"`
	Events        []WorldEvent   `json:"events,omitempty"`
	Announcements []Announcement `json:"announcements,omitempty"`
}

type WorldEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type Announcement struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Priority int    `json:"priority"`
}

type CombatState struct {
	InCombat     bool        `json:"in_combat"`
	CurrentEnemy *Enemy      `json:"current_enemy,omitempty"`
	BattleLog    []CombatLog `json:"battle_log,omitempty"`
	TurnNumber   int         `json:"turn_number"`
}

type Enemy struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Level     int     `json:"level"`
	Health    int     `json:"health"`
	MaxHealth int     `json:"max_health"`
	Attack    int     `json:"attack"`
	Defense   int     `json:"defense"`
	Skills    []Skill `json:"skills,omitempty"`
}

type CombatLog struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
}

type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Damage      int    `json:"damage"`
	EnergyCost  int    `json:"energy_cost"`
	Description string `json:"description"`
}

type SocialInfo struct {
	Friends  []Friend        `json:"friends,omitempty"`
	Guild    *Guild          `json:"guild,omitempty"`
	Messages []Message       `json:"messages,omitempty"`
	Requests []FriendRequest `json:"friend_requests,omitempty"`
}

type Friend struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Level    int       `json:"level"`
	Online   bool      `json:"online"`
	LastSeen time.Time `json:"last_seen"`
}

type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Members int    `json:"members"`
	Leader  string `json:"leader"`
}

type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
	Read       bool      `json:"read"`
}

type FriendRequest struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	FromName  string    `json:"from_name"`
	Timestamp time.Time `json:"timestamp"`
}

type Settings struct {
	AudioVolume       float64 `json:"audio_volume"`
	MusicVolume       float64 `json:"music_volume"`
	ShowDamageNumbers bool    `json:"show_damage_numbers"`
	AutoPlay          bool    `json:"auto_play"`
	ShowFPS           bool    `json:"show_fps"`
	Language          string  `json:"language"`
}

type AppState struct {
	CurrentWindow string
	IsLoggedIn    bool
	User          *Entity
	Character     *Character
	World         *WorldState
	Combat        *CombatState
	Social        *SocialInfo
	Settings      *Settings
}

// ===== WebSocket =====

type WSMessageType string

const (
	WSMessageTypeCombatUpdate WSMessageType = "combat_update"
	WSMessageTypeWorldUpdate  WSMessageType = "world_update"
	WSMessageTypeSocialUpdate WSMessageType = "social_update"
	WSMessageTypeNewMessage   WSMessageType = "new_message"
	WSMessageTypeError        WSMessageType = "error"
	WSMessageTypeOperation    WSMessageType = "operation"
	WSMessageTypeOpResult     WSMessageType = "op_result"
)

type WSMessage struct {
	Type      WSMessageType          `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp int64                  `json:"timestamp,omitempty"`
}

// OperationResult 操作结果
type OperationResult struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Effects   map[string]interface{} `json:"effects,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}
