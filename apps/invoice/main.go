package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/config"
	"github.com/smallbiznis/valora/internal/observability"
	"github.com/smallbiznis/valora/internal/payment"
	paymentprovider "github.com/smallbiznis/valora/internal/providers/payment"
	"github.com/smallbiznis/valora/internal/publicinvoice"
	"github.com/smallbiznis/valora/internal/ratelimit"
	"github.com/smallbiznis/valora/internal/server"
	"github.com/smallbiznis/valora/pkg/db"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module,
		observability.Module,
		fx.Provide(RegisterSnowflake),
		db.Module,

		// Invoice Service dependencies
		publicinvoice.Module,
		payment.Module,         // For payment intent creation
		paymentprovider.Module, // For interfacing with Stripe/etc
		ratelimit.Module,       // Public endpoints are rate limited

		fx.Provide(server.NewEngine),
		fx.Provide(server.NewServer),
		fx.Invoke(func(s *server.Server) {
			s.RegisterPublicRoutes()
			s.RegisterFallback() // If this service also serves a public UI
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
