package product

import (
	"github.com/smallbiznis/valora/internal/product/repository"
	"github.com/smallbiznis/valora/internal/product/service"
	"go.uber.org/fx"
)

var Module = fx.Module("product.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.New),
)
