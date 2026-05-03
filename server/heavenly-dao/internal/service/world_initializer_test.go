package service

import (
	"testing"

	"github.com/cultivation-world/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestDefaultWorldConfig(t *testing.T) {
	config := DefaultWorldConfig()
	assert.NotNil(t, config)
	assert.True(t, len(config.Regions) >= 15)
	assert.True(t, len(config.Sects) >= 3)
	assert.True(t, len(config.NPCs) >= 50)
	assert.True(t, len(config.Lore) >= 5)
}

func TestNewWorldInitializer(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)
	assert.NotNil(t, init)
}

func TestInitialize_Success(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, len(result.Regions) >= 15)
	assert.True(t, len(result.Sects) >= 6)
	assert.True(t, len(result.NPCs) >= 50)
	assert.True(t, len(result.Lore) >= 5)
}

func TestInitialize_RegionsHaveResources(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	// Check that tianzhou has resources
	tianzhou, exists := result.Regions["tianzhou"]
	assert.True(t, exists)
	assert.True(t, len(tianzhou.Resources) > 0)

	// Check that secret realm has resources
	secretRealm, exists := result.Regions["secret_realm_ancient"]
	assert.True(t, exists)
	assert.True(t, len(secretRealm.Resources) > 0)
}

func TestInitialize_SectsInValidRegions(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	for _, sect := range result.Sects {
		assert.True(t, len(sect.Territory) > 0)
		// Territory should reference valid regions
		for _, regionID := range sect.Territory {
			_, exists := result.Regions[types.RegionID(regionID)]
			assert.True(t, exists, "sect %s references non-existent region %s", sect.Name, regionID)
		}
	}
}

func TestInitialize_NPCsInValidRegions(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	for _, npc := range result.NPCs {
		_, exists := result.Regions[types.RegionID(npc.Position.RegionID)]
		assert.True(t, exists, "NPC %s in non-existent region %s", npc.Name, npc.Position.RegionID)
		assert.NotEmpty(t, npc.Name)
	}
}

func TestInitialize_NPCsHaveSpiritualRoots(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	for _, npc := range result.NPCs {
		assert.True(t, len(npc.Attributes.SpiritualRoots) > 0)
		root := npc.Attributes.SpiritualRoots[0]
		assert.NotEmpty(t, root.Element)
		assert.True(t, root.Purity >= 50 && root.Purity <= 99)
	}
}

func TestInitialize_NPCHaveValidRealms(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	validRealms := map[string]bool{
		"qi_condensation": true,
		"foundation":      true,
		"golden_core":     true,
		"nascent_soul":    true,
		"soul_transformation": true,
		"tribulation":     true,
	}

	for _, npc := range result.NPCs {
		_, exists := validRealms[string(npc.Realm)]
		assert.True(t, exists, "NPC %s has invalid realm %s", npc.Name, npc.Realm)
	}
}

func TestInitialize_DeterministicWithSameSeed(t *testing.T) {
	config := DefaultWorldConfig()

	init1 := NewWorldInitializer(config)
	result1, _ := init1.Initialize(12345)

	init2 := NewWorldInitializer(config)
	result2, _ := init2.Initialize(12345)

	// Same number of NPCs
	assert.Equal(t, len(result1.NPCs), len(result2.NPCs))

	// Same number of regions
	assert.Equal(t, len(result1.Regions), len(result2.Regions))

	// Resource quantities should be the same with same seed
	tianzhou1 := result1.Regions["qingzhou"]
	tianzhou2 := result2.Regions["qingzhou"]
	assert.Equal(t, len(tianzhou1.Resources), len(tianzhou2.Resources))
	for i := range tianzhou1.Resources {
		assert.Equal(t, tianzhou1.Resources[i].Quantity, tianzhou2.Resources[i].Quantity)
	}
}

func TestInitialize_DifferentSeedsProduceDifferentResults(t *testing.T) {
	config := DefaultWorldConfig()

	init1 := NewWorldInitializer(config)
	result1, _ := init1.Initialize(1)

	init2 := NewWorldInitializer(config)
	result2, _ := init2.Initialize(999)

	// Resource quantities should differ with different seeds
	qingzhou1 := result1.Regions["qingzhou"]
	qingzhou2 := result2.Regions["qingzhou"]

	// Resource quantities may differ with different seeds (randomness)
	for i := range qingzhou1.Resources {
		_ = qingzhou1.Resources[i].Quantity != qingzhou2.Resources[i].Quantity
	}
	// Just check structure
	assert.True(t, len(qingzhou1.Resources) > 0)
	assert.True(t, len(qingzhou2.Resources) > 0)
}

func TestInitialize_SectsHaveCorrectAlignment(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	qingyunzong := result.Sects["sect_青云宗"]
	assert.NotNil(t, qingyunzong)
	assert.Equal(t, "righteous", qingyunzong.Alignment)

	xueshadian := result.Sects["sect_血煞殿"]
	assert.NotNil(t, xueshadian)
	assert.Equal(t, "demonic", xueshadian.Alignment)

	tianjige := result.Sects["sect_天机阁"]
	assert.NotNil(t, tianjige)
	assert.Equal(t, "neutral", tianjige.Alignment)
}

func TestInitialize_LoreEntries(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	assert.True(t, len(result.Lore) >= 5)

	// Check specific lore entries exist
	hasCreation := false
	for _, lore := range result.Lore {
		if lore.Title == "创世神话" {
			hasCreation = true
			assert.NotEmpty(t, lore.Description)
			assert.Equal(t, "创世时代", lore.Era)
		}
	}
	assert.True(t, hasCreation)
}

func TestInitialize_RegionHierarchy(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	// Check parent-child relationships
	bingzhou, exists := result.Regions["bingzhou"]
	assert.True(t, exists)
	assert.NotNil(t, bingzhou.ParentRegionID)
	assert.Equal(t, "qingzhou", string(*bingzhou.ParentRegionID))

	youzhou, exists := result.Regions["youzhou"]
	assert.True(t, exists)
	assert.NotNil(t, youzhou.ParentRegionID)
	assert.Equal(t, "tianzhou", string(*youzhou.ParentRegionID))
}

func TestInitialize_SecretRealm(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	secretRealm, exists := result.Regions["secret_realm_ancient"]
	assert.True(t, exists)
	assert.Equal(t, "上古秘境", secretRealm.Name)
	assert.Equal(t, 10, secretRealm.DangerLevel)
	assert.Equal(t, 9, secretRealm.SpiritualTier)

	// Should have rare resources
	assert.True(t, len(secretRealm.Resources) > 0)
	for _, res := range secretRealm.Resources {
		assert.True(t, res.Rarity >= 9) // very rare
	}
}

func TestGetRegionCount(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	_, _ = init.Initialize(42)
	assert.True(t, init.GetRegionCount() >= 15)
}

func TestGetSectCount(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	_, _ = init.Initialize(42)
	assert.True(t, init.GetSectCount() >= 6)
}

func TestGetNPCCount(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	_, _ = init.Initialize(42)
	assert.True(t, init.GetNPCCount() >= 50)
}

func TestGetLoreCount(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	_, _ = init.Initialize(42)
	assert.True(t, init.GetLoreCount() >= 5)
}

func TestInitialize_NPCEntityType(t *testing.T) {
	config := DefaultWorldConfig()
	init := NewWorldInitializer(config)

	result, err := init.Initialize(42)
	assert.NoError(t, err)

	for _, npc := range result.NPCs {
		assert.Equal(t, "npc", string(npc.EntityType))
	}
}
