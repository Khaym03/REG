package service

import (
	"context"
	"path/filepath"

	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/Khaym03/REG/utils"
	"github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/workflow/app"
	"github.com/Khaym03/REG/internal/workflow/command/guide"
	"github.com/Khaym03/REG/internal/workflow/command/inventory"
	"github.com/Khaym03/REG/internal/workflow/command/reception"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
)

type CleanUpFunc func()

const stateFilename = "state.json"

func NewApplication(
	ctx context.Context,
	eventBus event.Bus,
	sm mediator.SessionMediator,
) (*app.Application, error) {

	logger := logrus.NewEntry(logrus.StandardLogger())

	stateFilepath := filepath.Join(utils.BaseDir(), stateFilename)

	var persistance repo.Persistence[repo.RepositoryData] = repo.NewJSONPersistence(
		stateFilepath,
		func() repo.RepositoryData {
			return repo.RepositoryData{
				Months:         make(map[string][]domain.Guide),
				Rubros:         make(map[string]domain.Rubro),
				ReceptionState: make(map[string]domain.ReceptionResult),
			}
		})

	store := repo.NewJSONStore(persistance)
	guideRepo := repo.NewJSONGuideRepository(store)
	receptionRepo := repo.NewJSONReceptionRepository(store)
	rubroRepo := repo.NewJSONRubroRepository(store)

	scraperSvc := guide.NewGuidesScraper()
	worker := guide.NewRodRubroWorker(1, sm)

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
		inventory.NewInventoryScraper(sm),
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
			SessionMediator: sm,
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
