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
	OrganizationID     string         `json:"organization_id"`
	PriceID            string         `json:"price_id"`
	MeterID            *string        `json:"meter_id"`
	Currency           string         `json:"currency"`
	UnitAmountCents    int64          `json:"unit_amount_cents"`
	MinimumAmountCents *int64         `json:"minimum_amount_cents"`
	MaximumAmountCents *int64         `json:"maximum_amount_cents"`
	Metadata           map[string]any `json:"metadata"`
}

type Response struct {
	ID                 string    `json:"id"`
	OrganizationID     string    `json:"organization_id"`
	PriceID            string    `json:"price_id"`
	MeterID            *string   `json:"meter_id,omitempty"`
	Currency           string    `json:"currency"`
	UnitAmountCents    int64     `json:"unit_amount_cents"`
	MinimumAmountCents *int64    `json:"minimum_amount_cents,omitempty"`
	MaximumAmountCents *int64    `json:"maximum_amount_cents,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

var (
	ErrInvalidOrganization = errors.New("invalid_organization")
	ErrInvalidPrice        = errors.New("invalid_price")
	ErrInvalidCurrency     = errors.New("invalid_currency")
	ErrInvalidUnitAmount   = errors.New("invalid_unit_amount")
	ErrInvalidMinAmount    = errors.New("invalid_minimum_amount")
	ErrInvalidMaxAmount    = errors.New("invalid_maximum_amount")
	ErrInvalidMeterID      = errors.New("invalid_meter_id")
	ErrInvalidID           = errors.New("invalid_id")
	ErrNotFound            = errors.New("not_found")
)
