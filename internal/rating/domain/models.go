// Package domain contains persistence models for rating outputs.
package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/datatypes"
)

// RatingResult captures the aggregated charge outcome for a billing cycle.
type RatingResult struct {
	ID             snowflake.ID      `gorm:"primaryKey"`
	OrgID          snowflake.ID      `gorm:"not null;index"`
	BillingCycleID snowflake.ID      `gorm:"not null;index;uniqueIndex:ux_rating_billing_cycle"`
	SubscriptionID snowflake.ID      `gorm:"not null;index"`
	TotalAmount    int64             `gorm:"not null;default:0"`
	Currency       string            `gorm:"type:text;not null"`
	Metadata       datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt      time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (RatingResult) TableName() string { return "rating_results" }

// RatingResultItem stores line-level rating breakdown.
type RatingResultItem struct {
	ID                 snowflake.ID      `gorm:"primaryKey"`
	OrgID              snowflake.ID      `gorm:"not null;index"`
	RatingResultID     snowflake.ID      `gorm:"not null;index"`
	SubscriptionItemID snowflake.ID      `gorm:"not null;index"`
	Quantity           int64             `gorm:"not null"`
	UnitAmount         int64             `gorm:"not null"`
	Amount             int64             `gorm:"not null"`
	TierBreakdown      datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	Metadata           datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt          time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (RatingResultItem) TableName() string { return "rating_result_items" }
