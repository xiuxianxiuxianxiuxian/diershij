package repository

import (
	"context"
	"time"
)

type NPCProfileRow struct {
	NPCID           string    `db:"npc_id"`
	EntityID        string    `db:"entity_id"`
	PersonalityType string    `db:"personality_type"`
	MoralAlignment  string    `db:"moral_alignment"`
	AmbitionLevel   int       `db:"ambition_level"`
	RiskTolerance   float64   `db:"risk_tolerance"`
	BackgroundStory string    `db:"background_story"`
	CurrentGoal     string    `db:"current_goal"`
	CurrentRegion   string    `db:"current_region"`
	Realm           string    `db:"realm"`
	Status          string    `db:"status"`
	LastActiveAt    time.Time `db:"last_active_at"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type NPCMemoryRow struct {
	ID               int       `db:"id"`
	NPCID            string    `db:"npc_id"`
	MemoryType       string    `db:"memory_type"`
	MemoryKey        string    `db:"memory_key"`
	Content          string    `db:"content"`
	Importance       float64   `db:"importance"`
	RelatedEntityID  string    `db:"related_entity_id"`
	RelatedEntityName string   `db:"related_entity_name"`
	CreatedAt        time.Time `db:"created_at"`
	ExpiresAt        *time.Time `db:"expires_at"`
}

type NPCRelationshipRow struct {
	NPCID             string    `db:"npc_id"`
	TargetID          string    `db:"target_id"`
	TargetName        string    `db:"target_name"`
	RelationshipType  string    `db:"relationship_type"`
	Affinity          int       `db:"affinity"`
	Familiarity       int       `db:"familiarity"`
	LastInteractionAt time.Time `db:"last_interaction_at"`
	InteractionCount  int       `db:"interaction_count"`
	Notes             string    `db:"notes"`
}

type NPCRepository struct {
	db *Database
}

func NewNPCRepository(db *Database) *NPCRepository {
	return &NPCRepository{db: db}
}

func (r *NPCRepository) SaveProfile(ctx context.Context, profile *NPCProfileRow) error {
	query := `INSERT INTO npc_profiles
		(npc_id, entity_id, personality_type, moral_alignment, ambition_level, risk_tolerance,
		 background_story, current_goal, current_region, realm, status, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW())
		ON CONFLICT (npc_id) DO UPDATE SET
			personality_type=EXCLUDED.personality_type,
			moral_alignment=EXCLUDED.moral_alignment,
			ambition_level=EXCLUDED.ambition_level,
			risk_tolerance=EXCLUDED.risk_tolerance,
			background_story=EXCLUDED.background_story,
			current_goal=EXCLUDED.current_goal,
			current_region=EXCLUDED.current_region,
			realm=EXCLUDED.realm,
			status=EXCLUDED.status,
			updated_at=NOW()`
	_, err := r.db.Pool().Exec(ctx, query,
		profile.NPCID, profile.EntityID, profile.PersonalityType, profile.MoralAlignment,
		profile.AmbitionLevel, profile.RiskTolerance, profile.BackgroundStory,
		profile.CurrentGoal, profile.CurrentRegion, profile.Realm, profile.Status)
	return err
}

func (r *NPCRepository) GetProfile(ctx context.Context, npcID string) (*NPCProfileRow, error) {
	query := `SELECT npc_id, entity_id, personality_type, moral_alignment, ambition_level,
		risk_tolerance, background_story, current_goal, current_region, realm, status,
		last_active_at, created_at, updated_at
		FROM npc_profiles WHERE npc_id = $1`
	row := r.db.Pool().QueryRow(ctx, query, npcID)
	p := &NPCProfileRow{}
	err := row.Scan(&p.NPCID, &p.EntityID, &p.PersonalityType, &p.MoralAlignment,
		&p.AmbitionLevel, &p.RiskTolerance, &p.BackgroundStory, &p.CurrentGoal,
		&p.CurrentRegion, &p.Realm, &p.Status,
		&p.LastActiveAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *NPCRepository) GetAllActiveProfiles(ctx context.Context) ([]*NPCProfileRow, error) {
	query := `SELECT npc_id, entity_id, personality_type, moral_alignment, ambition_level,
		risk_tolerance, background_story, current_goal, current_region, realm, status,
		last_active_at, created_at, updated_at
		FROM npc_profiles WHERE status != 'inactive' ORDER BY last_active_at DESC`
	rows, err := r.db.Pool().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*NPCProfileRow
	for rows.Next() {
		p := &NPCProfileRow{}
		if err := rows.Scan(&p.NPCID, &p.EntityID, &p.PersonalityType, &p.MoralAlignment,
			&p.AmbitionLevel, &p.RiskTolerance, &p.BackgroundStory, &p.CurrentGoal,
			&p.CurrentRegion, &p.Realm, &p.Status,
			&p.LastActiveAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *NPCRepository) UpdateStatus(ctx context.Context, npcID string, status string, region string) error {
	query := `UPDATE npc_profiles SET status=$1, current_region=$2, last_active_at=NOW(), updated_at=NOW() WHERE npc_id=$3`
	_, err := r.db.Pool().Exec(ctx, query, status, region, npcID)
	return err
}

func (r *NPCRepository) UpdateGoal(ctx context.Context, npcID string, goal string) error {
	_, err := r.db.Pool().Exec(ctx, `UPDATE npc_profiles SET current_goal=$1, updated_at=NOW() WHERE npc_id=$2`, goal, npcID)
	return err
}

func (r *NPCRepository) UpdateRealm(ctx context.Context, npcID string, realm string) error {
	_, err := r.db.Pool().Exec(ctx, `UPDATE npc_profiles SET realm=$1, updated_at=NOW() WHERE npc_id=$2`, realm, npcID)
	return err
}

func (r *NPCRepository) DeleteProfile(ctx context.Context, npcID string) error {
	_, err := r.db.Pool().Exec(ctx, `DELETE FROM npc_profiles WHERE npc_id=$1`, npcID)
	return err
}

func (r *NPCRepository) SaveMemory(ctx context.Context, m *NPCMemoryRow) error {
	query := `INSERT INTO npc_memory
		(npc_id, memory_type, memory_key, content, importance, related_entity_id, related_entity_name, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.Pool().Exec(ctx, query,
		m.NPCID, m.MemoryType, m.MemoryKey, m.Content, m.Importance,
		m.RelatedEntityID, m.RelatedEntityName, m.ExpiresAt)
	return err
}

// ClearMemories deletes all memories for an NPC before re-saving the current set.
func (r *NPCRepository) ClearMemories(ctx context.Context, npcID string) error {
	_, err := r.db.Pool().Exec(ctx, `DELETE FROM npc_memory WHERE npc_id=$1`, npcID)
	return err
}

// SaveMemoriesBatch clears old memories then saves the current set (atomic-like replace).
func (r *NPCRepository) SaveMemoriesBatch(ctx context.Context, npcID string, memories []*NPCMemoryRow) error {
	tx, err := r.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM npc_memory WHERE npc_id=$1`, npcID); err != nil {
		return err
	}

	for _, m := range memories {
		_, err := tx.Exec(ctx,
			`INSERT INTO npc_memory
			 (npc_id, memory_type, memory_key, content, importance, related_entity_id, related_entity_name, expires_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			m.NPCID, m.MemoryType, m.MemoryKey, m.Content, m.Importance,
			m.RelatedEntityID, m.RelatedEntityName, m.ExpiresAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *NPCRepository) GetMemories(ctx context.Context, npcID string, memoryType string, limit int) ([]*NPCMemoryRow, error) {
	query := `SELECT id, npc_id, memory_type, memory_key, content, importance,
		related_entity_id, related_entity_name, created_at, expires_at
		FROM npc_memory
		WHERE npc_id=$1 AND (memory_type=$2 OR $2='') AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY importance DESC, created_at DESC LIMIT $3`
	rows, err := r.db.Pool().Query(ctx, query, npcID, memoryType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*NPCMemoryRow
	for rows.Next() {
		m := &NPCMemoryRow{}
		if err := rows.Scan(&m.ID, &m.NPCID, &m.MemoryType, &m.MemoryKey, &m.Content,
			&m.Importance, &m.RelatedEntityID, &m.RelatedEntityName, &m.CreatedAt, &m.ExpiresAt); err != nil {
			return nil, err
		}
		memories = append(memories, m)
	}
	return memories, nil
}

func (r *NPCRepository) CleanExpiredMemories(ctx context.Context, npcID string) error {
	_, err := r.db.Pool().Exec(ctx, `DELETE FROM npc_memory WHERE npc_id=$1 AND expires_at < NOW()`, npcID)
	return err
}

func (r *NPCRepository) GetRelationship(ctx context.Context, npcID string, targetID string) (*NPCRelationshipRow, error) {
	row := r.db.Pool().QueryRow(ctx,
		`SELECT npc_id, target_id, target_name, relationship_type, affinity, familiarity,
		 last_interaction_at, interaction_count, notes
		 FROM npc_relationships WHERE npc_id=$1 AND target_id=$2`, npcID, targetID)
	rel := &NPCRelationshipRow{}
	err := row.Scan(&rel.NPCID, &rel.TargetID, &rel.TargetName, &rel.RelationshipType,
		&rel.Affinity, &rel.Familiarity, &rel.LastInteractionAt, &rel.InteractionCount, &rel.Notes)
	if err != nil {
		return nil, err
	}
	return rel, nil
}

func (r *NPCRepository) UpsertRelationship(ctx context.Context, rel *NPCRelationshipRow) error {
	query := `INSERT INTO npc_relationships
		(npc_id, target_id, target_name, relationship_type, affinity, familiarity, last_interaction_at, interaction_count, notes)
		VALUES ($1,$2,$3,$4,$5,$6,NOW(),$7,$8)
		ON CONFLICT (npc_id, target_id) DO UPDATE SET
			target_name=EXCLUDED.target_name,
			affinity=EXCLUDED.affinity,
			familiarity=EXCLUDED.familiarity,
			last_interaction_at=NOW(),
			interaction_count=EXCLUDED.interaction_count,
			notes=EXCLUDED.notes`
	_, err := r.db.Pool().Exec(ctx, query,
		rel.NPCID, rel.TargetID, rel.TargetName, rel.RelationshipType,
		rel.Affinity, rel.Familiarity, rel.InteractionCount, rel.Notes)
	return err
}

func (r *NPCRepository) GetRelationships(ctx context.Context, npcID string) ([]*NPCRelationshipRow, error) {
	rows, err := r.db.Pool().Query(ctx,
		`SELECT npc_id, target_id, target_name, relationship_type, affinity, familiarity,
		 last_interaction_at, interaction_count, notes
		 FROM npc_relationships WHERE npc_id=$1 ORDER BY ABS(affinity) DESC, familiarity DESC`, npcID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []*NPCRelationshipRow
	for rows.Next() {
		rel := &NPCRelationshipRow{}
		if err := rows.Scan(&rel.NPCID, &rel.TargetID, &rel.TargetName, &rel.RelationshipType,
			&rel.Affinity, &rel.Familiarity, &rel.LastInteractionAt, &rel.InteractionCount, &rel.Notes); err != nil {
			return nil, err
		}
		rels = append(rels, rel)
	}
	return rels, nil
}
