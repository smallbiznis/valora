package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pricetierdomain "github.com/smallbiznis/valora/internal/pricetier/domain"
)

func (s *Server) ListPriceTiers(c *gin.Context) {
	resp, err := s.priceTierSvc.List(c.Request.Context())
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) GetPriceTierByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	resp, err := s.priceTierSvc.Get(c.Request.Context(), id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) CreatePriceTier(c *gin.Context) {
	var req pricetierdomain.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	req.PriceID = strings.TrimSpace(req.PriceID)
	req.Unit = strings.TrimSpace(req.Unit)

	resp, err := s.priceTierSvc.Create(c.Request.Context(), req)
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
