package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/smallbiznis/valora/internal/auth/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

const (
	argonTime    uint32 = 1
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
	argonSaltLen        = 16

	sessionTokenBytes = 32
	sessionTTL        = 24 * time.Hour
)

type Service struct {
	log         *zap.Logger
	repo        domain.Repository
	sessionRepo domain.SessionRepository
	genID       *snowflake.Node
}

func New(log *zap.Logger, repo domain.Repository, sessionRepo domain.SessionRepository, genID *snowflake.Node) domain.Service {
	return &Service{
		log:         log.Named("auth.service"),
		repo:        repo,
		sessionRepo: sessionRepo,
		genID:       genID,
	}
}

func (s *Service) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, domain.ErrInvalidCredentials
	}

	if _, err := s.repo.FindByUsername(ctx, req.Username); err == nil {
		return nil, domain.ErrUserExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	hashed, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:                  s.genID.Generate(),
		Username:            req.Username,
		Email:               req.Email,
		PasswordHash:        &hashed,
		IsDefault:           false,
		LastPasswordChanged: &now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResult, error) {
	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if user.PasswordHash == nil || !verifyPassword(req.Password, *user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	rawToken, err := newSessionToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	session := &domain.Session{
		ID:               s.genID.Generate(),
		UserID:           user.ID,
		SessionTokenHash: hashToken(rawToken),
		UserAgent:        strings.TrimSpace(req.UserAgent),
		IPAddress:        strings.TrimSpace(req.IPAddress),
		ExpiresAt:        now.Add(sessionTTL),
		CreatedAt:        now,
		LastSeenAt:       now,
	}
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	passwordState := "rotated"
	if user.IsDefault || user.LastPasswordChanged == nil {
		passwordState = "default"
	}

	return &domain.LoginResult{
		Session: &domain.SessionView{
			Metadata: map[string]any{
				"user_id":               user.ID.String(),
				"username":              user.Username,
				"email":                 user.Email,
				"is_default":            user.IsDefault,
				"last_password_changed": user.LastPasswordChanged,
				"password_state":        passwordState,
				"auth_provider":         "local",
			},
		},
		RawToken:  rawToken,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *Service) Logout(ctx context.Context, rawToken string) error {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return domain.ErrInvalidSession
	}

	session, err := s.sessionRepo.GetSessionByTokenHash(ctx, hashToken(token))
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return domain.ErrInvalidSession
		}
		return err
	}

	now := time.Now().UTC()
	return s.sessionRepo.RevokeSession(ctx, session.ID, now)
}

func (s *Service) Authenticate(ctx context.Context, rawToken string) (*domain.Session, error) {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return nil, domain.ErrInvalidSession
	}

	session, err := s.sessionRepo.GetSessionByTokenHash(ctx, hashToken(token))
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil, domain.ErrInvalidSession
		}
		return nil, err
	}

	now := time.Now().UTC()
	if session.RevokedAt != nil {
		return nil, domain.ErrSessionRevoked
	}
	if now.After(session.ExpiresAt) {
		return nil, domain.ErrSessionExpired
	}

	if err := s.sessionRepo.UpdateLastSeen(ctx, session.ID, now); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	if newPassword == "" {
		return domain.ErrInvalidCredentials
	}

	id, err := snowflake.ParseString(userID)
	if err != nil {
		return err
	}

	_, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	hashed, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	fields := map[string]any{
		"password_hash":         hashed,
		"last_password_changed": &now,
		"is_default":            false,
		"updated_at":            now,
	}

	return s.repo.UpdateFields(ctx, id, fields)
}

func (s *Service) CurrentUser(ctx context.Context) (*domain.User, error) {
	_ = ctx
	return nil, domain.ErrUserNotFound
}

func newSessionToken() (string, error) {
	buf := make([]byte, sessionTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads, saltB64, hashB64), nil
}

func verifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false
	}

	var memory uint32
	var timeCost uint32
	var threads uint8
	{
		params := strings.Split(parts[3], ",")
		if len(params) != 3 {
			return false
		}

		m, ok := strings.CutPrefix(params[0], "m=")
		if !ok {
			return false
		}
		t, ok := strings.CutPrefix(params[1], "t=")
		if !ok {
			return false
		}
		p, ok := strings.CutPrefix(params[2], "p=")
		if !ok {
			return false
		}

		m64, err := strconv.ParseUint(m, 10, 32)
		if err != nil {
			return false
		}
		t64, err := strconv.ParseUint(t, 10, 32)
		if err != nil {
			return false
		}
		p64, err := strconv.ParseUint(p, 10, 8)
		if err != nil {
			return false
		}

		memory = uint32(m64)
		timeCost = uint32(t64)
		threads = uint8(p64)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	check := argon2.IDKey([]byte(password), salt, timeCost, memory, threads, uint32(len(hash)))
	return subtle.ConstantTimeCompare(hash, check) == 1
}
