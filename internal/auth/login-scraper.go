package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/internal/browser"
	c "github.com/Khaym03/REG/internal/constants"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/session"
	"github.com/go-rod/rod"
	"github.com/mustafaturan/bus/v3"
)

var _ AuthService = (*LoginScraper)(nil)

type LoginScraper struct {
	eventBus *bus.Bus
}

func NewLoginScraper(eventBus *bus.Bus) *LoginScraper {
	return &LoginScraper{eventBus: eventBus}
}

func (l *LoginScraper) Login(
	ctx context.Context,
	s session.Session,
	user User,
) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	login := func(p *rod.Page) (err error) {
		if err := l.eventBus.Emit(ctx, string(event.Login), struct{}{}); err != nil {
			return fmt.Errorf("emit login event: %w", err)
		}

		loginPage := NewLoginPage(p)

		if err := loginPage.Open(); err != nil {
			return fmt.Errorf("open login page: %w", err)
		}

		if err := loginPage.EnterCredentials(user); err != nil {
			return fmt.Errorf("enter credentials: %w", err)
		}

		if err := loginPage.Submit(); err != nil {
			return fmt.Errorf("submit login: %w", err)
		}

		return nil
	}

	login = browser.WithRetry(ctx, 3, time.Second*10)(login)

	if err := s.Do(ctx, login); err != nil {
		return fmt.Errorf("session do login: %w", err)
	}

	return nil
}

func (l *LoginScraper) Logout(
	ctx context.Context,
	s session.Session,
) (err error) {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	logout := func(p *rod.Page) error {
		if err := l.eventBus.Emit(ctx, string(event.Logout), struct{}{}); err != nil {
			return fmt.Errorf("emit logout event: %w", err)
		}

		logoutPage := NewLogoutPage(p)

		if err := logoutPage.Open(); err != nil {
			return fmt.Errorf("open logout page: %w", err)
		}

		if err := logoutPage.Logout(); err != nil {
			return fmt.Errorf("perform logout: %w", err)
		}

		return nil
	}

	logout = browser.WithRetry(ctx, 3, c.DefaultTimeout)(logout)

	if err := s.Do(ctx, logout); err != nil {
		return fmt.Errorf("session do logout: %w", err)
	}

	return nil
}
