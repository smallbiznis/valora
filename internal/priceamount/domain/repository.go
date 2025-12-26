package domain

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type Repository interface {
	Insert(ctx context.Context, db *gorm.DB, amount *PriceAmount) error
	FindByID(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) (*PriceAmount, error)
	List(ctx context.Context, db *gorm.DB, f PriceAmount) ([]PriceAmount, error)
}
