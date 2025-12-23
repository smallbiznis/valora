package auth

import (
	"github.com/smallbiznis/valora/internal/auth/repository"
	"github.com/smallbiznis/valora/internal/auth/service"
	"go.uber.org/fx"
)

var Module = fx.Module("auth.service",
	fx.Provide(repository.New),
	fx.Provide(service.New),
)
