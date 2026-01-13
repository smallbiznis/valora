package tax

import (
	"github.com/smallbiznis/valora/internal/tax/repository"
	"github.com/smallbiznis/valora/internal/tax/service"
	"go.uber.org/fx"
)

var Module = fx.Module("tax.service",
	fx.Provide(repository.NewRepository),
	fx.Provide(service.NewResolver),
	fx.Provide(service.NewService),
)
