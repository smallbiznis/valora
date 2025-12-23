package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	priceamountdomain "github.com/smallbiznis/valora/internal/priceamount/domain"
	"gorm.io/gorm"
)

type repo struct{}

func Provide() priceamountdomain.Repository {
	return &repo{}
}

func (r *repo) Insert(ctx context.Context, db *gorm.DB, amount *priceamountdomain.PriceAmount) error {
	return db.WithContext(ctx).Exec(
		`INSERT INTO price_amounts (
			id, org_id, price_id, meter_id, currency, unit_amount_cents, minimum_amount_cents, maximum_amount_cents,
			metadata, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		amount.ID,
		amount.OrgID,
		amount.PriceID,
		amount.MeterID,
		amount.Currency,
		amount.UnitAmountCents,
		amount.MinimumAmountCents,
		amount.MaximumAmountCents,
		amount.Metadata,
		amount.CreatedAt,
		amount.UpdatedAt,
	).Error
}

func (r *repo) FindByID(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) (*priceamountdomain.PriceAmount, error) {
	var amount priceamountdomain.PriceAmount
	err := db.WithContext(ctx).Raw(
		`SELECT id, org_id, price_id, meter_id, currency, unit_amount_cents, minimum_amount_cents, maximum_amount_cents,
		 metadata, created_at, updated_at
		 FROM price_amounts WHERE org_id = ? AND id = ?`,
		orgID,
		id,
	).Scan(&amount).Error
	if err != nil {
		return nil, err
	}
	if amount.ID == 0 {
		return nil, nil
	}
	return &amount, nil
}

func (r *repo) List(ctx context.Context, db *gorm.DB, orgID snowflake.ID) ([]priceamountdomain.PriceAmount, error) {
	var items []priceamountdomain.PriceAmount
	err := db.WithContext(ctx).Raw(
		`SELECT id, org_id, price_id, meter_id, currency, unit_amount_cents, minimum_amount_cents, maximum_amount_cents,
		 metadata, created_at, updated_at
		 FROM price_amounts WHERE org_id = ? ORDER BY created_at ASC`,
		orgID,
	).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
