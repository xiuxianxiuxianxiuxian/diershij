package types

// CultivationMethod represents a cultivation technique with 60+ attributes.
type CultivationMethod struct {
	// ===== Basic information =====
	ID              string `json:"id"`
	Name            string `json:"name"`
	CreatorID       string `json:"creator_id"`
	OriginSect      string `json:"origin_sect"`
	Rank            string `json:"rank"`               // 天/地/玄/黄 x 上/中/下/极品 (12 levels)
	Category        string `json:"category"`           // 主修功法/秘术/身法/神识/辅助/生活
	ElementAffinity string `json:"element_affinity"`   // 金木水火土风雷冰光暗/无
	Description     string `json:"description"`
	Version         int    `json:"version"`

	// ===== Core cultivation bonuses =====
	CultivationSpeedMult    float64 `json:"cultivation_speed_mult"`     // e.g. 1.5x
	SpiritualPowerCapMult   float64 `json:"spiritual_power_cap_mult"`  // spiritual power limit multiplier
	QiCapMult               float64 `json:"qi_cap_mult"`               // qi limit multiplier
	DivineSenseCapMult      float64 `json:"divine_sense_cap_mult"`     // divine sense limit multiplier
	LifespanBonus           int     `json:"lifespan_bonus"`            // extra years
	RecoverySpeedMult       float64 `json:"recovery_speed_mult"`       // recovery speed multiplier

	// ===== Combat bonuses =====
	AttackBonuses  map[string]float64 `json:"attack_bonuses"`  // e.g. {"fire_damage": 1.2, "penetration": 0.5}
	DefenseBonuses map[string]float64 `json:"defense_bonuses"` // e.g. {"magic_resist": 0.3, "damage_reduction": 0.1}
	UtilityBonuses map[string]float64 `json:"utility_bonuses"` // e.g. {"loot_rate": 0.1, "alchemy_success": 0.2}

	// ===== Special effects =====
	PassiveEffects []string `json:"passive_effects"`  // e.g. "mana_shield", "life_steal", "reflect"
	ActiveSkills   []Skill  `json:"active_skills"`    // skills provided by this method
	UltimateSkill  *Skill   `json:"ultimate_skill"`   // ultimate skill unlocked at mastery

	// ===== Law and Dao affinity =====
	LawAffinities         []string `json:"law_affinities"`          // laws accelerated by this method
	LawComprehensionBonus float64  `json:"law_comprehension_bonus"` // law comprehension speed multiplier
	DaoCompatibility      []string `json:"dao_compatibility"`       // compatible daos (e.g. "sword_dao" for sword cultivators)

	// ===== Restrictions =====
	RequiredRoots      []string `json:"required_roots"`      // required spiritual roots (e.g. ["fire", "metal"])
	RequiredPhysique   []string `json:"required_physique"`   // required physiques (e.g. ["pure_yang_body"])
	RealmRequirement   string   `json:"realm_requirement"`   // minimum realm requirement
	AlignmentRestriction string `json:"alignment_restriction"` // 正/魔/中立/无
	KarmaThreshold     int      `json:"karma_threshold"`     // karma threshold (cannot cultivate if exceeded)
	GenderRestriction  string   `json:"gender_restriction"`  // 男/女/无

	// ===== Inheritance and evolution =====
	ParentMethodID  string   `json:"parent_method_id"`  // derived from which method
	EvolutionPath   []string `json:"evolution_path"`    // evolution directions (upgraded versions)
	TransmissionMode string  `json:"transmission_mode"` // 玉简/口授/血脉/神念
	CanModify       bool     `json:"can_modify"`        // whether it can be modified by successors
	Complexity      int      `json:"complexity"`        // complexity score (affects learning difficulty and backlash risk)

	// ===== Evaluation attributes =====
	PowerScore int `json:"power_score"` // comprehensive power score (calculated by system)
	Potential  int `json:"potential"`   // potential score (affects later evolution ceiling)
	Popularity int `json:"popularity"`  // popularity (how many entities learned it)
}

// Skill represents an active or passive skill associated with a cultivation method.
type Skill struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"` // active, passive, ultimate
	Element     string  `json:"element"`
	DamageMult  float64 `json:"damage_mult"`  // damage multiplier
	Cooldown    int     `json:"cooldown"`     // cooldown in seconds
	ManaCost    float64 `json:"mana_cost"`    // spiritual power cost
	Range       float64 `json:"range"`        // effect range in meters
	Duration    int     `json:"duration"`     // effect duration in seconds
	Effects     []SkillEffect `json:"effects"`
}

// SkillEffect represents a specific effect of a skill.
type SkillEffect struct {
	Type      string  `json:"type"`      // damage, heal, buff, debuff, shield, teleport, etc.
	Value     float64 `json:"value"`     // effect value
	Target    string  `json:"target"`    // self, enemy, area, ally
	Duration  int     `json:"duration"`  // effect duration
	Condition string  `json:"condition"` // trigger condition
}

// EntityMethod represents a method that an entity has learned and is cultivating.
type EntityMethod struct {
	MethodID       string  `json:"method_id"`
	EntityID       string  `json:"entity_id"`
	MasteryLevel   float64 `json:"mastery_level"`   // 0-100%
	IsMainMethod   bool    `json:"is_main_method"`   // whether this is the primary cultivation method
	LearnedAt      int64   `json:"learned_at"`       // timestamp when learned
	LastPracticed  int64   `json:"last_practiced"`   // timestamp of last practice
	BacklashRisk   float64 `json:"backlash_risk"`    // current backlash risk (0-1)
	Modified       bool    `json:"modified"`         // whether the entity has modified this method
	ModifiedNotes  string  `json:"modified_notes"`   // notes about modifications
}
