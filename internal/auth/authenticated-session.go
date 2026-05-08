package auth

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var _ Session = (*AuthenticatedSession)(nil)

type AuthenticatedSession struct {
	base Session
	auth AuthService
	user User
}

func NewAuthenticatedSession(
	ctx context.Context,
	base Session,
	auth AuthService,
	user User,
) (*AuthenticatedSession, error) {
	if err := auth.Login(ctx, base, user); err != nil {
		return nil, err
	}

	return &AuthenticatedSession{
		base: base,
		auth: auth,
	}, nil
}

func (s *AuthenticatedSession) Do(ctx context.Context, fn PageFunc) error {
	return s.base.Do(ctx, fn)
}

func (s *AuthenticatedSession) NewIsolated(ctx context.Context) (Session, error) {
	return s.base.NewIsolated(ctx)
}

func (s *AuthenticatedSession) Close() error {
	// logout first, then close
	if err := s.auth.Logout(context.Background(), s.base); err != nil {
		log.Error(err)
	}

	return s.base.Close()
}
