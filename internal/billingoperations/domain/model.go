package domain

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/datatypes"
)

type BillingActionRecord struct {
	ID             snowflake.ID
	OrgID          snowflake.ID
	EntityType     string
	EntityID       snowflake.ID
	ActionType     string
	ActionBucket   time.Time
	IdempotencyKey string
	Metadata       datatypes.JSONMap
	ActorType      string
	ActorID        string
	CreatedAt      time.Time
}

type BillingAssignmentRecord struct {
	ID                  snowflake.ID
	OrgID               snowflake.ID
	EntityType          string
	EntityID            snowflake.ID
	AssignedTo          string
	AssignedAt          time.Time
	AssignmentExpiresAt time.Time
	Status              string
	ReleasedAt          sql.NullTime
	ReleasedBy          sql.NullString
	ReleaseReason       sql.NullString
	ResolvedAt          sql.NullTime
	ResolvedBy          sql.NullString
	BreachedAt          sql.NullTime
	BreachLevel         sql.NullString
	LastActionAt        sql.NullTime
	SnapshotMetadata    datatypes.JSON
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (BillingAssignmentRecord) TableName() string {
	return "billing_operation_assignments"
}

type BillingActionLookup struct {
	ID snowflake.ID `gorm:"column:id"`
}
