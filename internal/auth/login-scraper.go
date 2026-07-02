package auth

import (
	"context"
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
	return &LoginScraper{
		eventBus: eventBus,
	}
}

func (l *LoginScraper) Login(
	ctx context.Context,
	s session.Session,
	user User,
) (err error) {
	login := func(p *rod.Page) (err error) {
		err = l.eventBus.Emit(ctx, string(event.Login), struct{}{})
		if err != nil {
			return err
		}

		loginPage := NewLoginPage(p)

		if err = loginPage.Open(); err != nil {
			return
		}

		if err = loginPage.EnterCredentials(user); err != nil {
			return
		}

		if err = loginPage.Submit(); err != nil {
			return
		}

		return
	}

	login = browser.WithRetry(ctx, 3, time.Second*10)(login)

	return s.Do(ctx, login)
}

func (l *LoginScraper) Logout(
	ctx context.Context,
	s session.Session,
) (err error) {
	logout := func(p *rod.Page) error {
		err = l.eventBus.Emit(ctx, string(event.Logout), struct{}{})
		if err != nil {
			return err
		}

		logoutPage := NewLogoutPage(p)

		if err = logoutPage.Open(); err != nil {
			return err
		}

		return logoutPage.Logout()
	}

	logout = browser.WithRetry(ctx, 3, c.DefaultTimeout)(logout)

	return s.Do(ctx, logout)
}
