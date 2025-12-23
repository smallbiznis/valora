package invoice

import (
	"github.com/smallbiznis/valora/internal/invoice/service"
	"go.uber.org/fx"
)

var Module = fx.Module("invoice.service",
	fx.Provide(service.NewService),
)
