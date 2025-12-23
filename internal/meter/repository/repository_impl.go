package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	meterdomain "github.com/smallbiznis/valora/internal/meter/domain"
	"gorm.io/gorm"
)

type repo struct{}

func Provide() meterdomain.Repository {
	return &repo{}
}

func (r *repo) Insert(ctx context.Context, db *gorm.DB, m *meterdomain.Meter) error {
	return db.WithContext(ctx).Exec(
		`INSERT INTO meters (id, org_id, code, name, aggregation, unit, active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID,
		m.OrgID,
		m.Code,
		m.Name,
		m.Aggregation,
		m.Unit,
		m.Active,
		m.CreatedAt,
		m.UpdatedAt,
	).Error
}

func (r *repo) Update(ctx context.Context, db *gorm.DB, m *meterdomain.Meter) error {
	return db.WithContext(ctx).Exec(
		`UPDATE meters
		 SET name = ?, aggregation = ?, unit = ?, active = ?, updated_at = ?
		 WHERE org_id = ? AND id = ?`,
		m.Name,
		m.Aggregation,
		m.Unit,
		m.Active,
		m.UpdatedAt,
		m.OrgID,
		m.ID,
	).Error
}

func (r *repo) Delete(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) error {
	return db.WithContext(ctx).Exec(
		`DELETE FROM meters WHERE org_id = ? AND id = ?`,
		orgID,
		id,
	).Error
}

func (r *repo) FindByID(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) (*meterdomain.Meter, error) {
	var meter meterdomain.Meter
	err := db.WithContext(ctx).Raw(
		`SELECT id, org_id, code, name, aggregation, unit, active, created_at, updated_at
		 FROM meters WHERE org_id = ? AND id = ?`,
		orgID,
		id,
	).Scan(&meter).Error
	if err != nil {
		return nil, err
	}
	if meter.ID == 0 {
		return nil, nil
	}
	return &meter, nil
}

func (r *repo) FindByCode(ctx context.Context, db *gorm.DB, orgID snowflake.ID, code string) (*meterdomain.Meter, error) {
	var meter meterdomain.Meter
	err := db.WithContext(ctx).Raw(
		`SELECT id, org_id, code, name, aggregation, unit, active, created_at, updated_at
		 FROM meters WHERE org_id = ? AND code = ?`,
		orgID,
		code,
	).Scan(&meter).Error
	if err != nil {
		return nil, err
	}
	if meter.ID == 0 {
		return nil, nil
	}
	return &meter, nil
}

func (r *repo) List(ctx context.Context, db *gorm.DB, orgID snowflake.ID) ([]meterdomain.Meter, error) {
	var meters []meterdomain.Meter
	err := db.WithContext(ctx).Raw(
		`SELECT id, org_id, code, name, aggregation, unit, active, created_at, updated_at
		 FROM meters WHERE org_id = ? ORDER BY created_at ASC`,
		orgID,
	).Scan(&meters).Error
	if err != nil {
		return nil, err
	}
	return meters, nil
}
