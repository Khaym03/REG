package command

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type LoginCommand struct {
	User domain.User
	Page *rod.Page
}

type LoginHandler struct {
	auth domain.AuthService
}

func NewLoginHandler(auth domain.AuthService) LoginHandler {
	return LoginHandler{auth: auth}
}

func (h LoginHandler) Handle(ctx context.Context, cmd LoginCommand) error {
	return h.auth.Login(ctx, cmd.Page, cmd.User)
}
