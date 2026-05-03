package types

import "testing"

func TestSectInitialization(t *testing.T) {
	sect := Sect{
		ID:         "sect_001",
		Name:       "Qingyun Sect",
		FounderID:  "founder_001",
		Philosophy: "修身齐家",
		EntryRequirements: map[string]any{
			"min_realm":      "mortal",
			"min_age":        12,
			"max_age":        30,
			"alignment":      "good",
		},
		Territory: []string{"qingyun_town", "spirit_mist_mountain"},
		Rules: map[string]any{
			"tax_rate":        0.05,
			"combat_allowed":  false,
			"training_bonus":  1.2,
		},
		Alignment:   "good",
		CreatedAt:   1000000,
		MemberCount: 50,
		Prestige:    800,
		Wealth:      100000,
		FacilityScore: 75,
		CultivationResources: []string{"spirit_vein_1", "training_hall"},
	}

	if sect.ID != "sect_001" {
		t.Errorf("Sect ID mismatch")
	}
	if sect.Name != "Qingyun Sect" {
		t.Errorf("Sect Name mismatch, expected Qingyun Sect, got %s", sect.Name)
	}
	if sect.Philosophy != "修身齐家" {
		t.Errorf("Sect Philosophy mismatch")
	}
	if len(sect.Territory) != 2 {
		t.Errorf("Sect Territory length mismatch, expected 2, got %d", len(sect.Territory))
	}
	if sect.Alignment != "good" {
		t.Errorf("Sect Alignment mismatch, expected good, got %s", sect.Alignment)
	}
	if sect.MemberCount != 50 {
		t.Errorf("Sect MemberCount mismatch, expected 50, got %d", sect.MemberCount)
	}
	if sect.Prestige != 800 {
		t.Errorf("Sect Prestige mismatch, expected 800, got %d", sect.Prestige)
	}
	if sect.Wealth != 100000 {
		t.Errorf("Sect Wealth mismatch, expected 100000, got %d", sect.Wealth)
	}
}

func TestSectMemberInitialization(t *testing.T) {
	member := SectMember{
		SectID:       "sect_001",
		EntityID:     "entity_001",
		Rank:         "elder",
		Contribution: 5000.0,
		JoinedAt:     1000000,
		Privileges:   []string{"access_library", "use_training_hall", "trade_discount"},
	}

	if member.SectID != "sect_001" {
		t.Errorf("SectMember SectID mismatch")
	}
	if member.Rank != "elder" {
		t.Errorf("SectMember Rank mismatch, expected elder, got %s", member.Rank)
	}
	if member.Contribution != 5000.0 {
		t.Errorf("SectMember Contribution mismatch, expected 5000.0, got %f", member.Contribution)
	}
	if len(member.Privileges) != 3 {
		t.Errorf("SectMember Privileges length mismatch, expected 3, got %d", len(member.Privileges))
	}
}

func TestRelationshipInitialization(t *testing.T) {
	rel := Relationship{
		ID:               "rel_001",
		EntityAID:        "entity_001",
		EntityBID:        "entity_002",
		RelationshipType: "mentor_disciple",
		Strength:         80.0,
		History:          "Met at Qingyun Sect, master took disciple under wing",
		CreatedAt:        1000000,
	}

	if rel.ID != "rel_001" {
		t.Errorf("Relationship ID mismatch")
	}
	if rel.RelationshipType != "mentor_disciple" {
		t.Errorf("Relationship Type mismatch, expected mentor_disciple, got %s", rel.RelationshipType)
	}
	if rel.Strength != 80.0 {
		t.Errorf("Relationship Strength mismatch, expected 80.0, got %f", rel.Strength)
	}
}

func TestNPCPersonalityInitialization(t *testing.T) {
	personality := NPCPersonality{
		NPCID:            "npc_001",
		PersonalityType:  "ambitious",
		MoralAlignment:   "neutral",
		AmbitionLevel:    85,
		RiskTolerance:    0.7,
		SocialPreference: "extrovert",
		BackgroundStory:  "Orphan raised by wandering cultivator",
		CurrentGoal:      "reach foundation realm",
		HiddenSecrets:    []string{"has ancient bloodline"},
		LLMSystemPrompt:  "You are a young cultivator...",
		BehaviorTreeConfig: map[string]any{
			"priority": "cultivation",
			"social_frequency": 0.3,
		},
		InitialActions: []string{"cultivate", "explore", "gather"},
	}

	if personality.NPCID != "npc_001" {
		t.Errorf("NPCPersonality NPCID mismatch")
	}
	if personality.PersonalityType != "ambitious" {
		t.Errorf("NPCPersonality PersonalityType mismatch")
	}
	if personality.MoralAlignment != "neutral" {
		t.Errorf("NPCPersonality MoralAlignment mismatch")
	}
	if personality.AmbitionLevel != 85 {
		t.Errorf("NPCPersonality AmbitionLevel mismatch, expected 85, got %d", personality.AmbitionLevel)
	}
	if personality.RiskTolerance != 0.7 {
		t.Errorf("NPCPersonality RiskTolerance mismatch, expected 0.7, got %f", personality.RiskTolerance)
	}
	if len(personality.HiddenSecrets) != 1 {
		t.Errorf("NPCPersonality HiddenSecrets length mismatch, expected 1, got %d", len(personality.HiddenSecrets))
	}
	if len(personality.InitialActions) != 3 {
		t.Errorf("NPCPersonality InitialActions length mismatch, expected 3, got %d", len(personality.InitialActions))
	}
}

func TestNPCDecisionLogInitialization(t *testing.T) {
	log := NPCDecisionLog{
		ID:           "log_001",
		NPCID:        "npc_001",
		DecisionType: "cultivate",
		Context: map[string]any{
			"location":     "spirit_mist_mountain",
			"qi_level":     45.0,
			"time_of_day":  "morning",
		},
		ActionTaken: map[string]any{
			"action":   "cultivate",
			"duration": 3600,
		},
		Reasoning: "High spiritual density and morning qi surge make this optimal for cultivation",
		ModelUsed: "deepseek-chat",
		Source:    "llm",
		TokenCost: 150.0,
		Timestamp: 1000100,
	}

	if log.ID != "log_001" {
		t.Errorf("NPCDecisionLog ID mismatch")
	}
	if log.DecisionType != "cultivate" {
		t.Errorf("NPCDecisionLog DecisionType mismatch")
	}
	if log.Source != "llm" {
		t.Errorf("NPCDecisionLog Source mismatch, expected llm, got %s", log.Source)
	}
	if log.TokenCost != 150.0 {
		t.Errorf("NPCDecisionLog TokenCost mismatch, expected 150.0, got %f", log.TokenCost)
	}
}

func TestRelationshipTypes(t *testing.T) {
	relTypes := []string{
		"mentor",
		"disciple",
		"enemy",
		"ally",
		"lover",
		"sworn_sibling",
		"trade_partner",
		"rival",
		"benefactor",
	}

	for _, rt := range relTypes {
		rel := Relationship{
			RelationshipType: rt,
			Strength:         50.0,
		}
		if rel.RelationshipType != rt {
			t.Errorf("Relationship type mismatch for %s", rt)
		}
	}
}

func TestSectDefaultValues(t *testing.T) {
	sect := Sect{}

	if sect.MemberCount != 0 {
		t.Errorf("Default MemberCount should be 0, got %d", sect.MemberCount)
	}
	if sect.Prestige != 0 {
		t.Errorf("Default Prestige should be 0, got %d", sect.Prestige)
	}
	if sect.FacilityScore != 0 {
		t.Errorf("Default FacilityScore should be 0, got %d", sect.FacilityScore)
	}
}

func TestNPCPersonalityDefaultValues(t *testing.T) {
	p := NPCPersonality{}

	if p.AmbitionLevel != 0 {
		t.Errorf("Default AmbitionLevel should be 0, got %d", p.AmbitionLevel)
	}
	if p.RiskTolerance != 0 {
		t.Errorf("Default RiskTolerance should be 0, got %f", p.RiskTolerance)
	}
}
