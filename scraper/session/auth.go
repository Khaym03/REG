package session

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type AuthService interface {
	Login(context.Context, *rod.Page, domain.User) error
	Logout(context.Context, *rod.Page) error
}
