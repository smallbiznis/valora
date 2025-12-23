package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usagedomain "github.com/smallbiznis/valora/internal/usage/domain"
)

func (s *Server) IngestUsage(c *gin.Context) {

	var req usagedomain.CreateIngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, err)
		return
	}

	usage, err := s.usagesvc.Ingest(c.Request.Context(), req)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, usage)
}
