package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/cultivation-world/shared/types"
)

type WorldRepository struct {
	db *sql.DB
}

func NewWorldRepository(host string, port int, user, password, dbname string) (*WorldRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &WorldRepository{db: db}, nil
}

func (r *WorldRepository) Close() error {
	return r.db.Close()
}

// SaveWorldState 持久化世界状态（epoch + balance metrics）
func (r *WorldRepository) SaveWorldState(epoch int64, metrics types.BalanceMetrics) error {
	data, err := json.Marshal(map[string]interface{}{
		"power_distribution":   metrics.PowerDistribution,
		"resource_circulation": metrics.ResourceCirculation,
		"sect_diversity":       metrics.SectDiversity,
		"karma_distribution":   metrics.KarmaDistribution,
	})
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`DELETE FROM world_state`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO world_state (epoch, state_data, updated_at)
		VALUES ($1, $2, NOW())
	`, epoch, string(data))
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// LoadWorldState 加载最新世界状态
func (r *WorldRepository) LoadWorldState() (int64, *types.BalanceMetrics, error) {
	var epoch int64
	var stateData string

	err := r.db.QueryRow(`
		SELECT epoch, state_data FROM world_state ORDER BY id DESC LIMIT 1
	`).Scan(&epoch, &stateData)
	if err == sql.ErrNoRows {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}

	var parsed struct {
		PowerDistribution   float64 `json:"power_distribution"`
		ResourceCirculation float64 `json:"resource_circulation"`
		SectDiversity       float64 `json:"sect_diversity"`
		KarmaDistribution   float64 `json:"karma_distribution"`
	}
	if err := json.Unmarshal([]byte(stateData), &parsed); err != nil {
		return 0, nil, err
	}

	metrics := &types.BalanceMetrics{
		PowerDistribution:   parsed.PowerDistribution,
		ResourceCirculation: parsed.ResourceCirculation,
		SectDiversity:       parsed.SectDiversity,
		KarmaDistribution:   parsed.KarmaDistribution,
	}
	return epoch, metrics, nil
}

// SaveRegionResources 持久化所有区域的资源数量
func (r *WorldRepository) SaveRegionResources(regions map[string]*types.Region) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for regionID, region := range regions {
		for _, res := range region.Resources {
			var lastHarvested interface{}
			if res.LastHarvested != nil {
				lastHarvested = res.LastHarvested
			}
			_, err := tx.Exec(`
				INSERT INTO region_resources (region_id, resource_id, quantity, max_quantity, last_harvested, updated_at)
				VALUES ($1, $2, $3, $4, $5, NOW())
				ON CONFLICT (region_id, resource_id)
				DO UPDATE SET quantity = $3, last_harvested = $5, updated_at = NOW()
			`, regionID, res.ID, res.Quantity, 100, lastHarvested)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// LoadRegionResources 加载持久化的资源数量
func (r *WorldRepository) LoadRegionResources(regionID string) (map[string]int, error) {
	rows, err := r.db.Query(`
		SELECT resource_id, quantity FROM region_resources WHERE region_id = $1
	`, regionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var resID string
		var qty int
		if err := rows.Scan(&resID, &qty); err != nil {
			return nil, err
		}
		result[resID] = qty
	}
	return result, rows.Err()
}
