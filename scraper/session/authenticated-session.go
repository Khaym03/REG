package session

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/domain"
)

var _ domain.Session = (*AuthenticatedSession)(nil)

type AuthenticatedSession struct {
	base domain.Session
	auth AuthService
	user domain.User
}

func NewAuthenticatedSession(
	ctx context.Context,
	base domain.Session,
	auth AuthService,
	user domain.User,
) (*AuthenticatedSession, error) {

	if err := auth.Login(ctx, base, user); err != nil {
		return nil, err
	}

	return &AuthenticatedSession{
		base: base,
		auth: auth,
	}, nil
}

func (s *AuthenticatedSession) Do(ctx context.Context, fn domain.PageFunc) error {
	return s.base.Do(ctx, fn)
}

func (s *AuthenticatedSession) NewIsolated(ctx context.Context) (domain.Session, error) {
	return s.base.NewIsolated(ctx)
}

func (s *AuthenticatedSession) Close() error {
	// logout first, then close
	if err := s.auth.Logout(context.Background(), s.base); err != nil {
		log.Error(err)
	}

	return s.base.Close()
}
