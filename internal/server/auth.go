package server

import (
	"net/http"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	authdomain "github.com/smallbiznis/valora/internal/auth/domain"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) Login(c *gin.Context) {

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	result, err := s.authsvc.Login(c.Request.Context(), authdomain.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: c.Request.UserAgent(),
		IPAddress: c.ClientIP(),
	})
	if err != nil {
		AbortWithError(c, err)
		return
	}

	s.sessions.Set(c, result.RawToken, result.ExpiresAt)

	s.enrichSessionMetadata(c, result)

	c.JSON(http.StatusOK, result.Session)
}

func (s *Server) Logout(c *gin.Context) {
	token, ok := s.sessions.ReadToken(c)
	if !ok {
		AbortWithError(c, ErrUnauthorized)
		return
	}

	if err := s.authsvc.Logout(c.Request.Context(), token); err != nil {
		AbortWithError(c, err)
		return
	}

	s.sessions.Clear(c)
	c.Status(http.StatusNoContent)
}

func (s *Server) Me(c *gin.Context) {
	AbortWithError(c, ErrServiceUnavailable)
}

func (s *Server) Forgot(c *gin.Context) {
	AbortWithError(c, ErrServiceUnavailable)
}

func (s *Server) enrichSessionMetadata(c *gin.Context, result *authdomain.LoginResult) {
	if result == nil || result.Session == nil {
		return
	}

	rawUserID, ok := result.Session.Metadata["user_id"].(string)
	if !ok {
		return
	}

	parsedUserID, err := snowflake.ParseString(rawUserID)
	if err != nil {
		return
	}

	orgIDs, err := s.loadUserOrgIDs(c.Request.Context(), parsedUserID)
	if err != nil {
		return
	}

	if err := s.authsvc.UpdateSessionOrgContext(c.Request.Context(), result.SessionID, nil, orgIDs); err != nil {
		return
	}

	result.Session.Metadata["org_ids"] = toOrgIDStrings(orgIDs)
}
