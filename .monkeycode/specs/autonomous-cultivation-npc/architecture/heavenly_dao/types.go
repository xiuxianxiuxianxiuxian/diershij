package heavenlydao

import "time"

type EntityID string
type RegionID string

type ActionType string

const (
	ActionCultivate    ActionType = "cultivate"
	ActionBreakthrough ActionType = "breakthrough"
	ActionCombat       ActionType = "combat"
	ActionTrade        ActionType = "trade"
	ActionExplore      ActionType = "explore"
	ActionCreateMethod ActionType = "create_method"
)

type WorldTime struct {
	WorldDay   int
	WorldHour  int
	EpochName  string
	TideLevel  string
	ServerTime time.Time
}

type LocationInfo struct {
	RegionID          RegionID
	RegionName        string
	SpiritualDensity  float64
	SpiritualTier     int
	DangerLevel       int
	TerrainTags       []string
	ActiveRuleEffects map[string]float64
}

type EntitySnapshot struct {
	ID                    EntityID
	Name                  string
	Realm                 string
	RealmLevel            int
	Karma                 float64
	Merit                 float64
	Luck                  float64
	MentalStability       float64
	Comprehension         float64
	AttackPower           float64
	Defense               float64
	CritRate              float64
	CritDamage            float64
	CurrentHP             float64
	MaxHP                 float64
	CurrentSP             float64
	MaxSP                 float64
	SpiritualRoots        map[string]int
	SpiritStoneBalances   SpiritStoneBalances
	ActiveMethodID        string
	ActiveMethodQuality   float64
	ConsumedBonuses       map[string]float64
	FactionID             string
	RelationshipTags      []string
	RecentBreakthroughs7d int
}

type SpiritStoneBalances struct {
	LowGrade     int64
	MediumGrade  int64
	HighGrade    int64
	PremiumGrade int64
}

type RuleContext struct {
	ActorID     EntityID
	TargetID    EntityID
	ActionType  ActionType
	Params      map[string]any
	Location    *LocationInfo
	WorldTime   WorldTime
	ActorState  *EntitySnapshot
	TargetState *EntitySnapshot
	Seed        int64
}

type DamageResult struct {
	BaseDamage       float64
	FinalDamage      float64
	RealmModifier    float64
	ElementModifier  float64
	DefenseReduction float64
	CritTriggered    bool
	CritMultiplier   float64
}

type BreakthroughResult struct {
	Success                 bool
	SuccessProbability      float64
	TribulationTriggered    bool
	TribulationProbability  float64
	TribulationStrength     float64
	CultivationLossPercent  float64
	CooldownHours           int
	MentalDamage            float64
	Reason                  string
}

type KarmaResult struct {
	BaseValue          float64
	ContextMultiplier  float64
	RelationshipFactor float64
	FinalDelta         float64
}

type RuleEvent struct {
	Name      string
	ActorID   EntityID
	TargetID  EntityID
	Payload   map[string]any
	CreatedAt time.Time
}
