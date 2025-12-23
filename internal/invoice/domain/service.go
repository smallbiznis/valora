package domain

import (
	"context"

	"github.com/smallbiznis/valora/pkg/db/pagination"
)

type ListInvoiceRequest struct{
	OrgID string
}

type ListInvoiceResponse struct {
	pagination.PageInfo
	Invoices []Invoice `json:"invoices"`
}

type Service interface {
	List(context.Context, ListInvoiceRequest) (ListInvoiceResponse, error)
	GetByID(ctx context.Context, id string) (Invoice, error)
}
