package main

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/audit"
	"github.com/smallbiznis/valora/internal/authorization"
	"github.com/smallbiznis/valora/internal/billingdashboard/rollup"
	"github.com/smallbiznis/valora/internal/billingoperations"
	"github.com/smallbiznis/valora/internal/clock"
	"github.com/smallbiznis/valora/internal/config"
	"github.com/smallbiznis/valora/internal/feature"
	"github.com/smallbiznis/valora/internal/invoice"
	"github.com/smallbiznis/valora/internal/invoicetemplate"
	"github.com/smallbiznis/valora/internal/ledger"
	"github.com/smallbiznis/valora/internal/meter"
	"github.com/smallbiznis/valora/internal/observability"
	"github.com/smallbiznis/valora/internal/price"
	"github.com/smallbiznis/valora/internal/priceamount"
	"github.com/smallbiznis/valora/internal/pricetier"
	"github.com/smallbiznis/valora/internal/product"
	"github.com/smallbiznis/valora/internal/productfeature"
	"github.com/smallbiznis/valora/internal/rating"
	"github.com/smallbiznis/valora/internal/scheduler"
	"github.com/smallbiznis/valora/internal/subscription"
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

		// Domain services required by scheduler
		scheduler.Module,
		rating.Module,
		invoice.Module,
		ledger.Module,
		subscription.Module,
		audit.Module,
		authorization.Module,
		billingoperations.Module,
		rollup.Module,
		
		// Transitive dependencies (invoice needs product/price etc)
		product.Module,
		productfeature.Module,
		feature.Module,
		price.Module,
		priceamount.Module,
		pricetier.Module,
		invoicetemplate.Module,
		meter.Module,

		// No server module!
		fx.Invoke(StartScheduler),
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

func StartScheduler(lc fx.Lifecycle, s *scheduler.Scheduler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.RunForever(context.Background())
			return nil
		},
	})
}
