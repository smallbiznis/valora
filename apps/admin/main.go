package main

import (
	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/apikey"
	"github.com/smallbiznis/valora/internal/audit"
	"github.com/smallbiznis/valora/internal/auth"
	authlocal "github.com/smallbiznis/valora/internal/auth/local"
	authoauth2provider "github.com/smallbiznis/valora/internal/auth/oauth2provider"
	"github.com/smallbiznis/valora/internal/auth/session"
	"github.com/smallbiznis/valora/internal/authorization"
	"github.com/smallbiznis/valora/internal/billingdashboard"
	"github.com/smallbiznis/valora/internal/billingoperations"
	"github.com/smallbiznis/valora/internal/billingoverview"
	"github.com/smallbiznis/valora/internal/cloudmetrics"
	"github.com/smallbiznis/valora/internal/config"
	"github.com/smallbiznis/valora/internal/customer"
	"github.com/smallbiznis/valora/internal/events"
	"github.com/smallbiznis/valora/internal/feature"
	"github.com/smallbiznis/valora/internal/invoice"
	"github.com/smallbiznis/valora/internal/invoicetemplate"
	"github.com/smallbiznis/valora/internal/ledger"
	"github.com/smallbiznis/valora/internal/meter"
	"github.com/smallbiznis/valora/internal/observability"
	"github.com/smallbiznis/valora/internal/organization"
	"github.com/smallbiznis/valora/internal/payment"
	"github.com/smallbiznis/valora/internal/price"
	"github.com/smallbiznis/valora/internal/priceamount"
	"github.com/smallbiznis/valora/internal/pricetier"
	"github.com/smallbiznis/valora/internal/product"
	"github.com/smallbiznis/valora/internal/productfeature"
	"github.com/smallbiznis/valora/internal/providers/email"
	paymentprovider "github.com/smallbiznis/valora/internal/providers/payment"
	"github.com/smallbiznis/valora/internal/providers/pdf"
	"github.com/smallbiznis/valora/internal/ratelimit"
	"github.com/smallbiznis/valora/internal/rating"
	"github.com/smallbiznis/valora/internal/reference"
	"github.com/smallbiznis/valora/internal/server"
	"github.com/smallbiznis/valora/internal/subscription"
	"github.com/smallbiznis/valora/internal/usage"
	"github.com/smallbiznis/valora/pkg/db"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module,
		cloudmetrics.Module,
		observability.Module,
		fx.Provide(RegisterSnowflake),
		db.Module,

		// Admin needs almost everything
		authorization.Module,
		audit.Module,
		events.Module,
		auth.Module,
		authlocal.Module,
		authoauth2provider.Module,
		session.Module,
		apikey.Module,
		customer.Module,
		billingdashboard.Module,
		billingoperations.Module,
		email.Module,
		pdf.Module,
		billingoverview.Module,
		invoice.Module,
		invoicetemplate.Module,
		ledger.Module,
		meter.Module,
		organization.Module,
		price.Module,
		priceamount.Module,
		pricetier.Module,
		product.Module,
		productfeature.Module,
		feature.Module,
		payment.Module,
		paymentprovider.Module,
		reference.Module,
		rating.Module,
		ratelimit.Module,
		subscription.Module,
		usage.Module, // Needed for dashboard stats, but maybe not ingestion

		fx.Provide(server.NewEngine),
		fx.Provide(server.NewServer),
		fx.Invoke(func(s *server.Server) {
			s.RegisterAuthRoutes()
			s.RegisterAdminRoutes()
			s.RegisterUIRoutes() // Monolith style: serve the react app
			s.RegisterFallback()
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
