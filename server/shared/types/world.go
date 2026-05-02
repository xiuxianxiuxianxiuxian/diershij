package types

import "time"

type RegionID string

type Region struct {
    ID               RegionID   `json:"id"`
    Name             string     `json:"name"`
    ParentRegionID   *RegionID  `json:"parent_region_id,omitempty"`
    SpiritualDensity float64    `json:"spiritual_density"`
    SpiritualTier    int        `json:"spiritual_tier"`
    DangerLevel      int        `json:"danger_level"`
    Resources        []Resource `json:"resources"`
    Rules            RegionRules `json:"rules"`
    Description      string     `json:"description"`
    Lore             string     `json:"lore"`
}

type Resource struct {
    ID           string  `json:"id"`
    Name         string  `json:"name"`
    Type         string  `json:"type"`
    Rarity       int     `json:"rarity"`
    Quantity     int     `json:"quantity"`
    RespawnRate  float64 `json:"respawn_rate"`
    LastHarvested *time.Time `json:"last_harvested,omitempty"`
}

type RegionRules struct {
    IsRestricted    bool     `json:"is_restricted"`
    RestrictedBy    string   `json:"restricted_by,omitempty"`
    TaxRate         float64  `json:"tax_rate"`
    ForbiddenActions []string `json:"forbidden_actions"`
}

type WorldEvent struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Description string    `json:"description"`
    RegionID    RegionID  `json:"region_id"`
    StartTime   time.Time `json:"start_time"`
    EndTime     *time.Time `json:"end_time,omitempty"`
    Participants []EntityID `json:"participants"`
    Status      string    `json:"status"`
}

type WorldState struct {
    Epoch          int64                `json:"epoch"`
    Regions        map[RegionID]Region  `json:"regions"`
    ActiveEvents   []WorldEvent         `json:"active_events"`
    BalanceMetrics BalanceMetrics       `json:"balance_metrics"`
    LastUpdated    time.Time            `json:"last_updated"`
}

type BalanceMetrics struct {
    PowerDistribution  float64 `json:"power_distribution"`
    ResourceCirculation float64 `json:"resource_circulation"`
    SectDiversity      float64 `json:"sect_diversity"`
    KarmaDistribution  float64 `json:"karma_distribution"`
}
