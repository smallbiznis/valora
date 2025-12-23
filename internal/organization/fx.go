package organization

import (
	"github.com/smallbiznis/valora/internal/organization/event"
	"github.com/smallbiznis/valora/internal/organization/repository"
	"github.com/smallbiznis/valora/internal/organization/service"
	"go.uber.org/fx"
)

var Module = fx.Module("organization.service",
	fx.Provide(repository.NewRepository),
	fx.Provide(event.NewOutboxPublisher),
	fx.Provide(service.NewService),
)
