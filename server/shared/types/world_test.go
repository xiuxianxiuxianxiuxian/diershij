package types

import (
    "testing"
)

func TestRegionInitialization(t *testing.T) {
    region := &Region{
        ID: "test-region",
        Name: "Test Region",
        Description: "Test Description",
        SpiritualDensity: 50.0,
        DangerLevel: 2,
    }
    
    if region.ID != "test-region" {
        t.Errorf("Region.ID mismatch, expected test-region, got %s", region.ID)
    }
    
    if region.DangerLevel != 2 {
        t.Errorf("Region.DangerLevel mismatch, expected 2, got %d", region.DangerLevel)
    }
}

func TestWorldStateInitialization(t *testing.T) {
    world := &WorldState{
        Epoch: 1000,
    }
    
    if world.Epoch != 1000 {
        t.Errorf("WorldState.Epoch mismatch, expected 1000, got %d", world.Epoch)
    }
}

func TestResourceInitialization(t *testing.T) {
    resource := &Resource{
        ID: "test-resource",
        Type: "spiritual_herb",
        Name: "Test Herb",
        Rarity: 3,
    }
    
    if resource.Type != "spiritual_herb" {
        t.Errorf("Resource.Type mismatch, expected spiritual_herb, got %s", resource.Type)
    }
    
    if resource.Rarity != 3 {
        t.Errorf("Resource.Rarity mismatch, expected 3, got %d", resource.Rarity)
    }
}
