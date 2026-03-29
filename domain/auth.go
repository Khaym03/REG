package domain

import (
	"context"

	"github.com/go-rod/rod"
)

type AuthService interface {
	Login(context.Context, *rod.Page, User) error
	Logout(context.Context, *rod.Page) error
}
