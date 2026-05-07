package auth

import (
	"context"

	"github.com/Khaym03/REG/domain"
)

type AuthService interface {
	Login(context.Context, domain.Session, domain.User) error
	Logout(context.Context, domain.Session) error
}
