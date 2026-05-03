package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Sect struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	FounderID         string            `json:"founder_id"`
	Philosophy        string            `json:"philosophy"`
	EntryRequirements map[string]interface{} `json:"entry_requirements"`
	Territory         map[string]interface{} `json:"territory"`
	Rules             map[string]interface{} `json:"rules"`
	Alignment         string            `json:"alignment"`
	CreatedAt         time.Time         `json:"created_at"`
}

type SectMember struct {
	SectID       string    `json:"sect_id"`
	EntityID     string    `json:"entity_id"`
	Rank         string    `json:"rank"`
	Contribution float64   `json:"contribution"`
	JoinedAt     time.Time `json:"joined_at"`
}

type SectRepository struct {
	db *Database
}

func NewSectRepository(db *Database) *SectRepository {
	return &SectRepository{db: db}
}

func (r *SectRepository) Create(ctx context.Context, sect *Sect) error {
	query := `
		INSERT INTO sects (id, name, founder_id, philosophy, entry_requirements, territory, rules, alignment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Pool().Exec(ctx, query,
		sect.ID, sect.Name, sect.FounderID, sect.Philosophy,
		sect.EntryRequirements, sect.Territory, sect.Rules,
		sect.Alignment, sect.CreatedAt,
	)
	return err
}

func (r *SectRepository) GetByID(ctx context.Context, id string) (*Sect, error) {
	query := `SELECT id, name, founder_id, philosophy, entry_requirements, territory, rules, alignment, created_at FROM sects WHERE id = $1`

	var sect Sect
	err := r.db.Pool().QueryRow(ctx, query, id).Scan(
		&sect.ID, &sect.Name, &sect.FounderID, &sect.Philosophy,
		&sect.EntryRequirements, &sect.Territory, &sect.Rules,
		&sect.Alignment, &sect.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &sect, err
}

func (r *SectRepository) GetByName(ctx context.Context, name string) (*Sect, error) {
	query := `SELECT id, name, founder_id, philosophy, entry_requirements, territory, rules, alignment, created_at FROM sects WHERE name = $1`

	var sect Sect
	err := r.db.Pool().QueryRow(ctx, query, name).Scan(
		&sect.ID, &sect.Name, &sect.FounderID, &sect.Philosophy,
		&sect.EntryRequirements, &sect.Territory, &sect.Rules,
		&sect.Alignment, &sect.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &sect, err
}

func (r *SectRepository) AddMember(ctx context.Context, member *SectMember) error {
	query := `
		INSERT INTO sect_members (sect_id, entity_id, rank, contribution, joined_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (sect_id, entity_id) DO UPDATE SET
			rank = EXCLUDED.rank,
			contribution = EXCLUDED.contribution
	`
	_, err := r.db.Pool().Exec(ctx, query,
		member.SectID, member.EntityID, member.Rank, member.Contribution, member.JoinedAt,
	)
	return err
}

func (r *SectRepository) GetMembers(ctx context.Context, sectID string) ([]*SectMember, error) {
	query := `SELECT sect_id, entity_id, rank, contribution, joined_at FROM sect_members WHERE sect_id = $1`

	rows, err := r.db.Pool().Query(ctx, query, sectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*SectMember
	for rows.Next() {
		var m SectMember
		if err := rows.Scan(&m.SectID, &m.EntityID, &m.Rank, &m.Contribution, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}

func (r *SectRepository) GetMember(ctx context.Context, sectID string, entityID string) (*SectMember, error) {
	query := `SELECT sect_id, entity_id, rank, contribution, joined_at FROM sect_members WHERE sect_id = $1 AND entity_id = $2`

	var m SectMember
	err := r.db.Pool().QueryRow(ctx, query, sectID, entityID).Scan(
		&m.SectID, &m.EntityID, &m.Rank, &m.Contribution, &m.JoinedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &m, err
}

func (r *SectRepository) RemoveMember(ctx context.Context, sectID string, entityID string) error {
	query := `DELETE FROM sect_members WHERE sect_id = $1 AND entity_id = $2`
	_, err := r.db.Pool().Exec(ctx, query, sectID, entityID)
	return err
}
