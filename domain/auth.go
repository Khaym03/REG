package domain

import "context"

type AuthService interface {
	Login(context.Context, User) error
	Logout(context.Context) error
}
