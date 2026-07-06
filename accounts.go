package main

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/repo"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	usersFilepath = "users.json"
)

type AccountsAPI struct {
	ctx            context.Context
	service        *auth.AccountService
	sessionManager mediator.SessionMediator
}

func NewAccountsAPI(sm mediator.SessionMediator, eventBus event.Bus) *AccountsAPI {
	var p repo.Persistence[[]auth.RegisterUsers] = repo.NewJSONPersistence(
		usersFilepath,
		func() []auth.RegisterUsers {
			return nil
		},
	)

	return &AccountsAPI{
		service: auth.NewAccountService(
			auth.NewLoginScraper(eventBus),
			p,
		),
		sessionManager: sm,
	}
}

func (api *AccountsAPI) ServiceStartup(
	ctx context.Context,
	_ application.ServiceOptions,
) error {

	api.ctx = ctx

	return nil
}

func (api *AccountsAPI) AuthUser(user auth.User) error {
	if api.service.KnownUser(user) {
		return nil
	}

	s, err := api.sessionManager.Create(api.ctx)
	if err != nil {
		return err
	}

	err = api.service.AuthUser(api.ctx, user, s)
	if err != nil {
		return err
	}

	defer func() {
		if err = api.sessionManager.Logout(api.ctx, s); err != nil {
			log.Error(err)
		}
	}()

	return nil
}

func (api *AccountsAPI) GetRegisterUsers() ([]auth.RegisterUsers, error) {
	return api.service.GetRegisterUsers()
}

func (api *AccountsAPI) GetUserPassword(username string) (auth.User, error) {
	return api.service.GetUserPassword(username)
}

func (api *AccountsAPI) UpdateUser(user auth.RegisterUsers) error {
	return api.service.UpdateUser(user)
}

func (api *AccountsAPI) CurrentUser() *auth.RegisterUsers {
	return api.service.CurrentUser()
}
