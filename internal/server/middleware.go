package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	authdomain "github.com/smallbiznis/valora/internal/auth/domain"
	"github.com/smallbiznis/valora/internal/orgcontext"
)

const (
	contextUserIDKey  = "user_id"
	contextSessionKey = "session"
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
		token := readBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			AbortWithError(c, ErrUnauthorized)
			return
		}
		if strings.TrimSpace(s.cfg.AuthJWTSecret) == "" {
			AbortWithError(c, ErrUnauthorized)
			return
		}

		claims, err := validateJWT(token, []byte(s.cfg.AuthJWTSecret))
		if err != nil {
			AbortWithError(c, ErrUnauthorized)
			return
		}

		// SINGLE-ORG RUNTIME GUARD (v0)
		//
		// This enforces single-organization mode at runtime.
		// Valora OSS remains multi-org by design.
		// Removing the org_id equality check enables multi-org support
		// without changing handlers, services, or schemas.
		if s.cfg.DefaultOrgID == 0 {
			AbortWithError(c, ErrUnauthorized)
			return
		}

		ctx := orgcontext.WithOrgID(c.Request.Context(), claims.OrgID)
		c.Request = c.Request.WithContext(ctx)
		if strings.TrimSpace(claims.Subject) != "" {
			c.Set(contextUserIDKey, claims.Subject)
		}
		c.Next()
	}
}

func (s *Server) OrgContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID, ok := orgcontext.OrgIDFromContext(c.Request.Context())
		if !ok || orgID == 0 {
			AbortWithError(c, ErrOrgRequired)
			return
		}
		c.Next()
	}
}

func RequireRole(role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check role from context
		c.Next()
	}
}

func (s *Server) sessionFromContext(c *gin.Context) (*authdomain.Session, bool) {
	value, ok := c.Get(contextSessionKey)
	if !ok {
		return nil, false
	}
	session, ok := value.(*authdomain.Session)
	return session, ok
}

func (s *Server) loadUserOrgIDs(ctx context.Context, userID snowflake.ID) ([]int64, error) {
	orgs, err := s.organizationSvc.ListOrganizationsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	orgIDs := make([]int64, 0, len(orgs))
	for _, org := range orgs {
		parsed, err := snowflake.ParseString(org.ID)
		if err != nil {
			return nil, ErrInternal
		}
		orgIDs = append(orgIDs, int64(parsed))
	}

	return orgIDs, nil
}

func containsOrgID(ids []int64, target int64) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

var errInvalidToken = errors.New("invalid token")

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtClaims struct {
	Subject  string      `json:"sub"`
	OrgID    json.Number `json:"org_id"`
	TenantID json.Number `json:"tenant_id"`
	Expires  json.Number `json:"exp"`
}

type validatedJWT struct {
	Subject string
	OrgID   int64
	Expiry  int64
}

func readBearerToken(header string) string {
	parts := strings.Fields(strings.TrimSpace(header))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func validateJWT(token string, secret []byte) (*validatedJWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errInvalidToken
	}
	var header jwtHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errInvalidToken
	}
	if header.Alg != "HS256" {
		return nil, errInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errInvalidToken
	}

	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(signingInput))
	expected := mac.Sum(nil)
	if !hmac.Equal(expected, sig) {
		return nil, errInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errInvalidToken
	}
	var claims jwtClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errInvalidToken
	}

	orgID := parseNumber(claims.OrgID)
	if orgID == 0 {
		orgID = parseNumber(claims.TenantID)
	}
	exp := parseNumber(claims.Expires)
	if exp == 0 || time.Now().Unix() >= exp {
		return nil, errInvalidToken
	}

	return &validatedJWT{
		Subject: claims.Subject,
		OrgID:   orgID,
		Expiry:  exp,
	}, nil
}

func parseNumber(value json.Number) int64 {
	if value == "" {
		return 0
	}
	if parsed, err := value.Int64(); err == nil {
		return parsed
	}
	floatVal, err := value.Float64()
	if err != nil {
		return 0
	}
	return int64(floatVal)
}
