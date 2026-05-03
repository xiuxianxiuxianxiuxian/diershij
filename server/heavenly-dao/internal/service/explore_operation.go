package service

import (
	"fmt"
	"time"

	"github.com/cultivation-world/shared/types"
)

// ExploreOperation handles the region exploration operation.
type ExploreOperation struct {
	fortuneRule    *FortuneRule
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
}

// NewExploreOperation creates a new ExploreOperation.
func NewExploreOperation(cooldownPeriod time.Duration) *ExploreOperation {
	return &ExploreOperation{
		fortuneRule:    NewFortuneRule(),
		cooldownMap:    make(map[string]time.Time),
		cooldownPeriod: cooldownPeriod,
	}
}

// ExploreInput holds inputs for an exploration operation.
type ExploreInput struct {
	EntityID        string
	RegionName      string
	RegionLevel     int
	SpiritualDensity float64
	IsDangerZone    bool
	FortuneScore    float64
	Luck            int
	Realm           types.CultivationRealm
	MovementSpeed   float64
}

// ExploreResult holds the outcome of an exploration.
type ExploreResult struct {
	Success       bool
	Events        []ExploreEvent
	ResourcesGained []ResourceGain
	TimeElapsed   int // in-game hours
	Message       string
}

// ExploreEvent represents a random event during exploration.
type ExploreEvent struct {
	Type        string // encounter / discovery / trap / nothing / epiphany
	Title       string
	Description string
	Reward      *ResourceGain
	Danger      float64
}

// ResourceGain represents a resource obtained from exploration.
type ResourceGain struct {
	Name     string
	Type     string // herb / ore / crystal / stone
	Quantity int
	Quality  string
}

// ExecuteExplore runs an exploration operation.
func (op *ExploreOperation) ExecuteExplore(input ExploreInput, now time.Time, randFloat func() float64) (*ExploreResult, error) {
	// Check cooldown
	if lastTime, ok := op.cooldownMap[input.EntityID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("explore cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	// Determine exploration duration based on region level and movement speed
	duration := calculateExploreDuration(input.RegionLevel, input.MovementSpeed)

	// Generate events based on region properties and fortune
	events := op.generateExploreEvents(input, duration, randFloat)

	// Collect resources from events
	var resources []ResourceGain
	totalDanger := 0.0
	for _, event := range events {
		if event.Reward != nil {
			resources = append(resources, *event.Reward)
		}
		totalDanger += event.Danger
	}

	// Check if danger exceeded threshold
	if totalDanger > float64(duration)*0.3 {
		events = append(events, ExploreEvent{
			Type:        "trap",
			Title:       "危险区域",
			Description: "探索过程中遭遇危险，被迫撤退",
			Danger:      1.0,
		})
	}

	op.cooldownMap[input.EntityID] = now

	success := totalDanger <= float64(duration)*0.3
	message := fmt.Sprintf("在 %s 探索了 %d 小时", input.RegionName, duration)
	if !success {
		message = "探索遭遇危险，被迫撤退"
	}

	return &ExploreResult{
		Success:       success,
		Events:        events,
		ResourcesGained: resources,
		TimeElapsed:   duration,
		Message:       message,
	}, nil
}

func calculateExploreDuration(regionLevel int, movementSpeed float64) int {
	// Base duration scales with region level
	base := regionLevel * 4 + 8 // 12-48 hours

	// Movement speed reduces duration
	if movementSpeed > 0 {
		duration := float64(base) / movementSpeed
		return int(duration)
	}
	return base
}

func (op *ExploreOperation) generateExploreEvents(input ExploreInput, duration int, randFloat func() float64) []ExploreEvent {
	var events []ExploreEvent

	// Number of potential events scales with duration
	eventCount := duration / 6
	if eventCount < 1 {
		eventCount = 1
	}
	if eventCount > 8 {
		eventCount = 8
	}

	for i := 0; i < eventCount; i++ {
		event := op.generateSingleEvent(input, randFloat)
		events = append(events, event)
	}

	return events
}

func (op *ExploreOperation) generateSingleEvent(input ExploreInput, randFloat func() float64) ExploreEvent {
	r := randFloat()

	// Fortune modifier: higher fortune = better events
	fortuneMult := 1.0 + (input.FortuneScore-50.0)/100.0

	// Event type weights
	if r < 0.30*fortuneMult {
		return op.generateDiscoveryEvent(input, randFloat)
	}
	if r < 0.45 {
		return op.generateEncounterEvent(input, randFloat)
	}
	if r < 0.55 && input.IsDangerZone {
		return op.generateTrapEvent(input, randFloat)
	}
	if r < 0.65*fortuneMult {
		return op.generateEpiphanyEvent(input, randFloat)
	}

	return ExploreEvent{
		Type:        "nothing",
		Title:       "平静",
		Description: "一路平安，没有特别发现",
		Danger:      0,
	}
}

func (op *ExploreOperation) generateDiscoveryEvent(input ExploreInput, randFloat func() float64) ExploreEvent {
	herbs := []string{"千年灵芝", "龙涎草", "紫金花", "冰心莲", "天灵果"}
	ores := []string{"玄铁矿", "灵石矿", "星辰砂", "寒玉髓", "赤金"}

	r := randFloat()
	var name, resType string
	if r < 0.6 {
		name = herbs[int(randFloat()*float64(len(herbs)))]
		resType = "herb"
	} else {
		name = ores[int(randFloat()*float64(len(ores)))]
		resType = "ore"
	}

	quantity := 1 + int(randFloat()*float64(input.RegionLevel))
	quality := "common"
	if randFloat() < 0.2 {
		quality = "rare"
	}

	return ExploreEvent{
		Type:        "discovery",
		Title:       "发现资源",
		Description: fmt.Sprintf("发现了 %s", name),
		Reward: &ResourceGain{
			Name:     name,
			Type:     resType,
			Quantity: quantity,
			Quality:  quality,
		},
		Danger: 0,
	}
}

func (op *ExploreOperation) generateEncounterEvent(input ExploreInput, randFloat func() float64) ExploreEvent {
	encounters := []string{
		"流浪修士", "散修", "妖兽", "秘境守护者", "神秘老者",
	}
	idx := int(randFloat() * float64(len(encounters)))
	encounter := encounters[idx]

	friendly := randFloat() < 0.4*input.FortuneScore/100.0
	if friendly {
		return ExploreEvent{
			Type:        "encounter",
			Title:       "友好相遇",
			Description: fmt.Sprintf("遇到了%s，获得了一些指点", encounter),
			Danger:      0,
		}
	}

	return ExploreEvent{
		Type:        "encounter",
		Title:       "遭遇",
		Description: fmt.Sprintf("遇到了%s", encounter),
		Danger:      0.3 + randFloat()*0.3,
	}
}

func (op *ExploreOperation) generateTrapEvent(input ExploreInput, randFloat func() float64) ExploreEvent {
	traps := []string{
		"幻阵", "毒瘴", "陷阱", "妖兽巢穴", "禁制",
	}
	idx := int(randFloat() * float64(len(traps)))

	return ExploreEvent{
		Type:        "trap",
		Title:       "触发危险",
		Description: fmt.Sprintf("误入了%s", traps[idx]),
		Danger:      0.5 + randFloat()*0.5,
	}
}

func (op *ExploreOperation) generateEpiphanyEvent(input ExploreInput, randFloat func() float64) ExploreEvent {
	insights := []string{
		"领悟了一丝剑意",
		"对功法有了新的理解",
		"感受到了天地法则的波动",
		"心境得到了提升",
	}
	idx := int(randFloat() * float64(len(insights)))

	return ExploreEvent{
		Type:        "epiphany",
		Title:       "顿悟",
		Description: insights[idx],
		Danger:      0,
	}
}

// GetCooldownRemaining returns remaining cooldown time for an entity.
func (op *ExploreOperation) GetCooldownRemaining(entityID string, now time.Time) time.Duration {
	if lastTime, ok := op.cooldownMap[entityID]; ok {
		remaining := op.cooldownPeriod - now.Sub(lastTime)
		if remaining > 0 {
			return remaining
		}
	}
	return 0
}

// ClearCooldown removes cooldown for an entity.
func (op *ExploreOperation) ClearCooldown(entityID string) {
	delete(op.cooldownMap, entityID)
}
