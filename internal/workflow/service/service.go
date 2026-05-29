package service

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/mustafaturan/bus/v3"
	"github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/workflow/app"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
)

type CleanUpFunc func()

func NewApplication(
	ctx context.Context,
	eventBus *bus.Bus,
) (*app.Application, error) {

	logger := logrus.NewEntry(logrus.StandardLogger())

	store := repo.NewJSONStore("state.json")
	guideRepo := repo.NewJSONGuideRepository(store)
	receptionRepo := repo.NewJSONReceptionRepository(store)
	rubroRepo := repo.NewJSONRubroRepository(store)

	authService := auth.NewLoginScraper()
	sessionProvider := auth.NewProvider(authService, eventBus)

	scraperSvc := guide.NewGuidesScraper()
	worker := guide.NewRodRubroWorker(1)

	statsHandler := stats.NewStatsHandler(logger, eventBus)
	gatherHandler := guide.NewGatherGuidesHandler(
		guideRepo,
		rubroRepo,
		scraperSvc,
		worker,
		logger,
		eventBus,
	)
	inventoryHandler := inventory.NewInventoryHandler(
		rubroRepo,
		inventory.NewInventoryScraper(),
		logger,
		eventBus,
	)
	receptionHandler := reception.NewReceptionistHandler(
		receptionRepo,
		reception.NewReceptionistScraper(),
		logger,
		eventBus,
	)

	return &app.Application{
			EventBus:        eventBus,
			SessionProvider: sessionProvider,
			Commands: app.Commands{
				GatherGuides:  gatherHandler,
				SyncInventory: inventoryHandler,
				Receptionist:  receptionHandler,
			},
			Queries: app.Queries{
				Stats: statsHandler,
			},
		},
		nil
}
