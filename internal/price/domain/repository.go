package domain

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)

type Repository interface {
	Insert(ctx context.Context, db *gorm.DB, price *Price) error
	FindByID(ctx context.Context, db *gorm.DB, orgID, id snowflake.ID) (*Price, error)
	List(ctx context.Context, db *gorm.DB, orgID snowflake.ID) ([]Price, error)
}
