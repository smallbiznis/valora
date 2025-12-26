package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/cloudmetrics"
	meterdomain "github.com/smallbiznis/valora/internal/meter/domain"
	"github.com/smallbiznis/valora/internal/orgcontext"
	subscriptiondomain "github.com/smallbiznis/valora/internal/subscription/domain"
	usagedomain "github.com/smallbiznis/valora/internal/usage/domain"
	"github.com/smallbiznis/valora/pkg/db/option"
	"github.com/smallbiznis/valora/pkg/db/pagination"
	"github.com/smallbiznis/valora/pkg/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ServiceParam struct {
	fx.In

	DB       *gorm.DB
	Log      *zap.Logger
	GenID    *snowflake.Node
	MeterSvc meterdomain.Service
	SubSvc   subscriptiondomain.Service
}

type Service struct {
	db  *gorm.DB
	log *zap.Logger

	genID     *snowflake.Node
	metersvc  meterdomain.Service
	subSvc    subscriptiondomain.Service
	usagerepo repository.Repository[usagedomain.UsageRecord]
}

func NewService(p ServiceParam) usagedomain.Service {
	return &Service{
		db:  p.DB,
		log: p.Log.Named("usage.service"),

		genID:     p.GenID,
		metersvc:  p.MeterSvc,
		subSvc:    p.SubSvc,
		usagerepo: repository.ProvideStore[usagedomain.UsageRecord](p.DB),
	}
}

func (s *Service) Ingest(ctx context.Context, req usagedomain.CreateIngestRequest) (*usagedomain.UsageRecord, error) {
	orgID, err := s.orgIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	customerID, err := s.parseID(req.CustomerID, usagedomain.ErrInvalidCustomer)
	if err != nil {
		return nil, err
	}

	meterCode := strings.TrimSpace(req.MeterCode)
	if meterCode == "" {
		return nil, usagedomain.ErrInvalidMeterCode
	}

	meter, err := s.metersvc.GetByCode(ctx, meterCode)
	if err != nil {
		switch {
		case errors.Is(err, meterdomain.ErrInvalidCode), errors.Is(err, meterdomain.ErrNotFound):
			return nil, usagedomain.ErrInvalidMeterCode
		default:
			return nil, err
		}
	}

	subscription, err := s.subSvc.GetActiveByCustomerID(ctx, subscriptiondomain.GetActiveByCustomerIDRequest{
		CustomerID: req.CustomerID,
	})
	if err != nil {
		switch {
		case errors.Is(err, subscriptiondomain.ErrSubscriptionNotFound):
			return nil, usagedomain.ErrInvalidSubscription
		case errors.Is(err, subscriptiondomain.ErrInvalidCustomer):
			return nil, usagedomain.ErrInvalidCustomer
		default:
			return nil, err
		}
	}

	subscriptionItem, err := s.subSvc.GetSubscriptionItem(ctx, subscriptiondomain.GetSubscriptionItemRequest{
		SubscriptionID: subscription.ID.String(),
		MeterID:        meter.ID,
	})
	if err != nil {
		switch {
		case errors.Is(err, subscriptiondomain.ErrSubscriptionItemNotFound):
			return nil, usagedomain.ErrInvalidSubscriptionItem
		case errors.Is(err, subscriptiondomain.ErrInvalidMeterID):
			return nil, usagedomain.ErrInvalidMeter
		case errors.Is(err, subscriptiondomain.ErrInvalidMeterCode):
			return nil, usagedomain.ErrInvalidMeterCode
		case errors.Is(err, subscriptiondomain.ErrInvalidSubscription):
			return nil, usagedomain.ErrInvalidSubscription
		default:
			return nil, err
		}
	}

	meterID, err := s.parseID(meter.ID, usagedomain.ErrInvalidMeter)
	if err != nil {
		return nil, err
	}

	if subscriptionItem.MeterID != nil {
		return nil, usagedomain.ErrInvalidMeter
	}

	if req.RecordedAt.IsZero() {
		return nil, usagedomain.ErrInvalidRecordedAt
	}

	if math.IsNaN(req.Value) || math.IsInf(req.Value, 0) {
		return nil, usagedomain.ErrInvalidValue
	}

	var idempotencyKey *string
	if req.IdempotencyKey != nil {
		key := strings.TrimSpace(*req.IdempotencyKey)
		if key != "" {
			idempotencyKey = &key
		}
	}

	record := &usagedomain.UsageRecord{
		ID:                 s.genID.Generate(),
		OrgID:              orgID,
		CustomerID:         customerID,
		SubscriptionID:     subscription.ID,
		SubscriptionItemID: subscriptionItem.ID,
		MeterID:            meterID,
		MeterCode:          meterCode,
		Value:              req.Value,
		RecordedAt:         req.RecordedAt,
		IdempotencyKey:     idempotencyKey,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}
	if req.Metadata != nil {
		record.Metadata = datatypes.JSONMap(req.Metadata)
	}

	if err := s.usagerepo.Create(ctx, record); err != nil {
		return nil, err
	}

	cloudmetrics.RecordUsageEvent(orgID.String(), meterCode)
	return record, nil
}

func (s *Service) List(ctx context.Context, req usagedomain.ListUsageRequest) (usagedomain.ListUsageResponse, error) {
	orgID, err := s.orgIDFromContext(ctx)
	if err != nil {
		return usagedomain.ListUsageResponse{}, err
	}

	filter := &usagedomain.UsageRecord{
		OrgID: orgID,
	}

	if req.CustomerID != "" {
		customerID, err := s.parseID(req.CustomerID, usagedomain.ErrInvalidCustomer)
		if err != nil {
			return usagedomain.ListUsageResponse{}, err
		}
		filter.CustomerID = customerID
	}

	if req.SubscriptionID != "" {
		subscriptionID, err := s.parseID(req.SubscriptionID, usagedomain.ErrInvalidSubscription)
		if err != nil {
			return usagedomain.ListUsageResponse{}, err
		}
		filter.SubscriptionID = subscriptionID
	}

	if req.MeterID != "" {
		meterID, err := s.parseID(req.MeterID, usagedomain.ErrInvalidMeter)
		if err != nil {
			return usagedomain.ListUsageResponse{}, err
		}
		filter.MeterID = meterID
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}

	items, err := s.usagerepo.Find(ctx, filter,
		option.ApplyPagination(pagination.Pagination{
			PageToken: req.PageToken,
			PageSize:  int(pageSize),
		}),
		option.WithSortBy(option.QuerySortBy{Allow: map[string]bool{"created_at": true}}),
	)
	if err != nil {
		return usagedomain.ListUsageResponse{}, err
	}

	pageInfo := pagination.BuildCursorPageInfo(items, pageSize, func(record *usagedomain.UsageRecord) string {
		token, err := pagination.EncodeCursor(pagination.Cursor{
			ID:        record.ID.String(),
			CreatedAt: record.CreatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return ""
		}
		return token
	})
	if pageInfo != nil && pageInfo.HasMore && len(items) > int(pageSize) {
		items = items[:pageSize]
	}

	records := make([]usagedomain.UsageRecord, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		records = append(records, *item)
	}

	resp := usagedomain.ListUsageResponse{
		UsageRecords: records,
	}
	if pageInfo != nil {
		resp.PageInfo = *pageInfo
	}

	return resp, nil
}

func (s *Service) parseID(value string, invalidErr error) (snowflake.ID, error) {
	id, err := snowflake.ParseString(strings.TrimSpace(value))
	if err != nil || id == 0 {
		return 0, invalidErr
	}
	return id, nil
}

func (s *Service) orgIDFromContext(ctx context.Context) (snowflake.ID, error) {
	orgID, ok := orgcontext.OrgIDFromContext(ctx)
	if !ok || orgID == 0 {
		return 0, usagedomain.ErrInvalidOrganization
	}
	return snowflake.ID(orgID), nil
}
