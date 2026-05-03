package service

import (
	"fmt"
	"strings"
)

// TemplateType represents the type of template.
type TemplateType string

const (
	TemplateTypeBehavior  TemplateType = "behavior"
	TemplateTypeDialogue  TemplateType = "dialogue"
	TemplateTypeDecision  TemplateType = "decision"
)

// BehaviorTemplate represents a reusable behavior pattern.
type BehaviorTemplate struct {
	ID          string
	Type        TemplateType
	Pattern     string            // pattern to match
	Action      string            // action to take
	Parameters  map[string]string // action parameters
	Priority    int               // higher = more important
	UsageCount  int               // how many times used
	SuccessRate float64           // success rate (0-1)
	Tags        []string          // categorization tags
}

// TemplateLibrary stores and manages behavior templates.
type TemplateLibrary struct {
	templates []*BehaviorTemplate
	// Type -> pattern -> template index
	index map[TemplateType]map[string]int
}

// NewTemplateLibrary creates a new template library.
func NewTemplateLibrary() *TemplateLibrary {
	return &TemplateLibrary{
		index: make(map[TemplateType]map[string]int),
	}
}

// AddTemplate adds a template to the library.
func (lib *TemplateLibrary) AddTemplate(tmpl *BehaviorTemplate) {
	lib.templates = append(lib.templates, tmpl)

	tmplType := tmpl.Type
	if lib.index[tmplType] == nil {
		lib.index[tmplType] = make(map[string]int)
	}
	lib.index[tmplType][tmpl.Pattern] = len(lib.templates) - 1
}

// MatchAndDecide finds the best matching template for the current context.
func (lib *TemplateLibrary) MatchAndDecide(ctx *NPCContext) *DecisionResult {
	bestMatch := lib.FindBestMatch(ctx)
	if bestMatch == nil {
		return nil
	}

	// Increment usage count
	bestMatch.UsageCount++

	return &DecisionResult{
		Action:     bestMatch.Action,
		Parameters: convertMap(bestMatch.Parameters),
		Reasoning:  "模板匹配决策",
		Confidence: bestMatch.SuccessRate,
	}
}

func convertMap(src map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range src {
		result[k] = v
	}
	return result
}

// FindBestMatch finds the best matching template for the current context.
func (lib *TemplateLibrary) FindBestMatch(ctx *NPCContext) *BehaviorTemplate {
	var best *BehaviorTemplate
	bestScore := 0.0

	for _, tmpl := range lib.templates {
		score := lib.calculateSimilarity(ctx, tmpl)
		if score > bestScore {
			bestScore = score
			best = tmpl
		}
	}

	return best
}

// calculateSimilarity calculates how well a template matches the context.
func (lib *TemplateLibrary) calculateSimilarity(ctx *NPCContext, tmpl *BehaviorTemplate) float64 {
	score := 0.0
	pattern := strings.ToLower(tmpl.Pattern)

	// Check pattern keywords against context
	if strings.Contains(pattern, "combat") && ctx.IsInCombat {
		score += 2.0
	}
	if strings.Contains(pattern, "danger") && ctx.IsInDanger {
		score += 2.0
	}
	if strings.Contains(pattern, "injured") && ctx.IsInjured() {
		score += 2.0
	}
	if strings.Contains(pattern, "healthy") && ctx.IsHealthy() {
		score += 1.0
	}
	if strings.Contains(pattern, "low_qi") && ctx.IsLowOnQi() {
		score += 2.0
	}
	if strings.Contains(pattern, "has_target") && ctx.HasTarget {
		score += 1.5
	}

	// Check tags
	for _, tag := range tmpl.Tags {
		if strings.Contains(pattern, tag) {
			score += 0.5
		}
	}

	// Boost by priority
	score += float64(tmpl.Priority) * 0.1

	// Boost by success rate
	score += tmpl.SuccessRate * 0.5

	return score
}

// AddLLMResult adds an LLM-generated result to the template library.
func (lib *TemplateLibrary) AddLLMResult(pattern string, action string, parameters map[string]string, tmplType TemplateType) {
	lib.AddTemplate(&BehaviorTemplate{
		ID:       lib.generateID(),
		Type:     tmplType,
		Pattern:  pattern,
		Action:   action,
		Parameters: parameters,
		Priority: 5,
		SuccessRate: 0.7, // Initial success rate for LLM-generated templates
		Tags:     []string{"llm_generated"},
	})
}

// UpdateTemplateSuccess updates a template's success rate.
func (lib *TemplateLibrary) UpdateTemplateSuccess(templateID string, success bool) {
	for _, tmpl := range lib.templates {
		if tmpl.ID == templateID {
			if success {
				tmpl.SuccessRate = min(1.0, tmpl.SuccessRate+0.05)
			} else {
				tmpl.SuccessRate = max(0.0, tmpl.SuccessRate-0.05)
			}
			return
		}
	}
}

func (lib *TemplateLibrary) generateID() string {
	return generateTemplateID(len(lib.templates))
}

func generateTemplateID(count int) string {
	return fmt.Sprintf("tmpl_%d", count)
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// GetTemplateStats returns statistics about the template library.
func (lib *TemplateLibrary) GetTemplateStats() map[string]int {
	stats := make(map[string]int)
	stats["total"] = len(lib.templates)

	for _, tmpl := range lib.templates {
		stats[string(tmpl.Type)]++
	}

	return stats
}

// GetTemplatesByType returns all templates of a specific type.
func (lib *TemplateLibrary) GetTemplatesByType(tmplType TemplateType) []*BehaviorTemplate {
	var result []*BehaviorTemplate
	for _, tmpl := range lib.templates {
		if tmpl.Type == tmplType {
			result = append(result, tmpl)
		}
	}
	return result
}

// SearchTemplates searches for templates by keyword.
func (lib *TemplateLibrary) SearchTemplates(keyword string) []*BehaviorTemplate {
	var result []*BehaviorTemplate
	keyword = strings.ToLower(keyword)

	for _, tmpl := range lib.templates {
		if strings.Contains(strings.ToLower(tmpl.Pattern), keyword) ||
		   strings.Contains(strings.ToLower(tmpl.Action), keyword) {
			result = append(result, tmpl)
		}
	}

	return result
}

// DefaultBehaviorTemplates returns a set of default behavior templates.
func DefaultBehaviorTemplates() []*BehaviorTemplate {
	return []*BehaviorTemplate{
		// Combat behaviors
		{ID: "b001", Type: TemplateTypeBehavior, Pattern: "combat has_target is_healthy", Action: "attack", Priority: 10, SuccessRate: 0.85, Tags: []string{"combat", "offensive"}},
		{ID: "b002", Type: TemplateTypeBehavior, Pattern: "combat has_target injured", Action: "flee", Priority: 9, SuccessRate: 0.90, Tags: []string{"combat", "defensive"}},
		{ID: "b003", Type: TemplateTypeBehavior, Pattern: "combat low_qi", Action: "retreat", Priority: 8, SuccessRate: 0.80, Tags: []string{"combat", "resource"}},

		// Cultivation behaviors
		{ID: "c001", Type: TemplateTypeBehavior, Pattern: "cultivate healthy no_danger", Action: "cultivate", Priority: 5, SuccessRate: 0.75, Tags: []string{"cultivation", "routine"}},
		{ID: "c002", Type: TemplateTypeBehavior, Pattern: "cultivate low_qi", Action: "meditate", Priority: 7, SuccessRate: 0.70, Tags: []string{"cultivation", "recovery"}},
		{ID: "c003", Type: TemplateTypeBehavior, Pattern: "breakthrough ready", Action: "breakthrough", Priority: 10, SuccessRate: 0.60, Tags: []string{"cultivation", "advancement"}},

		// Survival behaviors
		{ID: "s001", Type: TemplateTypeBehavior, Pattern: "danger injured", Action: "flee", Priority: 10, SuccessRate: 0.95, Tags: []string{"survival", "emergency"}},
		{ID: "s002", Type: TemplateTypeBehavior, Pattern: "injured no_danger", Action: "rest", Priority: 8, SuccessRate: 0.85, Tags: []string{"survival", "recovery"}},
		{ID: "s003", Type: TemplateTypeBehavior, Pattern: "low_qi no_danger", Action: "meditate", Priority: 7, SuccessRate: 0.80, Tags: []string{"survival", "recovery"}},

		// Resource behaviors
		{ID: "r001", Type: TemplateTypeBehavior, Pattern: "gather need_herbs", Action: "gather", Priority: 4, SuccessRate: 0.70, Tags: []string{"resource", "gathering"}},
		{ID: "r002", Type: TemplateTypeBehavior, Pattern: "craft need_pills", Action: "craft", Priority: 5, SuccessRate: 0.65, Tags: []string{"resource", "crafting"}},

		// Decision templates
		{ID: "d001", Type: TemplateTypeDecision, Pattern: "strong_enemy nearby", Action: "avoid", Priority: 8, SuccessRate: 0.80, Tags: []string{"decision", "cautious"}},
		{ID: "d002", Type: TemplateTypeDecision, Pattern: "weak_enemy valuable_loot", Action: "engage", Priority: 7, SuccessRate: 0.75, Tags: []string{"decision", "aggressive"}},
	}
}

// DefaultDialogueTemplates returns a set of default dialogue templates.
func DefaultDialogueTemplates() []*BehaviorTemplate {
	return []*BehaviorTemplate{
		{ID: "dlg001", Type: TemplateTypeDialogue, Pattern: "greeting friendly", Action: "greet", Priority: 5, SuccessRate: 0.90, Tags: []string{"social", "greeting"}},
		{ID: "dlg002", Type: TemplateTypeDialogue, Pattern: "trade interested", Action: "offer_trade", Priority: 6, SuccessRate: 0.70, Tags: []string{"social", "trade"}},
		{ID: "dlg003", Type: TemplateTypeDialogue, Pattern: "combat hostile", Action: "threaten", Priority: 8, SuccessRate: 0.60, Tags: []string{"social", "combat"}},
		{ID: "dlg004", Type: TemplateTypeDialogue, Pattern: "help needy", Action: "offer_help", Priority: 4, SuccessRate: 0.80, Tags: []string{"social", "helpful"}},
	}
}
