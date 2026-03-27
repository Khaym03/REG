package app

import (
	"context"
	"log"
	"time"

	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
)

type WorkFlowInput struct {
	User domain.User
}

type ReceptionWorkflow struct {
	loginHandler  dcommand.CommandHandler[command.LoginCommand]
	logoutHandler dcommand.CommandHandler[command.LogoutCommand]
	gatherHandler dcommand.CommandHandler[command.GatherGuidesCommand]
}

func NewReceptionWorkflow(
	loginH dcommand.CommandHandler[command.LoginCommand],
	logoutH dcommand.CommandHandler[command.LogoutCommand],
	gatherH dcommand.CommandHandler[command.GatherGuidesCommand],
) *ReceptionWorkflow {
	return &ReceptionWorkflow{
		loginHandler:  loginH,
		logoutHandler: logoutH,
		gatherHandler: gatherH,
	}
}

func (w *ReceptionWorkflow) Run(ctx context.Context, input WorkFlowInput) (err error) {
	err = w.loginHandler.Handle(ctx, command.LoginCommand{
		User: input.User,
	})
	if err != nil {
		return err
	}

	defer func() {
		logoutErr := w.logoutHandler.Handle(ctx, command.LogoutCommand{})

		if logoutErr != nil {
			if err != nil {
				log.Printf("logout failed: %v", logoutErr)
				return
			}

			err = logoutErr
		}
	}()

	err = w.gatherHandler.Handle(ctx, command.GatherGuidesCommand{
		From: time.Now().AddDate(-1, 0, 0),
		To:   time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
