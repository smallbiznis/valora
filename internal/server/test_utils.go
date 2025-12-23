package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type testCleanupRequest struct {
	Prefix string `json:"prefix"`
}

func (s *Server) TestCleanup(c *gin.Context) {
	if s.cfg.Environment == "production" {
		AbortWithError(c, ErrNotFound)
		return
	}

	var req testCleanupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		AbortWithError(c, invalidRequestError())
		return
	}

	prefix := strings.TrimSpace(req.Prefix)
	if prefix == "" {
		AbortWithError(c, newValidationError("prefix", "required", "prefix is required"))
		return
	}

	ctx := c.Request.Context()
	like := prefix + "%"

	var orgIDs []int64
	if err := s.db.WithContext(ctx).
		Table("organizations").
		Select("id").
		Where("name LIKE ?", like).
		Scan(&orgIDs).Error; err != nil {
		AbortWithError(c, err)
		return
	}

	if len(orgIDs) > 0 {
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM price_tiers WHERE org_id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM price_amounts WHERE org_id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM prices WHERE org_id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM products WHERE org_id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM customers WHERE org_id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM organization_members WHERE org_id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM organizations WHERE id IN ?`, orgIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
	}

	var userIDs []int64
	if err := s.db.WithContext(ctx).
		Table("users").
		Select("id").
		Where("username LIKE ?", like).
		Scan(&userIDs).Error; err != nil {
		AbortWithError(c, err)
		return
	}

	if len(userIDs) > 0 {
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM sessions WHERE user_id IN ?`, userIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM organization_members WHERE user_id IN ?`, userIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
		if err := s.db.WithContext(ctx).Exec(
			`DELETE FROM users WHERE id IN ?`, userIDs,
		).Error; err != nil {
			AbortWithError(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
