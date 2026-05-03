package service

import (
	"math"
	"testing"
)

func TestDefaultKarmaConfig(t *testing.T) {
	cfg := DefaultKarmaConfig()

	if cfg.KarmaDecayRate != 0.01 {
		t.Errorf("Default decay rate mismatch, expected 0.01, got %f", cfg.KarmaDecayRate)
	}
	if cfg.KarmaCap != 10000 {
		t.Errorf("Default karma cap mismatch, expected 10000, got %d", cfg.KarmaCap)
	}
	if len(cfg.ActionKarmaMap) == 0 {
		t.Error("ActionKarmaMap should not be empty")
	}
	if cfg.ActionKarmaMap["kill_innocent"] != 500 {
		t.Errorf("kill_innocent karma mismatch, expected 500, got %d", cfg.ActionKarmaMap["kill_innocent"])
	}
	if cfg.ActionKarmaMap["save_life"] != -100 {
		t.Errorf("save_life karma mismatch, expected -100, got %d", cfg.ActionKarmaMap["save_life"])
	}
}

func TestNewKarmaRule(t *testing.T) {
	rule := NewKarmaRule(nil)
	if rule == nil {
		t.Error("NewKarmaRule(nil) should not return nil")
	}
	if rule.cfg.KarmaDecayRate != 0.01 {
		t.Error("Default config should have decay rate 0.01")
	}

	customCfg := &KarmaRuleConfig{
		KarmaDecayRate: 0.05,
		KarmaCap:       5000,
	}
	rule2 := NewKarmaRule(customCfg)
	if rule2.cfg.KarmaDecayRate != 0.05 {
		t.Errorf("Custom config decay rate mismatch, expected 0.05, got %f", rule2.cfg.KarmaDecayRate)
	}
	if rule2.cfg.KarmaCap != 5000 {
		t.Errorf("Custom config karma cap mismatch, expected 5000, got %d", rule2.cfg.KarmaCap)
	}
}

func TestCalculateKarmaChange_KillInnocent(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctx := &KarmaContext{
		ActionType: "kill_innocent",
		ActorKarma: 50,
	}

	result := rule.CalculateKarmaChange(ctx)

	if result.BaseKarma != 500 {
		t.Errorf("Base karma mismatch, expected 500, got %d", result.BaseKarma)
	}
	if result.KarmaChange <= 0 {
		t.Errorf("Kill innocent should increase karma (positive), got %d", result.KarmaChange)
	}
	if result.NewHeavenlyMark == "" {
		t.Error("NewHeavenlyMark should not be empty")
	}
}

func TestCalculateKarmaChange_SaveLife(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctx := &KarmaContext{
		ActionType: "save_life",
		ActorKarma: 100,
	}

	result := rule.CalculateKarmaChange(ctx)

	if result.BaseKarma != -100 {
		t.Errorf("Base karma mismatch, expected -100, got %d", result.BaseKarma)
	}
	if result.KarmaChange >= 0 {
		t.Errorf("Save life should decrease karma (negative/merit), got %d", result.KarmaChange)
	}
}

func TestCalculateKarmaChange_ContextMultiplier_DiminishingReturns(t *testing.T) {
	rule := NewKarmaRule(nil)

	// Low karma actor
	ctxLow := &KarmaContext{
		ActionType: "kill_innocent",
		ActorKarma: 0,
	}
	resultLow := rule.CalculateKarmaChange(ctxLow)

	// High karma actor (near cap)
	ctxHigh := &KarmaContext{
		ActionType: "kill_innocent",
		ActorKarma: 9000,
	}
	resultHigh := rule.CalculateKarmaChange(ctxHigh)

	// High karma actor should have smaller karma change due to diminishing returns
	if resultHigh.KarmaChange >= resultLow.KarmaChange {
		t.Errorf("High karma actor should have smaller karma change due to diminishing returns: low=%d, high=%d",
			resultLow.KarmaChange, resultHigh.KarmaChange)
	}
}

func TestCalculateKarmaChange_TargetRealmModifier(t *testing.T) {
	rule := NewKarmaRule(nil)

	// Killing someone of equal realm
	ctxEqual := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		ActorRealm:   2,
		TargetRealm:  2,
	}
	resultEqual := rule.CalculateKarmaChange(ctxEqual)

	// Killing someone of higher realm
	ctxHigher := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		ActorRealm:   2,
		TargetRealm:  5,
	}
	resultHigher := rule.CalculateKarmaChange(ctxHigher)

	if resultHigher.KarmaChange <= resultEqual.KarmaChange {
		t.Errorf("Killing higher realm target should produce more karma: equal=%d, higher=%d",
			resultEqual.KarmaChange, resultHigher.KarmaChange)
	}
}

func TestCalculateKarmaChange_SelfDefense(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctxNormal := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		ActorRealm:   2,
		TargetRealm:  2,
	}
	resultNormal := rule.CalculateKarmaChange(ctxNormal)

	ctxSelfDefense := &KarmaContext{
		ActionType:    "kill_cultivator",
		ActorKarma:    0,
		ActorRealm:    2,
		TargetRealm:   2,
		IsSelfDefense: true,
	}
	resultSelfDefense := rule.CalculateKarmaChange(ctxSelfDefense)

	if resultSelfDefense.KarmaChange >= resultNormal.KarmaChange {
		t.Errorf("Self-defense should reduce karma: normal=%d, self_defense=%d",
			resultNormal.KarmaChange, resultSelfDefense.KarmaChange)
	}
}

func TestCalculateKarmaChange_Relationship_Mentor(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctxNormal := &KarmaContext{
		ActionType: "kill_cultivator",
		ActorKarma: 0,
	}
	resultNormal := rule.CalculateKarmaChange(ctxNormal)

	ctxMentor := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		Relationship: "mentor",
	}
	resultMentor := rule.CalculateKarmaChange(ctxMentor)

	// Mentor relationship should double karma
	expectedKarma := resultNormal.KarmaChange * 2
	if resultMentor.KarmaChange != expectedKarma {
		t.Errorf("Mentor betrayal should double karma: normal=%d, expected=%d, got=%d",
			resultNormal.KarmaChange, expectedKarma, resultMentor.KarmaChange)
	}
}

func TestCalculateKarmaChange_Relationship_Enemy(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctxNormal := &KarmaContext{
		ActionType: "kill_cultivator",
		ActorKarma: 0,
	}
	resultNormal := rule.CalculateKarmaChange(ctxNormal)

	ctxEnemy := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		Relationship: "enemy",
	}
	resultEnemy := rule.CalculateKarmaChange(ctxEnemy)

	// Enemy relationship should halve karma
	expectedKarma := resultNormal.KarmaChange / 2
	if resultEnemy.KarmaChange != expectedKarma {
		t.Errorf("Killing enemy should halve karma: normal=%d, expected=%d, got=%d",
			resultNormal.KarmaChange, expectedKarma, resultEnemy.KarmaChange)
	}
}

func TestCalculateKarmaChange_Relationship_Disciple(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctxNormal := &KarmaContext{
		ActionType: "kill_cultivator",
		ActorKarma: 0,
	}
	resultNormal := rule.CalculateKarmaChange(ctxNormal)

	ctxDisciple := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		Relationship: "disciple",
	}
	resultDisciple := rule.CalculateKarmaChange(ctxDisciple)

	// Disciple relationship should triple karma
	expectedKarma := resultNormal.KarmaChange * 3
	if resultDisciple.KarmaChange != expectedKarma {
		t.Errorf("Killing disciple should triple karma: normal=%d, expected=%d, got=%d",
			resultNormal.KarmaChange, expectedKarma, resultDisciple.KarmaChange)
	}
}

func TestCalculateKarmaChange_Relationship_Lover(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctxNormal := &KarmaContext{
		ActionType: "kill_cultivator",
		ActorKarma: 0,
	}
	resultNormal := rule.CalculateKarmaChange(ctxNormal)

	ctxLover := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		Relationship: "lover",
	}
	resultLover := rule.CalculateKarmaChange(ctxLover)

	// Lover relationship should triple karma
	expectedKarma := resultNormal.KarmaChange * 3
	if resultLover.KarmaChange != expectedKarma {
		t.Errorf("Killing lover should triple karma: normal=%d, expected=%d, got=%d",
			resultNormal.KarmaChange, expectedKarma, resultLover.KarmaChange)
	}
}

func TestCalculateKarmaChange_Relationship_SectMember(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctxNormal := &KarmaContext{
		ActionType: "kill_cultivator",
		ActorKarma: 0,
	}
	resultNormal := rule.CalculateKarmaChange(ctxNormal)

	ctxSect := &KarmaContext{
		ActionType:   "kill_cultivator",
		ActorKarma:   0,
		Relationship: "sect_member",
	}
	resultSect := rule.CalculateKarmaChange(ctxSect)

	// Sect member relationship should multiply karma by 1.5
	expectedKarma := int(float64(resultNormal.KarmaChange) * 1.5)
	if resultSect.KarmaChange != expectedKarma {
		t.Errorf("Killing sect member should multiply karma by 1.5: normal=%d, expected=%d, got=%d",
			resultNormal.KarmaChange, expectedKarma, resultSect.KarmaChange)
	}
}

func TestCalculateKarmaChange_KarmaCap(t *testing.T) {
	cfg := &KarmaRuleConfig{
		KarmaDecayRate: 0.01,
		KarmaCap:       500,
		ActionKarmaMap: map[string]int{"kill_innocent": 500},
	}
	rule := NewKarmaRule(cfg)

	ctx := &KarmaContext{
		ActionType: "kill_innocent",
		ActorKarma: 400,
	}
	result := rule.CalculateKarmaChange(ctx)

	// New karma should not exceed cap
	if result.NewKarma > cfg.KarmaCap {
		t.Errorf("New karma %d should not exceed cap %d", result.NewKarma, cfg.KarmaCap)
	}
}

func TestCalculateKarmaChange_EmptyContext(t *testing.T) {
	rule := NewKarmaRule(nil)

	result := rule.CalculateKarmaChange(nil)

	if result.Reason != "empty context" {
		t.Errorf("Empty context should return 'empty context' reason, got %s", result.Reason)
	}
}

func TestCalculateKarmaChange_UnknownAction(t *testing.T) {
	rule := NewKarmaRule(nil)

	ctx := &KarmaContext{
		ActionType: "unknown_action",
		ActorKarma: 0,
	}
	result := rule.CalculateKarmaChange(ctx)

	// Unknown action should have 0 base karma
	if result.BaseKarma != 0 {
		t.Errorf("Unknown action should have 0 base karma, got %d", result.BaseKarma)
	}
	if result.KarmaChange != 0 {
		t.Errorf("Unknown action should have 0 karma change, got %d", result.KarmaChange)
	}
}

func TestCalculateHeavenlyMark(t *testing.T) {
	rule := NewKarmaRule(nil)

	tests := []struct {
		karma        int
		expectedMark string
	}{
		{0, "clear"},
		{50, "clear"},
		{99, "clear"},
		{100, "slight"},
		{200, "slight"},
		{499, "slight"},
		{500, "heavy"},
		{800, "heavy"},
		{999, "heavy"},
		{1000, "notorious"},
		{3000, "notorious"},
		{4999, "notorious"},
		{5000, "heaven_fury"},
		{9999, "heaven_fury"},
		{10000, "heaven_fury"},
	}

	for _, tt := range tests {
		mark := rule.CalculateHeavenlyMark(tt.karma)
		if mark != tt.expectedMark {
			t.Errorf("Karma %d: expected mark %s, got %s", tt.karma, tt.expectedMark, mark)
		}
	}
}

func TestApplyKarmaDecay_NoDecay(t *testing.T) {
	rule := NewKarmaRule(nil)

	// Zero hours = no decay
	result := rule.ApplyKarmaDecay(1000, 0)
	if result != 1000 {
		t.Errorf("Zero hours should not decay: expected 1000, got %d", result)
	}

	// Negative karma = no decay
	result = rule.ApplyKarmaDecay(-100, 10)
	if result != -100 {
		t.Errorf("Negative karma should not decay: expected -100, got %d", result)
	}
}

func TestApplyKarmaDecay_PositiveDecay(t *testing.T) {
	rule := NewKarmaRule(&KarmaRuleConfig{
		KarmaDecayRate: 0.01,
		KarmaCap:       10000,
	})

	// After 10 hours at 1% decay rate
	result := rule.ApplyKarmaDecay(1000, 10)

	// Expected: 1000 * exp(-0.01 * 10) = 1000 * exp(-0.1) ≈ 1000 * 0.9048 ≈ 905
	expected := int(float64(1000) * math.Exp(-0.01*10))
	if result != expected {
		t.Errorf("Decay after 10h: expected %d, got %d", expected, result)
	}
}

func TestApplyKarmaDecay_LongTime(t *testing.T) {
	rule := NewKarmaRule(&KarmaRuleConfig{
		KarmaDecayRate: 0.01,
		KarmaCap:       10000,
	})

	// After 1000 hours, karma should approach 0
	result := rule.ApplyKarmaDecay(1000, 1000)

	// Expected: 1000 * exp(-0.01 * 1000) = 1000 * exp(-10) ≈ 1000 * 0.000045 ≈ 0
	if result > 1 {
		t.Errorf("After 1000h karma should approach 0, got %d", result)
	}
}

func TestApplyKarmaDecay_CustomRate(t *testing.T) {
	rule := NewKarmaRule(&KarmaRuleConfig{
		KarmaDecayRate: 0.05,
		KarmaCap:       10000,
	})

	result := rule.ApplyKarmaDecay(500, 10)

	// Expected: 500 * exp(-0.05 * 10) = 500 * exp(-0.5) ≈ 500 * 0.6065 ≈ 303
	expected := int(float64(500) * math.Exp(-0.05*10))
	if result != expected {
		t.Errorf("Custom rate decay: expected %d, got %d", expected, result)
	}
}

func TestGetKarmaThresholds(t *testing.T) {
	rule := NewKarmaRule(nil)
	thresholds := rule.GetKarmaThresholds()

	if len(thresholds) != 5 {
		t.Errorf("Expected 5 thresholds, got %d", len(thresholds))
	}
	if thresholds["clear"] != 100 {
		t.Errorf("clear threshold mismatch, expected 100, got %d", thresholds["clear"])
	}
	if thresholds["heavy"] != 1000 {
		t.Errorf("heavy threshold mismatch, expected 1000, got %d", thresholds["heavy"])
	}
	if thresholds["heaven_fury"] != 10000 {
		t.Errorf("heaven_fury threshold mismatch, expected 10000, got %d", thresholds["heaven_fury"])
	}
}

func TestCalculateKarmaChange_ComplexScenario(t *testing.T) {
	rule := NewKarmaRule(nil)

	// Scenario: A disciple (actor karma 200) kills their master (higher realm) in self-defense
	ctx := &KarmaContext{
		ActionType:    "kill_innocent",
		ActorKarma:    200,
		ActorRealm:    2,
		TargetRealm:   5,
		Relationship:  "mentor",
		IsSelfDefense: true,
	}

	result := rule.CalculateKarmaChange(ctx)

	// Base karma for kill_innocent is 500
	// Context: diminishing returns (karma 200/10000) + realm diff (5-2=3, +60%) + self-defense (0.5)
	// Relation: mentor (2.0)
	// Should still be positive karma overall
	if result.BaseKarma != 500 {
		t.Errorf("Base karma should be 500 for kill_innocent, got %d", result.BaseKarma)
	}
	if result.NewKarma > rule.cfg.KarmaCap {
		t.Errorf("New karma should not exceed cap")
	}
	if result.Reason == "" {
		t.Error("Reason should not be empty")
	}
}

func TestCalculateKarmaChange_NegativeKarma(t *testing.T) {
	rule := NewKarmaRule(nil)

	// Saving someone from great evil
	ctx := &KarmaContext{
		ActionType: "save_life",
		ActorKarma: 500,
	}

	result := rule.CalculateKarmaChange(ctx)

	// Negative karma action should decrease karma
	if result.KarmaChange >= 0 {
		t.Errorf("Saving life should decrease karma, got %d", result.KarmaChange)
	}

	// New karma should not go below 0
	if result.NewKarma < 0 {
		t.Errorf("New karma should not be negative, got %d", result.NewKarma)
	}
}

func TestCalculateKarmaChange_AllActionTypes(t *testing.T) {
	rule := NewKarmaRule(nil)

	// Test that all defined action types work without errors
	for actionType := range DefaultKarmaConfig().ActionKarmaMap {
		ctx := &KarmaContext{
			ActionType: actionType,
			ActorKarma: 0,
		}
		result := rule.CalculateKarmaChange(ctx)

		// Should not panic and should return a result
		if result == nil {
			t.Errorf("CalculateKarmaChange(%s) returned nil", actionType)
		}
	}
}
