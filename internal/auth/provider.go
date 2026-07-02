package auth

// import (
// 	"context"

// 	"github.com/go-rod/rod"
// 	"github.com/mustafaturan/bus/v3"
// )

// type Provider struct {
// 	auth     AuthService
// 	eventBus *bus.Bus
// }

// func NewProvider(
// 	auth AuthService,
// 	eventBus *bus.Bus,
// ) *Provider {
// 	return &Provider{
// 		auth:     auth,
// 		eventBus: eventBus,
// 	}
// }

// func (p *Provider) Start(
// 	ctx context.Context,
// 	user User,
// 	browser *rod.Browser,
// ) (Session, error) {
// 	base, err := NewRodSession(browser, p.eventBus, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	sess, err := NewAuthenticatedSession(ctx, base, p.auth, user, p.eventBus)
// 	if err != nil {
// 		base.Close(ctx)
// 		return nil, err
// 	}

// 	return sess, nil
// }
