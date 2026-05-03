package types

import "time"

// 实体（用户/角色）
type Entity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Realm       string `json:"realm"`
	EntityType  string `json:"entity_type"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录响应
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Entity  Entity `json:"entity"`
	Message string `json:"message,omitempty"`
}

// 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 注册响应
type RegisterResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Entity  Entity `json:"entity"`
	Message string `json:"message,omitempty"`
}

// 用户
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Level     int       `json:"level"`
	Exp       int       `json:"exp"`
	CultivationLevel string `json:"cultivation_level,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// 角色
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
}

// 世界状态
type WorldState struct {
	CurrentMap     string            `json:"current_map"`
	PlayersOnline  int               `json:"players_online"`
	Events         []WorldEvent      `json:"events,omitempty"`
	Announcements  []Announcement    `json:"announcements,omitempty"`
}

// 世界事件
type WorldEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

// 公告
type Announcement struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Priority int   `json:"priority"`
}

// 战斗状态
type CombatState struct {
	InCombat      bool         `json:"in_combat"`
	CurrentEnemy  *Enemy       `json:"current_enemy,omitempty"`
	BattleLog     []CombatLog  `json:"battle_log,omitempty"`
	TurnNumber    int          `json:"turn_number"`
}

// 敌人
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

// 战斗日志
type CombatLog struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
}

// 技能
type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Damage      int    `json:"damage"`
	EnergyCost  int    `json:"energy_cost"`
	Description string `json:"description"`
}

// 社交信息
type SocialInfo struct {
	Friends     []Friend         `json:"friends,omitempty"`
	Guild       *Guild           `json:"guild,omitempty"`
	Messages    []Message        `json:"messages,omitempty"`
	Requests    []FriendRequest  `json:"friend_requests,omitempty"`
}

// 好友
type Friend struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Level    int       `json:"level"`
	Online   bool      `json:"online"`
	LastSeen time.Time `json:"last_seen"`
}

// 门派/公会
type Guild struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Members int    `json:"members"`
	Leader  string `json:"leader"`
}

// 消息
type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
	Read       bool      `json:"read"`
}

// 好友请求
type FriendRequest struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	FromName  string    `json:"from_name"`
	Timestamp time.Time `json:"timestamp"`
}

// 设置
type Settings struct {
	AudioVolume       float64 `json:"audio_volume"`
	MusicVolume       float64 `json:"music_volume"`
	ShowDamageNumbers bool    `json:"show_damage_numbers"`
	AutoPlay          bool    `json:"auto_play"`
	ShowFPS           bool    `json:"show_fps"`
	Language          string  `json:"language"`
}

// 应用状态
type AppState struct {
	CurrentWindow string
	IsLoggedIn    bool
	User          *User
	Character     *Character
	World         *WorldState
	Combat        *CombatState
	Social        *SocialInfo
	Settings      *Settings
}

// WebSocket 消息类型
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

// WebSocket 消息
type WSMessage struct {
	Type    WSMessageType          `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// OperationResult 操作结果
type OperationResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
