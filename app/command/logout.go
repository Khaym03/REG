package command

import (
	"context"

	"github.com/Khaym03/REG/domain"
)

type LogoutCommand struct{}

type LogoutHandler struct {
	auth domain.AuthService
}

func NewLogoutHandler(auth domain.AuthService) LogoutHandler {
	return LogoutHandler{auth: auth}
}

func (h LogoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	return h.auth.Logout(ctx)
}
