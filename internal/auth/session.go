package auth

import (
	"context"
	"errors"
	"sync"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/event"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/mustafaturan/bus/v3"
	"github.com/sirupsen/logrus"
)

var _ Session = (*RodSession)(nil)

var SessionClosed = errors.New("session closed")

type RodSession struct {
	browser *rod.Browser
	page    *rod.Page
	mu      sync.Mutex
	closed  bool

	eventBus *bus.Bus
}

// This implementation only allow 1 page
func NewRodSession(browser *rod.Browser, eventBus *bus.Bus) (*RodSession, error) {
	page, err := browser.Page(proto.TargetCreateTarget{URL: ""})
	if err != nil {
		return nil, err
	}

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	return &RodSession{
		browser:  browser,
		page:     page,
		mu:       sync.Mutex{},
		eventBus: eventBus,
	}, nil
}

func (s *RodSession) Do(ctx context.Context, fn browser.PageFunc) error {
	if s.isClosed() {
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

func (s *RodSession) NewIsolated(ctx context.Context) (Session, error) {
	if s.isClosed() {
		return nil, SessionClosed
	}

	incognito, err := s.browser.Incognito()
	if err != nil {
		return nil, err
	}

	return NewRodSession(incognito, s.eventBus)
}

func (s *RodSession) Close(ctx context.Context) error {
	if s.isClosed() {
		return nil
	}

	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()

	if err := s.eventBus.Emit(ctx, event.DestroyingBrowser, struct{}{}); err != nil {
		logrus.Error(err)
	}
	return s.browser.Close()
}

func (s *RodSession) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.closed
}
