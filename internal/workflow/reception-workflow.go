package workflow

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/workflow/app"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"

	"github.com/Khaym03/REG/internal/domain"
)

type WorkFlowInput struct {
	User                   auth.User        `json:"user"`
	Date                   domain.DateRange `json:"date"`
	ReceiveGuidesInTransit bool             `json:"receive_guides_in_transit,omitempty"`
}

type ReceptionWorkflow struct {
	app *app.Application
}

func NewReceptionWorkflow(
	app *app.Application,
) *ReceptionWorkflow {
	return &ReceptionWorkflow{
		app: app,
	}
}

func (w *ReceptionWorkflow) Run(ctx context.Context, input WorkFlowInput) error {
	session, err := w.app.SessionProvider.Start(ctx, input.User)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := session.Close(); cerr != nil {
			log.Error("cleanup error:", cerr)
		}

		log.Info("Done")
	}()

	stats, err := w.app.Queries.Stats.Handle(ctx, session, stats.StatsQuery{})
	if err != nil {
		return err
	}

	if !stats.HasActionableGuides(input.ReceiveGuidesInTransit) {
		log.Info("No actionable guides found")
		return nil
	}

	err = w.app.Commands.GatherGuides.Handle(ctx, session, guide.GatherGuidesCommand{
		DateRange: input.Date,
	})
	if err != nil {
		return err
	}

	err = w.app.Commands.SyncInventory.Handle(ctx, session, inventory.SyncInventoryCommand{})
	if err != nil {
		return err
	}

	err = w.app.Commands.Receptionist.Handle(ctx, session, reception.ReceptionistCommand{
		Date:                   input.Date,
		ReceiveGuidesInTransit: input.ReceiveGuidesInTransit,
	})
	if err != nil {
		return err
	}

	return nil
}
