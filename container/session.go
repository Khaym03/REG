package container

import (
	"github.com/Khaym03/REG/session"
	"github.com/go-rod/rod"
)

func buildSessionProvider(
	browser *rod.Browser,
	auth session.AuthService,
) *session.Provider {
	return session.NewProvider(browser, auth)
}
