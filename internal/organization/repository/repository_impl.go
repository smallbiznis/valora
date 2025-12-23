package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/organization/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) WithTx(tx *gorm.DB) domain.Repository {
	return &repository{db: tx}
}

func (r *repository) CreateOrganization(ctx context.Context, org domain.Organization) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO organizations (id, name, slug, country_code, timezone_name, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		org.ID,
		org.Name,
		org.Slug,
		org.CountryCode,
		org.TimezoneName,
		org.CreatedAt,
	).Error
}

func (r *repository) AddMember(ctx context.Context, member domain.OrganizationMember) error {
	return r.db.WithContext(ctx).Exec(
		`INSERT INTO organization_members (id, org_id, user_id, role, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		member.ID,
		member.OrgID,
		member.UserID,
		member.Role,
		member.CreatedAt,
	).Error
}

func (r *repository) ListOrganizationsByUser(ctx context.Context, userID snowflake.ID) ([]domain.OrganizationListItem, error) {
	var items []domain.OrganizationListItem
	err := r.db.WithContext(ctx).Raw(
		`SELECT o.id, o.name, m.role, o.created_at
		 FROM organizations o
		 JOIN organization_members m ON m.org_id = o.id
		 WHERE m.user_id = ?
		 ORDER BY o.created_at ASC`,
		userID,
	).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
