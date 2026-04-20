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
	repo := adapters.NewJSONGuideRepository("state.json")

	authService := buildAuthService()
	sessionProvider := buildSessionProvider(browser, authService)

	gatherHandler := buildGatherGuidesHandler(repo)
	inventoryHandler := buildInventoryHandler(repo)
	receptionHandler := buildReceptionHandler(repo)

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
