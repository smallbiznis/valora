package server

import "github.com/gin-gonic/gin"

func (s *Server) RunRatingJob(c *gin.Context) {
	AbortWithError(c, ErrServiceUnavailable)
}
