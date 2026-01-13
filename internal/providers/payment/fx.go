package payment

import (
	"github.com/smallbiznis/valora/internal/providers/payment/repository"
	"github.com/smallbiznis/valora/internal/providers/payment/service"
	"go.uber.org/fx"
)

var Module = fx.Module("paymentprovider.service",
	fx.Provide(repository.Provide),
	fx.Provide(service.New),
)
