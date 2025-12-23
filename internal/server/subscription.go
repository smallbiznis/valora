package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	subscriptiondomain "github.com/smallbiznis/valora/internal/subscription/domain"
)

func (s *Server) CreateSubscription(c *gin.Context) {
	var req subscriptiondomain.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.subscriptionSvc.Create(c.Request.Context(), subscriptiondomain.CreateSubscriptionRequest{
		OrganizationID: strings.TrimSpace(req.OrganizationID),
		CustomerID:     strings.TrimSpace(req.CustomerID),
		Items:          normalizeSubscriptionItems(req.Items),
		Metadata:       req.Metadata,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) ListSubscriptions(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	status := strings.TrimSpace(c.Query("status"))
	pageToken := strings.TrimSpace(c.Query("page_token"))
	pageSize := int32(0)
	if pageSizeRaw := strings.TrimSpace(c.Query("page_size")); pageSizeRaw != "" {
		parsedSize, err := strconv.Atoi(pageSizeRaw)
		if err != nil {
			AbortWithError(c, newValidationError("page_size", "invalid_page_size", "invalid page size"))
			return
		}
		pageSize = int32(parsedSize)
	}

	resp, err := s.subscriptionSvc.List(c.Request.Context(), subscriptiondomain.ListSubscriptionRequest{
		OrgID:     orgID,
		Status:    status,
		PageToken: pageToken,
		PageSize:  pageSize,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp.Subscriptions, "page_info": resp.PageInfo})
}

func (s *Server) GetSubscriptionByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if _, err := snowflake.ParseString(id); err != nil {
		AbortWithError(c, newValidationError("id", "invalid_id", "invalid id"))
		return
	}

	item, err := s.subscriptionSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (s *Server) CancelSubscription(c *gin.Context) {
	AbortWithError(c, ErrServiceUnavailable)
}

func normalizeSubscriptionItems(items []subscriptiondomain.CreateSubscriptionItemRequest) []subscriptiondomain.CreateSubscriptionItemRequest {
	if len(items) == 0 {
		return nil
	}

	normalized := make([]subscriptiondomain.CreateSubscriptionItemRequest, 0, len(items))
	for _, item := range items {
		normalized = append(normalized, subscriptiondomain.CreateSubscriptionItemRequest{
			PriceID:  strings.TrimSpace(item.PriceID),
			Quantity: item.Quantity,
		})
	}
	return normalized
}

func isSubscriptionValidationError(err error) bool {
	switch {
	case errors.Is(err, subscriptiondomain.ErrInvalidOrganization),
		errors.Is(err, subscriptiondomain.ErrInvalidCustomer),
		errors.Is(err, subscriptiondomain.ErrInvalidSubscription),
		errors.Is(err, subscriptiondomain.ErrInvalidMeterID),
		errors.Is(err, subscriptiondomain.ErrInvalidMeterCode),
		errors.Is(err, subscriptiondomain.ErrInvalidStatus),
		errors.Is(err, subscriptiondomain.ErrInvalidCollectionMode),
		errors.Is(err, subscriptiondomain.ErrInvalidBillingCycleType),
		errors.Is(err, subscriptiondomain.ErrInvalidStartAt),
		errors.Is(err, subscriptiondomain.ErrInvalidPeriod),
		errors.Is(err, subscriptiondomain.ErrInvalidItems),
		errors.Is(err, subscriptiondomain.ErrInvalidQuantity),
		errors.Is(err, subscriptiondomain.ErrInvalidPrice),
		errors.Is(err, subscriptiondomain.ErrMultipleFlatPrices):
		return true
	default:
		return false
	}
}
