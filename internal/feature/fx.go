package feature

import (
	"github.com/smallbiznis/valora/internal/feature/repository"
	"github.com/smallbiznis/valora/internal/feature/service"
	"go.uber.org/fx"
)

var Module = fx.Module("feature.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.New),
)
