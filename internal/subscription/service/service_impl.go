package service

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	billingcycledomain "github.com/smallbiznis/valora/internal/billingcycle/domain"
	pricedomain "github.com/smallbiznis/valora/internal/price/domain"
	subscriptiondomain "github.com/smallbiznis/valora/internal/subscription/domain"
	"github.com/smallbiznis/valora/pkg/db/option"
	"github.com/smallbiznis/valora/pkg/db/pagination"
	"github.com/smallbiznis/valora/pkg/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	log *zap.Logger

	genID            *snowflake.Node
	repo             subscriptiondomain.Repository
	billingCycleRepo repository.Repository[billingcycledomain.BillingCycle]
	subscriptionRepo repository.Repository[subscriptiondomain.Subscription]

	pricesvc pricedomain.Service
}

type ServiceParam struct {
	fx.In

	DB    *gorm.DB
	Log   *zap.Logger
	GenID *snowflake.Node
	Repo  subscriptiondomain.Repository

	Pricesvc pricedomain.Service
}

func NewService(p ServiceParam) subscriptiondomain.Service {
	return &Service{
		db:  p.DB,
		log: p.Log.Named("subscription.service"),

		genID:            p.GenID,
		repo:             p.Repo,
		billingCycleRepo: repository.ProvideStore[billingcycledomain.BillingCycle](p.DB),
		subscriptionRepo: repository.ProvideStore[subscriptiondomain.Subscription](p.DB),

		pricesvc: p.Pricesvc,
	}
}

// GetActiveByCustomerID implements domain.Service.
func (s *Service) GetActiveByCustomerID(ctx context.Context, req subscriptiondomain.GetActiveByCustomerIDRequest) (subscriptiondomain.Subscription, error) {
	orgID, err := s.parseID(req.OrgID, subscriptiondomain.ErrInvalidOrganization)
	if err != nil {
		return subscriptiondomain.Subscription{}, err
	}

	customerID, err := s.parseID(req.CustomerID, subscriptiondomain.ErrInvalidCustomer)
	if err != nil {
		return subscriptiondomain.Subscription{}, err
	}

	statuses := []subscriptiondomain.SubscriptionStatus{
		subscriptiondomain.SubscriptionStatusActive,
		subscriptiondomain.SubscriptionStatusTrialing,
	}

	item, err := s.repo.FindActiveByCustomerID(ctx, s.db, orgID, customerID, statuses)
	if err != nil {
		return subscriptiondomain.Subscription{}, err
	}
	if item == nil {
		return subscriptiondomain.Subscription{}, subscriptiondomain.ErrSubscriptionNotFound
	}

	return *item, nil
}

// GetSubscriptionItem implements domain.Service.
func (s *Service) GetSubscriptionItem(ctx context.Context, req subscriptiondomain.GetSubscriptionItemRequest) (subscriptiondomain.SubscriptionItem, error) {
	orgID, err := s.parseID(req.OrgID, subscriptiondomain.ErrInvalidOrganization)
	if err != nil {
		return subscriptiondomain.SubscriptionItem{}, err
	}

	subscriptionID, err := s.parseID(req.SubscriptionID, subscriptiondomain.ErrInvalidSubscription)
	if err != nil {
		return subscriptiondomain.SubscriptionItem{}, err
	}

	meterID := strings.TrimSpace(req.MeterID)
	if meterID != "" {
		parsedMeterID, err := s.parseID(meterID, subscriptiondomain.ErrInvalidMeterID)
		if err != nil {
			return subscriptiondomain.SubscriptionItem{}, err
		}

		item, err := s.repo.FindSubscriptionItemByMeterID(ctx, s.db, orgID, subscriptionID, parsedMeterID)
		if err != nil {
			return subscriptiondomain.SubscriptionItem{}, err
		}
		if item == nil {
			return subscriptiondomain.SubscriptionItem{}, subscriptiondomain.ErrSubscriptionItemNotFound
		}

		return *item, nil
	}

	meterCode := strings.TrimSpace(req.MeterCode)
	if meterCode == "" {
		return subscriptiondomain.SubscriptionItem{}, subscriptiondomain.ErrInvalidMeterCode
	}

	item, err := s.repo.FindSubscriptionItemByMeterCode(ctx, s.db, orgID, subscriptionID, meterCode)
	if err != nil {
		return subscriptiondomain.SubscriptionItem{}, err
	}
	if item == nil {
		return subscriptiondomain.SubscriptionItem{}, subscriptiondomain.ErrSubscriptionItemNotFound
	}

	return *item, nil
}

func (s *Service) List(ctx context.Context, req subscriptiondomain.ListSubscriptionRequest) (subscriptiondomain.ListSubscriptionResponse, error) {
	orgID, err := s.parseID(req.OrgID, subscriptiondomain.ErrInvalidOrganization)
	if err != nil {
		return subscriptiondomain.ListSubscriptionResponse{}, err
	}

	filter := &subscriptiondomain.Subscription{
		OrgID: orgID,
	}

	statusFilter, err := parseStatusFilter(req.Status)
	if err != nil {
		return subscriptiondomain.ListSubscriptionResponse{}, err
	}
	if statusFilter != nil {
		filter.Status = *statusFilter
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}

	items, err := s.subscriptionRepo.Find(ctx, filter,
		option.ApplyPagination(pagination.Pagination{
			PageToken: req.PageToken,
			PageSize:  int(pageSize),
		}),
		option.WithSortBy(option.WithQuerySortBy("created_at", "desc", map[string]bool{"created_at": true})),
	)
	if err != nil {
		return subscriptiondomain.ListSubscriptionResponse{}, err
	}

	pageInfo := pagination.BuildCursorPageInfo(items, pageSize, func(item *subscriptiondomain.Subscription) string {
		token, err := pagination.EncodeCursor(pagination.Cursor{
			ID:        item.ID.String(),
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return ""
		}
		return token
	})
	if pageInfo != nil && pageInfo.HasMore && len(items) > int(pageSize) {
		items = items[:pageSize]
	}

	subscriptions := make([]subscriptiondomain.Subscription, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		subscriptions = append(subscriptions, *item)
	}

	resp := subscriptiondomain.ListSubscriptionResponse{
		Subscriptions: subscriptions,
	}
	if pageInfo != nil {
		resp.PageInfo = *pageInfo
	}

	return resp, nil
}

func (s *Service) Create(ctx context.Context, req subscriptiondomain.CreateSubscriptionRequest) (subscriptiondomain.CreateSubscriptionResponse, error) {
	orgID, err := s.parseID(req.OrganizationID, subscriptiondomain.ErrInvalidOrganization)
	if err != nil {
		return subscriptiondomain.CreateSubscriptionResponse{}, err
	}

	customerID, err := s.parseID(req.CustomerID, subscriptiondomain.ErrInvalidCustomer)
	if err != nil {
		return subscriptiondomain.CreateSubscriptionResponse{}, err
	}

	if len(req.Items) == 0 {
		return subscriptiondomain.CreateSubscriptionResponse{}, subscriptiondomain.ErrInvalidItems
	}

	now := time.Now().UTC()
	subscription := subscriptiondomain.Subscription{
		ID:             s.genID.Generate(),
		OrgID:          orgID,
		CustomerID:     customerID,
		Status:         subscriptiondomain.SubscriptionStatusActive,
		CollectionMode: subscriptiondomain.SendInvoice,
		StartAt:        now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if req.Metadata != nil {
		subscription.Metadata = datatypes.JSONMap(req.Metadata)
	}

	priceCache := make(map[string]*pricedomain.Response, len(req.Items))
	flatCount := 0
	subscriptionItems := make([]subscriptiondomain.SubscriptionItem, 0, len(req.Items))

	for _, item := range req.Items {
		priceID := strings.TrimSpace(item.PriceID)
		price, ok := priceCache[priceID]
		if !ok {
			loadedPrice, err := s.pricesvc.Get(ctx, orgID.String(), priceID)
			if err != nil {
				return subscriptiondomain.CreateSubscriptionResponse{}, err
			}
			if !loadedPrice.Active {
				return subscriptiondomain.CreateSubscriptionResponse{}, subscriptiondomain.ErrInvalidPrice
			}
			price = loadedPrice
			priceCache[priceID] = loadedPrice
		}

		if price.PricingModel == pricedomain.Flat {
			flatCount++
			if flatCount > 1 {
				return subscriptiondomain.CreateSubscriptionResponse{}, subscriptiondomain.ErrMultipleFlatPrices
			}
		}

		quantity := item.Quantity
		if quantity == 0 {
			quantity = 1
		} else if quantity < 0 {
			return subscriptiondomain.CreateSubscriptionResponse{}, subscriptiondomain.ErrInvalidQuantity
		}

		parsedPriceID, err := s.parseID(price.ID, subscriptiondomain.ErrInvalidPrice)
		if err != nil {
			return subscriptiondomain.CreateSubscriptionResponse{}, err
		}

		priceCode := price.Code
		var priceCodePtr *string
		if priceCode != "" {
			priceCodePtr = &priceCode
		}

		subscriptionItems = append(subscriptionItems, subscriptiondomain.SubscriptionItem{
			ID:               s.genID.Generate(),
			OrgID:            orgID,
			SubscriptionID:   subscription.ID,
			PriceID:          parsedPriceID,
			PriceCode:        priceCodePtr,
			Quantity:         quantity,
			BillingMode:      string(price.BillingMode),
			BillingThreshold: price.BillingThreshold,
			CreatedAt:        now,
			UpdatedAt:        now,
		})
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.Insert(ctx, tx, &subscription); err != nil {
			return err
		}
		if err := s.repo.InsertItems(ctx, tx, subscriptionItems); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return subscriptiondomain.CreateSubscriptionResponse{}, err
	}

	return s.toCreateResponse(&subscription, subscriptionItems), nil
}

func (s *Service) GetByID(ctx context.Context, id string) (subscriptiondomain.Subscription, error) {
	subscriptionID, err := snowflake.ParseString(strings.TrimSpace(id))
	if err != nil {
		return subscriptiondomain.Subscription{}, err
	}

	item, err := s.repo.FindByID(ctx, s.db, subscriptionID)
	if err != nil {
		return subscriptiondomain.Subscription{}, err
	}
	if item == nil {
		return subscriptiondomain.Subscription{}, gorm.ErrRecordNotFound
	}

	return *item, nil
}

func (s *Service) parseID(value string, invalidErr error) (snowflake.ID, error) {
	id, err := snowflake.ParseString(strings.TrimSpace(value))
	if err != nil || id == 0 {
		return 0, invalidErr
	}
	return id, nil
}

func parseStatusFilter(value string) (*subscriptiondomain.SubscriptionStatus, error) {
	status := strings.TrimSpace(value)
	if status == "" {
		return nil, nil
	}

	status = strings.ToUpper(status)
	switch subscriptiondomain.SubscriptionStatus(status) {
	case subscriptiondomain.SubscriptionStatusDraft,
		subscriptiondomain.SubscriptionStatusActive,
		subscriptiondomain.SubscriptionStatusTrialing,
		subscriptiondomain.SubscriptionStatusPastDue,
		subscriptiondomain.SubscriptionStatusCanceled,
		subscriptiondomain.SubscriptionStatusEnded:
		parsed := subscriptiondomain.SubscriptionStatus(status)
		return &parsed, nil
	default:
		return nil, subscriptiondomain.ErrInvalidStatus
	}
}

func parseStatus(value string) (subscriptiondomain.SubscriptionStatus, error) {
	status := strings.TrimSpace(value)
	if status == "" {
		return subscriptiondomain.SubscriptionStatusActive, nil
	}

	switch subscriptiondomain.SubscriptionStatus(status) {
	case subscriptiondomain.SubscriptionStatusDraft,
		subscriptiondomain.SubscriptionStatusActive,
		subscriptiondomain.SubscriptionStatusTrialing,
		subscriptiondomain.SubscriptionStatusPastDue,
		subscriptiondomain.SubscriptionStatusCanceled,
		subscriptiondomain.SubscriptionStatusEnded:
		return subscriptiondomain.SubscriptionStatus(status), nil
	default:
		return "", subscriptiondomain.ErrInvalidStatus
	}
}

func parseCollectionMode(value string) (subscriptiondomain.SubscriptionCollectionMode, error) {
	mode := strings.TrimSpace(value)
	if mode == "" {
		return "", subscriptiondomain.ErrInvalidCollectionMode
	}

	switch subscriptiondomain.SubscriptionCollectionMode(mode) {
	case subscriptiondomain.SendInvoice,
		subscriptiondomain.ChargeAutomatically:
		return subscriptiondomain.SubscriptionCollectionMode(mode), nil
	default:
		return "", subscriptiondomain.ErrInvalidCollectionMode
	}
}

func (s *Service) toCreateResponse(subscription *subscriptiondomain.Subscription, items []subscriptiondomain.SubscriptionItem) subscriptiondomain.CreateSubscriptionResponse {
	respItems := make([]subscriptiondomain.CreateSubscriptionItemResponse, 0, len(items))
	for _, item := range items {
		var meterID *string
		if item.MeterID != nil {
			value := item.MeterID.String()
			meterID = &value
		}
		var meterCode *string
		if item.MeterCode != nil && strings.TrimSpace(*item.MeterCode) != "" {
			value := strings.TrimSpace(*item.MeterCode)
			meterCode = &value
		}

		respItems = append(respItems, subscriptiondomain.CreateSubscriptionItemResponse{
			ID:                item.ID.String(),
			PriceID:           item.PriceID.String(),
			PriceCode:         item.PriceCode,
			MeterID:           meterID,
			MeterCode:         meterCode,
			Quantity:          item.Quantity,
			BillingMode:       item.BillingMode,
			UsageBehavior:     item.UsageBehavior,
			BillingThreshold:  item.BillingThreshold,
			ProrationBehavior: item.ProrationBehavior,
		})
	}

	var metadata map[string]any
	if subscription.Metadata != nil {
		metadata = map[string]any(subscription.Metadata)
	}

	return subscriptiondomain.CreateSubscriptionResponse{
		ID:             subscription.ID.String(),
		OrganizationID: subscription.OrgID.String(),
		CustomerID:     subscription.CustomerID.String(),
		Status:         subscription.Status,
		CollectionMode: subscription.CollectionMode,
		StartAt:        subscription.StartAt,
		Items:          respItems,
		Metadata:       metadata,
	}
}
