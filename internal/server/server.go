package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	authdomain "github.com/smallbiznis/valora/internal/auth/domain"
	"github.com/smallbiznis/valora/internal/auth/session"
	"github.com/smallbiznis/valora/internal/billingprovisioning"
	"github.com/smallbiznis/valora/internal/cloudmetrics"
	"github.com/smallbiznis/valora/internal/config"
	"github.com/smallbiznis/valora/internal/customer"
	customerdomain "github.com/smallbiznis/valora/internal/customer/domain"
	"github.com/smallbiznis/valora/internal/invoice"
	invoicedomain "github.com/smallbiznis/valora/internal/invoice/domain"
	"github.com/smallbiznis/valora/internal/meter"
	meterdomain "github.com/smallbiznis/valora/internal/meter/domain"
	"github.com/smallbiznis/valora/internal/organization"
	organizationdomain "github.com/smallbiznis/valora/internal/organization/domain"
	"github.com/smallbiznis/valora/internal/price"
	pricedomain "github.com/smallbiznis/valora/internal/price/domain"
	"github.com/smallbiznis/valora/internal/priceamount"
	priceamountdomain "github.com/smallbiznis/valora/internal/priceamount/domain"
	"github.com/smallbiznis/valora/internal/pricetier"
	pricetierdomain "github.com/smallbiznis/valora/internal/pricetier/domain"
	"github.com/smallbiznis/valora/internal/product"
	productdomain "github.com/smallbiznis/valora/internal/product/domain"
	"github.com/smallbiznis/valora/internal/reference"
	referencedomain "github.com/smallbiznis/valora/internal/reference/domain"
	signupdomain "github.com/smallbiznis/valora/internal/signup/domain"
	"github.com/smallbiznis/valora/internal/subscription"
	subscriptiondomain "github.com/smallbiznis/valora/internal/subscription/domain"
	"github.com/smallbiznis/valora/internal/usage"
	usagedomain "github.com/smallbiznis/valora/internal/usage/domain"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module("http.server",
	config.Module,
	cloudmetrics.Module,
	fx.Provide(registerGin),
	billingprovisioning.Module,
	customer.Module,
	invoice.Module,
	meter.Module,
	organization.Module,
	price.Module,
	priceamount.Module,
	pricetier.Module,
	product.Module,
	reference.Module,
	subscription.Module,
	usage.Module,
	fx.Invoke(NewServer),
	fx.Invoke(run),
)

func registerGin() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(ErrorHandlingMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}

func run(lc fx.Lifecycle, r *gin.Engine) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := r.Run(":8080"); err != nil {
						panic(err)
					}
				}()
				return nil
			},
		},
	)
}

type Server struct {
	engine          *gin.Engine
	cfg             config.Config
	db              *gorm.DB
	authsvc         authdomain.Service
	sessions        *session.Manager
	invoiceSvc      invoicedomain.Service
	meterSvc        meterdomain.Service
	organizationSvc organizationdomain.Service
	customerSvc     customerdomain.Service
	priceSvc        pricedomain.Service
	priceAmountSvc  priceamountdomain.Service
	priceTierSvc    pricetierdomain.Service
	productSvc      productdomain.Service
	refrepo         referencedomain.Repository
	signupsvc       signupdomain.Service
	subscriptionSvc subscriptiondomain.Service
	usagesvc        usagedomain.Service
}

type ServerParams struct {
	fx.In

	Gin             *gin.Engine
	Cfg             config.Config
	DB              *gorm.DB
	InvoiceSvc      invoicedomain.Service
	MeterSvc        meterdomain.Service
	OrganizationSvc organizationdomain.Service
	CustomerSvc     customerdomain.Service
	PriceSvc        pricedomain.Service
	PriceAmountSvc  priceamountdomain.Service
	PriceTierSvc    pricetierdomain.Service
	ProductSvc      productdomain.Service
	Refrepo         referencedomain.Repository
	SubscriptionSvc subscriptiondomain.Service
	Usagesvc        usagedomain.Service
}

func NewServer(p ServerParams) *Server {
	svc := &Server{
		engine:          p.Gin,
		cfg:             p.Cfg,
		db:              p.DB,
		invoiceSvc:      p.InvoiceSvc,
		meterSvc:        p.MeterSvc,
		organizationSvc: p.OrganizationSvc,
		customerSvc:     p.CustomerSvc,
		priceSvc:        p.PriceSvc,
		priceAmountSvc:  p.PriceAmountSvc,
		priceTierSvc:    p.PriceTierSvc,
		productSvc:      p.ProductSvc,
		refrepo:         p.Refrepo,
		subscriptionSvc: p.SubscriptionSvc,
		usagesvc:        p.Usagesvc,
	}

	svc.registerAPIRoutes()

	return svc
}

func (s *Server) registerAPIRoutes() {
	api := s.engine.Group("/api")

	secured := api.Group("")

	// --- global middlewares ---
	secured.Use(RequestID())
	secured.Use(s.AuthRequired())
	secured.Use(s.OrgContext())

	secured.GET("/countries", s.ListCountries)
	secured.GET("/timezones", s.ListTimezones)
	secured.GET("/currencies", s.ListCurrencies)

	// -------- Meters --------
	secured.GET("/meters", s.ListMeters)
	secured.POST("/meters", s.CreateMeter)
	secured.GET("/meters/:id", s.GetMeterByID)
	secured.PATCH("/meters/:id", s.UpdateMeter)
	secured.DELETE("/meters/:id", s.DeleteMeter)

	// -------- Product --------
	secured.GET("/products", s.ListProducts)
	secured.POST("/products", s.CreateProduct)
	secured.GET("/products/:id", s.GetProductByID)

	// -------- Pricing --------
	secured.GET("/pricings", s.ListPricings)
	secured.POST("/pricings", s.CreatePricing)
	secured.GET("/pricings/:id", s.GetPricingByID)

	// -------- Prices --------
	secured.GET("/prices", s.ListPrices)
	secured.POST("/prices", s.CreatePrice)
	secured.GET("/prices/:id", s.GetPriceByID)

	// -------- Price Amounts --------
	secured.GET("/price_amounts", s.ListPriceAmounts)
	secured.POST("/price_amounts", s.CreatePriceAmount)
	secured.GET("/price_amounts/:id", s.GetPriceAmountByID)

	// -------- Tiers ---------
	secured.GET("/price_tiers", s.ListPriceTiers)
	secured.POST("/price_tiers", s.CreatePriceTier)
	secured.GET("/price_tiers/:id", s.GetPriceTierByID)

	// -------- Subscriptions --------
	secured.GET("/subscriptions", s.ListSubscriptions)
	secured.POST("/subscriptions", s.CreateSubscription)
	secured.GET("/subscriptions/:id", s.GetSubscriptionByID)
	secured.POST("/subscriptions/:id/cancel", s.CancelSubscription)

	// -------- Usage / Rating --------
	secured.POST("/usage", s.IngestUsage)
	secured.POST("/rating/run", s.RunRatingJob)

	// -------- Invoices --------
	secured.GET("/invoices", s.ListInvoices)
	secured.GET("/invoices/:id", s.GetInvoiceByID)

	// -------- Customers --------
	secured.GET("/customers", s.ListCustomers)
	secured.POST("/customers", s.CreateCustomer)
	secured.GET("/customers/:id", s.GetCustomerByID)

	if s.cfg.Environment != "production" {
		secured.POST("/test/cleanup", s.TestCleanup)
	}
}
