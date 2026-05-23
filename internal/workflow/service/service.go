package service

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/repo"
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
	conf config.BrowserConfig,
) (*app.Application, CleanUpFunc, error) {

	browser, err := browser.BuildBrowser(ctx, conf)
	if err != nil {
		return nil, nil, err
	}

	logger := logrus.NewEntry(logrus.StandardLogger())

	store := repo.NewJSONStore("state.json")
	guideRepo := repo.NewJSONGuideRepository(store)
	receptionRepo := repo.NewJSONReceptionRepository(store)
	rubroRepo := repo.NewJSONRubroRepository(store)

	authService := auth.NewLoginScraper()
	sessionProvider := auth.NewProvider(browser, authService)

	scraperSvc := guide.NewGuidesScraper()
	worker := guide.NewRodRubroWorker(1)

	statsHandler := stats.NewStatsHandler(logger)
	gatherHandler := guide.NewGatherGuidesHandler(
		guideRepo,
		rubroRepo,
		scraperSvc,
		worker,
		logger,
	)
	inventoryHandler := inventory.NewInventoryHandler(
		rubroRepo,
		inventory.NewInventoryScraper(),
		logger,
	)
	receptionHandler := reception.NewReceptionistHandler(
		receptionRepo,
		reception.NewReceptionistScraper(),
		logger,
	)

	return &app.Application{
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
		func() {
			browser.MustClose()
		},
		nil
}
