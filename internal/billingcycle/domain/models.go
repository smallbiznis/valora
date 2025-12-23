package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/datatypes"
)

// BillingCycleStatus represents rating/invoicing progress for a cycle.
type BillingCycleStatus string

const (
	BillingCycleStatusOpen   BillingCycleStatus = "OPEN"
	BillingCycleStatusClosed BillingCycleStatus = "CLOSED"
)

// BillingCycle represents a billing period for a subscription.
type BillingCycle struct {
	ID             snowflake.ID       `gorm:"primaryKey"`
	OrgID          snowflake.ID       `gorm:"not null;index"`
	SubscriptionID snowflake.ID       `gorm:"not null;index;uniqueIndex:ux_billing_cycle_period,priority:1"`
	PeriodStart    time.Time          `gorm:"not null;uniqueIndex:ux_billing_cycle_period,priority:2"`
	PeriodEnd      time.Time          `gorm:"not null;uniqueIndex:ux_billing_cycle_period,priority:3"`
	Status         BillingCycleStatus `gorm:"type:text;not null;default:'OPEN'"`
	RatedAt        *time.Time         `gorm:""`
	InvoicedAt     *time.Time         `gorm:""`
	ClosedAt       *time.Time         `gorm:""`
	Metadata       datatypes.JSONMap  `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt      time.Time          `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time          `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (BillingCycle) TableName() string { return "billing_cycles" }