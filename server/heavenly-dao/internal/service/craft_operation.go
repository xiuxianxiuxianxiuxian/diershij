package service

import (
	"fmt"
	"time"
)

// CraftOperation handles crafting operations (alchemy, forging, formation).
type CraftOperation struct {
	alchemyRule    *AlchemyRule
	forgingRule    *ForgingRule
	formationRule  *FormationRule
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
}

// NewCraftOperation creates a new CraftOperation.
func NewCraftOperation(cooldownPeriod time.Duration) *CraftOperation {
	return &CraftOperation{
		alchemyRule:    NewAlchemyRule(),
		forgingRule:    NewForgingRule(),
		formationRule:  NewFormationRule(),
		cooldownMap:    make(map[string]time.Time),
		cooldownPeriod: cooldownPeriod,
	}
}

// CraftType represents the type of crafting operation.
type CraftType string

const (
	CraftTypeAlchemy    CraftType = "alchemy"
	CraftTypeForging    CraftType = "forging"
	CraftTypeFormation  CraftType = "formation"
)

// CraftInput holds inputs for a crafting operation.
type CraftInput struct {
	EntityID      string
	CraftType     CraftType
	MentalStability int
	LocationBonus   float64
	EquipmentBonus  float64
	ElementAffinity float64
	Luck            int
	SkillLevel      int
	SpiritStoneCount int
}

// CraftResult holds the outcome of a crafting operation.
type CraftResult struct {
	Success        bool
	CraftType      CraftType
	Quality        string
	QualityValue   float64
	BonusDrop      bool
	Message        string
}

// ExecuteCraft dispatches to the appropriate crafting algorithm.
func (op *CraftOperation) ExecuteCraft(input CraftInput, now time.Time, randFloat func() float64) (*CraftResult, error) {
	// Check cooldown
	if lastTime, ok := op.cooldownMap[input.EntityID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("craft cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	var result *CraftResult
	var err error

	switch input.CraftType {
	case CraftTypeAlchemy:
		result, err = op.executeAlchemy(input, randFloat)
	case CraftTypeForging:
		result, err = op.executeForging(input, randFloat)
	case CraftTypeFormation:
		result, err = op.executeFormation(input, randFloat)
	default:
		return nil, fmt.Errorf("unknown craft type: %s", input.CraftType)
	}

	if err != nil {
		return nil, err
	}

	// Set cooldown
	op.cooldownMap[input.EntityID] = now

	return result, nil
}

func (op *CraftOperation) executeAlchemy(input CraftInput, randFloat func() float64) (*CraftResult, error) {
	recipe := PillRecipe{
		Name:          "custom_pill",
		RequiredLevel: input.SkillLevel,
		BaseQuality:   2,
		Difficulty:    1.0,
	}

	alchemyInput := AlchemyInput{
		AlchemyLevel:    input.SkillLevel,
		MentalStability: input.MentalStability,
		LocationBonus:   input.LocationBonus,
		EquipmentBonus:  input.EquipmentBonus,
		ElementAffinity: input.ElementAffinity,
		Luck:            input.Luck,
	}

	pillResult := op.alchemyRule.ResolveAlchemy(recipe, alchemyInput, randFloat)

	if pillResult.Success {
		return &CraftResult{
			Success:      true,
			CraftType:    CraftTypeAlchemy,
			Quality:      string(pillResult.Grade),
			QualityValue: pillResult.Quality,
			BonusDrop:    pillResult.Grade == PillGradeExcellent || pillResult.Grade == PillGradePerfect,
			Message:      fmt.Sprintf("炼丹成功，品质：%s", pillResult.Grade),
		}, nil
	}

	return &CraftResult{
		Success:   false,
		CraftType: CraftTypeAlchemy,
		Message:   pillResult.Message,
	}, nil
}

func (op *CraftOperation) executeForging(input CraftInput, randFloat func() float64) (*CraftResult, error) {
	recipe := ArtifactRecipe{
		Name:          "custom_artifact",
		RequiredLevel: input.SkillLevel,
		BaseGrade:     2,
		Difficulty:    1.0,
	}

	forgingInput := ForgingInput{
		ArtificingLevel: input.SkillLevel,
		MentalStability: input.MentalStability,
		LocationBonus:   input.LocationBonus,
		EquipmentBonus:  input.EquipmentBonus,
		ElementAffinity: input.ElementAffinity,
		Luck:            input.Luck,
	}

	forgingResult := op.forgingRule.ResolveForging(recipe, forgingInput, randFloat)

	if forgingResult.Success {
		return &CraftResult{
			Success:      true,
			CraftType:    CraftTypeForging,
			Quality:      string(forgingResult.Grade),
			QualityValue: forgingResult.Quality,
			BonusDrop:    forgingResult.Grade == ArtifactGradeHeavenly || forgingResult.Grade == ArtifactGradeDivine,
			Message:      fmt.Sprintf("炼器成功，品质：%s", forgingResult.Grade),
		}, nil
	}

	return &CraftResult{
		Success:   false,
		CraftType: CraftTypeForging,
		Message:   forgingResult.Message,
	}, nil
}

func (op *CraftOperation) executeFormation(input CraftInput, randFloat func() float64) (*CraftResult, error) {
	formation := Formation{
		Name:          "custom_formation",
		RequiredLevel: input.SkillLevel,
		BasePower:     100.0,
		Difficulty:    1.0,
		Nodes: []FormationNode{
			{Power: 10.0, IsCore: true},
			{Power: 5.0, IsCore: false},
			{Power: 5.0, IsCore: false},
		},
	}

	formationInput := FormationInput{
		FormationLevel:   input.SkillLevel,
		MentalStability:  input.MentalStability,
		LocationBonus:    input.LocationBonus,
		EquipmentBonus:   input.EquipmentBonus,
		ElementAffinity:  input.ElementAffinity,
		SpiritStoneCount: input.SpiritStoneCount,
		Luck:             input.Luck,
	}

	power := op.formationRule.CalculateFormationPower(formation, formationInput)

	// Formation success if power exceeds a threshold
	threshold := float64(input.SkillLevel) * 50.0
	success := power >= threshold

	var quality string
	if power >= threshold*1.5 {
		quality = "excellent"
	} else if power >= threshold*1.2 {
		quality = "fine"
	} else if success {
		quality = "common"
	} else {
		quality = "failed"
	}

	return &CraftResult{
		Success:      success,
		CraftType:    CraftTypeFormation,
		Quality:      quality,
		QualityValue: power,
		BonusDrop:    power >= threshold*1.5,
		Message:      fmt.Sprintf("阵法布置%s，威力：%.0f", map[bool]string{true: "成功", false: "失败"}[success], power),
	}, nil
}

// GetCooldownRemaining returns remaining cooldown time for an entity.
func (op *CraftOperation) GetCooldownRemaining(entityID string, now time.Time) time.Duration {
	if lastTime, ok := op.cooldownMap[entityID]; ok {
		remaining := op.cooldownPeriod - now.Sub(lastTime)
		if remaining > 0 {
			return remaining
		}
	}
	return 0
}

// ClearCooldown removes cooldown for an entity.
func (op *CraftOperation) ClearCooldown(entityID string) {
	delete(op.cooldownMap, entityID)
}
