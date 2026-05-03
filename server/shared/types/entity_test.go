package types

import (
	"testing"
	"time"
)

func TestGenerateEntityID(t *testing.T) {
	id := GenerateEntityID()
	if id == "" {
		t.Error("GenerateEntityID() returned empty string")
	}

	id2 := GenerateEntityID()
	if id == id2 {
		t.Error("GenerateEntityID() returned duplicate ID")
	}
}

func TestGenerateOperationID(t *testing.T) {
	id := GenerateOperationID()
	if id == "" {
		t.Error("GenerateOperationID() returned empty string")
	}
}

func TestEntityInitialization(t *testing.T) {
	id := GenerateEntityID()
	now := time.Now()

	entity := &Entity{
		ID: id,
		EntityType: EntityTypePlayer,
		Name: "TestPlayer",
		Realm: RealmMortal,
		Position: WorldPosition{
			RegionID: "test",
			X:        0,
			Y:        0,
		},
		Attributes: Attributes{
			Qi:    100,
			MaxQi: 100,
		},
		Status:    StatusNormal,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if entity.ID != id {
		t.Errorf("Entity.ID mismatch, expected %s, got %s", id, entity.ID)
	}

	if entity.EntityType != EntityTypePlayer {
		t.Errorf("Entity.EntityType mismatch, expected %s, got %s", EntityTypePlayer, entity.EntityType)
	}

	if entity.Realm != RealmMortal {
		t.Errorf("Entity.Realm mismatch, expected %s, got %s", RealmMortal, entity.Realm)
	}
}

func TestAttributesInitialization(t *testing.T) {
	attr := Attributes{
		Qi:    50,
		MaxQi: 100,
	}

	if attr.Qi != 50 {
		t.Errorf("Qi mismatch, expected 50, got %f", attr.Qi)
	}

	if attr.MaxQi != 100 {
		t.Errorf("MaxQi mismatch, expected 100, got %f", attr.MaxQi)
	}
}

func TestKarmaInitialization(t *testing.T) {
	karma := Karma{
		KarmaValue:   10,
		Merit:        5,
		HeavenlyMark: "clear",
	}

	if karma.KarmaValue != 10 {
		t.Errorf("KarmaValue mismatch, expected 10, got %d", karma.KarmaValue)
	}

	if karma.Merit != 5 {
		t.Errorf("Merit mismatch, expected 5, got %d", karma.Merit)
	}
}

func TestSpiritualRootInitialization(t *testing.T) {
	root := SpiritualRoot{
		Element: "fire",
		Purity:  80,
	}

	if root.Element != "fire" {
		t.Errorf("SpiritualRoot.Element mismatch, expected fire, got %s", root.Element)
	}

	if root.Purity != 80 {
		t.Errorf("SpiritualRoot.Purity mismatch, expected 80, got %d", root.Purity)
	}
}

func TestAttributesExtendedInitialization(t *testing.T) {
	attr := Attributes{
		// Basic attributes
		Age:        20,
		Gender:     "male",
		Appearance: 75,
		Charisma:   60,

		// Combat attributes
		CritRate:        5.0,
		CritDamage:      150.0,
		DodgeRate:       3.0,
		HitRate:         95.0,
		Penetration:     10.0,
		DamageReduction: 15.0,

		// Spiritual roots
		SpiritualRoots: []SpiritualRoot{
			{Element: "fire", Purity: 80},
			{Element: "wood", Purity: 50},
		},
		RootPurity:   70,
		RootAwakened: false,

		// Mental state
		MentalStability:      80,
		ObsessionCount:       2,
		DaoHeart:             45,
		InnerDemonResistance: 60,
		Enlightenment:        10,

		// Life skills
		AlchemyLevel:    3,
		ArtificingLevel: 2,
		FormationLevel:  1,
		FireControl:     40,
		HerbKnowledge:   55,
		MiningSkill:     30,
		TalismanSkill:   2,
		BeastTaming:     1,

		// Social attributes
		Reputation:        100,
		SectContribution:  50,
		FactionStandings:  map[string]int{"qingyun_sect": 80, "blood_sect": -20},
		RelationshipCount: 5,
		DiscipleIDs:       []string{"disciple_1"},
		SwornSiblings:     []string{"sibling_1", "sibling_2"},
		Enemies:           []string{"enemy_1"},
		Lovers:            []string{},

		// Wealth attributes
		SpiritStones: SpiritStones{
			LowGrade:     1000,
			MediumGrade:  50,
			HighGrade:    2,
			PremiumGrade: 0,
		},
		PropertyValue:  5000,
		RealEstate:     []string{"cave_001"},
		BusinessIncome: 0,

		// Special attributes
		Bloodline:        "mortal",
		BloodlinePurity:  100,
		Physique:         "none",
		PhysiqueAwakened: false,
		Destiny:          50,
		WorldFavor:       0,

		// Law attributes
		Laws:           map[string]float64{"fire": 5.0, "metal": 2.0},
		LawResonance:   0,
		DomainPower:    0,
		DomainRange:    0,
		LawSuppression: 0,

		// Dao attributes
		DaoSeedType:          "none",
		DaoSeedLevel:         0,
		DaoSeedGrowth:        0,
		DaoMarks:             0,
		DaoHeartComprehension: 0,
		DestinyPath:          "mortal_path",

		// Lifespan attributes
		RemainingLifespan: 80,
		MaxLifespan:       80,
		AgingPenalty:      0,

		// Status effects
		Injuries:    []Injury{},
		Buffs:       []Buff{},
		Debuffs:     []Debuff{},
		PoisonLevel: 0,
		CurseLevel:  0,
	}

	// Verify basic attributes
	if attr.Age != 20 {
		t.Errorf("Age mismatch, expected 20, got %d", attr.Age)
	}
	if attr.Gender != "male" {
		t.Errorf("Gender mismatch, expected male, got %s", attr.Gender)
	}
	if attr.Appearance != 75 {
		t.Errorf("Appearance mismatch, expected 75, got %d", attr.Appearance)
	}

	// Verify combat attributes
	if attr.CritRate != 5.0 {
		t.Errorf("CritRate mismatch, expected 5.0, got %f", attr.CritRate)
	}
	if attr.CritDamage != 150.0 {
		t.Errorf("CritDamage mismatch, expected 150.0, got %f", attr.CritDamage)
	}
	if attr.DodgeRate != 3.0 {
		t.Errorf("DodgeRate mismatch, expected 3.0, got %f", attr.DodgeRate)
	}
	if attr.HitRate != 95.0 {
		t.Errorf("HitRate mismatch, expected 95.0, got %f", attr.HitRate)
	}
	if attr.Penetration != 10.0 {
		t.Errorf("Penetration mismatch, expected 10.0, got %f", attr.Penetration)
	}
	if attr.DamageReduction != 15.0 {
		t.Errorf("DamageReduction mismatch, expected 15.0, got %f", attr.DamageReduction)
	}

	// Verify spiritual roots
	if len(attr.SpiritualRoots) != 2 {
		t.Errorf("SpiritualRoots length mismatch, expected 2, got %d", len(attr.SpiritualRoots))
	}
	if attr.SpiritualRoots[0].Element != "fire" {
		t.Errorf("SpiritualRoots[0].Element mismatch, expected fire, got %s", attr.SpiritualRoots[0].Element)
	}
	if attr.RootPurity != 70 {
		t.Errorf("RootPurity mismatch, expected 70, got %d", attr.RootPurity)
	}
	if attr.RootAwakened != false {
		t.Errorf("RootAwakened mismatch, expected false, got %v", attr.RootAwakened)
	}

	// Verify mental state
	if attr.MentalStability != 80 {
		t.Errorf("MentalStability mismatch, expected 80, got %d", attr.MentalStability)
	}
	if attr.ObsessionCount != 2 {
		t.Errorf("ObsessionCount mismatch, expected 2, got %d", attr.ObsessionCount)
	}
	if attr.DaoHeart != 45 {
		t.Errorf("DaoHeart mismatch, expected 45, got %d", attr.DaoHeart)
	}
	if attr.InnerDemonResistance != 60 {
		t.Errorf("InnerDemonResistance mismatch, expected 60, got %d", attr.InnerDemonResistance)
	}

	// Verify life skills
	if attr.AlchemyLevel != 3 {
		t.Errorf("AlchemyLevel mismatch, expected 3, got %d", attr.AlchemyLevel)
	}
	if attr.ArtificingLevel != 2 {
		t.Errorf("ArtificingLevel mismatch, expected 2, got %d", attr.ArtificingLevel)
	}
	if attr.FormationLevel != 1 {
		t.Errorf("FormationLevel mismatch, expected 1, got %d", attr.FormationLevel)
	}
	if attr.FireControl != 40 {
		t.Errorf("FireControl mismatch, expected 40, got %d", attr.FireControl)
	}
	if attr.HerbKnowledge != 55 {
		t.Errorf("HerbKnowledge mismatch, expected 55, got %d", attr.HerbKnowledge)
	}
	if attr.MiningSkill != 30 {
		t.Errorf("MiningSkill mismatch, expected 30, got %d", attr.MiningSkill)
	}
	if attr.TalismanSkill != 2 {
		t.Errorf("TalismanSkill mismatch, expected 2, got %d", attr.TalismanSkill)
	}
	if attr.BeastTaming != 1 {
		t.Errorf("BeastTaming mismatch, expected 1, got %d", attr.BeastTaming)
	}

	// Verify social attributes
	if attr.Reputation != 100 {
		t.Errorf("Reputation mismatch, expected 100, got %d", attr.Reputation)
	}
	if attr.SectContribution != 50 {
		t.Errorf("SectContribution mismatch, expected 50, got %d", attr.SectContribution)
	}
	if attr.FactionStandings["qingyun_sect"] != 80 {
		t.Errorf("FactionStandings qingyun_sect mismatch, expected 80, got %d", attr.FactionStandings["qingyun_sect"])
	}
	if attr.FactionStandings["blood_sect"] != -20 {
		t.Errorf("FactionStandings blood_sect mismatch, expected -20, got %d", attr.FactionStandings["blood_sect"])
	}
	if attr.RelationshipCount != 5 {
		t.Errorf("RelationshipCount mismatch, expected 5, got %d", attr.RelationshipCount)
	}
	if len(attr.DiscipleIDs) != 1 || attr.DiscipleIDs[0] != "disciple_1" {
		t.Errorf("DiscipleIDs mismatch")
	}
	if len(attr.Enemies) != 1 || attr.Enemies[0] != "enemy_1" {
		t.Errorf("Enemies mismatch")
	}

	// Verify wealth attributes
	if attr.SpiritStones.LowGrade != 1000 {
		t.Errorf("SpiritStones.LowGrade mismatch, expected 1000, got %d", attr.SpiritStones.LowGrade)
	}
	if attr.SpiritStones.MediumGrade != 50 {
		t.Errorf("SpiritStones.MediumGrade mismatch, expected 50, got %d", attr.SpiritStones.MediumGrade)
	}
	if attr.SpiritStones.HighGrade != 2 {
		t.Errorf("SpiritStones.HighGrade mismatch, expected 2, got %d", attr.SpiritStones.HighGrade)
	}
	if attr.SpiritStones.PremiumGrade != 0 {
		t.Errorf("SpiritStones.PremiumGrade mismatch, expected 0, got %d", attr.SpiritStones.PremiumGrade)
	}
	if attr.PropertyValue != 5000 {
		t.Errorf("PropertyValue mismatch, expected 5000, got %d", attr.PropertyValue)
	}

	// Verify special attributes
	if attr.Bloodline != "mortal" {
		t.Errorf("Bloodline mismatch, expected mortal, got %s", attr.Bloodline)
	}
	if attr.BloodlinePurity != 100 {
		t.Errorf("BloodlinePurity mismatch, expected 100, got %d", attr.BloodlinePurity)
	}
	if attr.Physique != "none" {
		t.Errorf("Physique mismatch, expected none, got %s", attr.Physique)
	}
	if attr.Destiny != 50 {
		t.Errorf("Destiny mismatch, expected 50, got %d", attr.Destiny)
	}

	// Verify law attributes
	if len(attr.Laws) != 2 {
		t.Errorf("Laws length mismatch, expected 2, got %d", len(attr.Laws))
	}
	if attr.Laws["fire"] != 5.0 {
		t.Errorf("Laws fire mismatch, expected 5.0, got %f", attr.Laws["fire"])
	}

	// Verify dao attributes
	if attr.DaoSeedType != "none" {
		t.Errorf("DaoSeedType mismatch, expected none, got %s", attr.DaoSeedType)
	}
	if attr.DaoSeedLevel != 0 {
		t.Errorf("DaoSeedLevel mismatch, expected 0, got %d", attr.DaoSeedLevel)
	}

	// Verify lifespan attributes
	if attr.RemainingLifespan != 80 {
		t.Errorf("RemainingLifespan mismatch, expected 80, got %d", attr.RemainingLifespan)
	}
	if attr.AgingPenalty != 0 {
		t.Errorf("AgingPenalty mismatch, expected 0, got %f", attr.AgingPenalty)
	}

	// Verify status effects
	if attr.PoisonLevel != 0 {
		t.Errorf("PoisonLevel mismatch, expected 0, got %d", attr.PoisonLevel)
	}
	if attr.CurseLevel != 0 {
		t.Errorf("CurseLevel mismatch, expected 0, got %d", attr.CurseLevel)
	}
}

func TestInjuryCreation(t *testing.T) {
	injury := Injury{
		Type:        "internal",
		Severity:    3,
		Cause:       "combat",
		HealTime:    time.Now().Add(24 * time.Hour).Unix(),
		Description: "damaged meridians",
	}

	if injury.Type != "internal" {
		t.Errorf("Injury.Type mismatch, expected internal, got %s", injury.Type)
	}
	if injury.Severity != 3 {
		t.Errorf("Injury.Severity mismatch, expected 3, got %d", injury.Severity)
	}
	if injury.Cause != "combat" {
		t.Errorf("Injury.Cause mismatch, expected combat, got %s", injury.Cause)
	}
	if injury.HealTime == 0 {
		t.Error("Injury.HealTime should be set")
	}
}

func TestBuffCreation(t *testing.T) {
	buff := Buff{
		Name:       "spirit_boost",
		Effect:     "increase_spiritual_power",
		Value:      20.0,
		Source:     "pill",
		ExpiryTime: time.Now().Add(2 * time.Hour).Unix(),
	}

	if buff.Name != "spirit_boost" {
		t.Errorf("Buff.Name mismatch, expected spirit_boost, got %s", buff.Name)
	}
	if buff.Effect != "increase_spiritual_power" {
		t.Errorf("Buff.Effect mismatch, expected increase_spiritual_power, got %s", buff.Effect)
	}
	if buff.Value != 20.0 {
		t.Errorf("Buff.Value mismatch, expected 20.0, got %f", buff.Value)
	}
	if buff.Source != "pill" {
		t.Errorf("Buff.Source mismatch, expected pill, got %s", buff.Source)
	}
}

func TestDebuffCreation(t *testing.T) {
	debuff := Debuff{
		Name:       "poison",
		Effect:     "reduce_qi",
		Value:      5.0,
		Source:     "herb_trap",
		ExpiryTime: time.Now().Add(30 * time.Minute).Unix(),
	}

	if debuff.Name != "poison" {
		t.Errorf("Debuff.Name mismatch, expected poison, got %s", debuff.Name)
	}
	if debuff.Effect != "reduce_qi" {
		t.Errorf("Debuff.Effect mismatch, expected reduce_qi, got %s", debuff.Effect)
	}
	if debuff.Value != 5.0 {
		t.Errorf("Debuff.Value mismatch, expected 5.0, got %f", debuff.Value)
	}
}

func TestKarmaExtendedFields(t *testing.T) {
	karma := Karma{
		KarmaValue:   200,
		Merit:        50,
		KarmicDebt:   10,
		HeavenlyMark: "slight",
	}

	if karma.KarmicDebt != 10 {
		t.Errorf("KarmicDebt mismatch, expected 10, got %d", karma.KarmicDebt)
	}
}

func TestEntityStatusConstants(t *testing.T) {
	statuses := map[EntityStatus]bool{
		StatusNormal:     true,
		StatusCultivating: true,
		StatusCombat:     true,
		StatusResting:    true,
		StatusDead:       true,
		StatusExploring:  true,
		StatusCrafting:   true,
		StatusMeditating: true,
	}

	for s, expected := range statuses {
		if s == "" && expected {
			t.Errorf("Status constant %v should not be empty", s)
		}
	}
}

func TestCultivationRealmConstants(t *testing.T) {
	expectedRealms := []CultivationRealm{
		RealmMortal,
		RealmQiCondensation,
		RealmFoundation,
		RealmGoldenCore,
		RealmNascentSoul,
		RealmSoulTransform,
		RealmVoidRefinement,
		RealmIntegration,
		RealmMahayana,
		RealmTribulation,
	}

	for _, realm := range expectedRealms {
		if realm == "" {
			t.Errorf("Realm constant should not be empty")
		}
	}

	if len(expectedRealms) != 10 {
		t.Errorf("Expected 10 realms, got %d", len(expectedRealms))
	}
}

func TestSpiritualRootsList(t *testing.T) {
	roots := []SpiritualRoot{
		{Element: "gold", Purity: 60},
		{Element: "wood", Purity: 70},
		{Element: "water", Purity: 80},
		{Element: "fire", Purity: 90},
		{Element: "earth", Purity: 50},
	}

	if len(roots) != 5 {
		t.Errorf("Expected 5 spiritual roots, got %d", len(roots))
	}

	for _, root := range roots {
		if root.Element == "" {
			t.Error("Spiritual root element should not be empty")
		}
		if root.Purity < 0 || root.Purity > 100 {
			t.Errorf("Spiritual root purity %d should be between 0 and 100", root.Purity)
		}
	}
}

func TestAttributesDefaultValues(t *testing.T) {
	attr := Attributes{}

	// All numeric fields should default to zero
	if attr.Qi != 0 {
		t.Errorf("Default Qi should be 0, got %f", attr.Qi)
	}
	if attr.MaxQi != 0 {
		t.Errorf("Default MaxQi should be 0, got %f", attr.MaxQi)
	}
	if attr.CritRate != 0 {
		t.Errorf("Default CritRate should be 0, got %f", attr.CritRate)
	}
	if attr.RootPurity != 0 {
		t.Errorf("Default RootPurity should be 0, got %d", attr.RootPurity)
	}
	if attr.MentalStability != 0 {
		t.Errorf("Default MentalStability should be 0, got %d", attr.MentalStability)
	}
	if attr.AlchemyLevel != 0 {
		t.Errorf("Default AlchemyLevel should be 0, got %d", attr.AlchemyLevel)
	}
	if attr.Reputation != 0 {
		t.Errorf("Default Reputation should be 0, got %d", attr.Reputation)
	}
	if attr.SpiritStones.LowGrade != 0 {
		t.Errorf("Default LowGrade should be 0, got %d", attr.SpiritStones.LowGrade)
	}
	if attr.BloodlinePurity != 0 {
		t.Errorf("Default BloodlinePurity should be 0, got %d", attr.BloodlinePurity)
	}
	if attr.DomainPower != 0 {
		t.Errorf("Default DomainPower should be 0, got %f", attr.DomainPower)
	}
	if attr.DaoSeedLevel != 0 {
		t.Errorf("Default DaoSeedLevel should be 0, got %d", attr.DaoSeedLevel)
	}
	if attr.RemainingLifespan != 0 {
		t.Errorf("Default RemainingLifespan should be 0, got %d", attr.RemainingLifespan)
	}
	if attr.PoisonLevel != 0 {
		t.Errorf("Default PoisonLevel should be 0, got %d", attr.PoisonLevel)
	}

	// Slice and map fields should be nil/empty
	if attr.SpiritualRoots != nil {
		t.Error("Default SpiritualRoots should be nil")
	}
	if attr.FactionStandings != nil {
		t.Error("Default FactionStandings should be nil")
	}
	if attr.Laws != nil {
		t.Error("Default Laws should be nil")
	}
	if attr.DiscipleIDs != nil {
		t.Error("Default DiscipleIDs should be nil")
	}
	if attr.Injuries != nil {
		t.Error("Default Injuries should be nil")
	}
	if attr.Buffs != nil {
		t.Error("Default Buffs should be nil")
	}
	if attr.Debuffs != nil {
		t.Error("Default Debuffs should be nil")
	}
}

func TestSpiritStonesGrades(t *testing.T) {
	stones := SpiritStones{
		LowGrade:     10000,
		MediumGrade:  100,
		HighGrade:    1,
		PremiumGrade: 0,
	}

	if stones.LowGrade != 10000 {
		t.Errorf("LowGrade mismatch, expected 10000, got %d", stones.LowGrade)
	}
	if stones.MediumGrade != 100 {
		t.Errorf("MediumGrade mismatch, expected 100, got %d", stones.MediumGrade)
	}
	if stones.HighGrade != 1 {
		t.Errorf("HighGrade mismatch, expected 1, got %d", stones.HighGrade)
	}
	if stones.PremiumGrade != 0 {
		t.Errorf("PremiumGrade mismatch, expected 0, got %d", stones.PremiumGrade)
	}
}

func TestCompleteEntityWithAllAttributes(t *testing.T) {
	id := GenerateEntityID()
	now := time.Now()

	entity := &Entity{
		ID:         id,
		EntityType: EntityTypeNPC,
		Name:       "TestNPC",
		Realm:      RealmQiCondensation,
		Position: WorldPosition{
			RegionID: "qingyun_town",
			X:        100,
			Y:        200,
		},
		Attributes: Attributes{
			Age:        25,
			Gender:     "female",
			Appearance: 85,
			Charisma:   70,
			Qi:         150,
			MaxQi:      200,
			SpiritualPower:    120,
			MaxSpiritualPower: 200,
			DivineSense:       15,
			Comprehension:     65,
			Constitution:      55,
			Luck:              40,
			CultivationProgress: 25.5,
			CritRate:        8.0,
			CritDamage:      180.0,
			DodgeRate:       5.0,
			HitRate:         90.0,
			Penetration:     15.0,
			DamageReduction: 10.0,
			SpiritualRoots: []SpiritualRoot{
				{Element: "water", Purity: 75},
			},
			RootPurity:   75,
			RootAwakened: false,
			MutatedRoot:  "",
			MentalStability:      90,
			ObsessionCount:       1,
			DaoHeart:             30,
			InnerDemonResistance: 50,
			Enlightenment:        5,
			AlchemyLevel:    2,
			ArtificingLevel: 1,
			FormationLevel:  1,
			FireControl:     30,
			HerbKnowledge:   40,
			MiningSkill:     10,
			TalismanSkill:   1,
			BeastTaming:     0,
			Reputation:        50,
			SectContribution:  20,
			FactionStandings:  map[string]int{"qingyun_sect": 60},
			RelationshipCount: 2,
			MentorID:          "mentor_001",
			DiscipleIDs:       []string{},
			SwornSiblings:     []string{"sibling_a"},
			Enemies:           []string{"enemy_x"},
			Lovers:            []string{},
			SpiritStones: SpiritStones{
				LowGrade:     5000,
				MediumGrade:  20,
				HighGrade:    0,
				PremiumGrade: 0,
			},
			PropertyValue:  2000,
			RealEstate:     []string{},
			BusinessIncome: 0,
			Bloodline:        "mortal",
			BloodlinePurity:  100,
			Physique:         "none",
			PhysiqueAwakened: false,
			Destiny:          50,
			WorldFavor:       0,
			Laws:           map[string]float64{"water": 3.0},
			LawResonance:   0,
			DomainPower:    0,
			DomainRange:    0,
			LawSuppression: 0,
			DaoSeedType:          "none",
			DaoSeedLevel:         0,
			DaoSeedGrowth:        0,
			DaoMarks:             0,
			DaoHeartComprehension: 0,
			DestinyPath:          "mortal_path",
			RemainingLifespan: 120,
			MaxLifespan:       120,
			AgingPenalty:      0,
			Injuries:          []Injury{},
			Buffs:             []Buff{},
			Debuffs:           []Debuff{},
			PoisonLevel:       0,
			CurseLevel:        0,
		},
		Karma: Karma{
			KarmaValue:   0,
			Merit:        10,
			KarmicDebt:   0,
			HeavenlyMark: "clear",
		},
		Status:    StatusNormal,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if entity.ID != id {
		t.Errorf("Entity.ID mismatch")
	}
	if entity.EntityType != EntityTypeNPC {
		t.Errorf("EntityType should be npc")
	}
	if entity.Realm != RealmQiCondensation {
		t.Errorf("Realm mismatch, expected qi_condensation, got %s", entity.Realm)
	}
	if entity.Attributes.Age != 25 {
		t.Errorf("Age mismatch, expected 25, got %d", entity.Attributes.Age)
	}
	if entity.Attributes.SpiritualRoots[0].Element != "water" {
		t.Errorf("First spiritual root should be water")
	}
	if entity.Karma.KarmaValue != 0 {
		t.Errorf("Karma value mismatch, expected 0, got %d", entity.Karma.KarmaValue)
	}
	if entity.Karma.Merit != 10 {
		t.Errorf("Merit mismatch, expected 10, got %d", entity.Karma.Merit)
	}
}
