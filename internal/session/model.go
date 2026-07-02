package session

import (
	"context"
	"io"

	"github.com/Khaym03/REG/internal/browser"
)

// Session provides controlled access to browser pages.
// Pages are short-lived and managed by the Session.
type Session interface {
	// Do executes fn with a managed page.
	// The page must not be used outside fn.
	Do(ctx context.Context, fn PageFunc) error

	ID() SessionID

	io.Closer
}

type SessionFactory interface {
	Create(ctx context.Context) (Session, error)
	CreateIsolated(ctx context.Context, parent Session) (Session, error)
}

type PageFunc = browser.PageFunc

type SessionID string
