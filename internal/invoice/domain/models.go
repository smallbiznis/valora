// Package domain contains persistence models for invoicing.
package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/datatypes"
)

// InvoiceStatus represents invoice lifecycle states.
type InvoiceStatus string

const (
	InvoiceStatusDraft         InvoiceStatus = "DRAFT"
	InvoiceStatusOpen          InvoiceStatus = "OPEN"
	InvoiceStatusPartiallyPaid InvoiceStatus = "PARTIALLY_PAID"
	InvoiceStatuusOverdue      InvoiceStatus = "OVERDUE"
	InvoiceStatusPaid          InvoiceStatus = "PAID"
	InvoiceStatusVoid          InvoiceStatus = "VOID"
	InvoiceStatusUncollectible InvoiceStatus = "UNCOLLECTIBLE"
)

// Invoice represents a generated invoice.
type Invoice struct {
	ID             snowflake.ID      `gorm:"primaryKey"`
	OrgID          snowflake.ID      `gorm:"not null;index"`
	BillingCycleID snowflake.ID      `gorm:"not null;index;uniqueIndex:ux_invoice_billing_cycle"`
	SubscriptionID snowflake.ID      `gorm:"not null;index"`
	CustomerID     snowflake.ID      `gorm:"not null;index"`
	Status         InvoiceStatus     `gorm:"type:text;not null;default:'DRAFT'"`
	TotalAmount    int64             `gorm:"not null;default:0"`
	Currency       string            `gorm:"type:text;not null"`
	IssuedAt       *time.Time        `gorm:""`
	DueAt          *time.Time        `gorm:""`
	VoidedAt       *time.Time        `gorm:""`
	PaidAt         *time.Time        `gorm:""`
	Metadata       datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt      time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (Invoice) TableName() string { return "invoices" }

// RecalculateState updates the invoice status based on the total paid amount and current time.
func (i *Invoice) RecalculateState(now time.Time, totalPaid int64) {
	// VOID is terminal
	if i.VoidedAt != nil {
		i.Status = InvoiceStatusVoid
		return
	}

	// Fully paid
	if totalPaid >= i.TotalAmount {
		i.Status = InvoiceStatusPaid
		if i.PaidAt == nil {
			i.PaidAt = &now
		}
		return
	}

	// Partial payment
	if totalPaid > 0 {
		if i.DueAt != nil && now.After(*i.DueAt) {
			i.Status = InvoiceStatuusOverdue
		} else {
			i.Status = InvoiceStatusPartiallyPaid
		}
		return
	}

	// No payment yet
	if i.DueAt != nil && now.After(*i.DueAt) {
		i.Status = InvoiceStatuusOverdue
		return
	}

	i.Status = InvoiceStatusOpen
}

// InvoiceItem represents a line on an invoice.
type InvoiceItem struct {
	ID                 snowflake.ID      `gorm:"primaryKey"`
	OrgID              snowflake.ID      `gorm:"not null;index"`
	InvoiceID          snowflake.ID      `gorm:"not null;index"`
	RatingResultItemID *snowflake.ID     `gorm:"index"`
	SubscriptionItemID *snowflake.ID     `gorm:"index"`
	Description        string            `gorm:"type:text"`
	Quantity           int64             `gorm:"not null"`
	UnitAmount         int64             `gorm:"not null"`
	Amount             int64             `gorm:"not null"`
	Metadata           datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt          time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (InvoiceItem) TableName() string { return "invoice_items" }
