package types

import "time"

type SpellID string

type SpellType string

const (
	SpellTypeAttack SpellType = "attack"
	SpellTypeHeal   SpellType = "heal"
	SpellTypeBuff   SpellType = "buff"
	SpellTypeDebuff SpellType = "debuff"
)

type SpellElement string

const (
	ElementFire   SpellElement = "fire"
	ElementWater  SpellElement = "water"
	ElementEarth  SpellElement = "earth"
	ElementMetal  SpellElement = "metal"
	ElementWood   SpellElement = "wood"
	ElementWind   SpellElement = "wind"
	ElementThunder SpellElement = "thunder"
	ElementIce    SpellElement = "ice"
	ElementLight  SpellElement = "light"
	ElementDark   SpellElement = "dark"
)

type Spell struct {
	ID               SpellID      `json:"id"`
	Name             string       `json:"name"`
	Type             SpellType    `json:"type"`
	Element          SpellElement `json:"element"`
	Cost             int          `json:"cost"`
	BaseDamage       int          `json:"base_damage"`
	BaseHeal         int          `json:"base_heal"`
	Duration         int          `json:"duration"`
	Cooldown         int          `json:"cooldown"`
	Description      string       `json:"description"`
	RealmRequirement CultivationRealm `json:"realm_requirement"`
	CreatedAt        time.Time    `json:"created_at"`
}

type EntitySpell struct {
	EntityID    EntityID  `json:"entity_id"`
	SpellID     SpellID   `json:"spell_id"`
	Spell       *Spell    `json:"spell,omitempty"`
	Proficiency int       `json:"proficiency"`
	LearnedAt   time.Time `json:"learned_at"`
	LastCastAt  *time.Time `json:"last_cast_at,omitempty"`
}

func (es *EntitySpell) CanCast(now time.Time) bool {
	if es.LastCastAt == nil {
		return true
	}
	if es.Spell == nil {
		return true
	}
	cooldownEnd := es.LastCastAt.Add(time.Duration(es.Spell.Cooldown) * time.Second)
	return now.After(cooldownEnd)
}

func (es *EntitySpell) GetCooldownRemaining(now time.Time) time.Duration {
	if es.LastCastAt == nil || es.Spell == nil {
		return 0
	}
	cooldownEnd := es.LastCastAt.Add(time.Duration(es.Spell.Cooldown) * time.Second)
	if now.After(cooldownEnd) {
		return 0
	}
	return cooldownEnd.Sub(now)
}
