package service

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	pricedomain "github.com/smallbiznis/valora/internal/price/domain"
	productdomain "github.com/smallbiznis/valora/internal/product/domain"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Params struct {
	fx.In

	DB          *gorm.DB
	Log         *zap.Logger
	GenID       *snowflake.Node
	Repo        pricedomain.Repository
	ProductRepo productdomain.Repository
}

type Service struct {
	db          *gorm.DB
	log         *zap.Logger
	genID       *snowflake.Node
	repo        pricedomain.Repository
	productRepo productdomain.Repository
}

func New(p Params) pricedomain.Service {
	return &Service{
		db:          p.DB,
		log:         p.Log.Named("price.service"),
		genID:       p.GenID,
		repo:        p.Repo,
		productRepo: p.ProductRepo,
	}
}

func (s *Service) Create(ctx context.Context, req pricedomain.CreateRequest) (*pricedomain.Response, error) {
	orgID, err := s.parseOrganizationID(req.OrganizationID)
	if err != nil {
		return nil, err
	}

	productID, err := parseID(req.ProductID)
	if err != nil {
		return nil, pricedomain.ErrInvalidProduct
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, pricedomain.ErrInvalidCode
	}

	pricingModel, err := parsePricingModel(req.PricingModel)
	if err != nil {
		return nil, err
	}

	billingMode, err := parseBillingMode(req.BillingMode)
	if err != nil {
		return nil, err
	}

	billingInterval, err := parseBillingInterval(req.BillingInterval)
	if err != nil {
		return nil, err
	}

	taxBehavior, err := parseTaxBehavior(req.TaxBehavior)
	if err != nil {
		return nil, err
	}

	if req.BillingIntervalCount <= 0 {
		return nil, pricedomain.ErrInvalidBillingIntervalCount
	}

	productExists, err := s.productExists(ctx, orgID, productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, pricedomain.ErrInvalidProduct
	}

	var aggregateUsagePtr *pricedomain.AggregateUsage
	if req.AggregateUsage != nil {
		aggregateUsage, err := parseAggregateUsage(*req.AggregateUsage)
		if err != nil {
			return nil, err
		}
		aggregateUsagePtr = &aggregateUsage
	}

	var billingUnitPtr *pricedomain.BillingUnit
	if req.BillingUnit != nil {
		billingUnit, err := parseBillingUnit(*req.BillingUnit)
		if err != nil {
			return nil, err
		}
		billingUnitPtr = &billingUnit
	}

	if err := validatePricingModelConfig(pricingModel, billingMode, aggregateUsagePtr, billingUnitPtr, req.BillingThreshold); err != nil {
		return nil, err
	}

	taxCode := strings.TrimSpace(ptrToString(req.TaxCode))
	var taxCodePtr *string
	if taxCode != "" {
		taxCodePtr = &taxCode
	}

	version := int32(1)
	if req.Version != nil {
		if *req.Version <= 0 {
			return nil, pricedomain.ErrInvalidVersion
		}
		version = *req.Version
	}

	isDefault := false
	if req.IsDefault != nil {
		isDefault = *req.IsDefault
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	now := time.Now().UTC()
	entity := &pricedomain.Price{
		ID:                   s.genID.Generate(),
		OrgID:                orgID,
		ProductID:            productID,
		Code:                 code,
		LookupKey:            &req.LookupKey,
		Name:                 req.Name,
		Description:          req.Description,
		PricingModel:         pricingModel,
		BillingMode:          billingMode,
		BillingInterval:      billingInterval,
		BillingIntervalCount: req.BillingIntervalCount,
		AggregateUsage:       aggregateUsagePtr,
		BillingUnit:          billingUnitPtr,
		BillingThreshold:     req.BillingThreshold,
		TaxBehavior:          taxBehavior,
		TaxCode:              taxCodePtr,
		Version:              version,
		IsDefault:            isDefault,
		Active:               active,
		RetiredAt:            req.RetiredAt,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if req.Metadata != nil {
		entity.Metadata = datatypes.JSONMap(req.Metadata)
	}

	if err := s.repo.Insert(ctx, s.db, entity); err != nil {
		return nil, err
	}

	return s.toResponse(entity), nil
}

func (s *Service) List(ctx context.Context, organizationID string) ([]pricedomain.Response, error) {
	orgID, err := s.parseOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.List(ctx, s.db, orgID)
	if err != nil {
		return nil, err
	}

	resp := make([]pricedomain.Response, 0, len(items))
	for i := range items {
		resp = append(resp, *s.toResponse(&items[i]))
	}

	return resp, nil
}

func (s *Service) Get(ctx context.Context, organizationID string, id string) (*pricedomain.Response, error) {
	orgID, err := s.parseOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	priceID, err := parseID(id)
	if err != nil {
		return nil, pricedomain.ErrInvalidID
	}

	entity, err := s.repo.FindByID(ctx, s.db, orgID, priceID)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, pricedomain.ErrNotFound
	}

	return s.toResponse(entity), nil
}

func (s *Service) parseOrganizationID(value string) (snowflake.ID, error) {
	orgID, err := snowflake.ParseString(strings.TrimSpace(value))
	if err != nil || orgID == 0 {
		return 0, pricedomain.ErrInvalidOrganization
	}
	return orgID, nil
}

func (s *Service) productExists(ctx context.Context, orgID, productID snowflake.ID) (bool, error) {
	item, err := s.productRepo.FindByID(ctx, s.db, orgID.Int64(), productID.Int64())
	if err != nil {
		return false, err
	}
	return item != nil, nil
}

func (s *Service) toResponse(p *pricedomain.Price) *pricedomain.Response {
	return &pricedomain.Response{
		ID:                   p.ID.String(),
		OrganizationID:       p.OrgID.String(),
		ProductID:            p.ProductID.String(),
		Code:                 p.Code,
		LookupKey:            p.LookupKey,
		Name:                 p.Name,
		Description:          p.Description,
		PricingModel:         p.PricingModel,
		BillingMode:          p.BillingMode,
		BillingInterval:      p.BillingInterval,
		BillingIntervalCount: p.BillingIntervalCount,
		AggregateUsage:       p.AggregateUsage,
		BillingUnit:          p.BillingUnit,
		BillingThreshold:     p.BillingThreshold,
		TaxBehavior:          p.TaxBehavior,
		TaxCode:              p.TaxCode,
		Version:              p.Version,
		IsDefault:            p.IsDefault,
		Active:               p.Active,
		RetiredAt:            p.RetiredAt,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}
}

func parseID(value string) (snowflake.ID, error) {
	return snowflake.ParseString(strings.TrimSpace(value))
}

func ptrToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func parsePricingModel(value pricedomain.PricingModel) (pricedomain.PricingModel, error) {
	switch strings.ToUpper(strings.TrimSpace(string(value))) {
	case string(pricedomain.Flat):
		return pricedomain.Flat, nil
	case string(pricedomain.PerUnit):
		return pricedomain.PerUnit, nil
	case string(pricedomain.TieredVolume):
		return pricedomain.TieredVolume, nil
	case string(pricedomain.TieredGraduated):
		return pricedomain.TieredGraduated, nil
	default:
		return "", pricedomain.ErrInvalidPricingModel
	}
}

func parseBillingMode(value pricedomain.BillingMode) (pricedomain.BillingMode, error) {
	switch strings.ToUpper(strings.TrimSpace(string(value))) {
	case string(pricedomain.Licensed):
		return pricedomain.Licensed, nil
	case string(pricedomain.Metered):
		return pricedomain.Metered, nil
	default:
		return "", pricedomain.ErrInvalidBillingMode
	}
}

func parseBillingInterval(value pricedomain.BillingInterval) (pricedomain.BillingInterval, error) {
	switch strings.ToUpper(strings.TrimSpace(string(value))) {
	case string(pricedomain.Day):
		return pricedomain.Day, nil
	case string(pricedomain.Week):
		return pricedomain.Week, nil
	case string(pricedomain.Month):
		return pricedomain.Month, nil
	case string(pricedomain.Year):
		return pricedomain.Year, nil
	default:
		return "", pricedomain.ErrInvalidBillingInterval
	}
}

func parseAggregateUsage(value pricedomain.AggregateUsage) (pricedomain.AggregateUsage, error) {
	switch strings.ToUpper(strings.TrimSpace(string(value))) {
	case string(pricedomain.SUM):
		return pricedomain.SUM, nil
	case string(pricedomain.MAX):
		return pricedomain.MAX, nil
	case string(pricedomain.LAST):
		return pricedomain.LAST, nil
	default:
		return "", pricedomain.ErrInvalidAggregateUsage
	}
}

func parseBillingUnit(value pricedomain.BillingUnit) (pricedomain.BillingUnit, error) {
	switch strings.ToUpper(strings.TrimSpace(string(value))) {
	case "API_CALL":
		return pricedomain.API_CALL, nil
	case "GB":
		return pricedomain.GB, nil
	case "GIB":
		return pricedomain.GiB, nil
	case "MB":
		return pricedomain.MB, nil
	case "MIB":
		return pricedomain.MiB, nil
	case "SECOND":
		return pricedomain.Second, nil
	case "MINUTE":
		return pricedomain.Minute, nil
	case "HOUR":
		return pricedomain.Hour, nil
	case "SEAT":
		return pricedomain.Seat, nil
	default:
		return "", pricedomain.ErrInvalidBillingUnit
	}
}

func parseTaxBehavior(value pricedomain.TaxBehavior) (pricedomain.TaxBehavior, error) {
	switch strings.ToUpper(strings.TrimSpace(string(value))) {
	case string(pricedomain.Inclusive):
		return pricedomain.Inclusive, nil
	case string(pricedomain.Exclusive):
		return pricedomain.Exclusive, nil
	case string(pricedomain.Inline):
		return pricedomain.Inline, nil
	default:
		return "", pricedomain.ErrInvalidTaxBehavior
	}
}

func validatePricingModelConfig(pricingModel pricedomain.PricingModel, billingMode pricedomain.BillingMode, aggregateUsage *pricedomain.AggregateUsage, billingUnit *pricedomain.BillingUnit, billingThreshold *float64) error {
	switch pricingModel {
	case pricedomain.Flat:
		if billingMode != pricedomain.Licensed {
			return pricedomain.ErrInvalidBillingMode
		}
		if aggregateUsage != nil {
			return pricedomain.ErrInvalidAggregateUsage
		}
		if billingUnit != nil {
			return pricedomain.ErrInvalidBillingUnit
		}
		if billingThreshold != nil {
			return pricedomain.ErrInvalidBillingThreshold
		}
	case pricedomain.PerUnit:
		if billingMode != pricedomain.Metered {
			return pricedomain.ErrInvalidBillingMode
		}
		if billingUnit == nil {
			return pricedomain.ErrInvalidBillingUnit
		}
		if aggregateUsage == nil || *aggregateUsage != pricedomain.SUM {
			return pricedomain.ErrInvalidAggregateUsage
		}
	case pricedomain.TieredVolume, pricedomain.TieredGraduated:
		if billingMode != pricedomain.Metered {
			return pricedomain.ErrInvalidBillingMode
		}
		if billingUnit == nil {
			return pricedomain.ErrInvalidBillingUnit
		}
		if aggregateUsage == nil || *aggregateUsage != pricedomain.SUM {
			return pricedomain.ErrInvalidAggregateUsage
		}
	default:
		return pricedomain.ErrInvalidPricingModel
	}

	return nil
}
