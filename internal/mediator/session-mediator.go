package mediator

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/session"
	log "github.com/sirupsen/logrus"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

//wails:internal
type SessionMediator interface {
	io.Closer
	SessionFactory

	auth.AuthService

	browser.Reconfigurer

	GetSession(id SessionID) (Session, bool)
	DeleteSession(id SessionID) error
}

var _ SessionMediator = (*sessionMediator)(nil)

type sessionMediator struct {
	pool     BrowserPool
	factory  SessionFactory
	registry SessionRegistry
	auth     auth.AuthService
}

func NewSessionMediator(
	eventBus event.Bus,
) *sessionMediator {

	pool := browser.NewBrowserPool(browser.BrowserPoolConfig{})

	return &sessionMediator{
		pool:     pool,
		factory:  session.NewSessionFactory(pool, eventBus),
		auth:     auth.NewLoginScraper(eventBus),
		registry: *session.NewSessionRegistry(),
	}
}

// Reconfigure implements [SessionMediator].
func (s *sessionMediator) Reconfigure(ctx context.Context, cfg config.BrowserConfig) error {
	b, err := s.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer func() {
		s.pool.Release(b)
	}()

	return b.Reconfigure(ctx, cfg)
}

// Close implements [SessionMediator].
func (s *sessionMediator) Close() error {
	return s.pool.Close()
}

// Create implements [SessionMediator].
func (s *sessionMediator) Create(ctx context.Context) (Session, error) {
	sess, err := s.factory.Create(ctx)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	log.Infof("added session: %s", sess.ID())
	s.registry.Add(sess)

	return sess, nil
}

// CreateIsolated implements [SessionMediator].
func (s *sessionMediator) CreateIsolated(
	ctx context.Context,
	parent Session,
) (Session, error) {

	// parent, ok := s.registry.Get(parent.ID())
	// if !ok {
	// 	return nil, fmt.Errorf("parent: %w", ErrSessionNotFound)
	// }

	sess, err := s.factory.CreateIsolated(ctx, parent)
	if err != nil {
		return nil, err
	}

	s.registry.Add(sess)

	return sess, nil
}

// Login implements [SessionMediator].
func (s *sessionMediator) Login(
	ctx context.Context,
	session Session,
	user auth.User,
) error {
	log.Infof("receive session: %s", session.ID())
	sess, ok := s.registry.Get(session.ID())
	if !ok {
		s.registry.Add(session)
		return s.auth.Login(ctx, session, user)
		// return fmt.Errorf("login err: %w", ErrSessionNotFound)
	}

	return s.auth.Login(ctx, sess, user)
}

// Logout implements [SessionMediator].
func (s *sessionMediator) Logout(ctx context.Context, session Session) error {
	sess, ok := s.registry.Get(session.ID())
	if !ok {
		s.registry.Add(session)
		return s.auth.Logout(ctx, session)
	}

	return s.auth.Logout(ctx, sess)
}

func (m *sessionMediator) GetSession(id SessionID) (Session, bool) {
	return m.registry.Get(id)
}

func (m *sessionMediator) DeleteSession(id SessionID) error {
	sess, ok := m.registry.Get(id)
	if !ok {
		return fmt.Errorf("delete: %w", ErrSessionNotFound)
	}

	if err := sess.Close(); err != nil {
		return err
	}

	m.registry.Remove(id)
	return nil
}

type Session = session.Session
type SessionID = session.SessionID
type BrowserPool = browser.BrowserPool
type SessionFactory = session.SessionFactory
type SessionRegistry = session.SessionRegistry
