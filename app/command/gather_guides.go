package command

import (
	"context"
	"time"

	"github.com/Khaym03/REG/domain"
)

type GatherGuidesCommand struct {
	domain.DateRange
}

type GatherGuidesHandler struct {
	guideRepo      domain.GuideRepository
	rubroRepo      domain.RubroRepository
	scraper        domain.GuideCollector
	rubroExtractor domain.RubroExtractor
}

func NewGatherGuidesHandler(
	guideRepo domain.GuideRepository,
	rubroRepo domain.RubroRepository,
	scraper domain.GuideCollector,
	rubroExtractor domain.RubroExtractor,
) *GatherGuidesHandler {
	return &GatherGuidesHandler{
		guideRepo:      guideRepo,
		rubroRepo:      rubroRepo,
		scraper:        scraper,
		rubroExtractor: rubroExtractor,
	}
}

func (h GatherGuidesHandler) Handle(
	ctx context.Context,
	session domain.Session,
	cmd GatherGuidesCommand,
) (err error) {
	dates := domain.MonthlyDateRanges(cmd.From, cmd.To, time.Now())

	for _, d := range dates {
		exist, err := h.guideRepo.Exists(ctx, d)

		if err != nil {
			return err
		}

		if exist {
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
