package app

import (
	"context"
	"log"
	"time"

	"github.com/Khaym03/REG/app/command"
	dcommand "github.com/Khaym03/REG/common/decorator/command"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/session"
	"github.com/Khaym03/REG/utils"
)

type WorkFlowInput struct {
	User domain.User
}

type ReceptionWorkflow struct {
	sessionProvider      *session.Provider
	gatherHandler        dcommand.CommandHandler[command.GatherGuidesCommand]
	syncInventoryHandler dcommand.CommandHandler[command.SyncInventoryCommand]
	receptionistHandler  dcommand.CommandHandler[command.ReceptionistCommand]
}

func NewReceptionWorkflow(
	sp *session.Provider,
	gatherH dcommand.CommandHandler[command.GatherGuidesCommand],
	syncInventoryH dcommand.CommandHandler[command.SyncInventoryCommand],
	receptionistH dcommand.CommandHandler[command.ReceptionistCommand],
) *ReceptionWorkflow {
	return &ReceptionWorkflow{
		sessionProvider:      sp,
		gatherHandler:        gatherH,
		syncInventoryHandler: syncInventoryH,
		receptionistHandler:  receptionistH,
	}
}

func (w *ReceptionWorkflow) Run(ctx context.Context, input WorkFlowInput) (err error) {
	session, err := w.sessionProvider.Start(ctx, input.User)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := session.Close(); closeErr != nil {
			if err != nil {
				log.Printf("session close failed: %v", closeErr)
			} else {
				err = closeErr
			}
		}

		if logoutErr := w.sessionProvider.End(ctx); logoutErr != nil {
			if err != nil {
				log.Printf("logout failed: %v", logoutErr)
			} else {
				err = logoutErr
			}
		}
	}()

	lastYearToPresent := utils.DateRange{
		From: time.Now().AddDate(-1, 0, 0),
		To:   time.Now(),
	}

	err = w.gatherHandler.Handle(ctx, command.GatherGuidesCommand{
		DateRange: lastYearToPresent,
	})
	if err != nil {
		return err
	}

	err = w.syncInventoryHandler.Handle(ctx, command.SyncInventoryCommand{})
	if err != nil {
		return err
	}

	err = w.receptionistHandler.Handle(ctx, command.ReceptionistCommand{
		DateRange: lastYearToPresent,
	})
	if err != nil {
		return err
	}

	return nil
}
