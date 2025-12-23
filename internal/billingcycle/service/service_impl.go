package service

import (
	"context"

	billingcycledomain "github.com/smallbiznis/valora/internal/billingcycle/domain"
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

	billingcyclerepo repository.Repository[billingcycledomain.BillingCycle]
}

func NewService(p ServiceParam) billingcycledomain.Service {
	return &Service{
		db:  p.DB,
		log: p.Log.Named("billingcycle.service"),

		billingcyclerepo: repository.ProvideStore[billingcycledomain.BillingCycle](p.DB),
	}
}

func (s *Service) List(ctx context.Context, req billingcycledomain.ListBillingCycleRequest) (billingcycledomain.ListBillingCycleResponse, error) {
	_ = ctx
	_ = req
	return billingcycledomain.ListBillingCycleResponse{}, nil
}
