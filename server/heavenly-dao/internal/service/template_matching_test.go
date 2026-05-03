package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplateLibrary(t *testing.T) {
	lib := NewTemplateLibrary()
	assert.NotNil(t, lib)
}

func TestAddTemplate(t *testing.T) {
	lib := NewTemplateLibrary()

	tmpl := &BehaviorTemplate{
		ID:       "test1",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat has_target",
		Action:   "attack",
		Priority: 10,
	}

	lib.AddTemplate(tmpl)
	assert.Equal(t, 1, len(lib.templates))

	stats := lib.GetTemplateStats()
	assert.Equal(t, 1, stats["total"])
	assert.Equal(t, 1, stats["behavior"])
}

func TestFindBestMatch_Combat(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{
		ID:       "b1",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat has_target healthy",
		Action:   "attack",
		Priority: 10,
		SuccessRate: 0.85,
	})
	lib.AddTemplate(&BehaviorTemplate{
		ID:       "b2",
		Type:     TemplateTypeBehavior,
		Pattern:  "injured no_danger",
		Action:   "rest",
		Priority: 8,
		SuccessRate: 0.90,
	})

	ctx := NewNPCContext("test_npc")
	ctx.IsInCombat = true
	ctx.HasTarget = true
	ctx.Health = 90

	best := lib.FindBestMatch(ctx)
	assert.Equal(t, "b1", best.ID)
	assert.Equal(t, "attack", best.Action)
}

func TestFindBestMatch_Injured(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{
		ID:       "b1",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat has_target healthy",
		Action:   "attack",
		Priority: 10,
		SuccessRate: 0.85,
	})
	lib.AddTemplate(&BehaviorTemplate{
		ID:       "b2",
		Type:     TemplateTypeBehavior,
		Pattern:  "injured no_danger",
		Action:   "rest",
		Priority: 8,
		SuccessRate: 0.90,
	})

	ctx := NewNPCContext("test_npc")
	ctx.Health = 30

	best := lib.FindBestMatch(ctx)
	assert.Equal(t, "b2", best.ID)
	assert.Equal(t, "rest", best.Action)
}

func TestMatchAndDecide(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{
		ID:       "d1",
		Type:     TemplateTypeDecision,
		Pattern:  "cultivate healthy",
		Action:   "cultivate",
		Parameters: map[string]string{"duration": "1h"},
		Priority: 5,
		SuccessRate: 0.75,
	})

	ctx := NewNPCContext("test_npc")
	ctx.Health = 90

	result := lib.MatchAndDecide(ctx)
	assert.NotNil(t, result)
	assert.Equal(t, "cultivate", result.Action)
	assert.Equal(t, "1h", result.Parameters["duration"])
	assert.InDelta(t, 0.75, result.Confidence, 0.01)
}

func TestMatchAndDecide_NoMatch(t *testing.T) {
	lib := NewTemplateLibrary()

	ctx := NewNPCContext("test_npc")
	result := lib.MatchAndDecide(ctx)
	assert.Nil(t, result)
}

func TestAddLLMResult(t *testing.T) {
	lib := NewTemplateLibrary()

	lib.AddLLMResult("danger low_health", "flee", map[string]string{"direction": "north"}, TemplateTypeBehavior)
	assert.Equal(t, 1, len(lib.templates))

	tmpl := lib.templates[0]
	assert.Equal(t, "flee", tmpl.Action)
	assert.Equal(t, "north", tmpl.Parameters["direction"])
	assert.Equal(t, 0.7, tmpl.SuccessRate)
	assert.Contains(t, tmpl.Tags, "llm_generated")
}

func TestUpdateTemplateSuccess(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{
		ID:          "t1",
		Type:        TemplateTypeBehavior,
		Pattern:     "test",
		Action:      "test_action",
		SuccessRate: 0.5,
	})

	// Success should increase rate
	lib.UpdateTemplateSuccess("t1", true)
	assert.Equal(t, 0.55, lib.templates[0].SuccessRate)

	// Failure should decrease rate
	lib.UpdateTemplateSuccess("t1", false)
	assert.Equal(t, 0.50, lib.templates[0].SuccessRate)
}

func TestGetTemplateStats(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{ID: "b1", Type: TemplateTypeBehavior, Pattern: "p1", Action: "a1"})
	lib.AddTemplate(&BehaviorTemplate{ID: "d1", Type: TemplateTypeDialogue, Pattern: "p2", Action: "a2"})
	lib.AddTemplate(&BehaviorTemplate{ID: "dc1", Type: TemplateTypeDecision, Pattern: "p3", Action: "a3"})

	stats := lib.GetTemplateStats()
	assert.Equal(t, 3, stats["total"])
	assert.Equal(t, 1, stats["behavior"])
	assert.Equal(t, 1, stats["dialogue"])
	assert.Equal(t, 1, stats["decision"])
}

func TestGetTemplatesByType(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{ID: "b1", Type: TemplateTypeBehavior, Pattern: "p1", Action: "a1"})
	lib.AddTemplate(&BehaviorTemplate{ID: "b2", Type: TemplateTypeBehavior, Pattern: "p2", Action: "a2"})
	lib.AddTemplate(&BehaviorTemplate{ID: "d1", Type: TemplateTypeDialogue, Pattern: "p3", Action: "a3"})

	behaviors := lib.GetTemplatesByType(TemplateTypeBehavior)
	assert.Len(t, behaviors, 2)

	dialogues := lib.GetTemplatesByType(TemplateTypeDialogue)
	assert.Len(t, dialogues, 1)
}

func TestSearchTemplates(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{ID: "b1", Type: TemplateTypeBehavior, Pattern: "combat attack", Action: "attack"})
	lib.AddTemplate(&BehaviorTemplate{ID: "b2", Type: TemplateTypeBehavior, Pattern: "cultivate meditate", Action: "cultivate"})

	results := lib.SearchTemplates("combat")
	assert.Len(t, results, 1)
	assert.Equal(t, "b1", results[0].ID)

	results = lib.SearchTemplates("cultivate")
	assert.Len(t, results, 1)
	assert.Equal(t, "b2", results[0].ID)
}

func TestDefaultBehaviorTemplates(t *testing.T) {
	templates := DefaultBehaviorTemplates()
	assert.True(t, len(templates) >= 10)

	// Check types
	for _, tmpl := range templates {
		assert.NotEmpty(t, tmpl.ID)
		assert.NotEmpty(t, tmpl.Pattern)
		assert.NotEmpty(t, tmpl.Action)
		assert.True(t, tmpl.SuccessRate > 0 && tmpl.SuccessRate <= 1.0)
	}
}

func TestDefaultDialogueTemplates(t *testing.T) {
	templates := DefaultDialogueTemplates()
	assert.True(t, len(templates) >= 4)

	for _, tmpl := range templates {
		assert.Equal(t, TemplateTypeDialogue, tmpl.Type)
	}
}

func TestTemplateUsageCount(t *testing.T) {
	lib := NewTemplateLibrary()
	lib.AddTemplate(&BehaviorTemplate{
		ID:          "t1",
		Type:        TemplateTypeBehavior,
		Pattern:     "combat has_target",
		Action:      "attack",
		SuccessRate: 0.8,
	})

	ctx := NewNPCContext("test_npc")
	ctx.IsInCombat = true
	ctx.HasTarget = true

	// First match
	lib.MatchAndDecide(ctx)
	assert.Equal(t, 1, lib.templates[0].UsageCount)

	// Second match
	lib.MatchAndDecide(ctx)
	assert.Equal(t, 2, lib.templates[0].UsageCount)
}

func TestCalculateSimilarity_PriorityBoost(t *testing.T) {
	lib := NewTemplateLibrary()

	tmpl1 := &BehaviorTemplate{
		ID:       "high_priority",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat",
		Action:   "attack",
		Priority: 10,
		SuccessRate: 0.5,
	}

	tmpl2 := &BehaviorTemplate{
		ID:       "low_priority",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat",
		Action:   "attack",
		Priority: 1,
		SuccessRate: 0.5,
	}

	ctx := NewNPCContext("test_npc")
	ctx.IsInCombat = true

	score1 := lib.calculateSimilarity(ctx, tmpl1)
	score2 := lib.calculateSimilarity(ctx, tmpl2)

	// Higher priority should get higher score
	assert.Greater(t, score1, score2)
}

func TestCalculateSimilarity_SuccessRateBoost(t *testing.T) {
	lib := NewTemplateLibrary()

	tmpl1 := &BehaviorTemplate{
		ID:       "high_success",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat",
		Action:   "attack",
		SuccessRate: 0.9,
	}

	tmpl2 := &BehaviorTemplate{
		ID:       "low_success",
		Type:     TemplateTypeBehavior,
		Pattern:  "combat",
		Action:   "attack",
		SuccessRate: 0.3,
	}

	ctx := NewNPCContext("test_npc")
	ctx.IsInCombat = true

	score1 := lib.calculateSimilarity(ctx, tmpl1)
	score2 := lib.calculateSimilarity(ctx, tmpl2)

	// Higher success rate should get higher score
	assert.Greater(t, score1, score2)
}

func TestTemplateLibraryWithDefaultTemplates(t *testing.T) {
	lib := NewTemplateLibrary()
	for _, tmpl := range DefaultBehaviorTemplates() {
		lib.AddTemplate(tmpl)
	}

	ctx := NewNPCContext("test_npc")
	ctx.IsInCombat = true
	ctx.HasTarget = true
	ctx.Health = 90

	result := lib.MatchAndDecide(ctx)
	assert.NotNil(t, result)
	assert.Equal(t, "attack", result.Action)
}
