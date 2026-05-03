package service

import (
	"math"
	"sort"
)

// WorldBalanceRule handles world health evaluation and balance adjustment.
type WorldBalanceRule struct {
	healthyGiniThreshold   float64 // max Gini coefficient for "healthy"
	unhealthyGiniThreshold float64 // threshold for "unhealthy"
	minSectDiversity       int     // minimum number of active sects
}

// NewWorldBalanceRule creates a new WorldBalanceRule with default configuration.
func NewWorldBalanceRule() *WorldBalanceRule {
	return &WorldBalanceRule{
		healthyGiniThreshold:   0.40,
		unhealthyGiniThreshold: 0.60,
		minSectDiversity:       3,
	}
}

// WorldMetrics holds aggregated world state data.
type WorldMetrics struct {
	EntitySpiritStones []float64 // spirit stone holdings per entity
	ResourceSpawnRate  float64   // current resource spawn rate
	ResourceConsumptionRate float64 // current resource consumption rate
	ActiveSectCount    int       // number of active sects
	SectSizes          []int     // member count per sect
	KarmaValues        []float64 // karma values across all entities
}

// WorldHealth represents the overall health assessment of the world.
type WorldHealth struct {
	GiniCoefficient    float64 // wealth inequality (0=perfect equality, 1=perfect inequality)
	ResourceRatio      float64 // spawn/consumption ratio
	SectDiversityScore float64 // diversity score based on sect count and size distribution
	KarmaStdDev        float64 // standard deviation of karma values
	OverallScore       float64 // composite health score (0-100)
	Status             string  // "healthy", "warning", "critical"
}

// EvaluateWorldHealth assesses the overall health of the cultivation world.
func (r *WorldBalanceRule) EvaluateWorldHealth(metrics WorldMetrics) WorldHealth {
	gini := r.calculateGiniCoefficient(metrics.EntitySpiritStones)
	resourceRatio := r.calculateResourceRatio(metrics.ResourceSpawnRate, metrics.ResourceConsumptionRate)
	sectDiversity := r.calculateSectDiversityScore(metrics.ActiveSectCount, metrics.SectSizes)
	karmaStdDev := r.calculateStdDev(metrics.KarmaValues)

	// Composite score (weighted average)
	// Lower Gini is better (invert), resource ratio near 1.0 is best
	giniScore := math.Max(0, 1.0-gini) * 100
	resourceScore := math.Max(0, 100-math.Abs(resourceRatio-1.0)*50)
	sectScore := math.Min(100, sectDiversity*100)
	karmaScore := math.Max(0, 100-karmaStdDev)

	overallScore := giniScore*0.25 + resourceScore*0.35 + sectScore*0.20 + karmaScore*0.20

	// Determine status
	var status string
	switch {
	case overallScore >= 70:
		status = "healthy"
	case overallScore >= 40:
		status = "warning"
	default:
		status = "critical"
	}

	return WorldHealth{
		GiniCoefficient:    gini,
		ResourceRatio:      resourceRatio,
		SectDiversityScore: sectDiversity,
		KarmaStdDev:        karmaStdDev,
		OverallScore:       overallScore,
		Status:             status,
	}
}

func (r *WorldBalanceRule) calculateGiniCoefficient(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	n := float64(len(sorted))
	sum := 0.0
	for _, v := range sorted {
		sum += v
	}
	if sum == 0 {
		return 0
	}

	// Gini formula: G = (2 * Σ(i * x_i)) / (n * Σx_i) - (n+1)/n
	weightedSum := 0.0
	for i, v := range sorted {
		weightedSum += float64(i+1) * v
	}

	gini := (2.0 * weightedSum) / (n * sum) - (n+1.0)/n
	return math.Max(0, math.Min(1, gini))
}

func (r *WorldBalanceRule) calculateResourceRatio(spawnRate, consumptionRate float64) float64 {
	if consumptionRate <= 0 {
		return 1.0 // no consumption, balanced
	}
	return spawnRate / consumptionRate
}

func (r *WorldBalanceRule) calculateSectDiversityScore(sectCount int, sizes []int) float64 {
	if sectCount == 0 {
		return 0
	}

	// Base score from count
	countScore := math.Min(1.0, float64(sectCount)/10.0)

	// Size diversity: use coefficient of variation
	if len(sizes) < 2 {
		return countScore * 0.5
	}

	mean := 0.0
	for _, s := range sizes {
		mean += float64(s)
	}
	mean /= float64(len(sizes))

	if mean == 0 {
		return countScore * 0.5
	}

	variance := 0.0
	for _, s := range sizes {
		diff := float64(s) - mean
		variance += diff * diff
	}
	variance /= float64(len(sizes))
	cv := math.Sqrt(variance) / mean

	// Lower CV = more balanced = better diversity score
	sizeScore := math.Max(0, 1.0-cv)

	return (countScore + sizeScore) / 2.0
}

func (r *WorldBalanceRule) calculateStdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	return math.Sqrt(variance)
}

// BalanceAdjustment represents a recommended adjustment to restore world balance.
type BalanceAdjustment struct {
	Type        string  // "resource_boost", "sect_promotion", "karma_redistribution", "wealth_tax"
	Target      string  // target entity or region
	Magnitude   float64 // adjustment magnitude
	Reason      string  // why this adjustment is needed
}

// ApplyBalanceAdjustment generates recommended adjustments based on world health.
func (r *WorldBalanceRule) ApplyBalanceAdjustment(health WorldHealth, metrics WorldMetrics) []BalanceAdjustment {
	var adjustments []BalanceAdjustment

	// High Gini: recommend wealth redistribution
	if health.GiniCoefficient > r.unhealthyGiniThreshold {
		adjustments = append(adjustments, BalanceAdjustment{
			Type:      "wealth_tax",
			Target:    "top_10_percent",
			Magnitude: health.GiniCoefficient * 0.1,
			Reason:    "财富差距过大，需要重新分配",
		})
	}

	// Resource shortage: boost spawn rate
	if health.ResourceRatio < 0.7 {
		adjustments = append(adjustments, BalanceAdjustment{
			Type:      "resource_boost",
			Target:    "all_regions",
			Magnitude: (0.7 - health.ResourceRatio) * 0.5,
			Reason:    "资源消耗大于产出，需要提升灵气密度",
		})
	}

	// Low sect diversity: promote smaller sects
	if metrics.ActiveSectCount < r.minSectDiversity {
		adjustments = append(adjustments, BalanceAdjustment{
			Type:      "sect_promotion",
			Target:    "new_or_small_sects",
			Magnitude: 0.3,
			Reason:    "宗门数量不足，需要扶持新宗门",
		})
	}

	// High karma variance: karma redistribution events
	if health.KarmaStdDev > 50 {
		adjustments = append(adjustments, BalanceAdjustment{
			Type:      "karma_redistribution",
			Target:    "karma_events",
			Magnitude: health.KarmaStdDev / 200.0,
			Reason:    "业力分布过于极端，需要天道干预",
		})
	}

	return adjustments
}
