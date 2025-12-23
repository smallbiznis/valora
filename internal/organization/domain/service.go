package domain

import (
	"context"
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
)

const (
	RoleOwner  = "OWNER"
	RoleMember = "MEMBER"
)

type Service interface {
	Create(ctx context.Context, userID snowflake.ID, req CreateOrganizationRequest) (*OrganizationResponse, error)
	GetByID(ctx context.Context, id string) (*OrganizationResponse, error)
	ListOrganizationsByUser(ctx context.Context, userID snowflake.ID) ([]OrganizationListResponseItem, error)
}

type CreateOrganizationRequest struct {
	Name            string
	CountryCode     string
	TimezoneName    string
}

type OrganizationResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Slug string `json:"slug"`
	CountryCode     string `json:"country_code"`
	TimezoneName    string `json:"timezone_name"`
}

type OrganizationListResponseItem struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrInvalidName         = errors.New("invalid_name")
	ErrInvalidCountry      = errors.New("invalid_country")
	ErrInvalidTimezone     = errors.New("invalid_timezone")
	ErrInvalidCurrency     = errors.New("invalid_currency")
	ErrInvalidUser         = errors.New("invalid_user")
	ErrInvalidOrganization = errors.New("invalid_organization")
)
