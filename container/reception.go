package container

import (
	"github.com/Khaym03/REG/adapters"
	"github.com/Khaym03/REG/app"
	"github.com/go-rod/rod"
)

type Container struct {
	Workflow *app.ReceptionWorkflow
}

func BuildContainer(browser *rod.Browser) *Container {
	store := adapters.NewJSONStore("state.json")
	guideRepo := adapters.NewJSONGuideRepository(store)
	receptionRepo := adapters.NewJSONReceptionRepository(store)
	rubroRepo := adapters.NewJSONRubroRepository(store)

	authService := buildAuthService()
	sessionProvider := buildSessionProvider(browser, authService)

	gatherHandler := buildGatherGuidesHandler(guideRepo, rubroRepo)
	inventoryHandler := buildInventoryHandler(rubroRepo)
	receptionHandler := buildReceptionHandler(receptionRepo)

	workflow := app.NewReceptionWorkflow(
		sessionProvider,
		gatherHandler,
		inventoryHandler,
		receptionHandler,
	)

	return &Container{
		Workflow: workflow,
	}
}
