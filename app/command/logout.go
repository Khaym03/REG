package command

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type LogoutCommand struct {
	Page *rod.Page
}

type LogoutHandler struct {
	auth domain.AuthService
}

func NewLogoutHandler(auth domain.AuthService) LogoutHandler {
	return LogoutHandler{auth: auth}
}

func (h LogoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	return h.auth.Logout(ctx, cmd.Page)
}
