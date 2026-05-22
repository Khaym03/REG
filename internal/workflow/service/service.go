package service

import (
	"context"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/repo"

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

	store := repo.NewJSONStore("state.json")
	guideRepo := repo.NewJSONGuideRepository(store)
	receptionRepo := repo.NewJSONReceptionRepository(store)
	rubroRepo := repo.NewJSONRubroRepository(store)

	authService := auth.NewLoginScraper()
	sessionProvider := auth.NewProvider(browser, authService)

	scraperSvc := guide.NewGuidesScraper()
	worker := guide.NewRodRubroWorker(1)

	statsHandler := stats.NewStatsHandler()
	gatherHandler := guide.NewGatherGuidesHandler(
		guideRepo,
		rubroRepo,
		scraperSvc,
		worker,
	)
	inventoryHandler := inventory.NewInventoryHandler(
		rubroRepo,
		inventory.NewInventoryScraper(),
	)
	receptionHandler := reception.NewReceptionistHandler(
		receptionRepo,
		reception.NewReceptionistScraper(),
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
