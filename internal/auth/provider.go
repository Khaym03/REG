package auth

import (
	"context"

	"github.com/go-rod/rod"
	"github.com/mustafaturan/bus/v3"
)

type Provider struct {
	browser *rod.Browser
	auth    AuthService

	eventBus *bus.Bus
}

func NewProvider(
	browser *rod.Browser,
	auth AuthService,
	eventBus *bus.Bus,
) *Provider {
	return &Provider{
		browser:  browser,
		auth:     auth,
		eventBus: eventBus,
	}
}

func (p *Provider) Start(
	ctx context.Context,
	user User,
) (Session, error) {
	base, err := NewRodSession(p.browser)
	if err != nil {
		return nil, err
	}

	sess, err := NewAuthenticatedSession(ctx, base, p.auth, user, p.eventBus)
	if err != nil {
		base.Close()
		return nil, err
	}

	return sess, nil
}
