package domain

import (
	"context"
	"time"
)

//go:generate mockgen -source=service.go -destination=../mocks/mock_service.go -package=mocks
type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResult, error)
	Logout(ctx context.Context, rawToken string) error
	Authenticate(ctx context.Context, rawToken string) (*Session, error)
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	CurrentUser(ctx context.Context) (*User, error)
}

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

type LoginRequest struct {
	Username  string
	Password  string
	UserAgent string
	IPAddress string
}

type LoginResult struct {
	Session   *SessionView
	RawToken  string
	ExpiresAt time.Time
}
