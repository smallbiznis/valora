package productfeature

import (
	"github.com/smallbiznis/valora/internal/productfeature/repository"
	"github.com/smallbiznis/valora/internal/productfeature/service"
	"go.uber.org/fx"
)

var Module = fx.Module("productfeature.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.New),
)
