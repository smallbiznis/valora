package subscription

import (
	"github.com/smallbiznis/valora/internal/subscription/repository"
	"github.com/smallbiznis/valora/internal/subscription/service"
	"go.uber.org/fx"
)

var Module = fx.Module("subscription.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.NewService),
)
