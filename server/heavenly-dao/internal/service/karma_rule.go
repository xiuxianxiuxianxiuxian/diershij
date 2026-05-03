package service

import (
	"fmt"
	"math"
)

// KarmaRuleConfig holds configuration for karma calculations.
type KarmaRuleConfig struct {
	KarmaDecayRate  float64            // hourly decay rate (default 0.01 = 1%)
	KarmaCap        int                // maximum karma value
	ActionKarmaMap  map[string]int     // base karma values for actions
	KarmaThresholds map[string]int     // thresholds for heavenly marks
}

// DefaultKarmaConfig returns the default karma configuration.
func DefaultKarmaConfig() *KarmaRuleConfig {
	return &KarmaRuleConfig{
		KarmaDecayRate: 0.01,
		KarmaCap:       10000,
		ActionKarmaMap: map[string]int{
			"kill_innocent":     500,
			"kill_cultivator":   200,
			"kill_demon":        -50,
			"save_life":         -100,
			"teach_method":      -200,
			"betray_master":     1000,
			"break_oath":        300,
			"destroy_sect":      800,
			"create_method":     -500,
			"steal":             150,
			"rob":               400,
			"help_stranger":     -30,
			"defeat_enemy":      50,
			"surrender":         -20,
			"spare_defeated":    -80,
			"torture":           600,
			"rescue_hostage":    -150,
			"pillaging":         300,
			"donate_resources":  -100,
			"protect_weak":      -120,
			"bully_weak":        200,
			"honor_agreement":   -50,
		},
		KarmaThresholds: map[string]int{
			"clear":       100,
			"slight":      500,
			"heavy":       1000,
			"notorious":   5000,
			"heaven_fury": 10000,
		},
	}
}

// KarmaContext provides the context for a karma calculation.
type KarmaContext struct {
	ActionType    string            // the action being performed
	ActorKarma    int               // current karma of the actor
	ActorRealm    int               // realm level of the actor
	TargetRealm   int               // realm level of the target (if applicable)
	Relationship  string            // relationship between actor and target (mentor, enemy, etc.)
	IsSelfDefense bool              // whether the action was in self-defense
	Location      string            // where the action took place
	Extras        map[string]any    // additional context
}

// KarmaResult is the result of a karma calculation.
type KarmaResult struct {
	BaseKarma      int     // base karma value for the action
	ContextMult    float64 // context multiplier
	RelationMult   float64 // relationship multiplier
	KarmaChange    int     // final karma change
	NewKarma       int     // karma after change (capped)
	NewHeavenlyMark string // updated heavenly mark
	Reason         string  // explanation of the calculation
}

// KarmaRule implements the karma algorithm from design doc section 8.1.
type KarmaRule struct {
	cfg *KarmaRuleConfig
}

// NewKarmaRule creates a new KarmaRule with the given config.
func NewKarmaRule(cfg *KarmaRuleConfig) *KarmaRule {
	if cfg == nil {
		cfg = DefaultKarmaConfig()
	}
	return &KarmaRule{cfg: cfg}
}

// CalculateKarmaChange computes the karma delta for an action.
// Formula: base_karma * context_multiplier * relationship_multiplier
func (r *KarmaRule) CalculateKarmaChange(ctx *KarmaContext) *KarmaResult {
	if ctx == nil {
		return &KarmaResult{Reason: "empty context"}
	}

	// 1. Get base karma for the action
	baseKarma := r.getActionKarmaBase(ctx.ActionType)

	// 2. Calculate context multiplier
	contextMult := r.calculateContextMultiplier(ctx)

	// 3. Calculate relationship multiplier
	relationMult := r.calculateRelationshipMultiplier(ctx)

	// 4. Compute final karma change
	karmaChange := int(float64(baseKarma) * contextMult * relationMult)

	// 5. Apply karma cap
	newKarma := ctx.ActorKarma + karmaChange
	if newKarma > r.cfg.KarmaCap {
		newKarma = r.cfg.KarmaCap
	}
	if newKarma < 0 {
		newKarma = 0
	}

	// 6. Determine new heavenly mark
	newMark := r.CalculateHeavenlyMark(newKarma)

	return &KarmaResult{
		BaseKarma:      baseKarma,
		ContextMult:    contextMult,
		RelationMult:   relationMult,
		KarmaChange:    karmaChange,
		NewKarma:       newKarma,
		NewHeavenlyMark: newMark,
		Reason:         r.formatReason(ctx, baseKarma, contextMult, relationMult),
	}
}

// getActionKarmaBase returns the base karma value for an action type.
func (r *KarmaRule) getActionKarmaBase(actionType string) int {
	if karma, ok := r.cfg.ActionKarmaMap[actionType]; ok {
		return karma
	}
	return 0
}

// calculateContextMultiplier computes the context modifier.
// Higher karma actors have diminishing returns on new karma.
func (r *KarmaRule) calculateContextMultiplier(ctx *KarmaContext) float64 {
	mult := 1.0

	// Diminishing returns: higher karma actors accumulate karma slower
	if r.cfg.KarmaCap > 0 {
		diminishFactor := 1.0 - float64(ctx.ActorKarma)/float64(r.cfg.KarmaCap)
		mult *= diminishFactor
	}

	// Target realm modifier: killing higher realm targets has more karma impact
	if ctx.TargetRealm > ctx.ActorRealm && ctx.TargetRealm > 0 {
		realmDiff := float64(ctx.TargetRealm - ctx.ActorRealm)
		mult *= (1.0 + realmDiff*0.2)
	}

	// Self-defense reduces karma impact
	if ctx.IsSelfDefense {
		mult *= 0.5
	}

	return mult
}

// calculateRelationshipMultiplier computes the relationship modifier.
func (r *KarmaRule) calculateRelationshipMultiplier(ctx *KarmaContext) float64 {
	switch ctx.Relationship {
	case "mentor":
		// Betraying a master doubles karma
		return 2.0
	case "enemy":
		// Killing an enemy halves karma
		return 0.5
	case "disciple":
		// Killing a disciple triples karma
		return 3.0
	case "sworn_sibling":
		// Betraying a sworn sibling doubles karma
		return 2.0
	case "lover":
		// Betraying a lover triples karma
		return 3.0
	case "sect_member":
		// Attacking sect member has 1.5x karma
		return 1.5
	default:
		return 1.0
	}
}

// CalculateHeavenlyMark determines the heavenly mark based on karma value.
func (r *KarmaRule) CalculateHeavenlyMark(karma int) string {
	if karma < r.cfg.KarmaThresholds["clear"] {
		return "clear"
	}
	if karma < r.cfg.KarmaThresholds["slight"] {
		return "slight"
	}
	if karma < r.cfg.KarmaThresholds["heavy"] {
		return "heavy"
	}
	if karma < r.cfg.KarmaThresholds["notorious"] {
		return "notorious"
	}
	return "heaven_fury"
}

// ApplyKarmaDecay applies time-based karma decay.
// Formula: new_karma = old_karma * exp(-decay_rate * hours_elapsed)
func (r *KarmaRule) ApplyKarmaDecay(oldKarma int, hoursElapsed float64) int {
	if oldKarma <= 0 || hoursElapsed <= 0 {
		return oldKarma
	}

	decayFactor := math.Exp(-r.cfg.KarmaDecayRate * hoursElapsed)
	newKarma := int(float64(oldKarma) * decayFactor)

	if newKarma < 0 {
		return 0
	}
	return newKarma
}

// GetKarmaThresholds returns the karma thresholds for heavenly marks.
func (r *KarmaRule) GetKarmaThresholds() map[string]int {
	return r.cfg.KarmaThresholds
}

// formatReason creates a human-readable explanation of the karma calculation.
func (r *KarmaRule) formatReason(ctx *KarmaContext, baseKarma int, contextMult, relationMult float64) string {
	reason := "action: " + ctx.ActionType + ", base: " + itoa(baseKarma)

	if contextMult != 1.0 {
		reason += ", context: " + ftoa(contextMult)
	}
	if relationMult != 1.0 {
		reason += ", relation(" + ctx.Relationship + "): " + ftoa(relationMult)
	}

	return reason
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func ftoa(f float64) string {
	return fmt.Sprintf("%.2f", f)
}
