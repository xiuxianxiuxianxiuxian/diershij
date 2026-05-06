package service

import (
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

type MemoryEntry struct {
	Content          string    `json:"content"`
	MemoryType       string    `json:"memory_type"` // short_term, long_term
	Importance       float64   `json:"importance"`  // 0.0 to 1.0
	RelatedEntityID  string    `json:"related_entity_id"`
	RelatedEntityName string   `json:"related_entity_name"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
}

type RelationshipMemory struct {
	TargetID        string `json:"target_id"`
	TargetName      string `json:"target_name"`
	RelationshipType string `json:"relationship_type"` // player, npc
	Affinity        int    `json:"affinity"`            // -100 to 100
	Familiarity     int    `json:"familiarity"`         // 0 to 100
	InteractionCount int   `json:"interaction_count"`
}

type NPCMemoryStore struct {
	mu sync.RWMutex

	npcID string

	shortTerm []*MemoryEntry
	longTerm  []*MemoryEntry

	relationships map[string]*RelationshipMemory

	// For memory consolidation
	lastConsolidation time.Time
}

const (
	maxShortTermMemories  = 20
	maxLongTermMemories   = 100
	consolidationInterval = 10 * time.Minute
)

func NewNPCMemoryStore(npcID string) *NPCMemoryStore {
	return &NPCMemoryStore{
		npcID:             npcID,
		shortTerm:         make([]*MemoryEntry, 0, maxShortTermMemories),
		longTerm:          make([]*MemoryEntry, 0, maxLongTermMemories),
		relationships:     make(map[string]*RelationshipMemory),
		lastConsolidation: time.Now(),
	}
}

// Remember adds a new short-term memory. High-importance memories auto-consolidate.
func (m *NPCMemoryStore) Remember(content string, importance float64, entityID string, entityName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := &MemoryEntry{
		Content:           content,
		MemoryType:        "short_term",
		Importance:        math.Max(0, math.Min(1, importance)),
		RelatedEntityID:   entityID,
		RelatedEntityName: entityName,
		CreatedAt:         time.Now(),
	}

	m.shortTerm = append(m.shortTerm, entry)

	// Consolidate immediately if importance is high
	if importance >= 0.7 {
		m.consolidateLocked()
	}

	// Trim if over capacity
	if len(m.shortTerm) > maxShortTermMemories {
		// Sort by importance, keep most important
		m.shortTerm = m.topMemories(m.shortTerm, maxShortTermMemories)
	}
}

// RememberInteraction records an interaction with a player/NPC and creates a memory.
func (m *NPCMemoryStore) RememberInteraction(entityID string, entityName string, content string, affinityDelta int) {
	m.Remember(content, 0.5, entityID, entityName)

	m.mu.Lock()
	rel, exists := m.relationships[entityID]
	if !exists {
		rel = &RelationshipMemory{
			TargetID:         entityID,
			TargetName:       entityName,
			RelationshipType: "player",
			Affinity:         0,
			Familiarity:      0,
			InteractionCount: 0,
		}
		m.relationships[entityID] = rel
	}
	rel.Affinity = clampInt(rel.Affinity+affinityDelta, -100, 100)
	rel.Familiarity = clampInt(rel.Familiarity+1, 0, 100)
	rel.InteractionCount++
	rel.TargetName = entityName
	m.mu.Unlock()
}

// GetRelationship returns a relationship memory for a target.
func (m *NPCMemoryStore) GetRelationship(targetID string) *RelationshipMemory {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rel, exists := m.relationships[targetID]
	if !exists {
		return nil
	}
	return rel
}

// GetAllRelationships returns all relationships sorted by familiarity.
func (m *NPCMemoryStore) GetAllRelationships() []*RelationshipMemory {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*RelationshipMemory, 0, len(m.relationships))
	for _, rel := range m.relationships {
		result = append(result, rel)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Familiarity > result[j].Familiarity
	})
	return result
}

// GetRecentMemories returns the most recent memories, with optional filtering.
func (m *NPCMemoryStore) GetRecentMemories(count int, entityID string) []*MemoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	all := make([]*MemoryEntry, 0, len(m.longTerm)+len(m.shortTerm))
	all = append(all, m.longTerm...)
	all = append(all, m.shortTerm...)
	if entityID != "" {
		filtered := make([]*MemoryEntry, 0, len(all))
		for _, mem := range all {
			if mem.RelatedEntityID == entityID {
				filtered = append(filtered, mem)
			}
		}
		all = filtered
	}

	// Sort by creation time, newest first
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	if count <= 0 || count > len(all) {
		count = len(all)
	}
	return all[:count]
}

// GetMemoryContext builds a prompt context string from memories.
func (m *NPCMemoryStore) GetMemoryContext() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var parts []string

	// Recent important memories
	important := m.topMemories(m.shortTerm, 5)
	for _, mem := range important {
		parts = append(parts, "- "+mem.Content)
	}

	// Active relationships
	if len(m.relationships) > 0 {
		var relParts []string
		for _, rel := range m.relationships {
			relParts = append(relParts, rel.TargetName)
		}
		parts = append(parts, "Known entities: "+strings.Join(relParts, ", "))
	}

	return strings.Join(parts, "\n")
}

// Consolidate promotes important short-term memories to long-term.
func (m *NPCMemoryStore) Consolidate() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.consolidateLocked()
}

func (m *NPCMemoryStore) consolidateLocked() {
	if time.Since(m.lastConsolidation) < consolidationInterval {
		return
	}
	m.lastConsolidation = time.Now()

	// Promote highly important short-term memories
	var remaining []*MemoryEntry
	for _, mem := range m.shortTerm {
		if mem.Importance >= 0.6 {
			mem.MemoryType = "long_term"
			m.longTerm = append(m.longTerm, mem)
		} else {
			remaining = append(remaining, mem)
		}
	}
	m.shortTerm = remaining

	// Trim long-term if over capacity
	if len(m.longTerm) > maxLongTermMemories {
		m.longTerm = m.topMemories(m.longTerm, maxLongTermMemories)
	}
}

// ToPersistableMemories returns memories that should be saved to the database.
func (m *NPCMemoryStore) ToPersistableMemories() []*MemoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Only return long-term + high-importance short-term for persistence
	result := make([]*MemoryEntry, 0, len(m.longTerm)+10)
	result = append(result, m.longTerm...)
	for _, mem := range m.shortTerm {
		if mem.Importance >= 0.5 {
			result = append(result, mem)
		}
	}
	return result
}

func (m *NPCMemoryStore) topMemories(entries []*MemoryEntry, n int) []*MemoryEntry {
	if len(entries) <= n {
		return entries
	}
	sorted := make([]*MemoryEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Importance != sorted[j].Importance {
			return sorted[i].Importance > sorted[j].Importance
		}
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})
	return sorted[:n]
}

func clampInt(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
