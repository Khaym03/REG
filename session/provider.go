package session

import (
	"context"

	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type Provider struct {
	browser       *rod.Browser
	session       Session
	loginHandler  dcommand.CommandHandler[command.LoginCommand]
	logoutHandler dcommand.CommandHandler[command.LogoutCommand]
}

func NewProvider(
	browser *rod.Browser,
	loginH dcommand.CommandHandler[command.LoginCommand],
	logoutH dcommand.CommandHandler[command.LogoutCommand],
) *Provider {
	return &Provider{
		browser:       browser,
		loginHandler:  loginH,
		logoutHandler: logoutH,
	}
}

func (p *Provider) Start(ctx context.Context, user domain.User) (Session, error) {
	session, err := NewRodSession(p.browser)
	if err != nil {
		return nil, err
	}

	err = p.loginHandler.Handle(ctx, command.LoginCommand{
		User: user,
		Page: session.MainPage(),
	})

	if err != nil {
		session.Close()
		return nil, err
	}

	p.session = session

	return session, nil
}

func (p *Provider) End(ctx context.Context) error {
	if p.session == nil {
		return nil
	}
	return p.logoutHandler.Handle(ctx, command.LogoutCommand{Page: p.session.MainPage()})
}

func (p *Provider) Get() Session {
	if p.session == nil {
		panic("session not started: call Provider.Start() first")
	}
	return p.session
}
