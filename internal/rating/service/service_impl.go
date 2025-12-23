package service

import (
	ratingdomain "github.com/smallbiznis/valora/internal/rating/domain"
	"github.com/smallbiznis/valora/pkg/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	log *zap.Logger

	ratingrepo     repository.Repository[ratingdomain.RatingResult]
	ratingitemrepo repository.Repository[ratingdomain.RatingResultItem]
}

type ServiceParam struct {
	fx.In

	DB  *gorm.DB
	Log *zap.Logger
}

func NewService(p ServiceParam) ratingdomain.Service {
	return &Service{
		db:  p.DB,
		log: p.Log.Named("rating.service"),

		ratingrepo:     repository.ProvideStore[ratingdomain.RatingResult](p.DB),
		ratingitemrepo: repository.ProvideStore[ratingdomain.RatingResultItem](p.DB),
	}
}
