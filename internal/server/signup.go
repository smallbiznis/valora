package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authdomain "github.com/smallbiznis/railzway/internal/auth/domain"
	signupdomain "github.com/smallbiznis/railzway/internal/signup/domain"
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
		s.sessions.Set(c, result.RawToken, result.ExpiresAt)
		s.enrichSessionFromToken(c, result.Session, result.RawToken)
	}

	c.JSON(http.StatusOK, result.Session)
}

func (s *Server) enrichSessionFromToken(c *gin.Context, view *authdomain.SessionView, rawToken string) {
	if view == nil || rawToken == "" {
		return
	}

	session, err := s.authsvc.Authenticate(c.Request.Context(), rawToken)
	if err != nil || session == nil {
		return
	}

	orgIDs, err := s.loadUserOrgIDs(c.Request.Context(), session.UserID)
	if err != nil {
		return
	}

	if err := s.authsvc.UpdateSessionOrgContext(c.Request.Context(), session.ID, nil, orgIDs); err != nil {
		return
	}

	if view.Metadata == nil {
		view.Metadata = map[string]any{}
	}
	view.Metadata["org_ids"] = toOrgIDStrings(orgIDs)
}
