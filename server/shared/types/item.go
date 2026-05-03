package types

import "time"

// ItemID is the unique identifier for items
type ItemID string

// ItemType represents the type of an item
type ItemType string

const (
	ItemTypeWeapon   ItemType = "weapon"
	ItemTypeArmor    ItemType = "armor"
	ItemTypePill     ItemType = "pill"
	ItemTypeMaterial ItemType = "material"
	ItemTypeTalisman ItemType = "talisman"
	ItemTypeArtifact ItemType = "artifact"
	ItemTypeTreasure ItemType = "treasure"
)

// Item represents an item in the game world
type Item struct {
	ID               ItemID                 `json:"id"`
	Name             string                 `json:"name"`
	Type             ItemType               `json:"type"`
	Rarity           int                    `json:"rarity"` // 1-5
	Description      string                 `json:"description"`
	Attributes       map[string]interface{} `json:"attributes"`
	Stackable        bool                   `json:"stackable"`
	MaxStack         int                    `json:"max_stack"`
	Usable           bool                   `json:"usable"`
	LevelRequirement int                    `json:"level_requirement"`
	RealmRequirement CultivationRealm       `json:"realm_requirement"`
	CreatedAt        time.Time              `json:"created_at"`
}

// InventoryItem represents an item in entity's inventory
type InventoryItem struct {
	ID         string   `json:"id"`
	EntityID   EntityID `json:"entity_id"`
	ItemID     ItemID   `json:"item_id"`
	Item       *Item    `json:"item,omitempty"`
	Quantity   int      `json:"quantity"`
	Equipped   bool     `json:"equipped"`
	Slot       string   `json:"slot"`
	Durability int      `json:"durability"`
	Bound      bool     `json:"bound"`
	AcquiredAt int64    `json:"acquired_at"`
}

// ItemTemplate represents a blueprint for an item in the world.
type ItemTemplate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`     // pill, artifact, talisman, material, treasure, etc.
	SubType     string `json:"sub_type"` // specific type within category
	Rank        string `json:"rank"`     // 天/地/玄/黄 x 上/中/下/极品
	Description string `json:"description"`
	BaseValue   int64  `json:"base_value"` // base value in low-grade spirit stones
	Stackable   bool   `json:"stackable"`
	MaxStack    int    `json:"max_stack"`
	Usable      bool   `json:"usable"`
	Consumable  bool   `json:"consumable"`
	Tradeable   bool   `json:"tradeable"`
	Droppable   bool   `json:"droppable"`
}

// EntityInventory represents an entity's inventory.
type EntityInventory struct {
	EntityID string          `json:"entity_id"`
	Items    []InventoryItem `json:"items"`
	Capacity int             `json:"capacity"` // max number of unique item types
}

// EntityEquipment represents equipment slots for an entity.
type EntityEquipment struct {
	EntityID   string         `json:"entity_id"`
	Weapon     *EquipmentItem `json:"weapon"`
	Armor      *EquipmentItem `json:"armor"`
	Helmet     *EquipmentItem `json:"helmet"`
	Boots      *EquipmentItem `json:"boots"`
	Necklace   *EquipmentItem `json:"necklace"`
	Ring1      *EquipmentItem `json:"ring_1"`
	Ring2      *EquipmentItem `json:"ring_2"`
	InnerArmor *EquipmentItem `json:"inner_armor"`
	Waist      *EquipmentItem `json:"waist"`
	Bracelet   *EquipmentItem `json:"bracelet"`
	TotalSlots int            `json:"total_slots"`
}

// EquipmentItem represents a single equipped item.
type EquipmentItem struct {
	TemplateID    string             `json:"template_id"`
	InstanceID    string             `json:"instance_id"`
	Name          string             `json:"name"`
	Slot          string             `json:"slot"`
	Rank          string             `json:"rank"`
	Quality       int                `json:"quality"`    // 1-100
	Level         int                `json:"level"`      // equipment level
	Durability    int                `json:"durability"` // current durability
	MaxDurability int                `json:"max_durability"`
	Stats         map[string]float64 `json:"stats"` // attribute bonuses
	Enchantments  []Enchantment      `json:"enchantments"`
	Soulbound     bool               `json:"soulbound"` // soulbound status
}

// Enchantment represents an enchantment on an equipment item.
type Enchantment struct {
	Name  string  `json:"name"`
	Stat  string  `json:"stat"`
	Value float64 `json:"value"`
	Tier  int     `json:"tier"` // enchantment tier (1-10)
}

// Pill represents a pill item.
type Pill struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Type          string  `json:"type"` // cultivation, healing, breakthrough, combat, etc.
	Rank          string  `json:"rank"`
	Quality       string  `json:"quality"`        // 下品/中品/上品/极品
	Effect        string  `json:"effect"`         // effect description
	EffectValue   float64 `json:"effect_value"`   // effect magnitude
	Duration      int     `json:"duration"`       // effect duration in seconds
	SuccessRate   float64 `json:"success_rate"`   // pill consumption success rate
	FailureEffect string  `json:"failure_effect"` // what happens on failure (toxicity, explosion)
	Toxicity      int     `json:"toxicity"`       // pill toxicity accumulation
	Cooldown      int     `json:"cooldown"`       // cooldown before another pill can be taken
}

// Artifact represents a magical artifact/treasure.
type Artifact struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Type           string             `json:"type"` // flying_sword, mirror, bell, pagoda, banner, etc.
	Rank           string             `json:"rank"`
	Grade          string             `json:"grade"` // 凡器/地器/天器/古宝
	AttackBonus    map[string]float64 `json:"attack_bonus"`
	DefenseBonus   map[string]float64 `json:"defense_bonus"`
	SpecialAbility string             `json:"special_ability"`
	Energy         float64            `json:"energy"` // artifact energy (used for abilities)
	MaxEnergy      float64            `json:"max_energy"`
	Level          int                `json:"level"`
	Experience     float64            `json:"experience"` // artifact EXP for leveling
	Refined        bool               `json:"refined"`    // whether refined by current owner
}

// Talisman represents a talisman item.
type Talisman struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // attack, defense, utility, movement, etc.
	Rank        string  `json:"rank"`
	Effect      string  `json:"effect"`
	EffectValue float64 `json:"effect_value"`
	Duration    int     `json:"duration"`
	Cooldown    int     `json:"cooldown"`
	Charges     int     `json:"charges"` // remaining uses
	MaxCharges  int     `json:"max_charges"`
}

// Recipe represents a crafting recipe.
type Recipe struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	Type               string           `json:"type"` // pill, artifact, talisman, array
	ResultItemID       string           `json:"result_item_id"`
	ResultQuantity     int              `json:"result_quantity"`
	RequiredLevel      int              `json:"required_level"`
	RequiredSkill      string           `json:"required_skill"` // alchemy, artificing, etc.
	RequiredSkillLevel int              `json:"required_skill_level"`
	Materials          []RecipeMaterial `json:"materials"`
	BaseSuccessRate    float64          `json:"base_success_rate"`
	CreatorID          string           `json:"creator_id"` // who discovered/created this recipe
	IsSecret           bool             `json:"is_secret"`  // whether recipe is a secret
}

// RecipeMaterial represents a material needed for a recipe.
type RecipeMaterial struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Quality  string `json:"quality"` // minimum quality required
}

// Material represents a raw material for crafting.
type Material struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Type       string             `json:"type"` // herb, ore, beast_part, essence, etc.
	Rank       string             `json:"rank"`
	Purity     int                `json:"purity"`     // 1-100
	Attributes map[string]float64 `json:"attributes"` // material-specific attributes
	Source     string             `json:"source"`     // where it can be found
	DropRate   float64            `json:"drop_rate"`  // base drop rate from source
}
