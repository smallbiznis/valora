package service

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/billingoperations/domain"
	"github.com/smallbiznis/valora/internal/orgcontext"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type FinOpsSnapshotRepository struct {
	db *gorm.DB
}

func NewFinOpsSnapshotRepository(db *gorm.DB) *FinOpsSnapshotRepository {
	return &FinOpsSnapshotRepository{db: db}
}

// FindByUser retrieves snapshots for a specific user within a period range.
func (r *FinOpsSnapshotRepository) FindByUser(ctx context.Context, orgID snowflake.ID, userID string, periodType string, start, end time.Time) ([]domain.FinOpsScoreSnapshot, error) {
	// Ensure OrgID matches context (best practice double check, though caller usually handles auth)
	if ctxOrgID, ok := orgcontext.OrgIDFromContext(ctx); ok && ctxOrgID != orgID {
		return nil, domain.ErrInvalidOrganization
	}

	var rows []finopsSnapshotRow
	if err := r.db.WithContext(ctx).Table("finops_performance_snapshots").
		Where("org_id = ? AND user_id = ? AND period_type = ? AND period_start >= ? AND period_start < ?",
			orgID, userID, periodType, start, end).
		Order("period_start ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	return mapRowsToSnapshots(rows), nil
}

// FindByUserWithLimit retrieves snapshots for a specific user with a limit.
// It orders by period_start DESC to get the most recent snapshots within the window.
func (r *FinOpsSnapshotRepository) FindByUserWithLimit(ctx context.Context, orgID snowflake.ID, userID string, periodType string, start, end time.Time, limit int) ([]domain.FinOpsScoreSnapshot, error) {
	if ctxOrgID, ok := orgcontext.OrgIDFromContext(ctx); ok && ctxOrgID != orgID {
		return nil, domain.ErrInvalidOrganization
	}

	var rows []finopsSnapshotRow
	query := r.db.WithContext(ctx).Table("finops_performance_snapshots").
		Where("org_id = ? AND user_id = ? AND period_type = ? AND period_start >= ? AND period_start < ?",
			orgID, userID, periodType, start, end).
		Order("period_start DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	return mapRowsToSnapshots(rows), nil
}

// FindByOrg retrieves snapshots for all users in an org within a period range (for team view).
// Ordered by UserID, then PeriodStart.
func (r *FinOpsSnapshotRepository) FindByOrg(ctx context.Context, orgID snowflake.ID, periodType string, start, end time.Time) ([]domain.FinOpsScoreSnapshot, error) {
	if ctxOrgID, ok := orgcontext.OrgIDFromContext(ctx); ok && ctxOrgID != orgID {
		return nil, domain.ErrInvalidOrganization
	}

	var rows []finopsSnapshotRow
	if err := r.db.WithContext(ctx).Table("finops_performance_snapshots").
		Where("org_id = ? AND period_type = ? AND period_start >= ? AND period_start < ?",
			orgID, periodType, start, end).
		Order("user_id ASC, period_start ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	return mapRowsToSnapshots(rows), nil
}

// Internal row structure to match DB schema (handling generic JSONB)
type finopsSnapshotRow struct {
	OrgID          snowflake.ID   `gorm:"column:org_id"`
	UserID         string         `gorm:"column:user_id"`
	PeriodType     string         `gorm:"column:period_type"`
	PeriodStart    time.Time      `gorm:"column:period_start"`
	PeriodEnd      time.Time      `gorm:"column:period_end"`
	ScoringVersion string         `gorm:"column:scoring_version"`
	Metrics        datatypes.JSON `gorm:"column:metrics"`
	Scores         datatypes.JSON `gorm:"column:scores"`
}

func mapRowsToSnapshots(rows []finopsSnapshotRow) []domain.FinOpsScoreSnapshot {
	snapshots := make([]domain.FinOpsScoreSnapshot, len(rows))
	for i, r := range rows {
		// We rely on domain JSON Unmarshal or manual if needed.
		// Since domain.FinOpsScoreSnapshot uses struct tags matching JSON, we can unmarshal.
		// Warning: datatypes.JSON is []byte.
		// We map to domain struct which has Metrics PerformanceMetrics.
		// We need to unmarshal the JSON content.

		// Helper to unmarshal safely
		var m domain.PerformanceMetrics
		var s domain.PerformanceScores
		// We suppress error here assuming DB data is valid JSON if inserted correctly.
		// In a real repo we might log error.
		_ = m.UnmarshalJSON(r.Metrics)
		_ = s.UnmarshalJSON(r.Scores)

		snapshots[i] = domain.FinOpsScoreSnapshot{
			OrgID:          r.OrgID.String(),
			UserID:         r.UserID,
			PeriodType:     r.PeriodType,
			PeriodStart:    r.PeriodStart,
			PeriodEnd:      r.PeriodEnd,
			ScoringVersion: r.ScoringVersion,
			Metrics:        m,
			Scores:         s,
		}
	}
	return snapshots
}
