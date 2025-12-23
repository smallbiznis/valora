package server

import (
	"net/http"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	invoicedomain "github.com/smallbiznis/valora/internal/invoice/domain"
)

func (s *Server) ListInvoices(c *gin.Context) {
	resp, err := s.invoiceSvc.List(c.Request.Context(), invoicedomain.ListInvoiceRequest{})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp.Invoices})
}

func (s *Server) GetInvoiceByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if _, err := snowflake.ParseString(id); err != nil {
		AbortWithError(c, newValidationError("id", "invalid_id", "invalid id"))
		return
	}

	item, err := s.invoiceSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}
