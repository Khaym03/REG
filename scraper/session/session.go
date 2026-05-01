package session

import "github.com/go-rod/rod"

type Session interface {
	MainPage() *rod.Page
	NewPage() *rod.Page
	Browser() *rod.Browser
	Close() error
}

type RodSession struct {
	browser *rod.Browser
	page    *rod.Page
}

func NewRodSession(browser *rod.Browser) (*RodSession, error) {
	page := browser.MustPage()
	page.MustWaitLoad()

	return &RodSession{
		browser: browser,
		page:    page,
	}, nil
}

func (s *RodSession) MainPage() *rod.Page {
	return s.page
}

func (s *RodSession) NewPage() *rod.Page {
	return s.browser.MustPage()
}

func (s *RodSession) Browser() *rod.Browser {
	return s.browser
}

func (s *RodSession) Close() error {
	return s.page.Close()
}
