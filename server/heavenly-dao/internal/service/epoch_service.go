package service

import (
	"fmt"
	"time"
)

// EpochPhase represents the current phase of the world.
type EpochPhase string

const (
	EpochSpiritualAge   EpochPhase = "spiritual_age"   // spiritual energy abundant
	EpochDharmaDecline  EpochPhase = "dharma_decline"  // spiritual energy declining
	EpochChaos          EpochPhase = "chaos"           // chaotic period, frequent events
	EpochRenewal        EpochPhase = "renewal"         // recovery and growth
)

// WorldEpoch represents a complete world epoch.
type WorldEpoch struct {
	EpochNumber  int
	Phase        EpochPhase
	StartTime    time.Time
	EndTime      time.Time
	Description  string

	// Epoch characteristics
	SpiritualBase   float64 // base spiritual density (0.0-1.0)
	ResourceCap     int     // maximum resources per region
	EventFrequency  float64 // events per day
	SecretRealmCount int    // number of active secret realms
}

// EpochTransitionCriteria defines when an epoch should transition.
type EpochTransitionCriteria struct {
	// Minimum duration before transition can occur
	MinimumDurationDays int

	// Spiritual density threshold for transition
	SpiritualThreshold float64

	// Event count trigger (too many events = chaos)
	EventCountThreshold int

	// Balance threshold (severe imbalance = decline)
	BalanceThreshold float64
}

// DefaultEpochTransitionCriteria returns the default criteria.
func DefaultEpochTransitionCriteria() *EpochTransitionCriteria {
	return &EpochTransitionCriteria{
		MinimumDurationDays:  365, // 1 year minimum
		SpiritualThreshold:   0.15, // transition if density drops below 15%
		EventCountThreshold:  50,   // transition if too many events
		BalanceThreshold:     0.85, // transition if severely unbalanced
	}
}

// EpochService manages world epoch transitions.
type EpochService struct {
	currentEpoch *WorldEpoch
	history      []*WorldEpoch
	criteria     *EpochTransitionCriteria
	epochCounter int
}

// NewEpochService creates a new EpochService.
func NewEpochService() *EpochService {
	return &EpochService{
		history:  []*WorldEpoch{},
		criteria: DefaultEpochTransitionCriteria(),
	}
}

// EpochInput holds inputs for creating a new epoch.
type EpochInput struct {
	Phase            EpochPhase
	Description      string
	SpiritualBase    float64
	ResourceCap      int
	EventFrequency   float64
	SecretRealmCount int
}

// StartNewEpoch creates and starts a new world epoch.
func (s *EpochService) StartNewEpoch(input EpochInput, now time.Time) (*WorldEpoch, error) {
	if input.SpiritualBase < 0 || input.SpiritualBase > 1 {
		return nil, fmt.Errorf("spiritual base must be between 0 and 1")
	}

	// End current epoch if exists
	if s.currentEpoch != nil {
		s.currentEpoch.EndTime = now
		s.history = append(s.history, s.currentEpoch)
	}

	s.epochCounter++
	epoch := &WorldEpoch{
		EpochNumber:      s.epochCounter,
		Phase:            input.Phase,
		StartTime:        now,
		EndTime:          now.Add(365 * 24 * time.Hour), // default 1 year
		Description:      input.Description,
		SpiritualBase:    input.SpiritualBase,
		ResourceCap:      input.ResourceCap,
		EventFrequency:   input.EventFrequency,
		SecretRealmCount: input.SecretRealmCount,
	}

	s.currentEpoch = epoch
	return epoch, nil
}

// EpochTransitionResult holds the outcome of a transition check.
type EpochTransitionResult struct {
	ShouldTransition bool
	Reason           string
	SuggestedPhase   EpochPhase
}

// CheckTransition checks if the world should transition to a new epoch.
func (s *EpochService) CheckTransition(
	spiritualDensity float64,
	eventCount int,
	balanceFactor float64,
	now time.Time,
) *EpochTransitionResult {
	if s.currentEpoch == nil {
		return &EpochTransitionResult{
			ShouldTransition: true,
			Reason:           "尚未开始任何纪元",
			SuggestedPhase:   EpochSpiritualAge,
		}
	}

	// Check minimum duration
	daysElapsed := int(now.Sub(s.currentEpoch.StartTime).Hours() / 24)
	if daysElapsed < s.criteria.MinimumDurationDays {
		return &EpochTransitionResult{
			ShouldTransition: false,
			Reason:           fmt.Sprintf("纪元持续时间不足 (%d/%d 天)", daysElapsed, s.criteria.MinimumDurationDays),
		}
	}

	// Check spiritual collapse
	if spiritualDensity < s.criteria.SpiritualThreshold {
		return &EpochTransitionResult{
			ShouldTransition: true,
			Reason:           fmt.Sprintf("灵气密度过低 (%.2f < %.2f)", spiritualDensity, s.criteria.SpiritualThreshold),
			SuggestedPhase:   EpochDharmaDecline,
		}
	}

	// Check event overload
	if eventCount >= s.criteria.EventCountThreshold {
		return &EpochTransitionResult{
			ShouldTransition: true,
			Reason:           fmt.Sprintf("事件过多 (%d >= %d)", eventCount, s.criteria.EventCountThreshold),
			SuggestedPhase:   EpochChaos,
		}
	}

	// Check severe imbalance
	if balanceFactor >= s.criteria.BalanceThreshold {
		return &EpochTransitionResult{
			ShouldTransition: true,
			Reason:           fmt.Sprintf("世界严重失衡 (%.2f >= %.2f)", balanceFactor, s.criteria.BalanceThreshold),
			SuggestedPhase:   EpochDharmaDecline,
		}
	}

	// Check if in decline but recovering
	if s.currentEpoch.Phase == EpochDharmaDecline && spiritualDensity > 0.5 {
		return &EpochTransitionResult{
			ShouldTransition: true,
			Reason:           fmt.Sprintf("灵气复苏 (%.2f > 0.50)", spiritualDensity),
			SuggestedPhase:   EpochRenewal,
		}
	}

	return &EpochTransitionResult{
		ShouldTransition: false,
		Reason:           "当前纪元稳定",
	}
}

// TransitionToNewEpoch performs the epoch transition.
func (s *EpochService) TransitionToNewEpoch(newPhase EpochPhase, now time.Time) (*WorldEpoch, error) {
	// Force transition regardless of check

	input := EpochInput{
		Phase: newPhase,
	}

	// Set epoch characteristics based on phase
	switch newPhase {
	case EpochSpiritualAge:
		input.Description = "灵气充沛，万物繁荣"
		input.SpiritualBase = 0.8
		input.ResourceCap = 150
		input.EventFrequency = 0.5
		input.SecretRealmCount = 5
	case EpochDharmaDecline:
		input.Description = "末法时代，灵气衰退"
		input.SpiritualBase = 0.2
		input.ResourceCap = 50
		input.EventFrequency = 0.2
		input.SecretRealmCount = 1
	case EpochChaos:
		input.Description = "混沌乱世，异变频发"
		input.SpiritualBase = 0.5
		input.ResourceCap = 80
		input.EventFrequency = 2.0
		input.SecretRealmCount = 3
	case EpochRenewal:
		input.Description = "万物复苏，新时代开启"
		input.SpiritualBase = 0.6
		input.ResourceCap = 100
		input.EventFrequency = 0.8
		input.SecretRealmCount = 4
	default:
		return nil, fmt.Errorf("未知的纪元阶段: %s", newPhase)
	}

	return s.StartNewEpoch(input, now)
}

// GetCurrentEpoch returns the current epoch.
func (s *EpochService) GetCurrentEpoch() (*WorldEpoch, bool) {
	if s.currentEpoch == nil {
		return nil, false
	}
	return s.currentEpoch, true
}

// GetEpochHistory returns all past epochs.
func (s *EpochService) GetEpochHistory() []*WorldEpoch {
	return s.history
}

// GetEpochDuration returns how long the current epoch has lasted.
func (s *EpochService) GetEpochDuration(now time.Time) (time.Duration, error) {
	if s.currentEpoch == nil {
		return 0, fmt.Errorf("当前没有活跃的纪元")
	}
	return now.Sub(s.currentEpoch.StartTime), nil
}

// GetCurrentPhase returns the current epoch phase.
func (s *EpochService) GetCurrentPhase() (EpochPhase, error) {
	if s.currentEpoch == nil {
		return "", fmt.Errorf("当前没有活跃的纪元")
	}
	return s.currentEpoch.Phase, nil
}
