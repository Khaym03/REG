package session

import (
	"context"
	"fmt"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/mustafaturan/bus/v3"
)

var _ SessionFactory = (*sessionFactory)(nil)

type sessionFactory struct {
	pool     browser.BrowserPool
	eventBus *bus.Bus
}

func NewSessionFactory(bp browser.BrowserPool, eventBus *bus.Bus) *sessionFactory {
	return &sessionFactory{
		pool:     bp,
		eventBus: eventBus,
	}
}

// Create implements [SessionFactory].
func (s *sessionFactory) Create(ctx context.Context) (Session, error) {
	bruntime, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		s.pool.Release(bruntime)
	}()

	b := bruntime.Browser()

	page, err := browser.CreatePageIfNotExist(b)
	if err != nil {
		return nil, fmt.Errorf("fail at CreatePageIfNotExist: %w", err)
	}

	return NewRodSession(page, s.eventBus), nil
}

// CreateIsolated implements [SessionFactory].
func (s *sessionFactory) CreateIsolated(
	ctx context.Context,
	parent Session,
) (Session, error) {
	bruntime, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		s.pool.Release(bruntime)
	}()

	b := bruntime.Browser()

	incognito, err := b.Incognito()
	if err != nil {
		return nil, err
	}

	err = browser.CopyCookies(b, incognito)
	if err != nil {
		return nil, err
	}

	page, err := browser.CreatePageIfNotExist(incognito)
	if err != nil {
		return nil, err
	}

	return NewRodSession(page, s.eventBus), nil
}
