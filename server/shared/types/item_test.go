package types

import "testing"

func TestItemTemplateInitialization(t *testing.T) {
	template := ItemTemplate{
		ID:          "item_001",
		Name:        "Spirit Gathering Pill",
		Type:        "pill",
		SubType:     "cultivation",
		Rank:        "mystic_upper",
		Description: "A pill that gathers spiritual energy",
		BaseValue:   500,
		Stackable:   true,
		MaxStack:    99,
		Usable:      true,
		Consumable:  true,
		Tradeable:   true,
		Droppable:   true,
	}

	if template.ID != "item_001" {
		t.Errorf("ItemTemplate ID mismatch, expected item_001, got %s", template.ID)
	}
	if template.Name != "Spirit Gathering Pill" {
		t.Errorf("ItemTemplate Name mismatch, expected Spirit Gathering Pill, got %s", template.Name)
	}
	if template.Type != "pill" {
		t.Errorf("ItemTemplate Type mismatch, expected pill, got %s", template.Type)
	}
	if template.BaseValue != 500 {
		t.Errorf("ItemTemplate BaseValue mismatch, expected 500, got %d", template.BaseValue)
	}
	if template.Stackable != true {
		t.Error("ItemTemplate Stackable should be true")
	}
	if template.MaxStack != 99 {
		t.Errorf("ItemTemplate MaxStack mismatch, expected 99, got %d", template.MaxStack)
	}
}

func TestInventoryItemInitialization(t *testing.T) {
	item := InventoryItem{
		ID:       "inst_001",
		ItemID:   "item_001",
		Quantity: 10,
		Bound:    false,
		Item: &Item{
			Name: "Spirit Gathering Pill",
			ID:   "item_001",
			Type: "pill",
		},
	}

	if item.ItemID != "item_001" {
		t.Errorf("InventoryItem ItemID mismatch")
	}
	if item.Quantity != 10 {
		t.Errorf("InventoryItem Quantity mismatch, expected 10, got %d", item.Quantity)
	}
	if item.Item.Name != "Spirit Gathering Pill" {
		t.Errorf("InventoryItem Name mismatch, expected Spirit Gathering Pill, got %s", item.Item.Name)
	}
	if item.Bound != false {
		t.Error("InventoryItem Bound should be false")
	}
}

func TestEntityEquipmentInitialization(t *testing.T) {
	equip := EntityEquipment{
		EntityID: "entity_001",
		Weapon: &EquipmentItem{
			TemplateID:    "weapon_001",
			InstanceID:    "w_inst_001",
			Name:          "Flame Sword",
			Slot:          "weapon",
			Rank:          "mystic_upper",
			Quality:       90,
			Level:         5,
			Durability:    100,
			MaxDurability: 100,
			Stats:         map[string]float64{"attack_power": 50.0, "fire_damage": 20.0},
			Soulbound:     true,
		},
		Armor: &EquipmentItem{
			TemplateID: "armor_001",
			Name:       "Iron Chestplate",
			Slot:       "armor",
			Rank:       "mystic_lower",
			Quality:    70,
			Level:      3,
			Stats:      map[string]float64{"defense": 30.0},
		},
		TotalSlots: 10,
	}

	if equip.EntityID != "entity_001" {
		t.Errorf("EntityEquipment EntityID mismatch")
	}
	if equip.Weapon == nil {
		t.Error("Weapon should not be nil")
	} else {
		if equip.Weapon.Name != "Flame Sword" {
			t.Errorf("Weapon Name mismatch, expected Flame Sword, got %s", equip.Weapon.Name)
		}
		if equip.Weapon.Stats["attack_power"] != 50.0 {
			t.Errorf("Weapon attack_power stat mismatch, expected 50.0, got %f", equip.Weapon.Stats["attack_power"])
		}
		if equip.Weapon.Soulbound != true {
			t.Error("Weapon should be soulbound")
		}
	}
	if equip.Armor == nil {
		t.Error("Armor should not be nil")
	}
	if equip.TotalSlots != 10 {
		t.Errorf("TotalSlots mismatch, expected 10, got %d", equip.TotalSlots)
	}
}

func TestEquipmentItemWithEnchantments(t *testing.T) {
	item := EquipmentItem{
		TemplateID:    "weapon_001",
		Name:          "Thunder Blade",
		Slot:          "weapon",
		Quality:       95,
		Level:         8,
		Durability:    80,
		MaxDurability: 100,
		Stats:         map[string]float64{"attack_power": 100.0},
		Enchantments: []Enchantment{
			{Name: "Thunder Strike", Stat: "thunder_damage", Value: 50.0, Tier: 3},
			{Name: "Swift", Stat: "speed", Value: 10.0, Tier: 2},
		},
	}

	if len(item.Enchantments) != 2 {
		t.Errorf("Enchantments length mismatch, expected 2, got %d", len(item.Enchantments))
	}
	if item.Enchantments[0].Name != "Thunder Strike" {
		t.Errorf("First enchantment Name mismatch")
	}
	if item.Enchantments[0].Value != 50.0 {
		t.Errorf("First enchantment Value mismatch, expected 50.0, got %f", item.Enchantments[0].Value)
	}
	if item.Enchantments[1].Tier != 2 {
		t.Errorf("Second enchantment Tier mismatch, expected 2, got %d", item.Enchantments[1].Tier)
	}
}

func TestPillInitialization(t *testing.T) {
	pill := Pill{
		ID:             "pill_001",
		Name:           "Foundation Pill",
		Type:           "breakthrough",
		Rank:           "heaven_lower",
		Quality:        "上品",
		Effect:         "increase_breakthrough_chance",
		EffectValue:    0.15,
		Duration:       0,
		SuccessRate:    0.9,
		FailureEffect:  "toxicity_increase",
		Toxicity:       5,
		Cooldown:       3600,
	}

	if pill.ID != "pill_001" {
		t.Errorf("Pill ID mismatch")
	}
	if pill.Type != "breakthrough" {
		t.Errorf("Pill Type mismatch, expected breakthrough, got %s", pill.Type)
	}
	if pill.EffectValue != 0.15 {
		t.Errorf("Pill EffectValue mismatch, expected 0.15, got %f", pill.EffectValue)
	}
	if pill.Toxicity != 5 {
		t.Errorf("Pill Toxicity mismatch, expected 5, got %d", pill.Toxicity)
	}
	if pill.Cooldown != 3600 {
		t.Errorf("Pill Cooldown mismatch, expected 3600, got %d", pill.Cooldown)
	}
}

func TestArtifactInitialization(t *testing.T) {
	artifact := Artifact{
		ID:    "artifact_001",
		Name:  "Heavenly Mirror",
		Type:  "mirror",
		Rank:  "heaven_upper",
		Grade: "天器",
		AttackBonus: map[string]float64{
			"divine_sense": 30.0,
		},
		DefenseBonus: map[string]float64{
			"mental_resist": 20.0,
		},
		SpecialAbility: "illusion_reflection",
		Energy:         100.0,
		MaxEnergy:      500.0,
		Level:          5,
		Experience:     250.0,
		Refined:        true,
	}

	if artifact.Name != "Heavenly Mirror" {
		t.Errorf("Artifact Name mismatch")
	}
	if artifact.Grade != "天器" {
		t.Errorf("Artifact Grade mismatch, expected 天器, got %s", artifact.Grade)
	}
	if artifact.AttackBonus["divine_sense"] != 30.0 {
		t.Errorf("Artifact divine_sense bonus mismatch")
	}
	if artifact.Energy != 100.0 {
		t.Errorf("Artifact Energy mismatch, expected 100.0, got %f", artifact.Energy)
	}
	if artifact.Refined != true {
		t.Error("Artifact should be refined")
	}
}

func TestTalismanInitialization(t *testing.T) {
	talisman := Talisman{
		ID:          "talisman_001",
		Name:        "Thunder Talisman",
		Type:        "attack",
		Rank:        "mystic_middle",
		Effect:      "thunder_damage",
		EffectValue: 200.0,
		Duration:    0,
		Cooldown:    0,
		Charges:     3,
		MaxCharges:  3,
	}

	if talisman.Name != "Thunder Talisman" {
		t.Errorf("Talisman Name mismatch")
	}
	if talisman.Charges != 3 {
		t.Errorf("Talisman Charges mismatch, expected 3, got %d", talisman.Charges)
	}
}

func TestRecipeInitialization(t *testing.T) {
	recipe := Recipe{
		ID:                 "recipe_001",
		Name:               "Foundation Pill Recipe",
		Type:               "pill",
		ResultItemID:       "pill_001",
		ResultQuantity:     3,
		RequiredLevel:      5,
		RequiredSkill:      "alchemy",
		RequiredSkillLevel: 4,
		Materials: []RecipeMaterial{
			{ItemID: "herb_001", Name: "Spirit Grass", Quantity: 10, Quality: "中品"},
			{ItemID: "herb_002", Name: "Fire Flower", Quantity: 5, Quality: "上品"},
			{ItemID: "water_001", Name: "Spirit Water", Quantity: 1, Quality: "下品"},
		},
		BaseSuccessRate: 0.6,
		CreatorID:       "creator_001",
		IsSecret:        false,
	}

	if recipe.Name != "Foundation Pill Recipe" {
		t.Errorf("Recipe Name mismatch")
	}
	if len(recipe.Materials) != 3 {
		t.Errorf("Recipe Materials length mismatch, expected 3, got %d", len(recipe.Materials))
	}
	if recipe.Materials[0].Quantity != 10 {
		t.Errorf("First material Quantity mismatch")
	}
	if recipe.BaseSuccessRate != 0.6 {
		t.Errorf("Recipe BaseSuccessRate mismatch, expected 0.6, got %f", recipe.BaseSuccessRate)
	}
}

func TestMaterialInitialization(t *testing.T) {
	material := Material{
		ID:     "herb_001",
		Name:   "Thousand Year Ginseng",
		Type:   "herb",
		Rank:   "earth_upper",
		Purity: 85,
		Attributes: map[string]float64{
			"qi_recovery": 100.0,
			"lifespan_bonus": 5.0,
		},
		Source:   "spirit_mist_mountain",
		DropRate: 0.05,
	}

	if material.Name != "Thousand Year Ginseng" {
		t.Errorf("Material Name mismatch")
	}
	if material.Purity != 85 {
		t.Errorf("Material Purity mismatch, expected 85, got %d", material.Purity)
	}
	if material.Attributes["qi_recovery"] != 100.0 {
		t.Errorf("Material qi_recovery attribute mismatch")
	}
	if material.DropRate != 0.05 {
		t.Errorf("Material DropRate mismatch, expected 0.05, got %f", material.DropRate)
	}
}

func TestEntityInventoryWithItems(t *testing.T) {
	inventory := EntityInventory{
		EntityID: "entity_001",
		Items: []InventoryItem{
			{ID: "inst_001", ItemID: "item_001", Item: &Item{Name: "Pill A"}, Quantity: 5, Slot: "0"},
			{ID: "inst_002", ItemID: "item_002", Item: &Item{Name: "Pill B"}, Quantity: 10, Slot: "1"},
			{ID: "inst_003", ItemID: "item_003", Item: &Item{Name: "Sword"}, Quantity: 1, Slot: "2", Bound: true},
		},
		Capacity: 50,
	}

	if len(inventory.Items) != 3 {
		t.Errorf("Inventory Items length mismatch, expected 3, got %d", len(inventory.Items))
	}
	if inventory.Capacity != 50 {
		t.Errorf("Inventory Capacity mismatch, expected 50, got %d", inventory.Capacity)
	}
	if !inventory.Items[2].Bound {
		t.Error("Third item should be bound")
	}
}

func TestItemTemplateDefaultValues(t *testing.T) {
	template := ItemTemplate{}

	if template.Stackable != false {
		t.Error("Default Stackable should be false")
	}
	if template.Usable != false {
		t.Error("Default Usable should be false")
	}
	if template.Tradeable != false {
		t.Error("Default Tradeable should be false")
	}
}
