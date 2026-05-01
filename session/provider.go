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

func (p *Provider) Start(ctx context.Context, user domain.User) (context.Context, error) {
	s, err := NewRodSession(p.browser)
	if err != nil {
		return nil, err
	}

	err = p.auth.Login(ctx, s.MainPage(), user)
	if err != nil {
		s.Close()
		return nil, err
	}

	ctx = WithSession(ctx, s)

	return ctx, nil
}

func (p *Provider) End(ctx context.Context) error {
	s := FromContext(ctx)

	err := p.auth.Logout(ctx, s.MainPage())

	// ensure cleanup even if logout fails
	closeErr := s.Close()

	if err != nil {
		return err
	}
	return closeErr
}
