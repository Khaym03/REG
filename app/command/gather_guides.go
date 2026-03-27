package command

import (
	"context"
	"time"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/utils"
)

type GatherGuidesCommand struct {
	From time.Time
	To   time.Time
}

type GatherGuidesHandler struct {
	repo    domain.GuideRepository
	scraper domain.GuideScraper
	workers domain.RubroWorker
}

func (h GatherGuidesHandler) Handle(ctx context.Context, cmd GatherGuidesCommand) error {
	dates := utils.MonthlyDateRanges(cmd.From, cmd.To)

	for _, d := range dates {

		if h.repo.Exists(d) {
			continue
		}

		ids, err := h.scraper.CollectIDs(ctx, d)
		if err != nil {
			return err
		}

		h.repo.SaveIDs(d, ids)

		rubros, err := h.workers.Process(ctx, ids)
		if err != nil {
			return err
		}

		h.repo.SaveRubros(rubros)
	}

	return nil
}
