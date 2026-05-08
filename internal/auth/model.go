package auth

import (
	"context"

	"github.com/Khaym03/REG/internal/browser"
)

// Session provides controlled access to browser pages.
// Pages are short-lived and managed by the Session.
type Session interface {
	// Do executes fn with a managed page.
	// The page must not be used outside fn.
	Do(ctx context.Context, fn PageFunc) error

	// NewIsolated creates a new Session with isolated browser state.
	// Inheriting its cookies.
	// The caller must call Close.
	NewIsolated(ctx context.Context) (Session, error)

	// Close cleans up the Session and its resources.
	Close() error
}

type AuthService interface {
	Login(context.Context, Session, User) error
	Logout(context.Context, Session) error
}

type User struct {
	Username string
	Password string
}

type (
	PageFunc = browser.PageFunc
)
