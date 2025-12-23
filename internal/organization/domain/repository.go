package domain

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
)


type OrganizationListItem struct {
	ID        snowflake.ID
	Name      string
	Role      string
	CreatedAt time.Time
}

type Repository interface {
	WithTx(tx *gorm.DB) Repository
	CreateOrganization(ctx context.Context, org Organization) error
	AddMember(ctx context.Context, member OrganizationMember) error
	ListOrganizationsByUser(ctx context.Context, userID snowflake.ID) ([]OrganizationListItem, error)
}
