package scraper

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/Khaym03/REG/scraper/session"
	"github.com/go-rod/rod"
)

var _ session.AuthService = (*LoginScraper)(nil)

type LoginScraper struct {
}

func NewLoginScraper() *LoginScraper {
	return &LoginScraper{}
}

func (l *LoginScraper) Login(
	ctx context.Context,
	s domain.Session,
	user domain.User,
) (err error) {
	return s.Do(ctx, func(p *rod.Page) (err error) {
		loginPage := pages.NewLoginPage(p)

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
	})
}

func (l *LoginScraper) Logout(
	ctx context.Context,
	s domain.Session,
) (err error) {
	return s.Do(ctx, func(p *rod.Page) error {
		logoutPage := pages.NewLogoutPage(p)

		if err = logoutPage.Open(); err != nil {
			return err
		}

		return logoutPage.Logout()
	})
}
