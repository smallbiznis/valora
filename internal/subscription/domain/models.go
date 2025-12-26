// Package domain contains persistence models for subscriptions and billing cycles.
package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/datatypes"
)

// SubscriptionStatus represents lifecycle states for a subscription.
type SubscriptionStatus string

const (
	SubscriptionStatusPending   SubscriptionStatus = "PENDING"
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusPastDue   SubscriptionStatus = "PAST_DUE"
	SubscriptionStatusSuspended SubscriptionStatus = "SUSPENDED"
	SubscriptionStatusCanceled  SubscriptionStatus = "CANCELED"
)

type SubscriptionCollectionMode string

var (
	SendInvoice         SubscriptionCollectionMode = "SEND_INVOICE"
	ChargeAutomatically SubscriptionCollectionMode = "CHARGE_AUTOMATICALLY"
)

// Subscription captures a customer's billing agreement.
type Subscription struct {
	ID                     snowflake.ID               `gorm:"primaryKey"`
	OrgID                  snowflake.ID               `gorm:"not null;index"`
	CustomerID             snowflake.ID               `gorm:"not null;index"`
	Status                 SubscriptionStatus         `gorm:"type:text;not null"`
	CollectionMode         SubscriptionCollectionMode `gorm:"type:text;not null"`
	StartAt                time.Time                  `gorm:"not null"`
	EndAt                  *time.Time                 `gorm:""`
	TrialStartsAt          *time.Time                 `gorm:""`
	TrialEndsAt            *time.Time                 `gorm:""`
	CancelAt               *time.Time                 `gorm:""`
	CancelAtPeriodEnd      bool                       `gorm:"not null;default:false"`
	CanceledAt             *time.Time                 `gorm:""`
	BillingAnchorDay       *int16                     `gorm:"type:smallint"`
	BillingCycleType       string                     `gorm:"type:text;not null"`
	DefaultPaymentTermDays *int                       `gorm:""`
	DefaultCurrency        *string                    `gorm:"type:text"`
	DefaultTaxBehavior     *string                    `gorm:"type:text"`
	Metadata               datatypes.JSONMap          `gorm:"type:jsonb"`
	CreatedAt              time.Time                  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt              time.Time                  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (Subscription) TableName() string { return "subscriptions" }

// SetTrial sets the trial period for the subscription.
func (s *Subscription) SetTrial(trialDurationDays int) *Subscription {
	now := time.Now().UTC()
	s.TrialStartsAt = &now
	trialEnd := now.AddDate(0, 0, trialDurationDays)
	s.TrialEndsAt = &trialEnd
	return s
}

// IsTrial checks if the subscription is currently in a trial period.
func (s *Subscription) IsTrial(now time.Time) bool {
	return s.TrialEndsAt != nil && now.Before(*s.TrialEndsAt)
}

// RecalculateState updates the subscription status based on current time and invoice status.
func (s *Subscription) RecalculateState(
	now time.Time,
	hasOverdueInvoice bool,
	overdueDuration time.Duration,
	gracePeriod time.Duration,
) {
	if s.CanceledAt != nil {
		s.Status = SubscriptionStatusCanceled
		return
	}

	if now.Before(s.StartAt) {
		s.Status = SubscriptionStatusPending
		return
	}

	if hasOverdueInvoice {
		if overdueDuration <= gracePeriod {
			s.Status = SubscriptionStatusPastDue
		} else {
			s.Status = SubscriptionStatusSuspended
		}
		return
	}

	s.Status = SubscriptionStatusActive
}

// SubscriptionItem associates subscriptions to prices/meters.
type SubscriptionItem struct {
	ID                snowflake.ID      `gorm:"primaryKey"`
	OrgID             snowflake.ID      `gorm:"not null;index"`
	SubscriptionID    snowflake.ID      `gorm:"not null;index"`
	PriceID           snowflake.ID      `gorm:"not null;index"`
	PriceCode         *string           `gorm:"type:text"`
	MeterID           *snowflake.ID     `gorm:"not null;index"`
	MeterCode         *string           `gorm:"type:text"`
	Quantity          int8              `gorm:"column:quantity"`
	BillingMode       string            `gorm:"type:text;not null"`
	UsageBehavior     *string           `gorm:"type:text"`
	BillingThreshold  *float64          `gorm:""`
	ProrationBehavior *string           `gorm:"type:text"`
	NextPeriodStart   *time.Time        `gorm:""`
	NextPeriodEnd     *time.Time        `gorm:""`
	Metadata          datatypes.JSONMap `gorm:"type:jsonb"`
	CreatedAt         time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (SubscriptionItem) TableName() string { return "subscription_items" }
