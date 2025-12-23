package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pricedomain "github.com/smallbiznis/valora/internal/price/domain"
)

func (s *Server) CreatePricing(c *gin.Context) {
	var req pricedomain.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.priceSvc.Create(c.Request.Context(), pricedomain.CreateRequest{
		OrganizationID: strings.TrimSpace(req.OrganizationID),
		Name:           strings.TrimSpace(req.Name),
		Description:    req.Description,
		Active:         req.Active,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) ListPricings(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	resp, err := s.priceSvc.List(c.Request.Context(), orgID)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) GetPricingByID(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	id := strings.TrimSpace(c.Param("id"))
	resp, err := s.priceSvc.Get(c.Request.Context(), orgID, id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func isPricingValidationError(err error) bool {
	switch err {
	case pricedomain.ErrInvalidOrganization,
		pricedomain.ErrInvalidProduct,
		pricedomain.ErrInvalidCode,
		pricedomain.ErrInvalidPricingModel,
		pricedomain.ErrInvalidBillingMode,
		pricedomain.ErrInvalidBillingInterval,
		pricedomain.ErrInvalidBillingIntervalCount,
		pricedomain.ErrInvalidAggregateUsage,
		pricedomain.ErrInvalidBillingUnit,
		pricedomain.ErrInvalidBillingThreshold,
		pricedomain.ErrInvalidTaxBehavior,
		pricedomain.ErrInvalidVersion,
		pricedomain.ErrInvalidID:
		return true
	default:
		return false
	}
}
