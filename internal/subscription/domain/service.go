package domain

import (
	"context"
	"errors"
	"time"

	"github.com/smallbiznis/valora/pkg/db/pagination"
)

type ListSubscriptionRequest struct {
	OrgID     string
	Status    string
	PageToken string
	PageSize  int32
}

type ListSubscriptionResponse struct {
	pagination.PageInfo
	Subscriptions []Subscription `json:"subscriptions"`
}

type CreateSubscriptionItemRequest struct {
	PriceID  string `json:"price_id"`
	Quantity int8   `json:"quantity,omitempty"`
}

type CreateSubscriptionRequest struct {
	OrganizationID string                          `json:"organization_id"`
	CustomerID     string                          `json:"customer_id"`
	Items          []CreateSubscriptionItemRequest `json:"items"`
	Metadata       map[string]any                  `json:"metadata,omitempty"`
}

type GetActiveByCustomerIDRequest struct {
	OrgID      string
	CustomerID string
}

type GetSubscriptionItemRequest struct {
	OrgID          string
	SubscriptionID string
	MeterID        string
	MeterCode      string
}

//go:generate mockgen -source=service.go -destination=./mocks/mock_service.go -package=mocks
type Service interface {
	List(context.Context, ListSubscriptionRequest) (ListSubscriptionResponse, error)
	Create(context.Context, CreateSubscriptionRequest) (CreateSubscriptionResponse, error)
	GetByID(context.Context, string) (Subscription, error)
	GetActiveByCustomerID(context.Context, GetActiveByCustomerIDRequest) (Subscription, error)
	GetSubscriptionItem(context.Context, GetSubscriptionItemRequest) (SubscriptionItem, error)
}

type CreateSubscriptionItemResponse struct {
	ID                string   `json:"id"`
	PriceID           string   `json:"price_id"`
	PriceCode         *string  `json:"price_code,omitempty"`
	MeterID           *string  `json:"meter_id,omitempty"`
	MeterCode         *string  `json:"meter_code,omitempty"`
	Quantity          int8     `json:"quantity"`
	BillingMode       string   `json:"billing_mode"`
	UsageBehavior     *string  `json:"usage_behavior,omitempty"`
	BillingThreshold  *float64 `json:"billing_threshold,omitempty"`
	ProrationBehavior *string  `json:"proration_behavior,omitempty"`
}

type CreateSubscriptionResponse struct {
	ID             string                           `json:"id"`
	OrganizationID string                           `json:"organization_id"`
	CustomerID     string                           `json:"customer_id"`
	Status         SubscriptionStatus               `json:"status"`
	CollectionMode SubscriptionCollectionMode       `json:"collection_mode"`
	StartAt        time.Time                        `json:"start_at"`
	Items          []CreateSubscriptionItemResponse `json:"items"`
	Metadata       map[string]any                   `json:"metadata,omitempty"`
}

var (
	ErrInvalidOrganization      = errors.New("invalid_organization")
	ErrInvalidCustomer          = errors.New("invalid_customer")
	ErrInvalidSubscription      = errors.New("invalid_subscription")
	ErrInvalidMeterID           = errors.New("invalid_meter_id")
	ErrInvalidMeterCode         = errors.New("invalid_meter_code")
	ErrInvalidStatus            = errors.New("invalid_status")
	ErrInvalidCollectionMode    = errors.New("invalid_collection_mode")
	ErrInvalidBillingCycleType  = errors.New("invalid_billing_cycle_type")
	ErrInvalidStartAt           = errors.New("invalid_start_at")
	ErrInvalidPeriod            = errors.New("invalid_period")
	ErrInvalidItems             = errors.New("invalid_items")
	ErrInvalidQuantity          = errors.New("invalid_quantity")
	ErrInvalidPrice             = errors.New("invalid_price")
	ErrMultipleFlatPrices       = errors.New("multiple_flat_prices_not_allowed")
	ErrSubscriptionNotFound     = errors.New("subscription_not_found")
	ErrSubscriptionItemNotFound = errors.New("subscription_item_not_found")
)
