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
            X: 0,
            Y: 0,
        },
        Attributes: Attributes{
            Qi: 100,
            MaxQi: 100,
        },
        Status: StatusNormal,
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
        Qi: 50,
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
        KarmaValue: 10,
        Merit: 5,
        HeavenlyMark: "clear",
    }
    
    if karma.KarmaValue != 10 {
        t.Errorf("KarmaValue mismatch, expected 10, got %d", karma.KarmaValue)
    }
    
    if karma.Merit != 5 {
        t.Errorf("Merit mismatch, expected 5, got %d", karma.Merit)
    }
}
