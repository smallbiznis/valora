package service

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	meterdomain "github.com/smallbiznis/valora/internal/meter/domain"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Params struct {
	fx.In

	DB    *gorm.DB
	Log   *zap.Logger
	GenID *snowflake.Node
	Repo  meterdomain.Repository
}

type Service struct {
	db    *gorm.DB
	log   *zap.Logger
	repo  meterdomain.Repository
	genID *snowflake.Node
}

func New(p Params) meterdomain.Service {
	return &Service{
		db:    p.DB,
		log:   p.Log.Named("meter.service"),
		repo:  p.Repo,
		genID: p.GenID,
	}
}

func (s *Service) Create(ctx context.Context, req meterdomain.CreateRequest) (*meterdomain.Response, error) {
	orgID, err := s.parseOrganizationID(req.OrganizationID)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		return nil, meterdomain.ErrInvalidCode
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, meterdomain.ErrInvalidName
	}

	aggregation := strings.TrimSpace(req.Aggregation)
	if aggregation == "" {
		return nil, meterdomain.ErrInvalidAggregation
	}

	unit := strings.TrimSpace(req.Unit)
	if unit == "" {
		return nil, meterdomain.ErrInvalidUnit
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	now := time.Now().UTC()
	m := &meterdomain.Meter{
		ID:          s.genID.Generate(),
		OrgID:       orgID,
		Code:        code,
		Name:        name,
		Aggregation: aggregation,
		Unit:        unit,
		Active:      active,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Insert(ctx, s.db, m); err != nil {
		return nil, err
	}

	return s.toResponse(m), nil
}

func (s *Service) List(ctx context.Context, organizationID string) ([]meterdomain.Response, error) {
	orgID, err := s.parseOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.List(ctx, s.db, orgID)
	if err != nil {
		return nil, err
	}

	resp := make([]meterdomain.Response, 0, len(items))
	for i := range items {
		resp = append(resp, *s.toResponse(&items[i]))
	}

	return resp, nil
}

func (s *Service) Update(ctx context.Context, req meterdomain.UpdateRequest) (*meterdomain.Response, error) {
	orgID, err := s.parseOrganizationID(req.OrganizationID)
	if err != nil {
		return nil, err
	}

	meterID, err := meterdomain.ParseID(strings.TrimSpace(req.ID))
	if err != nil {
		return nil, meterdomain.ErrInvalidID
	}

	item, err := s.repo.FindByID(ctx, s.db, orgID, meterID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, meterdomain.ErrNotFound
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, meterdomain.ErrInvalidName
		}
		item.Name = name
	}

	if req.Aggregation != nil {
		aggregation := strings.TrimSpace(*req.Aggregation)
		if aggregation == "" {
			return nil, meterdomain.ErrInvalidAggregation
		}
		item.Aggregation = aggregation
	}

	if req.Unit != nil {
		unit := strings.TrimSpace(*req.Unit)
		if unit == "" {
			return nil, meterdomain.ErrInvalidUnit
		}
		item.Unit = unit
	}

	if req.Active != nil {
		item.Active = *req.Active
	}

	item.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, s.db, item); err != nil {
		return nil, err
	}

	return s.toResponse(item), nil
}

func (s *Service) Delete(ctx context.Context, organizationID string, id string) error {
	orgID, err := s.parseOrganizationID(organizationID)
	if err != nil {
		return err
	}

	meterID, err := meterdomain.ParseID(strings.TrimSpace(id))
	if err != nil {
		return meterdomain.ErrInvalidID
	}

	item, err := s.repo.FindByID(ctx, s.db, orgID, meterID)
	if err != nil {
		return err
	}
	if item == nil {
		return meterdomain.ErrNotFound
	}

	return s.repo.Delete(ctx, s.db, orgID, meterID)
}

func (s *Service) GetByCode(ctx context.Context, organizationID string, code string) (*meterdomain.Response, error) {
	orgID, err := s.parseOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	item, err := s.repo.FindByCode(ctx, s.db, orgID, code)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, meterdomain.ErrNotFound
	}

	return s.toResponse(item), nil
}

func (s *Service) GetByID(ctx context.Context, organizationID string, id string) (*meterdomain.Response, error) {
	orgID, err := s.parseOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	meterID, err := meterdomain.ParseID(id)
	if err != nil {
		return nil, meterdomain.ErrInvalidID
	}

	item, err := s.repo.FindByID(ctx, s.db, orgID, meterID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, meterdomain.ErrNotFound
	}

	return s.toResponse(item), nil
}

func (s *Service) parseOrganizationID(value string) (snowflake.ID, error) {
	orgID, err := meterdomain.ParseID(strings.TrimSpace(value))
	if err != nil || orgID == 0 {
		return 0, meterdomain.ErrInvalidOrganization
	}
	return orgID, nil
}

func (s *Service) toResponse(m *meterdomain.Meter) *meterdomain.Response {
	return &meterdomain.Response{
		ID:             m.ID.String(),
		OrganizationID: m.OrgID.String(),
		Code:           m.Code,
		Name:           m.Name,
		Aggregation:    m.Aggregation,
		Unit:           m.Unit,
		Active:         m.Active,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
