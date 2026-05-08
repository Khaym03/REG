package app

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/auth"
	dcommand "github.com/Khaym03/REG/internal/common/decorator/command"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/guide"
	"github.com/Khaym03/REG/internal/inventory"
	"github.com/Khaym03/REG/internal/reception"
)

type WorkFlowInput struct {
	User auth.User
	Date domain.DateRange
}

type ReceptionWorkflow struct {
	sessionProvider      *auth.Provider
	gatherHandler        dcommand.CommandHandler[guide.GatherGuidesCommand]
	syncInventoryHandler dcommand.CommandHandler[inventory.SyncInventoryCommand]
	receptionistHandler  dcommand.CommandHandler[reception.ReceptionistCommand]
}

func NewReceptionWorkflow(
	sp *auth.Provider,
	gatherH dcommand.CommandHandler[guide.GatherGuidesCommand],
	syncInventoryH dcommand.CommandHandler[inventory.SyncInventoryCommand],
	receptionistH dcommand.CommandHandler[reception.ReceptionistCommand],
) *ReceptionWorkflow {
	return &ReceptionWorkflow{
		sessionProvider:      sp,
		gatherHandler:        gatherH,
		syncInventoryHandler: syncInventoryH,
		receptionistHandler:  receptionistH,
	}
}

func (w *ReceptionWorkflow) Run(ctx context.Context, input WorkFlowInput) error {
	session, err := w.sessionProvider.Start(ctx, input.User)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := session.Close(); cerr != nil {
			log.Error("cleanup error:", cerr)
		}

		log.Info("Done")
	}()

	err = w.gatherHandler.Handle(ctx, session, guide.GatherGuidesCommand{
		DateRange: input.Date,
	})
	if err != nil {
		return err
	}

	err = w.syncInventoryHandler.Handle(ctx, session, inventory.SyncInventoryCommand{})
	if err != nil {
		return err
	}

	err = w.receptionistHandler.Handle(ctx, session, reception.ReceptionistCommand{
		DateRange: input.Date,
	})
	if err != nil {
		return err
	}

	return nil
}
