package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authdomain "github.com/smallbiznis/valora/internal/auth/domain"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const sessionCookieName = "_sid"

func (s *Server) Login(c *gin.Context) {

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	result, err := s.authsvc.Login(c.Request.Context(), authdomain.LoginRequest{
		Username:  req.Username,
		Password:  req.Password,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	maxAge := int(time.Until(result.ExpiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	s.setSessionCookie(c, result.RawToken, maxAge)

	c.JSON(http.StatusOK, result.Session)
}

func (s *Server) Logout(c *gin.Context) {
	token, err := c.Cookie(sessionCookieName)
	if err != nil {
		AbortWithError(c, ErrUnauthorized)
		return
	}

	if err := s.authsvc.Logout(c.Request.Context(), token); err != nil {
		AbortWithError(c, err)
		return
	}

	s.setSessionCookie(c, "", -1)
	c.Status(http.StatusNoContent)
}

func (s *Server) Me(c *gin.Context) {
	AbortWithError(c, ErrServiceUnavailable)
}

func (s *Server) Forgot(c *gin.Context) {
	AbortWithError(c, ErrServiceUnavailable)
}

func (s *Server) setSessionCookie(c *gin.Context, value string, maxAge int) {
	secure := s.cfg.Environment == "production"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(sessionCookieName, value, maxAge, "/", "", secure, true)
}
