package service

import (
    "math"
)

// FortuneRule handles fortune score calculation and opportunity trigger probability.
type FortuneRule struct {
    baseOpportunityProb float64
    meritWeight         float64
    karmaWeight         float64
}

// NewFortuneRule creates a new FortuneRule with default configuration.
func NewFortuneRule() *FortuneRule {
    return &FortuneRule{
        baseOpportunityProb: 0.01,
        meritWeight:         0.1,
        karmaWeight:         0.05,
    }
}

// FortuneInput holds the inputs needed to calculate fortune score.
type FortuneInput struct {
    Luck              int     // base luck attribute (1-100)
    Merit             int     // merit/功德 value
    Karma             int     // karma/业力 value
    SaveLifeCount     int     // times saved others in last 7 days (+5 each)
    DestroyOppCount   int     // times destroyed opportunities (-10 each)
    HeavenlyTaskCount int     // times completed heavenly tasks (+20 each)
}

// CalculateFortuneScore computes the overall fortune score for an entity.
//
// Formula:
//
//	fortune = luck + merit × 0.1 - karma × 0.05 + action_bonus
//
// Where action_bonus = save_life×5 - destroy_opp×10 + heavenly_task×20
func (f *FortuneRule) CalculateFortuneScore(input FortuneInput) float64 {
    actionBonus := float64(input.SaveLifeCount)*5.0 -
        float64(input.DestroyOppCount)*10.0 +
        float64(input.HeavenlyTaskCount)*20.0

    fortune := float64(input.Luck) +
        float64(input.Merit)*f.meritWeight -
        float64(input.Karma)*f.karmaWeight +
        actionBonus

    return math.Max(0, fortune)
}

// FortuneGrade represents the fortune tier.
type FortuneGrade string

const (
    FortuneMisfortune FortuneGrade = "misfortune" // 厄运 < 20
    FortunePlain      FortuneGrade = "plain"      // 平淡 20-40
    FortuneNormal     FortuneGrade = "normal"     // 普通 40-60
    FortuneLucky      FortuneGrade = "lucky"      // 幸运 60-80
    FortuneDestined   FortuneGrade = "destined"   // 天命 > 80
)

// GetFortuneGrade returns the fortune grade for a given score.
func (f *FortuneRule) GetFortuneGrade(score float64) FortuneGrade {
    if score < 20 {
        return FortuneMisfortune
    }
    if score < 40 {
        return FortunePlain
    }
    if score < 60 {
        return FortuneNormal
    }
    if score < 80 {
        return FortuneLucky
    }
    return FortuneDestined
}

// GetFortuneMultiplier returns the opportunity multiplier based on fortune score.
func (f *FortuneRule) GetFortuneMultiplier(score float64) float64 {
    if score < 20 {
        return 0.5
    }
    if score < 40 {
        return 0.8
    }
    if score < 60 {
        return 1.0
    }
    if score < 80 {
        return 1.3
    }
    return 2.0
}

// OpportunityInput holds inputs for opportunity trigger probability calculation.
type OpportunityInput struct {
    Fortune            float64
    SpiritualDensity   int     // location spiritual density (0-100)
    ExplorationFactor  float64 // area exploration activity (0-1+)
    BalanceFactor      float64 // world balance modifier (0.5-1.5)
}

// CalculateOpportunityProbability computes the probability of triggering an opportunity.
//
// Formula:
//
//	P = 0.01 × fortune_factor × spiritual × exploration × balance
//
// Where:
//   - fortune_factor = 1.0 + fortune / 100
//   - spiritual = spiritual_density / 100
//   - exploration: area exploration activity coefficient
//   - balance: world balance modifier (0.5-1.5)
func (f *FortuneRule) CalculateOpportunityProbability(input OpportunityInput) float64 {
    fortuneFactor := 1.0 + input.Fortune/100.0

    spiritualFactor := float64(input.SpiritualDensity) / 100.0
    spiritualFactor = math.Max(0, math.Min(1, spiritualFactor))

    explorationFactor := math.Max(0, input.ExplorationFactor)

    balanceFactor := math.Max(0.5, math.Min(1.5, input.BalanceFactor))

    return f.baseOpportunityProb * fortuneFactor * spiritualFactor * explorationFactor * balanceFactor
}

// OpportunityType represents the type of opportunity encountered.
type OpportunityType string

const (
    OpHerbDiscovery OpportunityType = "herb_discovery"    // 灵草发现 40%
    OpElderEstate   OpportunityType = "elder_estate"      // 前辈遗府 5%
    OpEpiphany      OpportunityType = "epiphany"          // 功法顿悟 10%
    OpRootAwakening OpportunityType = "root_awakening"    // 灵根觉醒 1%
    OpHeavenlyBless OpportunityType = "heavenly_blessing" // 天道赐福 2%
    OpSecretRealm   OpportunityType = "secret_realm"      // 秘境入口 3%
    OpNobleHelp     OpportunityType = "noble_help"        // 贵人相助 15%
    OpNothing       OpportunityType = "nothing"           // 无事发生 24%
)

// opportunityWeights defines the base probability weights for each opportunity type.
var opportunityWeights = map[OpportunityType]float64{
    OpHerbDiscovery: 40,
    OpElderEstate:   5,
    OpEpiphany:      10,
    OpRootAwakening: 1,
    OpHeavenlyBless: 2,
    OpSecretRealm:   3,
    OpNobleHelp:     15,
    OpNothing:       24,
}

const totalOpportunityWeight = 100.0

// ResolveOpportunityType determines which type of opportunity is triggered.
// Uses weighted random selection based on the predefined distribution.
func (f *FortuneRule) ResolveOpportunityType(randFloat func() float64) OpportunityType {
    r := randFloat() * totalOpportunityWeight

    cumulative := 0.0
    for _, opType := range []OpportunityType{
        OpHerbDiscovery, OpElderEstate, OpEpiphany, OpRootAwakening,
        OpHeavenlyBless, OpSecretRealm, OpNobleHelp, OpNothing,
    } {
        cumulative += opportunityWeights[opType]
        if r < cumulative {
            return opType
        }
    }
    return OpNothing
}

// CalculateOpportunityQualityTier returns the quality tier (1-5) based on fortune score.
//
// Formula: quality_tier = clamp(fortune / 20, 1, 5)
func (f *FortuneRule) CalculateOpportunityQualityTier(fortune float64) int {
    tier := int(fortune / 20)
    if tier < 1 {
        return 1
    }
    if tier > 5 {
        return 5
    }
    return tier
}

// Destiny represents an entity's fate/destiny track.
type Destiny string

const (
    DestinyMortal        Destiny = "mortal"         // 凡途: default
    DestinyDefyHeaven    Destiny = "defy_heaven"    // 逆天改命: 3 consecutive breakthrough successes
    DestinyFollowHeaven  Destiny = "follow_heaven"  // 顺天应人: merit > 5000
    DestinyKillBreakWolf Destiny = "kill_break_wolf" // 杀破狼: kills > 100
    DestinyHeavenSecret  Destiny = "heaven_secret"  // 天机: discovered 5+ secret realms
    DestinyLoneStar      Destiny = "lone_star"      // 孤星: no master/spouse/sect
)

// DestinyInput holds inputs for destiny determination.
type DestinyInput struct {
    Merit                int
    Karma                int
    KillCount            int
    SecretRealmDiscovered int
    HasMaster            bool
    HasSpouse            bool
    HasSect              bool
    ConsecutiveBreakthroughs int
}

// DetermineDestiny determines the entity's destiny based on their achievements and state.
// Priority: highest matching destiny is returned.
func (f *FortuneRule) DetermineDestiny(input DestinyInput) Destiny {
    // 顺天应人: merit > 5000
    if input.Merit > 5000 {
        return DestinyFollowHeaven
    }

    // 杀破狼: kills > 100
    if input.KillCount > 100 {
        return DestinyKillBreakWolf
    }

    // 天机: discovered 5+ secret realms
    if input.SecretRealmDiscovered >= 5 {
        return DestinyHeavenSecret
    }

    // 逆天改命: 3 consecutive breakthrough successes
    if input.ConsecutiveBreakthroughs >= 3 {
        return DestinyDefyHeaven
    }

    // 孤星: no master, spouse, or sect
    if !input.HasMaster && !input.HasSpouse && !input.HasSect {
        return DestinyLoneStar
    }

    return DestinyMortal
}

// DestinyEffect represents the effects of a destiny.
type DestinyEffect struct {
    BreakthroughBonus float64
    TribulationReduction float64
    CombatBonus     float64
    KarmaPenalty    float64
    OpportunityBonus float64
    CultivationBonus float64
    MentalPenalty   float64
}

// GetDestinyEffects returns the stat modifications for a given destiny.
func (f *FortuneRule) GetDestinyEffects(destiny Destiny) DestinyEffect {
    switch destiny {
    case DestinyDefyHeaven:
        return DestinyEffect{BreakthroughBonus: 0.20}
    case DestinyFollowHeaven:
        return DestinyEffect{TribulationReduction: 0.30}
    case DestinyKillBreakWolf:
        return DestinyEffect{CombatBonus: 0.30, KarmaPenalty: 0.50}
    case DestinyHeavenSecret:
        return DestinyEffect{OpportunityBonus: 0.50}
    case DestinyLoneStar:
        return DestinyEffect{CultivationBonus: 0.20, MentalPenalty: 0.20}
    default:
        return DestinyEffect{}
    }
}
