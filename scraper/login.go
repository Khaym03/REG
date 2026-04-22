package scraper

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/go-rod/rod"
)

var _ domain.AuthService = (*LoginScraper)(nil)

type LoginScraper struct {
}

func NewLoginScraper() *LoginScraper {
	return &LoginScraper{}
}

func (l *LoginScraper) Login(ctx context.Context, page *rod.Page, user domain.User) (err error) {
	loginPage := pages.NewLoginPage(page.Context(ctx))

	if err = loginPage.Open(); err != nil {
		return
	}

	if err = loginPage.EnterCredentials(user); err != nil {
		return
	}

	if err = loginPage.Submit(); err != nil {
		return
	}

	return nil
}

func (l *LoginScraper) Logout(ctx context.Context, page *rod.Page) (err error) {
	logoutPage := pages.NewLogoutPage(page.Context(ctx))

	if err = logoutPage.Open(); err != nil {
		return err
	}

	return logoutPage.Logout()
}
