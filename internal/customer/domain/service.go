package domain

import (
	"context"
	"errors"

	"github.com/smallbiznis/valora/pkg/db/pagination"
)

type ListCustomerRequest struct {
	OrgID     string
	PageToken string
	PageSize  int32
}

type ListCustomerResponse struct {
	pagination.PageInfo
	Customers []Customer `json:"customers"`
}

type CreateCustomerRequest struct {
	OrganizationID string
	Name           string
	Email          string
}

type GetCustomerRequest struct {
	ID    string
	OrgID string
}

type Service interface {
	Create(context.Context, CreateCustomerRequest) (Customer, error)
	List(context.Context, ListCustomerRequest) (ListCustomerResponse, error)
	GetByID(context.Context, GetCustomerRequest) (Customer, error)
}

var (
	ErrInvalidOrganization = errors.New("invalid_organization")
	ErrInvalidName         = errors.New("invalid_name")
	ErrInvalidEmail        = errors.New("invalid_email")
	ErrInvalidID           = errors.New("invalid_id")
	ErrNotFound            = errors.New("not_found")
)
