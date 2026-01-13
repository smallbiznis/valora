package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/apikey"
	"github.com/smallbiznis/valora/internal/auth"
	"github.com/smallbiznis/valora/internal/clock"
	"github.com/smallbiznis/valora/internal/config"
	"github.com/smallbiznis/valora/internal/meter"
	"github.com/smallbiznis/valora/internal/observability"
	"github.com/smallbiznis/valora/internal/ratelimit"
	"github.com/smallbiznis/valora/internal/server"
	"github.com/smallbiznis/valora/internal/usage"
	"github.com/smallbiznis/valora/pkg/db"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module,
		observability.Module,
		fx.Provide(RegisterSnowflake),
		db.Module,
		clock.Module,

		// Core dependencies for API
		auth.Module,   // For API Key validation logic
		apikey.Module, // For API Key domain
		meter.Module,
		usage.Module,
		ratelimit.Module,

		fx.Provide(server.NewEngine),
		fx.Provide(server.NewServer),
		fx.Invoke(func(s *server.Server) {
			s.RegisterAPIRoutes()
		}),
		fx.Invoke(server.RunHTTP),
	)
	app.Run()
}

func RegisterSnowflake() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}
