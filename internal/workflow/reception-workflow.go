package workflow

import (
	"context"

	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/workflow/app"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	"github.com/Khaym03/REG/internal/workflow/service"

	"github.com/Khaym03/REG/internal/domain"
)

type WorkFlowInput struct {
	User                   auth.User        `json:"user"`
	Date                   domain.DateRange `json:"date"`
	ReceiveGuidesInTransit bool             `json:"receive_guides_in_transit,omitempty"`
	BrowserConf            config.BrowserConfig
}

type ReceptionWorkflow struct {
	app *app.Application
}

func NewReceptionWorkflow(
	ctx context.Context,
	eventBus *bus.Bus,

) (*ReceptionWorkflow, error) {
	application, err := service.NewApplication(
		ctx,
		eventBus,
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &ReceptionWorkflow{
			app: application,
		},
		nil
}

func (w *ReceptionWorkflow) Run(ctx context.Context, input WorkFlowInput) error {
	if err := w.app.EventBus.Emit(
		ctx,
		event.WorkflowStarted,
		struct{}{},
	); err != nil {
		log.Error(err)
	}

	defer func() {
		if err := w.app.EventBus.Emit(
			ctx,
			event.WorkflowFinished,
			struct{}{},
		); err != nil {
			log.Error(err)
		}
	}()

	if err := w.app.EventBus.Emit(
		ctx,
		event.BuildingBrowser,
		struct{}{},
	); err != nil {
		log.Error(err)
	}

	browser, err := browser.BuildBrowser(ctx, input.BrowserConf)
	if err != nil {
		return err
	}

	session, err := w.app.SessionProvider.Start(ctx, input.User, browser)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := session.Close(ctx); cerr != nil {
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
