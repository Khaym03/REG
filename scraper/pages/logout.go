package pages

import (
	"github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
)

type LogoutPage struct {
	page *rod.Page
}

func NewLogoutPage(p *rod.Page) *LogoutPage {
	return &LogoutPage{page: p}
}

func (lp *LogoutPage) Open() error {
	return navigate(lp.page, constants.BaseURL)
}

func (lp *LogoutPage) Logout() (err error) {
	if err = click(lp.page, profileDropdownSelector); err != nil {
		return err
	}

	return click(lp.page, logoutSelector)
}

const (
	profileDropdownSelector = `//*[@id="profileDropdown"]`
	logoutSelector          = `/html/body/div[1]/nav/div[2]/ul/li[3]/div/a`
)
