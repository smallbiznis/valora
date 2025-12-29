package usage

import (
	"github.com/smallbiznis/valora/internal/usage/service"
	"go.uber.org/fx"
)

var Module = fx.Module("usage.service",
	fx.Provide(service.NewService),
)
