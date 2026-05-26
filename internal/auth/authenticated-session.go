package auth

import (
	"context"

	"github.com/Khaym03/REG/internal/event"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"
)

var _ Session = (*AuthenticatedSession)(nil)

type AuthenticatedSession struct {
	base Session
	auth AuthService
	user User

	eventBus *bus.Bus
}

func NewAuthenticatedSession(
	ctx context.Context,
	base Session,
	auth AuthService,
	user User,
	eventBus *bus.Bus,
) (*AuthenticatedSession, error) {
	if err := eventBus.Emit(ctx, event.LogginTopic, struct{}{}); err != nil {
		log.Error(err)
	}

	if err := auth.Login(ctx, base, user); err != nil {
		return nil, err
	}

	return &AuthenticatedSession{
		base:     base,
		auth:     auth,
		eventBus: eventBus,
	}, nil
}

func (s *AuthenticatedSession) Do(ctx context.Context, fn PageFunc) error {
	return s.base.Do(ctx, fn)
}

func (s *AuthenticatedSession) NewIsolated(ctx context.Context) (Session, error) {
	return s.base.NewIsolated(ctx)
}

func (s *AuthenticatedSession) Close() error {
	s.eventBus.Emit(context.TODO(), event.LogoutTopic, struct{}{})
	// logout first, then close
	if err := s.auth.Logout(context.Background(), s.base); err != nil {
		log.Error(err)
	}

	return s.base.Close()
}
