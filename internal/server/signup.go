package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	signupdomain "github.com/smallbiznis/valora/internal/signup/domain"
)

type SignupRequest struct {
	OrgName  string `json:"org_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) Signup(c *gin.Context) {
	if !s.cfg.IsCloud() {
		AbortWithError(c, ErrNotFound)
		return
	}

	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	result, err := s.signupsvc.Signup(c.Request.Context(), signupdomain.Request{
		OrgName:   req.OrgName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	if result.Session != nil && result.RawToken != "" {
		maxAge := int(time.Until(result.ExpiresAt).Seconds())
		if maxAge < 0 {
			maxAge = 0
		}
		s.setSessionCookie(c, result.RawToken, maxAge)
	}

	c.JSON(http.StatusOK, result.Session)
}
