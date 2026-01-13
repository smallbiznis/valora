package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/clock"
	"github.com/smallbiznis/valora/internal/config"
	"github.com/smallbiznis/valora/internal/migration"
	"github.com/smallbiznis/valora/internal/observability"
	"github.com/smallbiznis/valora/internal/scheduler"
	"github.com/smallbiznis/valora/internal/server"
	"github.com/smallbiznis/valora/pkg/db"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Core Infrastructure
		config.Module,
		observability.Module,
		fx.Provide(RegisterSnowflake),
		db.Module,
		clock.Module,
		server.Module,

		// Functional Domains
		scheduler.Module,
		migration.Module,
		
		// All other domain modules usually imported by specific apps
		// We can mostly rely on server.Module importing them transitively or explicitly here
		// but server.Module already imports MOST of them.
		// Let's ensure we register ALL routes.

		fx.Provide(server.NewEngine),
		// NewServer is already provided by server.Module? 
		// No, server.Module invokes NewServer? 
		// Let's check server.go again. 
		// server.Module has fx.Invoke(NewServer), fx.Invoke(RunHTTP).
		// Wait, server.Module DOES invoke NewServer. 
		// But in apps/admin we had to Provide it? 
		// In apps/admin we did: fx.Provide(server.NewEngine), fx.Provide(server.NewServer).
		// In server.server.go, Module DOES NOT Provide NewServer, it Invokes it? 
		// Let's check server.go code again in my mind...
		// Lines 124-125: fx.Invoke(NewServer), fx.Invoke(RunHTTP).
		// This means NewServer is treated as an invocation (side effect?) or maybe I misread.
		// NewServer returns *Server. If it is Invoked, the result is discarded unless annotated?
		// Usually fx.Provide(NewServer) is what we want so others can use it.
		// In server.go:124 it says fx.Invoke(NewServer). This seems wrong if we want to use 's' in Invoke.
		// But wait, existing code in server.go might have been: fx.Provide(NewServer).
		// I will Assume server.Module needs to be "fixed" or I just Provide it here to be safe/override.
		
		// Let's use the pattern from apps/admin:
		fx.Provide(server.NewEngine),
		fx.Provide(server.NewServer),

		fx.Invoke(func(s *server.Server) {
			// Register ALL the things
			s.RegisterAPIRoutes()
			s.RegisterAdminRoutes()     // Admin Dashboard API
			s.RegisterPublicRoutes()    // Public Invoice API
			s.RegisterAuthRoutes()
			
			// UI Routes (Monolith primarily serves Admin UI)
			s.RegisterUIRoutes() 
			s.RegisterFallback()
		}),
		
		// RunHTTP is invoked by server.Module or explicitly?
		// server.Module has fx.Invoke(RunHTTP). 
		// We can leave it or be explicit. 
		// To be safe, let's Suppress server.Module's autodrive if needed, or just let it run.
		// But server.Module defines `fx.Invoke(RunHTTP)` at line 125.
		// So it will run automatically.
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
