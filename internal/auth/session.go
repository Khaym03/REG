package auth

import (
	"context"
	"errors"
	"sync"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var SessionClosed = errors.New("session closed")

type RodSession struct {
	browser *rod.Browser
	page    *rod.Page
	mu      sync.Mutex
	closed  bool
}

// This implementation only allow 1 page
func NewRodSession(browser *rod.Browser) (*RodSession, error) {
	page, err := browser.Page(proto.TargetCreateTarget{URL: ""})
	if err != nil {
		return nil, err
	}

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	return &RodSession{
		browser: browser,
		page:    page,
		mu:      sync.Mutex{},
	}, nil
}

func (s *RodSession) Do(ctx context.Context, fn browser.PageFunc) error {
	if s.isClosed() {
		return SessionClosed
	}

	return fn(s.page)
}

func (s *RodSession) NewIsolated(ctx context.Context) (domain.Session, error) {
	if s.isClosed() {
		return nil, SessionClosed
	}

	incognito, err := s.browser.Incognito()
	if err != nil {
		return nil, err
	}

	return NewRodSession(incognito)
}

func (s *RodSession) Close() error {
	if s.isClosed() {
		return nil
	}

	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()

	return s.browser.Close()
}

func (s *RodSession) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.closed
}
