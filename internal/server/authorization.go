package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
)

type actorContext struct {
	Type   string
	ID     snowflake.ID
	Scopes []string
}

func (s *Server) authorizeOrgAction(object string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := s.authorizeOrgActionWithContext(c, object, action); err != nil {
			AbortWithError(c, err)
			return
		}
		c.Next()
	}
}

func (s *Server) authorizeOrgActionWithContext(c *gin.Context, object string, action string) error {
	actor, ok := s.actorFromContext(c)
	if !ok {
		return ErrUnauthorized
	}

	orgID, err := s.orgIDFromRequest(c)
	if err != nil {
		return err
	}

	switch actor.Type {
	case "api_key":
		if !apiKeyScopeAllows(actor.Scopes, object, action) {
			return ErrForbidden
		}
		return nil
	case "user":
		if s.authzSvc == nil {
			return ErrForbidden
		}
		return s.authorizeForOrg(c, actor.subject(), orgID, object, action)
	default:
		return ErrUnauthorized
	}
}

func (s *Server) authorizeForOrg(c *gin.Context, actor string, orgID snowflake.ID, object string, action string) error {
	return s.authzSvc.Authorize(c.Request.Context(), actor, orgID.String(), strings.TrimSpace(object), strings.TrimSpace(action))
}

func (s *Server) actorFromContext(c *gin.Context) (actorContext, bool) {
	if c == nil {
		return actorContext{}, false
	}

	if authType, ok := c.Request.Context().Value(contextAuthTypeKey).(string); ok && strings.TrimSpace(authType) == "api_key" {
		apiKeyID, ok := apiKeyIDFromContext(c.Request.Context())
		if !ok {
			return actorContext{}, false
		}
		return actorContext{
			Type:   "api_key",
			ID:     apiKeyID,
			Scopes: apiKeyScopesFromContext(c.Request.Context()),
		}, true
	}

	userID, ok := s.userIDFromSession(c)
	if !ok {
		return actorContext{}, false
	}
	return actorContext{Type: "user", ID: userID}, true
}

func (a actorContext) subject() string {
	switch a.Type {
	case "user":
		return fmt.Sprintf("user:%s", a.ID.String())
	case "api_key":
		return fmt.Sprintf("api_key:%s", a.ID.String())
	default:
		return ""
	}
}

func apiKeyIDFromContext(ctx context.Context) (snowflake.ID, bool) {
	if ctx == nil {
		return 0, false
	}
	raw := ctx.Value(contextAPIKeyIDKey)
	switch value := raw.(type) {
	case int64:
		if value == 0 {
			return 0, false
		}
		return snowflake.ID(value), true
	case snowflake.ID:
		if value == 0 {
			return 0, false
		}
		return value, true
	case string:
		parsed, err := snowflake.ParseString(strings.TrimSpace(value))
		if err != nil || parsed == 0 {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func apiKeyScopesFromContext(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(contextAPIKeyScopesKey)
	scopes, ok := value.([]string)
	if !ok {
		return nil
	}
	return scopes
}

// API key scopes match exact actions or an object wildcard like "subscription.*" or "subscription:*".
func apiKeyScopeAllows(scopes []string, object string, action string) bool {
	if len(scopes) == 0 {
		return false
	}

	normalizedAction := strings.ToLower(strings.TrimSpace(action))
	if normalizedAction == "" {
		return false
	}
	normalizedObject := strings.ToLower(strings.TrimSpace(object))

	for _, scope := range scopes {
		normalizedScope := strings.ToLower(strings.TrimSpace(scope))
		if normalizedScope == "" {
			continue
		}
		if normalizedScope == "*" {
			return true
		}
		if normalizedScope == normalizedAction {
			return true
		}
		if normalizedObject != "" && (normalizedScope == normalizedObject+".*" || normalizedScope == normalizedObject+":*") {
			return true
		}
	}
	return false
}
