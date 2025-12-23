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
	PaidAt         *time.Time        `gorm:""`
	Metadata       datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt      time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (Invoice) TableName() string { return "invoices" }

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
