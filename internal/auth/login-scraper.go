package auth

import (
	"context"
	"time"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/go-rod/rod"
)

var _ AuthService = (*LoginScraper)(nil)

type LoginScraper struct{}

func NewLoginScraper() *LoginScraper {
	return &LoginScraper{}
}

func (l *LoginScraper) Login(
	ctx context.Context,
	s domain.Session,
	user domain.User,
) (err error) {
	login := func(p *rod.Page) (err error) {
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

	login = browser.WithRetry(3, time.Second*10)(login)

	return s.Do(ctx, login)
}

func (l *LoginScraper) Logout(
	ctx context.Context,
	s domain.Session,
) (err error) {
	logout := func(p *rod.Page) error {
		logoutPage := NewLogoutPage(p)

		if err = logoutPage.Open(); err != nil {
			return err
		}

		return logoutPage.Logout()
	}

	logout = browser.WithRetry(3, time.Second*10)(logout)

	return s.Do(ctx, logout)
}
