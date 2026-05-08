package domain

import (
	"context"

	"github.com/Khaym03/REG/internal/browser"
)

// Session provides controlled access to browser pages.
// Pages are short-lived and managed by the Session.
type Session interface {
	// Do executes fn with a managed page.
	// The page must not be used outside fn.
	Do(ctx context.Context, fn browser.PageFunc) error

	// NewIsolated creates a new Session with isolated browser state.
	// Inheriting its cookies.
	// The caller must call Close.
	NewIsolated(ctx context.Context) (Session, error)

	// Close cleans up the Session and its resources.
	Close() error
}

type GuideCollector interface {
	Collect(context.Context, Session, DateRange) ([]Guide, error)
}

type RubroExtractor interface {
	FromGuides(context.Context, Session, []Guide) ([]Rubro, error)
}

type ReceptionService interface {
	Receive(context.Context, Session, DateRange) (ReceptionResult, error)
}

type InventoryService interface {
	Snapshot(context.Context, Session) ([]Rubro, error)
	Insert(context.Context, Session, Rubro) error
}
