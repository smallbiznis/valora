package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	priceamountdomain "github.com/smallbiznis/valora/internal/priceamount/domain"
)

type createPriceAmountRequest struct {
	OrganizationID     string         `json:"organization_id"`
	PriceID            string         `json:"price_id"`
	MeterID            *string        `json:"meter_id"`
	Currency           string         `json:"currency"`
	UnitAmountCents    int64          `json:"unit_amount_cents"`
	MinimumAmountCents *int64         `json:"minimum_amount_cents"`
	MaximumAmountCents *int64         `json:"maximum_amount_cents"`
	Metadata           map[string]any `json:"metadata"`
}

func (s *Server) CreatePriceAmount(c *gin.Context) {
	var req createPriceAmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.priceAmountSvc.Create(c.Request.Context(), priceamountdomain.CreateRequest{
		OrganizationID:     strings.TrimSpace(req.OrganizationID),
		PriceID:            strings.TrimSpace(req.PriceID),
		MeterID:            req.MeterID,
		Currency:           strings.TrimSpace(req.Currency),
		UnitAmountCents:    req.UnitAmountCents,
		MinimumAmountCents: req.MinimumAmountCents,
		MaximumAmountCents: req.MaximumAmountCents,
		Metadata:           req.Metadata,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) ListPriceAmounts(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	resp, err := s.priceAmountSvc.List(c.Request.Context(), orgID)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) GetPriceAmountByID(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	id := strings.TrimSpace(c.Param("id"))
	resp, err := s.priceAmountSvc.Get(c.Request.Context(), orgID, id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func isPriceAmountValidationError(err error) bool {
	switch err {
	case priceamountdomain.ErrInvalidOrganization,
		priceamountdomain.ErrInvalidPrice,
		priceamountdomain.ErrInvalidCurrency,
		priceamountdomain.ErrInvalidUnitAmount,
		priceamountdomain.ErrInvalidMinAmount,
		priceamountdomain.ErrInvalidMaxAmount,
		priceamountdomain.ErrInvalidMeterID,
		priceamountdomain.ErrInvalidID:
		return true
	default:
		return false
	}
}
