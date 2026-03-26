package command

import (
	"context"

	"github.com/Khaym03/REG/domain"
)

type LoginCommand struct {
	User domain.User
}

type LoginHandler struct {
	auth domain.AuthService
}

func NewLoginHandler(auth domain.AuthService) LoginHandler {
	return LoginHandler{auth: auth}
}

func (h LoginHandler) Handle(ctx context.Context, cmd LoginCommand) error {
	return h.auth.Login(ctx, cmd.User)
}
