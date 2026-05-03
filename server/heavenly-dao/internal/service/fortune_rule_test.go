package service

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewFortuneRule(t *testing.T) {
    rule := NewFortuneRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 0.01, rule.baseOpportunityProb)
}

func TestCalculateFortuneScore_Baseline(t *testing.T) {
    rule := NewFortuneRule()

    input := FortuneInput{
        Luck:    50,
        Merit:   0,
        Karma:   0,
    }

    score := rule.CalculateFortuneScore(input)
    assert.Equal(t, 50.0, score)
}

func TestCalculateFortuneScore_WithMerit(t *testing.T) {
    rule := NewFortuneRule()

    input := FortuneInput{
        Luck:  50,
        Merit: 1000,
        Karma: 0,
    }

    // fortune = 50 + 1000*0.1 = 50 + 100 = 150
    score := rule.CalculateFortuneScore(input)
    assert.Equal(t, 150.0, score)
}

func TestCalculateFortuneScore_WithKarma(t *testing.T) {
    rule := NewFortuneRule()

    input := FortuneInput{
        Luck:  50,
        Merit: 0,
        Karma: 500,
    }

    // fortune = 50 - 500*0.05 = 50 - 25 = 25
    score := rule.CalculateFortuneScore(input)
    assert.Equal(t, 25.0, score)
}

func TestCalculateFortuneScore_WithActions(t *testing.T) {
    rule := NewFortuneRule()

    input := FortuneInput{
        Luck:              50,
        SaveLifeCount:     3,  // +15
        DestroyOppCount:   1,  // -10
        HeavenlyTaskCount: 2,  // +40
    }

    // fortune = 50 + 15 - 10 + 40 = 95
    score := rule.CalculateFortuneScore(input)
    assert.Equal(t, 95.0, score)
}

func TestCalculateFortuneScore_NegativeClamp(t *testing.T) {
    rule := NewFortuneRule()

    input := FortuneInput{
        Luck:    10,
        Karma:   1000, // -50
    }

    // fortune = 10 - 50 = -40 → clamped to 0
    score := rule.CalculateFortuneScore(input)
    assert.Equal(t, 0.0, score)
}

func TestCalculateFortuneScore_AllFactors(t *testing.T) {
    rule := NewFortuneRule()

    input := FortuneInput{
        Luck:              80,
        Merit:             500,  // +50
        Karma:             200,  // -10
        SaveLifeCount:     2,    // +10
        DestroyOppCount:   0,
        HeavenlyTaskCount: 1,    // +20
    }

    // fortune = 80 + 50 - 10 + 10 + 20 = 150
    score := rule.CalculateFortuneScore(input)
    assert.Equal(t, 150.0, score)
}

func TestGetFortuneGrade_AllTiers(t *testing.T) {
    rule := NewFortuneRule()

    tests := []struct {
        score  float64
        grade  FortuneGrade
        mult   float64
    }{
        {0, FortuneMisfortune, 0.5},
        {19.9, FortuneMisfortune, 0.5},
        {20, FortunePlain, 0.8},
        {39.9, FortunePlain, 0.8},
        {40, FortuneNormal, 1.0},
        {59.9, FortuneNormal, 1.0},
        {60, FortuneLucky, 1.3},
        {79.9, FortuneLucky, 1.3},
        {80, FortuneDestined, 2.0},
        {100, FortuneDestined, 2.0},
    }

    for _, tc := range tests {
        grade := rule.GetFortuneGrade(tc.score)
        assert.Equal(t, tc.grade, grade, "score: %.1f", tc.score)

        mult := rule.GetFortuneMultiplier(tc.score)
        assert.Equal(t, tc.mult, mult, "score: %.1f", tc.score)
    }
}

func TestCalculateOpportunityProbability_Baseline(t *testing.T) {
    rule := NewFortuneRule()

    input := OpportunityInput{
        Fortune:           50,  // fortune_factor = 1.5
        SpiritualDensity:  50,  // spiritual = 0.5
        ExplorationFactor: 1.0,
        BalanceFactor:     1.0,
    }

    // prob = 0.01 * 1.5 * 0.5 * 1.0 * 1.0 = 0.0075
    prob := rule.CalculateOpportunityProbability(input)
    assert.InDelta(t, 0.0075, prob, 0.0001)
}

func TestCalculateOpportunityProbability_HighFortune(t *testing.T) {
    rule := NewFortuneRule()

    input := OpportunityInput{
        Fortune:           100, // fortune_factor = 2.0
        SpiritualDensity:  80,  // spiritual = 0.8
        ExplorationFactor: 1.0,
        BalanceFactor:     1.0,
    }

    // prob = 0.01 * 2.0 * 0.8 = 0.016
    prob := rule.CalculateOpportunityProbability(input)
    assert.InDelta(t, 0.016, prob, 0.0001)
}

func TestCalculateOpportunityProbability_BalanceModifier(t *testing.T) {
    rule := NewFortuneRule()

    input := OpportunityInput{
        Fortune:           50,
        SpiritualDensity:  50,
        ExplorationFactor: 1.0,
        BalanceFactor:     1.5, // max balance boost
    }

    // prob = 0.01 * 1.5 * 0.5 * 1.0 * 1.5 = 0.01125
    prob := rule.CalculateOpportunityProbability(input)
    assert.InDelta(t, 0.01125, prob, 0.0001)
}

func TestCalculateOpportunityProbability_ZeroSpiritual(t *testing.T) {
    rule := NewFortuneRule()

    input := OpportunityInput{
        Fortune:           100,
        SpiritualDensity:  0,
        ExplorationFactor: 1.0,
        BalanceFactor:     1.0,
    }

    prob := rule.CalculateOpportunityProbability(input)
    assert.Equal(t, 0.0, prob) // spiritual = 0 → prob = 0
}

func TestCalculateOpportunityProbability_ClampedFactors(t *testing.T) {
    rule := NewFortuneRule()

    // Negative balance factor → clamped to 0.5
    input := OpportunityInput{
        Fortune:           50,
        SpiritualDensity:  50,
        ExplorationFactor: 1.0,
        BalanceFactor:     -1.0,
    }

    prob := rule.CalculateOpportunityProbability(input)
    // prob = 0.01 * 1.5 * 0.5 * 1.0 * 0.5 = 0.00375
    assert.InDelta(t, 0.00375, prob, 0.0001)

    // Excessive balance factor → clamped to 1.5
    input.BalanceFactor = 3.0
    prob = rule.CalculateOpportunityProbability(input)
    assert.InDelta(t, 0.01125, prob, 0.0001)
}

func TestResolveOpportunityType_Distribution(t *testing.T) {
    rule := NewFortuneRule()

    // The implementation does: r = randFloat() * 100, then compares against cumulative weights
    // Cumulative: herb=40, elder=45, epiphany=55, root=56, bless=58, realm=61, noble=76, nothing=100
    tests := []struct {
        randVal  float64
        expected OpportunityType
    }{
        {0.00, OpHerbDiscovery},    // r=0: herb (0-40)
        {0.30, OpHerbDiscovery},    // r=30: herb
        {0.39, OpHerbDiscovery},    // r=39: herb
        {0.40, OpElderEstate},      // r=40: elder (40-45)
        {0.44, OpElderEstate},      // r=44: elder
        {0.45, OpEpiphany},         // r=45: epiphany (45-55)
        {0.50, OpEpiphany},         // r=50: epiphany
        {0.55, OpRootAwakening},    // r=55: root (55-56)
        {0.56, OpHeavenlyBless},    // r=56: blessing (56-58)
        {0.57, OpHeavenlyBless},    // r=57: blessing
        {0.585, OpSecretRealm},     // r=58.5: realm (58-61)
        {0.60, OpSecretRealm},      // r=60: realm
        {0.61, OpNobleHelp},        // r=61: noble (61-76)
        {0.70, OpNobleHelp},        // r=70: noble
        {0.76, OpNothing},          // r=76: nothing (76-100)
        {0.90, OpNothing},          // r=90: nothing
    }

    for _, tc := range tests {
        result := rule.ResolveOpportunityType(func() float64 { return tc.randVal })
        assert.Equal(t, tc.expected, result, "randVal: %.3f", tc.randVal)
    }
}

func TestCalculateOpportunityQualityTier(t *testing.T) {
    rule := NewFortuneRule()

    tests := []struct {
        fortune float64
        tier    int
    }{
        {0, 1},
        {19.9, 1},
        {20, 1},
        {39.9, 1},
        {40, 2},
        {59.9, 2},
        {60, 3},
        {79.9, 3},
        {80, 4},
        {99.9, 4},
        {100, 5},
        {200, 5},
    }

    for _, tc := range tests {
        tier := rule.CalculateOpportunityQualityTier(tc.fortune)
        assert.Equal(t, tc.tier, tier, "fortune: %.1f", tc.fortune)
    }
}

func TestDetermineDestiny_Default(t *testing.T) {
    rule := NewFortuneRule()

    input := DestinyInput{
        HasMaster: true,
        HasSpouse: false,
        HasSect:   true,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyMortal, destiny)
}

func TestDetermineDestiny_FollowHeaven(t *testing.T) {
    rule := NewFortuneRule()

    input := DestinyInput{
        Merit:     5001,
        HasMaster: true,
        HasSect:   true,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyFollowHeaven, destiny)
}

func TestDetermineDestiny_KillBreakWolf(t *testing.T) {
    rule := NewFortuneRule()

    input := DestinyInput{
        KillCount: 101,
        HasSect:   true,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyKillBreakWolf, destiny)
}

func TestDetermineDestiny_HeavenSecret(t *testing.T) {
    rule := NewFortuneRule()

    input := DestinyInput{
        SecretRealmDiscovered: 5,
        HasSect:               true,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyHeavenSecret, destiny)
}

func TestDetermineDestiny_DefyHeaven(t *testing.T) {
    rule := NewFortuneRule()

    input := DestinyInput{
        ConsecutiveBreakthroughs: 3,
        HasSect:                  true,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyDefyHeaven, destiny)
}

func TestDetermineDestiny_LoneStar(t *testing.T) {
    rule := NewFortuneRule()

    input := DestinyInput{
        HasMaster: false,
        HasSpouse: false,
        HasSect:   false,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyLoneStar, destiny)
}

func TestDetermineDestiny_Priority(t *testing.T) {
    rule := NewFortuneRule()

    // Merit > 5000 AND kills > 100 → should be FollowHeaven (higher priority)
    input := DestinyInput{
        Merit:     6000,
        KillCount: 200,
        HasSect:   true,
    }

    destiny := rule.DetermineDestiny(input)
    assert.Equal(t, DestinyFollowHeaven, destiny)
}

func TestGetDestinyEffects_AllDestinies(t *testing.T) {
    rule := NewFortuneRule()

    tests := []struct {
        destiny      Destiny
        breakthrough float64
        tribulation  float64
        combat       float64
        karma        float64
        opportunity  float64
        cultivation  float64
        mental       float64
    }{
        {DestinyMortal, 0, 0, 0, 0, 0, 0, 0},
        {DestinyDefyHeaven, 0.20, 0, 0, 0, 0, 0, 0},
        {DestinyFollowHeaven, 0, 0.30, 0, 0, 0, 0, 0},
        {DestinyKillBreakWolf, 0, 0, 0.30, 0.50, 0, 0, 0},
        {DestinyHeavenSecret, 0, 0, 0, 0, 0.50, 0, 0},
        {DestinyLoneStar, 0, 0, 0, 0, 0, 0.20, 0.20},
    }

    for _, tc := range tests {
        effects := rule.GetDestinyEffects(tc.destiny)
        assert.Equal(t, tc.breakthrough, effects.BreakthroughBonus, "%s breakthrough", tc.destiny)
        assert.Equal(t, tc.tribulation, effects.TribulationReduction, "%s tribulation", tc.destiny)
        assert.Equal(t, tc.combat, effects.CombatBonus, "%s combat", tc.destiny)
        assert.Equal(t, tc.karma, effects.KarmaPenalty, "%s karma", tc.destiny)
        assert.Equal(t, tc.opportunity, effects.OpportunityBonus, "%s opportunity", tc.destiny)
        assert.Equal(t, tc.cultivation, effects.CultivationBonus, "%s cultivation", tc.destiny)
        assert.Equal(t, tc.mental, effects.MentalPenalty, "%s mental", tc.destiny)
    }
}
