package session

import (
	"context"

	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"

	"github.com/go-rod/rod"
)

type Provider struct {
	browser *rod.Browser
	// session       Session
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

func (p *Provider) Start(ctx context.Context, user domain.User) (context.Context, error) {
	s, err := NewRodSession(p.browser)
	if err != nil {
		return nil, err
	}

	err = p.loginHandler.Handle(ctx, command.LoginCommand{
		User: user,
		Page: s.MainPage(),
	})
	if err != nil {
		s.Close()
		return nil, err
	}

	ctx = WithSession(ctx, s)

	return ctx, nil
}

func (p *Provider) End(ctx context.Context) error {
	s := FromContext(ctx)

	err := p.logoutHandler.Handle(ctx, command.LogoutCommand{
		Page: s.MainPage(),
	})

	// ensure cleanup even if logout fails
	closeErr := s.Close()

	if err != nil {
		return err
	}
	return closeErr
}
