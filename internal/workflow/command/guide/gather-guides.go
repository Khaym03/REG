package guide

import (
	"context"
	"time"

	"github.com/Khaym03/REG/internal/common/decorator"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/mustafaturan/bus/v3"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type GatherGuidesCommand struct {
	DateRange
}

type GatherGuidesHandler decorator.CommandHandler[GatherGuidesCommand]

type gatherGuidesHandler struct {
	guideRepo      repo.GuideRepository
	rubroRepo      repo.RubroRepository
	scraper        GuideCollector
	rubroExtractor RubroExtractor

	eventBus *bus.Bus
}

func NewGatherGuidesHandler(
	guideRepo repo.GuideRepository,
	rubroRepo repo.RubroRepository,
	scraper GuideCollector,
	rubroExtractor RubroExtractor,
	logger *logrus.Entry,
	eventBus *bus.Bus,
) GatherGuidesHandler {

	return decorator.ApplyCommandDecorators(&gatherGuidesHandler{
		guideRepo:      guideRepo,
		rubroRepo:      rubroRepo,
		scraper:        scraper,
		rubroExtractor: rubroExtractor,
		eventBus:       eventBus,
	},
		logger,
	)
}

func (h gatherGuidesHandler) Handle(
	ctx context.Context,
	session Session,
	cmd GatherGuidesCommand,
) (err error) {

	if err := h.eventBus.Emit(ctx, event.GuidesGatherTopic, struct{}{}); err != nil {
		log.Error(err)
	}

	dates := domain.MonthlyDateRanges(cmd.From, cmd.To, time.Now())

	for _, d := range dates {
		exist, err := h.guideRepo.Exists(ctx, d)
		if err != nil {
			return err
		}

		isNotCurrentMonth := d.From.Month() != time.Now().Month()

		if exist && isNotCurrentMonth {
			continue
		}

		guides, err := h.scraper.Collect(ctx, session, d)
		if err != nil {
			return err
		}

		if err = h.guideRepo.Save(ctx, d, guides); err != nil {
			return err
		}

		rubros, err := h.rubroExtractor.FromGuides(ctx, session, guides)
		if err != nil {
			return err
		}

		if err = h.rubroRepo.Save(ctx, rubros); err != nil {
			return err
		}
	}

	return
}
