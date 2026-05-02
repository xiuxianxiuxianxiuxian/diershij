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

type Attributes struct {
    Qi                  float64 `json:"qi"`
    MaxQi               float64 `json:"max_qi"`
    SpiritualPower      float64 `json:"spiritual_power"`
    MaxSpiritualPower   float64 `json:"max_spiritual_power"`
    DivineSense         float64 `json:"divine_sense"`
    Comprehension       int     `json:"comprehension"`
    Constitution        int     `json:"constitution"`
    Luck                int     `json:"luck"`
    CultivationProgress float64 `json:"cultivation_progress"`
    AttackPower         float64 `json:"attack_power"`
    Defense             float64 `json:"defense"`
    Speed               float64 `json:"speed"`
    MentalStability     int     `json:"mental_stability"`
    RemainingLifespan   int     `json:"remaining_lifespan"`
    MaxLifespan         int     `json:"max_lifespan"`
}

type Karma struct {
    KarmaValue  int    `json:"karma_value"`
    Merit       int    `json:"merit"`
    HeavenlyMark string `json:"heavenly_mark"`
}

type EntityStatus string

const (
    StatusNormal   EntityStatus = "normal"
    StatusCultivating EntityStatus = "cultivating"
    StatusCombat   EntityStatus = "combat"
    StatusResting  EntityStatus = "resting"
    StatusDead     EntityStatus = "dead"
)

type SpiritStones struct {
    LowGrade     int64 `json:"low_grade"`
    MediumGrade  int64 `json:"medium_grade"`
    HighGrade    int64 `json:"high_grade"`
    PremiumGrade int64 `json:"premium_grade"`
}
