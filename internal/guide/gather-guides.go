package guide

import (
	"context"
	"time"

	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/repo"
)

type GatherGuidesCommand struct {
	DateRange
}

type GatherGuidesHandler struct {
	guideRepo      repo.GuideRepository
	rubroRepo      repo.RubroRepository
	scraper        GuideCollector
	rubroExtractor RubroExtractor
}

func NewGatherGuidesHandler(
	guideRepo repo.GuideRepository,
	rubroRepo repo.RubroRepository,
	scraper GuideCollector,
	rubroExtractor RubroExtractor,
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
	session Session,
	cmd GatherGuidesCommand,
) (err error) {
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
