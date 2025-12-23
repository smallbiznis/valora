package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type InvitationStatus string

var (
	Pending InvitationStatus = "PENDING"
	Completed InvitationStatus = "COMPLETED"
)

type Invitation struct {
	ID        snowflake.ID      `gorm:"primaryKey" json:"id"`
	OrgID     snowflake.ID      `gorm:"not null;index" json:"organization_id"`
	Email     string            `gorm:"not null" json:"email"`
	Role string `gorm:"column:role"`
	Code string `gorm:"column:code"`
	Status string `gorm:"column:status"`
	CreatedAt time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}