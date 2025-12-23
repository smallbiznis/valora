package domain

import (
	"context"

	"github.com/smallbiznis/valora/pkg/db/pagination"
)

type ListBillingCycleRequest struct{}

type ListBillingCycleResponse struct {
	pagination.PageInfo
	BillingCycle []BillingCycle `json:"billing_cycles"`
}

type Service interface {
	List(context.Context, ListBillingCycleRequest) (ListBillingCycleResponse, error)
}
