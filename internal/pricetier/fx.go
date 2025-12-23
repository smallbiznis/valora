package pricetier

import (
	"github.com/smallbiznis/valora/internal/pricetier/repository"
	"github.com/smallbiznis/valora/internal/pricetier/service"
	"go.uber.org/fx"
)

var Module = fx.Module("pricetier.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.New),
)
