package auth

import (
	"context"
	"io"
	"sync"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/config"
	"github.com/go-rod/rod"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type SessionProvider interface {
	Session(context.Context) (Session, error)
}

type Reconfigurer interface {
	Reconfigure(ctx context.Context, cfg config.BrowserConfig) error
}

type SessionManager interface {
	SessionProvider

	Reconfigurer

	Init(context.Context) error
	Login(context.Context, User) error
	Logout(context.Context) error

	io.Closer
}

type sessionManager struct {
	eventBus *bus.Bus
	service  AuthService

	browser     *rod.Browser
	browserConf config.BrowserConfig
	stateMu     sync.Mutex

	ready chan struct{}

	initOnce sync.Once
	initErr  error

	sessionMu sync.RWMutex
	session   Session
}

func NewSessionManager(eventBus *bus.Bus) *sessionManager {
	return &sessionManager{
		eventBus:    eventBus,
		service:     NewLoginScraper(),
		ready:       make(chan struct{}),
		stateMu:     sync.Mutex{},
		browserConf: config.BrowserConfFromENV(),
	}
}

func (s *sessionManager) ServiceStartup(
	ctx context.Context,
	_ application.ServiceOptions,
) error {

	go func() {
		err := s.Init(ctx)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func (s *sessionManager) ServiceShutdown() error {
	err := s.Close()
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (s *sessionManager) Init(ctx context.Context) error {
	s.initOnce.Do(func() {
		b, err := browser.BuildBrowser(ctx, config.BrowserConfFromENV())
		if err != nil {
			s.initErr = err
			close(s.ready)
			return
		}

		sess, err := NewRodSession(b, s.eventBus, false)
		if err != nil {
			_ = b.Close()

			s.initErr = err
			close(s.ready)
			return
		}

		s.browser = b

		s.sessionMu.Lock()
		s.session = sess
		s.sessionMu.Unlock()

		close(s.ready)
	})

	return s.initErr
}

func (s *sessionManager) Session(ctx context.Context) (Session, error) {
	select {
	case <-s.ready:

	case <-ctx.Done():
		return nil, context.Cause(ctx)
	}

	if s.initErr != nil {
		return nil, s.initErr
	}

	s.sessionMu.RLock()
	defer s.sessionMu.RUnlock()

	return s.session, nil
}

func (s *sessionManager) Login(ctx context.Context, user User) error {
	current, err := s.Session(ctx)
	if err != nil {
		return err
	}

	authenticated, err := NewAuthenticatedSession(
		ctx,
		current,
		s.service,
		user,
		s.eventBus,
	)
	if err != nil {
		return err
	}

	s.sessionMu.Lock()
	s.session = authenticated
	s.sessionMu.Unlock()

	return nil
}

func (s *sessionManager) Logout(ctx context.Context) error {
	current, err := s.Session(ctx)
	if err != nil {
		return err
	}

	if err := s.service.Logout(ctx, current); err != nil {
		return err
	}

	root, err := NewRodSession(s.browser, s.eventBus, false)
	if err != nil {
		return err
	}

	s.sessionMu.Lock()
	old := s.session
	s.session = root
	s.sessionMu.Unlock()

	// The browser remains alive. We only discard the previous session.
	if old != nil {
		_ = old.Close(ctx)
	}

	return nil
}

func (s *sessionManager) Close() error {
	<-s.ready

	s.sessionMu.RLock()
	sess := s.session
	s.sessionMu.RUnlock()

	if sess != nil {
		_ = sess.Close(context.Background())
	}

	return s.browser.Close()
}

func (s *sessionManager) Reconfigure(ctx context.Context, cfg config.BrowserConfig) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	log.Info(s.browserConf, cfg, s.browserConf.Equal(cfg))

	if s.browserConf.Equal(cfg) {
		return nil
	}

	if s.session != nil {
		_ = s.session.Close(ctx)
		s.session = nil
	}

	if s.browser != nil {
		_ = s.browser.Close()
		s.browser = nil
	}

	browser, err := browser.BuildBrowser(ctx, cfg)
	if err != nil {
		return err
	}

	session, err := NewRodSession(browser, s.eventBus, false)
	if err != nil {
		_ = browser.Close()
		return err
	}

	s.browser = browser
	s.browserConf = cfg
	s.session = session

	return nil
}
