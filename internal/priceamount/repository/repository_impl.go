package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	priceamountdomain "github.com/smallbiznis/valora/internal/priceamount/domain"
	"github.com/smallbiznis/valora/pkg/db/option"
	"gorm.io/gorm"
)

type repo struct{}

func Provide() priceamountdomain.Repository {
	return &repo{}
}

func (r *repo) Insert(ctx context.Context, db *gorm.DB, amount *priceamountdomain.PriceAmount) error {
	return db.WithContext(ctx).Exec(
		`INSERT INTO price_amounts (
			id, org_id, price_id, meter_id, currency, unit_amount_cents, minimum_amount_cents, maximum_amount_cents, effective_from, effective_to,
			metadata, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		amount.ID,
		amount.OrgID,
		amount.PriceID,
		amount.MeterID,
		amount.Currency,
		amount.UnitAmountCents,
		amount.MinimumAmountCents,
		amount.MaximumAmountCents,
		amount.EffectiveFrom,
		amount.EffectiveTo,
		amount.Metadata,
		amount.CreatedAt,
		amount.UpdatedAt,
	).Error
}

func (r *repo) FindByID(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) (*priceamountdomain.PriceAmount, error) {
	var amount priceamountdomain.PriceAmount
	err := db.WithContext(ctx).Raw(
		`SELECT id, org_id, price_id, meter_id, currency, unit_amount_cents, minimum_amount_cents, maximum_amount_cents, effective_from, effective_to,
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

func (r *repo) List(ctx context.Context, db *gorm.DB, f priceamountdomain.PriceAmount, opts ...option.QueryOption) ([]priceamountdomain.PriceAmount, error) {
	var items []priceamountdomain.PriceAmount
	stmt := db.WithContext(ctx).Model(&priceamountdomain.PriceAmount{})

	for _, opt := range opts {
		stmt = opt.Apply(stmt)
	}

	if err := stmt.Where(f).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repo) FindOne(ctx context.Context, db *gorm.DB, amount *priceamountdomain.PriceAmount) (*priceamountdomain.PriceAmount, error) {
	var result priceamountdomain.PriceAmount
	err := db.WithContext(ctx).Where(amount).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *repo) Update(ctx context.Context, db *gorm.DB, amount *priceamountdomain.PriceAmount) (*priceamountdomain.PriceAmount, error) {
	err := db.WithContext(ctx).Model(&priceamountdomain.PriceAmount{}).
		Where("org_id = ? AND id = ?", amount.OrgID, amount.ID).
		Updates(map[string]interface{}{
			"meter_id":             amount.MeterID,
			"currency":             amount.Currency,
			"unit_amount_cents":    amount.UnitAmountCents,
			"minimum_amount_cents": amount.MinimumAmountCents,
			"maximum_amount_cents": amount.MaximumAmountCents,
			"effective_from":       amount.EffectiveFrom,
			"effective_to":         amount.EffectiveTo,
			"metadata":             amount.Metadata,
			"updated_at":           amount.UpdatedAt,
		}).Error
	if err != nil {
		return nil, err
	}
	return amount, nil
}
