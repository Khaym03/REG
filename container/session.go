package container

import (
	"github.com/Khaym03/REG/app/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/session"
	"github.com/go-rod/rod"
)

func buildSessionProvider(
	browser *rod.Browser,
	auth domain.AuthService,
) *session.Provider {

	login := withRetryAndLogging(
		command.NewLoginHandler(auth),
	)

	logout := withRetryAndLogging(
		command.NewLogoutHandler(auth),
	)

	return session.NewProvider(browser, login, logout)
}
