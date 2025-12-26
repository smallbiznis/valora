package service

import (
	"context"
	"strings"

	"github.com/bwmarrin/snowflake"
	invoicedomain "github.com/smallbiznis/valora/internal/invoice/domain"
	"github.com/smallbiznis/valora/internal/orgcontext"
	"github.com/smallbiznis/valora/pkg/db/option"
	"github.com/smallbiznis/valora/pkg/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceParam struct {
	fx.In

	DB  *gorm.DB
	Log *zap.Logger
}

type Service struct {
	db  *gorm.DB
	log *zap.Logger

	invoicerepo repository.Repository[invoicedomain.Invoice]
}

func NewService(p ServiceParam) invoicedomain.Service {
	return &Service{
		db:          p.DB,
		log:         p.Log.Named("invoice.service"),
		invoicerepo: repository.ProvideStore[invoicedomain.Invoice](p.DB),
	}
}

func (s *Service) List(ctx context.Context, req invoicedomain.ListInvoiceRequest) (invoicedomain.ListInvoiceResponse, error) {
	_ = req
	orgID, err := s.orgIDFromContext(ctx)
	if err != nil {
		return invoicedomain.ListInvoiceResponse{}, err
	}

	items, err := s.invoicerepo.Find(ctx, &invoicedomain.Invoice{OrgID: orgID},
		option.WithSortBy(option.QuerySortBy{Allow: map[string]bool{"created_at": true}}),
	)
	if err != nil {
		return invoicedomain.ListInvoiceResponse{}, err
	}

	invoices := make([]invoicedomain.Invoice, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		invoices = append(invoices, *item)
	}

	return invoicedomain.ListInvoiceResponse{Invoices: invoices}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (invoicedomain.Invoice, error) {
	orgID, err := s.orgIDFromContext(ctx)
	if err != nil {
		return invoicedomain.Invoice{}, err
	}

	invoiceID, err := snowflake.ParseString(strings.TrimSpace(id))
	if err != nil {
		return invoicedomain.Invoice{}, err
	}

	item, err := s.invoicerepo.FindOne(ctx, &invoicedomain.Invoice{ID: invoiceID, OrgID: orgID})
	if err != nil {
		return invoicedomain.Invoice{}, err
	}
	if item == nil {
		return invoicedomain.Invoice{}, gorm.ErrRecordNotFound
	}

	return *item, nil
}

func (s *Service) orgIDFromContext(ctx context.Context) (snowflake.ID, error) {
	orgID, ok := orgcontext.OrgIDFromContext(ctx)
	if !ok || orgID == 0 {
		return 0, invoicedomain.ErrInvalidOrganization
	}
	return snowflake.ID(orgID), nil
}
