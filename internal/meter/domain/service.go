package domain

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Response, error)
	List(ctx context.Context, organizationID string) ([]Response, error)
	GetByID(ctx context.Context, organizationID string, id string) (*Response, error)
	GetByCode(ctx context.Context, organizationID string, code string) (*Response, error)
	Update(ctx context.Context, req UpdateRequest) (*Response, error)
	Delete(ctx context.Context, organizationID string, id string) error
}

type CreateRequest struct {
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Aggregation    string `json:"aggregation_type"`
	Unit           string `json:"unit"`
	Active         *bool  `json:"active"`
}

type UpdateRequest struct {
	OrganizationID string  `json:"organization_id"`
	ID             string  `json:"id"`
	Name           *string `json:"name,omitempty"`
	Aggregation    *string `json:"aggregation_type,omitempty"`
	Unit           *string `json:"unit,omitempty"`
	Active         *bool   `json:"active,omitempty"`
}

type Response struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Aggregation    string    `json:"aggregation"`
	Unit           string    `json:"unit"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var (
	ErrInvalidOrganization = errors.New("invalid_organization")
	ErrInvalidCode         = errors.New("invalid_code")
	ErrInvalidName         = errors.New("invalid_name")
	ErrInvalidAggregation  = errors.New("invalid_aggregation_type")
	ErrInvalidUnit         = errors.New("invalid_unit")
	ErrNotFound            = errors.New("not_found")
	ErrInvalidID           = errors.New("invalid_id")
)

func ParseID(value string) (snowflake.ID, error) {
	return snowflake.ParseString(value)
}
