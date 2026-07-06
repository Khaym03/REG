package session

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/event"
	"github.com/go-rod/rod"
	"github.com/google/uuid"
)

var _ Session = (*RodSession)(nil)

var SessionClosed = errors.New("session closed")

type RodSession struct {
	page *rod.Page
	mu   sync.Mutex

	eventBus event.Bus
	closed   atomic.Bool
}

// This implementation only allow 1 page
func NewRodSession(
	page *rod.Page,
	eventBus event.Bus,
) *RodSession {
	return &RodSession{
		page:     page,
		mu:       sync.Mutex{},
		eventBus: eventBus,
	}
}

func (s *RodSession) Do(ctx context.Context, fn browser.PageFunc) error {
	if s.closed.Load() {
		return SessionClosed
	}

	errCh := make(chan error, 1)

	page := s.page.Context(ctx)
	go func() {
		errCh <- fn(page)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

// ID implements [Session].
func (s *RodSession) ID() SessionID {
	return SessionID("rod-session-" + uuid.New().String())
}

func (s *RodSession) Close() error {
	if s.closed.Load() {
		return nil
	}

	s.closed.Store(true)

	return nil
}
