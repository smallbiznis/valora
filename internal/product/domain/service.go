package domain

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Response, error)
	List(ctx context.Context, organizationID string) ([]Response, error)
	Get(ctx context.Context, organizationID string, id string) (*Response, error)
}

type CreateRequest struct {
	OrganizationID string         `json:"organization_id"`
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	Description    *string        `json:"description"`
	Active         *bool          `json:"active"`
	Metadata       map[string]any `json:"metadata"`
}

type Response struct {
	ID             string         `json:"id"`
	OrganizationID string         `json:"organization_id"`
	Code           string         `json:"code"`
	Name           string         `json:"name"`
	Description    *string        `json:"description,omitempty"`
	Active         bool           `json:"active"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

var (
	ErrInvalidOrganization = errors.New("invalid_organization")
	ErrInvalidCode         = errors.New("invalid_code")
	ErrInvalidName         = errors.New("invalid_name")
	ErrNotFound            = errors.New("not_found")
	ErrInvalidID           = errors.New("invalid_id")
)
