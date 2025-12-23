package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) ListCountries(c *gin.Context) {
	countries, err := s.refrepo.ListCountries(c.Request.Context())
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": countries})
}

func (s *Server) ListTimezones(c *gin.Context) {
	country := strings.TrimSpace(c.Query("country"))
	if country == "" {
		AbortWithError(c, newValidationError("country", "invalid_country", "invalid country"))
		return
	}

	timezones, err := s.refrepo.ListTimezonesByCountry(c.Request.Context(), country)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": timezones})
}

func (s *Server) ListCurrencies(c *gin.Context) {
	currencies, err := s.refrepo.ListCurrencies(c.Request.Context())
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": currencies})
}
