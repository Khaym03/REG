package container

import (
	"github.com/Khaym03/REG/app"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/Khaym03/REG/internal/stats"
	"github.com/go-rod/rod"
)

type Container struct {
	Workflow *app.ReceptionWorkflow
}

func BuildContainer(browser *rod.Browser) *Container {
	store := repo.NewJSONStore("state.json")
	guideRepo := repo.NewJSONGuideRepository(store)
	receptionRepo := repo.NewJSONReceptionRepository(store)
	rubroRepo := repo.NewJSONRubroRepository(store)

	authService := buildAuthService()
	sessionProvider := buildSessionProvider(browser, authService)

	statsHandler := stats.NewStatsHandler()
	gatherHandler := buildGatherGuidesHandler(guideRepo, rubroRepo)
	inventoryHandler := buildInventoryHandler(rubroRepo)
	receptionHandler := buildReceptionHandler(receptionRepo)

	workflow := app.NewReceptionWorkflow(
		sessionProvider,
		statsHandler,
		gatherHandler,
		inventoryHandler,
		receptionHandler,
	)

	return &Container{
		Workflow: workflow,
	}
}
