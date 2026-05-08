package auth

import (
	"context"

	"github.com/Khaym03/REG/internal/domain"

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
	base, err := NewRodSession(p.browser)
	if err != nil {
		return nil, err
	}

	sess, err := NewAuthenticatedSession(ctx, base, p.auth, user)
	if err != nil {
		base.Close()
		return nil, err
	}

	return sess, nil
}
