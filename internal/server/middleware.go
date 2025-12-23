package server

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HeaderOrg        = "X-Org-ID"
	contextUserIDKey = "user_id"
)

func serveIndex(c *gin.Context) {
	c.File("./public/index.html")
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate / propagate request id
		c.Next()
	}
}

func (s *Server) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(sessionCookieName)
		fmt.Printf("sid: %v\n", sid)
		fmt.Printf("error: %v\n", err)
		if err != nil || strings.TrimSpace(sid) == "" {
			AbortWithError(c, ErrUnauthorized)
			return
		}

		session, err := s.authsvc.Authenticate(c.Request.Context(), sid)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.Set(contextUserIDKey, session.UserID.String())
		c.Next()
	}
}

func OrgContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// resolve org from token / header
		// inject to context
		c.Next()
	}
}

func RequireRole(role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check role from context
		c.Next()
	}
}
