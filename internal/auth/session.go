package auth

// import (
// 	"context"
// 	"errors"
// 	"sync"

// 	"github.com/Khaym03/REG/internal/browser"
// 	"github.com/Khaym03/REG/internal/event"
// 	"github.com/go-rod/rod"
// 	"github.com/go-rod/rod/lib/proto"
// 	"github.com/mustafaturan/bus/v3"
// 	"github.com/sirupsen/logrus"
// )

// var _ Session = (*RodSession)(nil)

// var SessionClosed = errors.New("session closed")

// type RodSession struct {
// 	browser *rod.Browser
// 	page    *rod.Page
// 	mu      sync.Mutex
// 	closed  bool

// 	// a work around while think in a redesign for the ownership of browser
// 	closeBrowser bool

// 	eventBus *bus.Bus
// }

// // This implementation only allow 1 page
// func NewRodSession(
// 	browser *rod.Browser,
// 	eventBus *bus.Bus,
// 	closeBroser bool,
// ) (*RodSession, error) {
// 	pages, err := browser.Pages()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var page *rod.Page

// 	if !pages.Empty() {
// 		page = pages.First()
// 	} else {
// 		page, err = browser.Page(proto.TargetCreateTarget{URL: ""})
// 		if err != nil {
// 			return nil, err
// 		}

// 		if err := page.WaitLoad(); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return &RodSession{
// 		browser:  browser,
// 		page:     page,
// 		mu:       sync.Mutex{},
// 		eventBus: eventBus,
// 	}, nil
// }

// func (s *RodSession) Do(ctx context.Context, fn browser.PageFunc) error {
// 	if s.isClosed() {
// 		return SessionClosed
// 	}

// 	errCh := make(chan error, 1)

// 	page := s.page.Context(ctx)
// 	go func() {
// 		errCh <- fn(page)
// 	}()

// 	select {
// 	case err := <-errCh:
// 		return err
// 	case <-ctx.Done():
// 		return context.Cause(ctx)
// 	}
// }

// func (s *RodSession) NewIsolated(ctx context.Context) (Session, error) {
// 	if s.isClosed() {
// 		return nil, SessionClosed
// 	}

// 	currentCookies, err := s.browser.GetCookies()
// 	if err != nil {
// 		return nil, err
// 	}

// 	incognito, err := s.browser.Incognito()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var params []*proto.NetworkCookieParam
// 	for _, cookie := range currentCookies {
// 		// Copiamos el SourcePort localmente para poder sacar su puntero de forma segura
// 		sourcePortCopy := cookie.SourcePort

// 		param := &proto.NetworkCookieParam{
// 			Name:         cookie.Name,
// 			Value:        cookie.Value,
// 			Domain:       cookie.Domain,
// 			Path:         cookie.Path,
// 			Secure:       cookie.Secure,
// 			HTTPOnly:     cookie.HTTPOnly,
// 			SameSite:     cookie.SameSite,
// 			Expires:      cookie.Expires,
// 			Priority:     cookie.Priority,
// 			SameParty:    cookie.SameParty,
// 			SourceScheme: cookie.SourceScheme,
// 			SourcePort:   &sourcePortCopy, // Puntero a la copia del int
// 			PartitionKey: cookie.PartitionKey,
// 			// URL se omite intencionalmente porque ya pasamos Domain y Path explícitos
// 		}

// 		params = append(params, param)
// 	}
// 	err = incognito.SetCookies(params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return NewRodSession(incognito, s.eventBus, true)
// }

// func (s *RodSession) Close(ctx context.Context) error {
// 	if s.isClosed() {
// 		return nil
// 	}

// 	s.mu.Lock()
// 	s.closed = true
// 	s.mu.Unlock()

// 	if err := s.eventBus.Emit(ctx, string(event.DestroyingBrowser), struct{}{}); err != nil {
// 		logrus.Error(err)
// 	}

// 	if ctx.Err() != nil {
// 		err := context.Cause(ctx)
// 		if !errors.Is(err, context.Canceled) {
// 			return nil
// 		}
// 	}

// 	if s.closeBrowser {
// 		_ = s.page.Close()
// 		return s.browser.Close()
// 	}
// 	return nil
// }

// func (s *RodSession) isClosed() bool {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	return s.closed
// }
