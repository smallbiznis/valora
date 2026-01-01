package domain

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/pkg/db/option"
	"gorm.io/gorm"
)

type Repository interface {
	Insert(ctx context.Context, db *gorm.DB, amount *PriceAmount) error
	FindOne(ctx context.Context, db *gorm.DB, amount *PriceAmount) (*PriceAmount, error)
	FindByID(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) (*PriceAmount, error)
	List(ctx context.Context, db *gorm.DB, f PriceAmount, opts ...option.QueryOption) ([]PriceAmount, error)
	Update(ctx context.Context, db *gorm.DB, amount *PriceAmount) (*PriceAmount, error)
}
