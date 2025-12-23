package customer

import (
	"github.com/smallbiznis/valora/internal/customer/repository"
	"github.com/smallbiznis/valora/internal/customer/service"
	"go.uber.org/fx"
)

var Module = fx.Module("customer.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.New),
)
