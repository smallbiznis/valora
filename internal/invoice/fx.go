package invoice

import (
	"github.com/smallbiznis/valora/internal/invoice/render"
	"github.com/smallbiznis/valora/internal/invoice/service"
	"github.com/smallbiznis/valora/internal/tax"
	"go.uber.org/fx"
)

var Module = fx.Module("invoice.service",
	tax.Module,
	fx.Provide(render.NewRenderer),
	fx.Provide(service.NewService),
)
