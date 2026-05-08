package container

import (
	"github.com/Khaym03/REG/internal/auth"
	"github.com/go-rod/rod"
)

func buildSessionProvider(
	browser *rod.Browser,
	service auth.AuthService,
) *auth.Provider {
	return auth.NewProvider(browser, service)
}
