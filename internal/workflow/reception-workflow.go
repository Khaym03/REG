package workflow

import (
	"context"

	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/auth"
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
	app           *app.Application
	sessinManager auth.SessionManager
}

func NewReceptionWorkflow(
	ctx context.Context,
	eventBus *bus.Bus,
	sm auth.SessionManager,

) (*ReceptionWorkflow, error) {
	application, err := service.NewApplication(
		ctx,
		eventBus,
		sm,
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &ReceptionWorkflow{
			app:           application,
			sessinManager: sm,
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

	err := w.sessinManager.Reconfigure(ctx, input.BrowserConf)
	if err != nil {
		return err
	}

	err = w.sessinManager.Login(ctx, input.User)
	if err != nil {
		return err
	}
	defer func() {
		err = w.sessinManager.Logout(ctx)
		if err != nil {
			log.Error(err)
		}
	}()

	session, err := w.sessinManager.Session(ctx)
	if err != nil {
		return err
	}

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
