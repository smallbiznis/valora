package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gosimple/slug"
	"github.com/smallbiznis/valora/internal/organization/domain"
	"github.com/smallbiznis/valora/internal/organization/event"
	referencedomain "github.com/smallbiznis/valora/internal/reference/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	db        *gorm.DB
	repo      domain.Repository
	ref       referencedomain.Repository
	genID     *snowflake.Node
	publisher event.EventPublisher
}

func NewService(db *gorm.DB, repo domain.Repository, ref referencedomain.Repository, genID *snowflake.Node, publisher event.EventPublisher) domain.Service {
	return &service{
		db:        db,
		repo:      repo,
		ref:       ref,
		genID:     genID,
		publisher: publisher,
	}
}

func (s *service) Create(ctx context.Context, userID snowflake.ID, req domain.CreateOrganizationRequest) (*domain.OrganizationResponse, error) {
	if userID == 0 {
		return nil, domain.ErrInvalidUser
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidName
	}

	countryCode := strings.TrimSpace(req.CountryCode)
	if countryCode == "" {
		return nil, domain.ErrInvalidCountry
	}

	timezoneName := strings.TrimSpace(req.TimezoneName)
	if timezoneName == "" {
		return nil, domain.ErrInvalidTimezone
	}

	countryOK, err := s.countryExists(ctx, countryCode)
	if err != nil {
		return nil, err
	}
	if !countryOK {
		return nil, domain.ErrInvalidCountry
	}

	timezoneOK, err := s.timezoneAllowed(ctx, countryCode, timezoneName)
	if err != nil {
		return nil, err
	}
	if !timezoneOK {
		return nil, domain.ErrInvalidTimezone
	}

	now := time.Now().UTC()
	orgID := s.genID.Generate()
	org := domain.Organization{
		ID:              orgID,
		Name:            name,
		Slug: slug.Make(name),
		CountryCode:     countryCode,
		TimezoneName:    timezoneName,
		CreatedAt:       now,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := s.repo.WithTx(tx)
		if err := repo.CreateOrganization(ctx, org); err != nil {
			return err
		}

		member := domain.OrganizationMember{
			ID:        s.genID.Generate(),
			OrgID:     orgID,
			UserID:    userID,
			Role:      domain.RoleOwner,
			CreatedAt: now,
		}

		if err := repo.AddMember(ctx, member); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.emitOrganizationCreated(ctx, org, userID)

	return &domain.OrganizationResponse{
		ID:           orgID.String(),
		Name:         name,
		Slug:         org.Slug,
		CountryCode:  countryCode,
		TimezoneName: timezoneName,
	}, nil
}

func (s *service) ListOrganizationsByUser(ctx context.Context, userID snowflake.ID) ([]domain.OrganizationListResponseItem, error) {
	if userID == 0 {
		return nil, domain.ErrInvalidUser
	}

	items, err := s.repo.ListOrganizationsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]domain.OrganizationListResponseItem, 0, len(items))
	for _, item := range items {
		resp = append(resp, domain.OrganizationListResponseItem{
			ID:        item.ID.String(),
			Name:      item.Name,
			Role:      item.Role,
			CreatedAt: item.CreatedAt,
		})
	}

	return resp, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*domain.OrganizationResponse, error) {
	raw := strings.TrimSpace(id)
	if raw == "" {
		return nil, domain.ErrInvalidOrganization
	}
	orgID, err := snowflake.ParseString(raw)
	if err != nil {
		return nil, domain.ErrInvalidOrganization
	}

	var org domain.Organization
	if err := s.db.WithContext(ctx).First(&org, "id = ?", orgID).Error; err != nil {
		return nil, err
	}

	return &domain.OrganizationResponse{
		ID:           org.ID.String(),
		Name:         org.Name,
		Slug:         org.Slug,
		CountryCode:  org.CountryCode,
		TimezoneName: org.TimezoneName,
	}, nil
}

func (s *service) countryExists(ctx context.Context, code string) (bool, error) {
	countries, err := s.ref.ListCountries(ctx)
	if err != nil {
		return false, err
	}
	for _, country := range countries {
		if country.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) timezoneAllowed(ctx context.Context, countryCode, timezoneName string) (bool, error) {
	timezones, err := s.ref.ListTimezonesByCountry(ctx, countryCode)
	if err != nil {
		return false, err
	}
	for _, tz := range timezones {
		if tz.Name == timezoneName {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) currencyExists(ctx context.Context, code string) (bool, error) {
	currencies, err := s.ref.ListCurrencies(ctx)
	if err != nil {
		return false, err
	}
	for _, currency := range currencies {
		if currency.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) emitOrganizationCreated(ctx context.Context, org domain.Organization, ownerUserID snowflake.ID) {
	if s.publisher == nil {
		return
	}

	payload := map[string]string{
		"organization_id":  org.ID.String(),
		"owner_user_id":    ownerUserID.String(),
		"country_code":     org.CountryCode,
		"timezone_name":    org.TimezoneName,
		"created_at":       org.CreatedAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		zap.L().Warn("failed to marshal organization.created payload", zap.Error(err))
		return
	}

	if err := s.publisher.Publish(ctx, event.OrganizationCreatedTopic, data); err != nil {
		zap.L().Warn("failed to publish organization.created", zap.Error(err))
	}
}
