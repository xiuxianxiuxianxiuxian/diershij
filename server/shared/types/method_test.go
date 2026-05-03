package types

import "testing"

func TestCultivationMethodInitialization(t *testing.T) {
	method := CultivationMethod{
		ID:              "method_001",
		Name:            "Nine Sun Scripture",
		CreatorID:       "creator_001",
		OriginSect:      "qingyun_sect",
		Rank:            "heaven_upper",
		Category:        "main",
		ElementAffinity: "fire",
		Description:     "A supreme fire cultivation method",
		Version:         1,

		CultivationSpeedMult:  1.5,
		SpiritualPowerCapMult: 1.3,
		QiCapMult:             1.4,
		DivineSenseCapMult:    1.1,
		LifespanBonus:         50,
		RecoverySpeedMult:     1.2,

		AttackBonuses:  map[string]float64{"fire_damage": 1.5},
		DefenseBonuses: map[string]float64{"fire_resist": 0.3},
		UtilityBonuses: map[string]float64{},

		PassiveEffects: []string{"fire_aura", "heat_resistance"},
		ActiveSkills: []Skill{
			{
				ID:          "skill_001",
				Name:        "Sun Flare",
				Category:    "active",
				Element:     "fire",
				DamageMult:  2.0,
				Cooldown:    10,
				ManaCost:    50.0,
			},
		},
		UltimateSkill: &Skill{
			ID:          "ult_001",
			Name:        "Nine Suns Annihilation",
			Category:    "ultimate",
			Element:     "fire",
			DamageMult:  10.0,
			Cooldown:    300,
			ManaCost:    500.0,
		},

		LawAffinities:         []string{"fire", "light"},
		LawComprehensionBonus: 1.3,
		DaoCompatibility:      []string{"sword_dao"},

		RequiredRoots:      []string{"fire"},
		RequiredPhysique:   []string{},
		RealmRequirement:   "foundation",
		AlignmentRestriction: "neutral",
		KarmaThreshold:     5000,
		GenderRestriction:  "none",

		ParentMethodID:  "",
		EvolutionPath:   []string{"method_002"},
		TransmissionMode: "jade_slip",
		CanModify:       true,
		Complexity:      8,

		PowerScore: 90,
		Potential:  95,
		Popularity: 3,
	}

	if method.ID != "method_001" {
		t.Errorf("Method ID mismatch, expected method_001, got %s", method.ID)
	}
	if method.Name != "Nine Sun Scripture" {
		t.Errorf("Method Name mismatch, expected Nine Sun Scripture, got %s", method.Name)
	}
	if method.Rank != "heaven_upper" {
		t.Errorf("Method Rank mismatch, expected heaven_upper, got %s", method.Rank)
	}
	if method.Category != "main" {
		t.Errorf("Method Category mismatch, expected main, got %s", method.Category)
	}
	if method.CultivationSpeedMult != 1.5 {
		t.Errorf("CultivationSpeedMult mismatch, expected 1.5, got %f", method.CultivationSpeedMult)
	}
	if method.SpiritualPowerCapMult != 1.3 {
		t.Errorf("SpiritualPowerCapMult mismatch, expected 1.3, got %f", method.SpiritualPowerCapMult)
	}
	if method.LifespanBonus != 50 {
		t.Errorf("LifespanBonus mismatch, expected 50, got %d", method.LifespanBonus)
	}
	if len(method.PassiveEffects) != 2 {
		t.Errorf("PassiveEffects length mismatch, expected 2, got %d", len(method.PassiveEffects))
	}
	if len(method.ActiveSkills) != 1 {
		t.Errorf("ActiveSkills length mismatch, expected 1, got %d", len(method.ActiveSkills))
	}
	if method.ActiveSkills[0].DamageMult != 2.0 {
		t.Errorf("First skill DamageMult mismatch, expected 2.0, got %f", method.ActiveSkills[0].DamageMult)
	}
	if method.UltimateSkill == nil {
		t.Error("UltimateSkill should not be nil")
	} else if method.UltimateSkill.Name != "Nine Suns Annihilation" {
		t.Errorf("UltimateSkill Name mismatch, expected Nine Suns Annihilation, got %s", method.UltimateSkill.Name)
	}
	if len(method.LawAffinities) != 2 {
		t.Errorf("LawAffinities length mismatch, expected 2, got %d", len(method.LawAffinities))
	}
	if method.LawComprehensionBonus != 1.3 {
		t.Errorf("LawComprehensionBonus mismatch, expected 1.3, got %f", method.LawComprehensionBonus)
	}
	if len(method.RequiredRoots) != 1 {
		t.Errorf("RequiredRoots length mismatch, expected 1, got %d", len(method.RequiredRoots))
	}
	if method.Complexity != 8 {
		t.Errorf("Complexity mismatch, expected 8, got %d", method.Complexity)
	}
	if method.PowerScore != 90 {
		t.Errorf("PowerScore mismatch, expected 90, got %d", method.PowerScore)
	}
	if method.Potential != 95 {
		t.Errorf("Potential mismatch, expected 95, got %d", method.Potential)
	}
	if method.Popularity != 3 {
		t.Errorf("Popularity mismatch, expected 3, got %d", method.Popularity)
	}
}

func TestSkillInitialization(t *testing.T) {
	skill := Skill{
		ID:          "skill_001",
		Name:        "Thunder Strike",
		Description: "A powerful lightning attack",
		Category:    "active",
		Element:     "thunder",
		DamageMult:  3.5,
		Cooldown:    15,
		ManaCost:    80.0,
		Range:       100.0,
		Duration:    0,
		Effects: []SkillEffect{
			{Type: "damage", Value: 300.0, Target: "enemy", Duration: 0},
			{Type: "stun", Value: 0, Target: "enemy", Duration: 3},
		},
	}

	if skill.ID != "skill_001" {
		t.Errorf("Skill ID mismatch, expected skill_001, got %s", skill.ID)
	}
	if skill.DamageMult != 3.5 {
		t.Errorf("Skill DamageMult mismatch, expected 3.5, got %f", skill.DamageMult)
	}
	if skill.ManaCost != 80.0 {
		t.Errorf("Skill ManaCost mismatch, expected 80.0, got %f", skill.ManaCost)
	}
	if len(skill.Effects) != 2 {
		t.Errorf("Skill Effects length mismatch, expected 2, got %d", len(skill.Effects))
	}
	if skill.Effects[0].Type != "damage" {
		t.Errorf("First effect Type mismatch, expected damage, got %s", skill.Effects[0].Type)
	}
	if skill.Effects[1].Type != "stun" {
		t.Errorf("Second effect Type mismatch, expected stun, got %s", skill.Effects[1].Type)
	}
	if skill.Effects[1].Duration != 3 {
		t.Errorf("Stun duration mismatch, expected 3, got %d", skill.Effects[1].Duration)
	}
}

func TestSkillEffectInitialization(t *testing.T) {
	effect := SkillEffect{
		Type:      "heal",
		Value:     500.0,
		Target:    "self",
		Duration:  0,
		Condition: "when_qi_below_30_percent",
	}

	if effect.Type != "heal" {
		t.Errorf("Effect Type mismatch, expected heal, got %s", effect.Type)
	}
	if effect.Value != 500.0 {
		t.Errorf("Effect Value mismatch, expected 500.0, got %f", effect.Value)
	}
	if effect.Target != "self" {
		t.Errorf("Effect Target mismatch, expected self, got %s", effect.Target)
	}
}

func TestEntityMethodInitialization(t *testing.T) {
	em := EntityMethod{
		MethodID:      "method_001",
		EntityID:      "entity_001",
		MasteryLevel:  45.5,
		IsMainMethod:  true,
		LearnedAt:     1000000,
		LastPracticed: 1000100,
		BacklashRisk:  0.15,
		Modified:      false,
	}

	if em.MethodID != "method_001" {
		t.Errorf("EntityMethod MethodID mismatch, expected method_001, got %s", em.MethodID)
	}
	if em.EntityID != "entity_001" {
		t.Errorf("EntityMethod EntityID mismatch, expected entity_001, got %s", em.EntityID)
	}
	if em.MasteryLevel != 45.5 {
		t.Errorf("EntityMethod MasteryLevel mismatch, expected 45.5, got %f", em.MasteryLevel)
	}
	if em.IsMainMethod != true {
		t.Error("EntityMethod IsMainMethod should be true")
	}
	if em.BacklashRisk != 0.15 {
		t.Errorf("EntityMethod BacklashRisk mismatch, expected 0.15, got %f", em.BacklashRisk)
	}
}

func TestCultivationMethodDefaultValues(t *testing.T) {
	method := CultivationMethod{}

	if method.CultivationSpeedMult != 0 {
		t.Errorf("Default CultivationSpeedMult should be 0, got %f", method.CultivationSpeedMult)
	}
	if method.Complexity != 0 {
		t.Errorf("Default Complexity should be 0, got %d", method.Complexity)
	}
	if method.CanModify != false {
		t.Error("Default CanModify should be false")
	}
	if method.AttackBonuses != nil {
		t.Error("Default AttackBonuses should be nil")
	}
	if method.PassiveEffects != nil {
		t.Error("Default PassiveEffects should be nil")
	}
	if method.ActiveSkills != nil {
		t.Error("Default ActiveSkills should be nil")
	}
}

func TestMethodAttackBonuses(t *testing.T) {
	method := CultivationMethod{
		AttackBonuses: map[string]float64{
			"fire_damage":  1.5,
			"penetration":  0.3,
			"crit_rate":    0.1,
		},
		DefenseBonuses: map[string]float64{
			"magic_resist":     0.3,
			"damage_reduction": 0.1,
		},
		UtilityBonuses: map[string]float64{
			"loot_rate":        0.1,
			"alchemy_success":  0.2,
		},
	}

	if len(method.AttackBonuses) != 3 {
		t.Errorf("AttackBonuses length mismatch, expected 3, got %d", len(method.AttackBonuses))
	}
	if method.AttackBonuses["fire_damage"] != 1.5 {
		t.Errorf("fire_damage bonus mismatch, expected 1.5, got %f", method.AttackBonuses["fire_damage"])
	}
	if len(method.DefenseBonuses) != 2 {
		t.Errorf("DefenseBonuses length mismatch, expected 2, got %d", len(method.DefenseBonuses))
	}
	if len(method.UtilityBonuses) != 2 {
		t.Errorf("UtilityBonuses length mismatch, expected 2, got %d", len(method.UtilityBonuses))
	}
}

func TestMethodMultipleSkills(t *testing.T) {
	method := CultivationMethod{
		ActiveSkills: []Skill{
			{ID: "s1", Name: "Fire Ball", Category: "active", DamageMult: 1.5},
			{ID: "s2", Name: "Fire Shield", Category: "active", DamageMult: 0},
			{ID: "s3", Name: "Fire Dash", Category: "active", DamageMult: 0.5},
		},
	}

	if len(method.ActiveSkills) != 3 {
		t.Errorf("Expected 3 active skills, got %d", len(method.ActiveSkills))
	}

	for _, skill := range method.ActiveSkills {
		if skill.Name == "" {
			t.Error("Skill name should not be empty")
		}
	}
}
