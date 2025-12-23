// Package domain contains core types for the auth service.
package domain

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/datatypes"
)

// User represents a system user account.
type User struct {
	ID                  snowflake.ID      `gorm:"primaryKey"`
	Username            string            `gorm:"type:text;not null;uniqueIndex"`
	Email               string            `gorm:"column:email;uniqueIndex"`
	PasswordHash        *string           `gorm:"type:text"`
	IsDefault           bool              `gorm:"column:is_default"`
	LastPasswordChanged *time.Time        `gorm:"column:last_password_changed"`
	Metadata            datatypes.JSONMap `gorm:"type:jsonb;not null;default:'{}'"`
	CreatedAt           time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time         `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (User) TableName() string { return "users" }

// Session represents a persisted login session.
type Session struct {
	ID               snowflake.ID `gorm:"primaryKey"`
	UserID           snowflake.ID `gorm:"column:user_id;not null;index"`
	SessionTokenHash string       `gorm:"column:session_token_hash;type:text;not null;uniqueIndex"`
	UserAgent        string       `gorm:"column:user_agent;type:text"`
	IPAddress        string       `gorm:"column:ip_address;type:text"`
	ExpiresAt        time.Time    `gorm:"column:expires_at;not null;index"`
	RevokedAt        *time.Time   `gorm:"column:revoked_at"`
	CreatedAt        time.Time    `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	LastSeenAt       time.Time    `gorm:"column:last_seen_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName sets the database table name.
func (Session) TableName() string { return "sessions" }

// SessionView is returned to clients without exposing token values.
type SessionView struct {
	Metadata map[string]any `json:"metadata"`
}
