package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	meterdomain "github.com/smallbiznis/valora/internal/meter/domain"
)

type createMeterRequest struct {
	OrganizationID  string  `json:"organization_id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	AggregationType string  `json:"aggregation_type"`
	Unit            string  `json:"unit"`
	Description     *string `json:"description"`
	Active          *bool   `json:"active"`
}

type updateMeterRequest struct {
	Name            *string `json:"name,omitempty"`
	AggregationType *string `json:"aggregation_type,omitempty"`
	Unit            *string `json:"unit,omitempty"`
	Active          *bool   `json:"active,omitempty"`
}

func (s *Server) CreateMeter(c *gin.Context) {
	var req createMeterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.meterSvc.Create(c.Request.Context(), meterdomain.CreateRequest{
		OrganizationID: strings.TrimSpace(req.OrganizationID),
		Code:           strings.TrimSpace(req.Code),
		Name:           strings.TrimSpace(req.Name),
		Aggregation:    strings.TrimSpace(req.AggregationType),
		Unit:           strings.TrimSpace(req.Unit),
		Active:         req.Active,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) ListMeters(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	resp, err := s.meterSvc.List(c.Request.Context(), orgID)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) GetMeterByID(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	id := strings.TrimSpace(c.Param("id"))
	resp, err := s.meterSvc.GetByID(c.Request.Context(), orgID, id)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) UpdateMeter(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	id := strings.TrimSpace(c.Param("id"))

	var req updateMeterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	resp, err := s.meterSvc.Update(c.Request.Context(), meterdomain.UpdateRequest{
		OrganizationID: orgID,
		ID:             id,
		Name:           trimStringPtr(req.Name),
		Aggregation:    trimStringPtr(req.AggregationType),
		Unit:           trimStringPtr(req.Unit),
		Active:         req.Active,
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (s *Server) DeleteMeter(c *gin.Context) {
	orgID := strings.TrimSpace(c.Query("organization_id"))
	id := strings.TrimSpace(c.Param("id"))

	if err := s.meterSvc.Delete(c.Request.Context(), orgID, id); err != nil {
		AbortWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func isMeterValidationError(err error) bool {
	switch err {
	case meterdomain.ErrInvalidOrganization,
		meterdomain.ErrInvalidCode,
		meterdomain.ErrInvalidName,
		meterdomain.ErrInvalidAggregation,
		meterdomain.ErrInvalidUnit,
		meterdomain.ErrInvalidID:
		return true
	default:
		return false
	}
}

func trimStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}
