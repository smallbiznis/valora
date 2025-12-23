package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/logger"
	"github.com/smallbiznis/valora/internal/server"
	"github.com/smallbiznis/valora/pkg/db"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		logger.Module,
		fx.Provide(func() *snowflake.Node {
			node, err := snowflake.NewNode(1)
			if err != nil {
				panic(err)
			}
			return node
		}),
		db.Module,
		server.Module,
	)
	app.Run()
}