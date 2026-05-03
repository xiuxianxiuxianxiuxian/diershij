package types

import "time"

type EntityID string

type EntityType string

const (
	EntityTypePlayer EntityType = "player"
	EntityTypeNPC    EntityType = "npc"
)

type CultivationRealm string

const (
	RealmMortal         CultivationRealm = "mortal"
	RealmQiCondensation CultivationRealm = "qi_condensation"
	RealmFoundation     CultivationRealm = "foundation"
	RealmGoldenCore     CultivationRealm = "golden_core"
	RealmNascentSoul    CultivationRealm = "nascent_soul"
	RealmSoulTransform  CultivationRealm = "soul_transformation"
	RealmVoidRefinement CultivationRealm = "void_refinement"
	RealmIntegration    CultivationRealm = "integration"
	RealmMahayana       CultivationRealm = "mahayana"
	RealmTribulation    CultivationRealm = "tribulation"
)

type Entity struct {
	ID           EntityID         `json:"id"`
	EntityType   EntityType       `json:"entity_type"`
	Name         string           `json:"name"`
	Realm        CultivationRealm `json:"realm"`
	Position     WorldPosition    `json:"position"`
	Attributes   Attributes       `json:"attributes"`
	Karma        Karma            `json:"karma"`
	Status       EntityStatus     `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type WorldPosition struct {
	RegionID string  `json:"region_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

// Attributes contains all 83+ core attributes organized by category.
type Attributes struct {
	// ===== Basic attributes =====
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Appearance int    `json:"appearance"`
	Charisma   int    `json:"charisma"`

	// ===== Cultivation attributes =====
	Qi                  float64 `json:"qi"`
	MaxQi               float64 `json:"max_qi"`
	SpiritualPower      float64 `json:"spiritual_power"`
	MaxSpiritualPower   float64 `json:"max_spiritual_power"`
	DivineSense         float64 `json:"divine_sense"`
	Comprehension       int     `json:"comprehension"`
	Constitution        int     `json:"constitution"`
	Luck                int     `json:"luck"`
	CultivationProgress float64 `json:"cultivation_progress"`

	// ===== Combat attributes =====
	AttackPower       float64 `json:"attack_power"`
	Defense           float64 `json:"defense"`
	Speed             float64 `json:"speed"`
	CritRate          float64 `json:"crit_rate"`
	CritDamage        float64 `json:"crit_damage"`
	DodgeRate         float64 `json:"dodge_rate"`
	HitRate           float64 `json:"hit_rate"`
	Penetration       float64 `json:"penetration"`
	DamageReduction   float64 `json:"damage_reduction"`

	// ===== Spiritual roots =====
	SpiritualRoots  []SpiritualRoot `json:"spiritual_roots"`
	RootPurity      int             `json:"root_purity"`
	RootAwakened    bool            `json:"root_awakened"`
	MutatedRoot     string          `json:"mutated_root"`

	// ===== Mental state =====
	MentalStability     int `json:"mental_stability"`
	ObsessionCount      int `json:"obsession_count"`
	DaoHeart            int `json:"dao_heart"`
	InnerDemonResistance int `json:"inner_demon_resistance"`
	Enlightenment       int `json:"enlightenment"`

	// ===== Life skills =====
	AlchemyLevel    int `json:"alchemy_level"`
	ArtificingLevel int `json:"artificing_level"`
	FormationLevel  int `json:"formation_level"`
	FireControl     int `json:"fire_control"`
	HerbKnowledge   int `json:"herb_knowledge"`
	MiningSkill     int `json:"mining_skill"`
	TalismanSkill   int `json:"talisman_skill"`
	BeastTaming     int `json:"beast_taming"`

	// ===== Social attributes =====
	Reputation        int      `json:"reputation"`
	SectContribution  int      `json:"sect_contribution"`
	FactionStandings  map[string]int `json:"faction_standings"`
	RelationshipCount int      `json:"relationship_count"`
	MentorID          string   `json:"mentor_id"`
	DiscipleIDs       []string `json:"disciple_ids"`
	SwornSiblings     []string `json:"sworn_siblings"`
	Enemies           []string `json:"enemies"`
	Lovers            []string `json:"lovers"`

	// ===== Wealth attributes =====
	SpiritStones   SpiritStones `json:"spirit_stones"`
	PropertyValue  int          `json:"property_value"`
	RealEstate     []string     `json:"real_estate"`
	BusinessIncome int          `json:"business_income"`

	// ===== Special attributes =====
	Bloodline        string `json:"bloodline"`
	BloodlinePurity  int    `json:"bloodline_purity"`
	Physique         string `json:"physique"`
	PhysiqueAwakened bool   `json:"physique_awakened"`
	Destiny          int    `json:"destiny"`
	WorldFavor       int    `json:"world_favor"`

	// ===== Law attributes =====
	Laws           map[string]float64 `json:"laws"`
	LawResonance   int                `json:"law_resonance"`
	DomainPower    float64            `json:"domain_power"`
	DomainRange    float64            `json:"domain_range"`
	LawSuppression float64            `json:"law_suppression"`

	// ===== Dao attributes =====
	DaoSeedType          string `json:"dao_seed_type"`
	DaoSeedLevel         int    `json:"dao_seed_level"`
	DaoSeedGrowth        float64 `json:"dao_seed_growth"`
	DaoMarks             int    `json:"dao_marks"`
	DaoHeartComprehension int   `json:"dao_heart_comprehension"`
	DestinyPath          string `json:"destiny_path"`

	// ===== Lifespan attributes =====
	RemainingLifespan int     `json:"remaining_lifespan"`
	MaxLifespan       int     `json:"max_lifespan"`
	AgingPenalty      float64 `json:"aging_penalty"`

	// ===== Status effects =====
	Injuries    []Injury    `json:"injuries"`
	Buffs       []Buff      `json:"buffs"`
	Debuffs     []Debuff    `json:"debuffs"`
	PoisonLevel int         `json:"poison_level"`
	CurseLevel  int         `json:"curse_level"`
}

// SpiritualRoot represents a spiritual root element.
type SpiritualRoot struct {
	Element string `json:"element"` // gold, wood, water, fire, earth, wind, thunder, ice, light, dark, etc.
	Purity  int    `json:"purity"`  // 1-100
}

// Injury represents an entity injury.
type Injury struct {
	Type        string `json:"type"`
	Severity    int    `json:"severity"`
	Cause       string `json:"cause"`
	HealTime    int64  `json:"heal_time"`
	Description string `json:"description"`
}

// Buff represents a positive status effect.
type Buff struct {
	Name        string  `json:"name"`
	Effect      string  `json:"effect"`
	Value       float64 `json:"value"`
	Source      string  `json:"source"`
	ExpiryTime  int64   `json:"expiry_time"`
}

// Debuff represents a negative status effect.
type Debuff struct {
	Name        string  `json:"name"`
	Effect      string  `json:"effect"`
	Value       float64 `json:"value"`
	Source      string  `json:"source"`
	ExpiryTime  int64   `json:"expiry_time"`
}

type Karma struct {
	KarmaValue   int    `json:"karma_value"`
	Merit        int    `json:"merit"`
	KarmicDebt   int    `json:"karmic_debt"`
	HeavenlyMark string `json:"heavenly_mark"`
}

type EntityStatus string

const (
	StatusNormal    EntityStatus = "normal"
	StatusCultivating EntityStatus = "cultivating"
	StatusCombat    EntityStatus = "combat"
	StatusResting   EntityStatus = "resting"
	StatusDead      EntityStatus = "dead"
	StatusExploring EntityStatus = "exploring"
	StatusCrafting  EntityStatus = "crafting"
	StatusMeditating EntityStatus = "meditating"
)

type SpiritStones struct {
	LowGrade     int64 `json:"low_grade"`
	MediumGrade  int64 `json:"medium_grade"`
	HighGrade    int64 `json:"high_grade"`
	PremiumGrade int64 `json:"premium_grade"`
}

// CultivationRealmLevel returns the numeric level of a cultivation realm.
// Mortal=0, QiCondensation=1, ..., Tribulation=9.
// Returns -1 for unknown realms.
func CultivationRealmLevel(realm CultivationRealm) int {
	switch realm {
	case RealmMortal:
		return 0
	case RealmQiCondensation:
		return 1
	case RealmFoundation:
		return 2
	case RealmGoldenCore:
		return 3
	case RealmNascentSoul:
		return 4
	case RealmSoulTransform:
		return 5
	case RealmVoidRefinement:
		return 6
	case RealmIntegration:
		return 7
	case RealmMahayana:
		return 8
	case RealmTribulation:
		return 9
	default:
		return -1
	}
}
