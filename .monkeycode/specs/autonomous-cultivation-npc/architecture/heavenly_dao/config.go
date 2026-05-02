package heavenlydao

type HeavenlyDaoConfig struct {
	Karma        KarmaConfig
	Tribulation  TribulationConfig
	Breakthrough BreakthroughConfig
	Combat       CombatConfig
}

type KarmaConfig struct {
	Cap           float64
	DecayRateHour float64
	BaseValues    map[string]float64
}

type TribulationConfig struct {
	BaseProbabilityByRealm map[string]float64
	StrengthPerKarma       float64
	RecentWindowDays       int
	RecentStrengthBonus    float64
	MeritFloorFactor       float64
	MinProbability         float64
	MaxProbability         float64
}

type BreakthroughConfig struct {
	BaseSuccessByRealm       map[string]float64
	MinSuccessRate           float64
	MaxSuccessRate           float64
	FailureCultivationLoss   float64
	FailureCooldownPerRealm  int
	FailureMentalDamage      float64
}

type CombatConfig struct {
	RealmSuppressionPerLevel float64
	BaseCritRate             float64
	BaseCritDamage           float64
	ElementCounters          map[string]float64
}

type ConfigLoader interface {
	Load() (*HeavenlyDaoConfig, error)
}

type StaticConfigLoader struct {
	Config *HeavenlyDaoConfig
}

func (l *StaticConfigLoader) Load() (*HeavenlyDaoConfig, error) {
	return l.Config, nil
}
