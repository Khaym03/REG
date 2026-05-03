package session

import (
	"context"

	"github.com/Khaym03/REG/domain"

	"github.com/go-rod/rod"
)

type Provider struct {
	browser *rod.Browser
	auth    AuthService
}

func NewProvider(
	browser *rod.Browser,
	auth AuthService,

) *Provider {
	return &Provider{
		browser: browser,
		auth:    auth,
	}
}

func (p *Provider) Start(
	ctx context.Context,
	user domain.User,
) (domain.Session, error) {
	s, err := NewRodSession(p.browser)
	if err != nil {
		return nil, err
	}

	err = p.auth.Login(ctx, s, user)
	if err != nil {
		s.Close()
		return nil, err
	}

	return s, nil
}

func (p *Provider) End(ctx context.Context, session domain.Session) error {
	err := p.auth.Logout(ctx, session)

	// ensure cleanup even if logout fails
	closeErr := session.Close()

	if err != nil {
		return err
	}

	return closeErr
}
