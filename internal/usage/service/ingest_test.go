package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/glebarez/sqlite"
	"github.com/smallbiznis/railzway/internal/meter/domain"
	meterdomain "github.com/smallbiznis/railzway/internal/meter/domain"
	"github.com/smallbiznis/railzway/internal/orgcontext"
	subscriptiondomain "github.com/smallbiznis/railzway/internal/subscription/domain"
	usagedomain "github.com/smallbiznis/railzway/internal/usage/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// -- Mocks --

type subscriptionMock struct {
	mock.Mock
}

func (m *subscriptionMock) GetActiveByCustomerID(ctx context.Context, req subscriptiondomain.GetActiveByCustomerIDRequest) (subscriptiondomain.Subscription, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(subscriptiondomain.Subscription), args.Error(1)
}

func (m *subscriptionMock) ValidateUsageEntitlement(ctx context.Context, subID, meterID snowflake.ID, at time.Time) error {
	args := m.Called(ctx, subID, meterID, at)
	return args.Error(0)
}

func (m *subscriptionMock) List(context.Context, subscriptiondomain.ListSubscriptionRequest) (subscriptiondomain.ListSubscriptionResponse, error) {
	return subscriptiondomain.ListSubscriptionResponse{}, nil
}
func (m *subscriptionMock) Create(context.Context, subscriptiondomain.CreateSubscriptionRequest) (subscriptiondomain.CreateSubscriptionResponse, error) {
	return subscriptiondomain.CreateSubscriptionResponse{}, nil
}
func (m *subscriptionMock) ReplaceItems(context.Context, subscriptiondomain.ReplaceSubscriptionItemsRequest) (subscriptiondomain.CreateSubscriptionResponse, error) {
	return subscriptiondomain.CreateSubscriptionResponse{}, nil
}
func (m *subscriptionMock) GetByID(context.Context, string) (subscriptiondomain.Subscription, error) {
	return subscriptiondomain.Subscription{}, nil
}
func (m *subscriptionMock) GetSubscriptionItem(context.Context, subscriptiondomain.GetSubscriptionItemRequest) (subscriptiondomain.SubscriptionItem, error) {
	return subscriptiondomain.SubscriptionItem{}, nil
}
func (m *subscriptionMock) TransitionSubscription(ctx context.Context, id string, status subscriptiondomain.SubscriptionStatus, reason subscriptiondomain.TransitionReason) error {
	return nil
}
func (m *subscriptionMock) ChangePlan(ctx context.Context, req subscriptiondomain.ChangePlanRequest) error {
	return nil
}

type meterMock struct {
	mock.Mock
}

func (m *meterMock) GetByCode(ctx context.Context, code string) (*domain.Response, error) {
	args := m.Called(ctx, code)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.(*domain.Response), args.Error(1)
}

func (m *meterMock) Create(ctx context.Context, req domain.CreateRequest) (*domain.Response, error) {
	return nil, nil
}
func (m *meterMock) List(ctx context.Context, req domain.ListRequest) ([]domain.Response, error) {
	return nil, nil
}
func (m *meterMock) GetByID(ctx context.Context, id string) (*domain.Response, error) {
	return nil, nil
}
func (m *meterMock) Update(ctx context.Context, req domain.UpdateRequest) (*domain.Response, error) {
	return nil, nil
}
func (m *meterMock) Delete(ctx context.Context, id string) error { return nil }

// -- Tests --

func TestIngest_EntitlementGating(t *testing.T) {
	// Setup In-Memory DB (needed for NewService and Ingest insert)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// Migrate usage_events table
	if err := db.AutoMigrate(&usagedomain.UsageEvent{}); err != nil {
		t.Fatal(err)
	}

	node, _ := snowflake.NewNode(1)
	genID := node
	logger := zap.NewNop()

	orgID := genID.Generate()
	customerID := genID.Generate()
	meterID := genID.Generate()
	subID := genID.Generate()

	tests := []struct {
		name         string
		req          usagedomain.CreateIngestRequest
		setupMocks   func(*subscriptionMock, *meterMock)
		expectedErr  error
		expectIngest bool // If we expect ingestion to succeed
	}{
		{
			name: "Success: Valid Entitlement",
			req: usagedomain.CreateIngestRequest{
				CustomerID:     customerID.String(),
				MeterCode:      "m1",
				Value:          10,
				RecordedAt:     time.Now(),
				IdempotencyKey: "k1",
			},
			setupMocks: func(s *subscriptionMock, m *meterMock) {
				// 1. Resolve Meter
				m.On("GetByCode", mock.Anything, "m1").Return(&meterdomain.Response{
					ID:   meterID.String(),
					Code: "m1",
				}, nil)
				// 2. Resolve Subscription
				s.On("GetActiveByCustomerID", mock.Anything, mock.MatchedBy(func(req subscriptiondomain.GetActiveByCustomerIDRequest) bool {
					return req.CustomerID == customerID.String()
				})).Return(subscriptiondomain.Subscription{ID: subID}, nil)
				// 3. Validate Entitlement -> Success
				s.On("ValidateUsageEntitlement", mock.Anything, subID, meterID, mock.Anything).Return(nil)
			},
			expectedErr:  nil,
			expectIngest: true,
		},
		{
			name: "Failure: Meter Not Found",
			req: usagedomain.CreateIngestRequest{
				CustomerID:     customerID.String(),
				MeterCode:      "invalid_meter",
				Value:          10,
				RecordedAt:     time.Now(),
				IdempotencyKey: "k2",
			},
			setupMocks: func(s *subscriptionMock, m *meterMock) {
				m.On("GetByCode", mock.Anything, "invalid_meter").Return(nil, meterdomain.ErrMeterNotFound)
			},
			expectedErr:  usagedomain.ErrInvalidMeter, // Assuming service maps it or returns ErrInvalidMeter if nil
			expectIngest: false,
		},
		{
			name: "Failure: No Active Subscription",
			req: usagedomain.CreateIngestRequest{
				CustomerID:     customerID.String(),
				MeterCode:      "m1",
				Value:          10,
				RecordedAt:     time.Now(),
				IdempotencyKey: "k3",
			},
			setupMocks: func(s *subscriptionMock, m *meterMock) {
				// 1. Resolve Subscription -> Not Found
				s.On("GetActiveByCustomerID", mock.Anything, mock.Anything).Return(subscriptiondomain.Subscription{}, subscriptiondomain.ErrSubscriptionNotFound)
			},
			expectedErr:  usagedomain.ErrInvalidSubscription,
			expectIngest: false,
		},
		{
			name: "Failure: Not Entitled",
			req: usagedomain.CreateIngestRequest{
				CustomerID:     customerID.String(),
				MeterCode:      "m1",
				Value:          10,
				RecordedAt:     time.Now(),
				IdempotencyKey: "k4",
			},
			setupMocks: func(s *subscriptionMock, m *meterMock) {
				m.On("GetByCode", mock.Anything, "m1").Return(&meterdomain.Response{
					ID:   meterID.String(),
					Code: "m1",
				}, nil)
				s.On("GetActiveByCustomerID", mock.Anything, mock.Anything).Return(subscriptiondomain.Subscription{ID: subID}, nil)
				// 3. Validate Entitlement -> FeatureNotEntitled
				s.On("ValidateUsageEntitlement", mock.Anything, subID, meterID, mock.Anything).Return(subscriptiondomain.ErrFeatureNotEntitled)
			},
			expectedErr:  errors.New("usage_rejected_feature_not_entitled"),
			expectIngest: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSub := new(subscriptionMock)
			mockMeter := new(meterMock)

			if tt.setupMocks != nil {
				tt.setupMocks(mockSub, mockMeter)
			}

			svc := NewService(ServiceParam{
				DB:       db,
				Log:      logger,
				GenID:    genID,
				MeterSvc: mockMeter,
				SubSvc:   mockSub,
			})

			ctx := WithTestOrgContext(context.Background(), orgID)
			t.Logf("Test %s: orgID=%s", tt.name, orgID)

			res, err := svc.Ingest(ctx, tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				if tt.expectIngest == false {
					assert.Nil(t, res)
				}
				if tt.expectedErr.Error() != "" {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				if tt.expectIngest {
					assert.NotNil(t, res)
				}
			}
		})
	}
}

func TestIngest_Idempotency_Strict(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&usagedomain.UsageEvent{})
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS ux_usage_events_idempotency ON usage_events(org_id, idempotency_key)")

	node, _ := snowflake.NewNode(1)
	orgID := node.Generate()
	customerID := node.Generate()
	subID := node.Generate()
	meterID := node.Generate()

	mockSub := new(subscriptionMock)
	mockMeter := new(meterMock)

	// Setup valid responses
	mockMeter.On("GetByCode", mock.Anything, "m1").Return(&meterdomain.Response{ID: meterID.String(), Code: "m1"}, nil)
	mockSub.On("GetActiveByCustomerID", mock.Anything, mock.MatchedBy(func(req subscriptiondomain.GetActiveByCustomerIDRequest) bool {
		return req.CustomerID == customerID.String()
	})).Return(subscriptiondomain.Subscription{ID: subID}, nil)

	// Entitlement passes initially
	mockSub.On("ValidateUsageEntitlement", mock.Anything, subID, meterID, mock.Anything).Return(nil)

	svc := NewService(ServiceParam{
		DB:       db,
		Log:      zap.NewNop(),
		GenID:    node,
		MeterSvc: mockMeter,
		SubSvc:   mockSub,
	})

	ctx := WithTestOrgContext(context.Background(), orgID)
	key := "strict_idempotency_key_1"
	req := usagedomain.CreateIngestRequest{
		CustomerID:     customerID.String(),
		MeterCode:      "m1",
		Value:          1,
		RecordedAt:     time.Now(),
		IdempotencyKey: key,
	}

	// 1. First Call: Success
	res1, err := svc.Ingest(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res1)
	assert.Equal(t, key, res1.IdempotencyKey)

	// 2. Second Call: Success (Same Result)
	res2, err := svc.Ingest(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, res1.ID, res2.ID, "Must return same ID")

	// 3. Third Call: With Entitlement Failure (Simulate cancelled sub)
	// We create a NEW service instance or update mocks if possible.
	// Since mocks are referenced, we can try to override or since we used "Return", it might be sticky.
	// Let's rely on the "Fast Path" preventing the mock call entirely.
	// If the code calls ValidateUsageEntitlement, it will return nil (pass) because of the mock above.
	// To verify logic, we should assert that ValidateUsageEntitlement is NOT called again?
	// Or we can construct a test where the fast path is the ONLY way to succeed.
}

func TestIngest_Idempotency_BypassEntitlementFailure(t *testing.T) {
	// dedicated test for the "Entitlement Revoked" Case
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&usagedomain.UsageEvent{})
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS ux_usage_events_idempotency ON usage_events(org_id, idempotency_key)")

	node, _ := snowflake.NewNode(1)
	orgID := node.Generate()
	// Pre-seed the DB with an event
	key := "pre_existing_key"
	existingEvent := usagedomain.UsageEvent{
		ID:             node.Generate(),
		OrgID:          orgID,
		IdempotencyKey: key,
		Status:         usagedomain.UsageStatusAccepted,
		RecordedAt:     time.Now(),
	}
	db.Create(&existingEvent)

	// Mocks that FAIL EVERYTHING
	mockSub := new(subscriptionMock)
	mockMeter := new(meterMock)
	// If the service calls these, they will panic or return error if we set them to.
	// We won't setup any expectations. If they are called, test fails (unexpected call) or return zero/error.
	// Actually testify mock returns error on unexpected call.
	// So if we don't setup On(), calls panic. Perfect.

	svc := NewService(ServiceParam{
		DB:       db,
		Log:      zap.NewNop(),
		GenID:    node,
		MeterSvc: mockMeter,
		SubSvc:   mockSub,
	})

	ctx := WithTestOrgContext(context.Background(), orgID)
	req := usagedomain.CreateIngestRequest{
		CustomerID:     "cust_1", // Doesn't matter, shouldn't reach resolve
		MeterCode:      "m1",
		Value:          10,
		RecordedAt:     time.Now(),
		IdempotencyKey: key,
	}

	// Call Ingest
	res, err := svc.Ingest(ctx, req)

	// Expect Success (Returning existing) WITHOUT calling mocks
	assert.NoError(t, err)
	assert.Equal(t, existingEvent.ID, res.ID)
}

// Helper for context
func WithTestOrgContext(ctx context.Context, orgID snowflake.ID) context.Context {
	return orgcontext.WithOrgID(ctx, int64(orgID))
}
