package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pricetierdomain "github.com/smallbiznis/valora/internal/pricetier/domain"
)

type createPriceTierRequest struct {
	OrganizationID  string         `json:"organization_id"`
	PriceID         string         `json:"price_id"`
	TierMode        int16          `json:"tier_mode"`
	StartQuantity   float64        `json:"start_quantity"`
	EndQuantity     *float64       `json:"end_quantity"`
	UnitAmountCents *int64         `json:"unit_amount_cents"`
	FlatAmountCents *int64         `json:"flat_amount_cents"`
	Unit            string         `json:"unit"`
	Metadata        map[string]any `json:"metadata"`
}

func (s *Server) ListPriceTiers(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	resp, err := s.priceTierSvc.List(c.Request.Context(), orgID)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) GetPriceTierByID(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	id := strings.TrimSpace(c.Param("id"))
	resp, err := s.priceTierSvc.Get(c.Request.Context(), orgID, id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) CreatePriceTier(c *gin.Context) {
	var req createPriceTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.priceTierSvc.Create(c.Request.Context(), pricetierdomain.CreateRequest{
		OrganizationID:  strings.TrimSpace(req.OrganizationID),
		PriceID:         strings.TrimSpace(req.PriceID),
		TierMode:        req.TierMode,
		StartQuantity:   req.StartQuantity,
		EndQuantity:     req.EndQuantity,
		UnitAmountCents: req.UnitAmountCents,
		FlatAmountCents: req.FlatAmountCents,
		Unit:            strings.TrimSpace(req.Unit),
		Metadata:        req.Metadata,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func isPriceTierValidationError(err error) bool {
	switch err {
	case pricetierdomain.ErrInvalidOrganization,
		pricetierdomain.ErrInvalidPrice,
		pricetierdomain.ErrInvalidTierMode,
		pricetierdomain.ErrInvalidStartQty,
		pricetierdomain.ErrInvalidEndQty,
		pricetierdomain.ErrInvalidUnitAmount,
		pricetierdomain.ErrInvalidFlatAmount,
		pricetierdomain.ErrInvalidUnit,
		pricetierdomain.ErrInvalidID:
		return true
	default:
		return false
	}
}
