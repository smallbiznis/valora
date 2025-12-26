package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	productdomain "github.com/smallbiznis/valora/internal/product/domain"
)

type createProductRequest struct {
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Active      *bool          `json:"active"`
	Metadata    map[string]any `json:"metadata"`
}

func (s *Server) CreateProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.productSvc.Create(c.Request.Context(), productdomain.CreateRequest{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: req.Description,
		Active:      req.Active,
		Metadata:    req.Metadata,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) ListProducts(c *gin.Context) {
	resp, err := s.productSvc.List(c.Request.Context())
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) GetProductByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	resp, err := s.productSvc.Get(c.Request.Context(), id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func isProductValidationError(err error) bool {
	switch err {
	case productdomain.ErrInvalidOrganization,
		productdomain.ErrInvalidCode,
		productdomain.ErrInvalidName,
		productdomain.ErrInvalidID:
		return true
	default:
		return false
	}
}
