package command

import (
	"context"
	"fmt"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/utils"
)

type ReceptionistCommand struct {
	utils.DateRange
}

type ReceptionistHandler struct {
	repo    domain.GuideRepository
	scraper domain.Receptionist
}

func NewReceptionistHandler(
	repo domain.GuideRepository,
	scraper domain.Receptionist,
) *ReceptionistHandler {
	return &ReceptionistHandler{repo: repo, scraper: scraper}
}

func (r *ReceptionistHandler) Handle(ctx context.Context, cmd ReceptionistCommand) error {
	dates := utils.MonthlyDateRanges(cmd.From, cmd.To)

	for _, d := range dates {
		if r.repo.IsReceptionCompleted(d) {
			continue
		}

		result, err := r.scraper.Receive(ctx, d)
		r.repo.SaveReceptionProgress(d, result)

		if err != nil {
			return fmt.Errorf(
				"receptionist failed for range %s - %s: %w",
				d.From, d.To, err,
			)
		}

		if result.Completed {
			r.repo.MarkReceptionCompleted(d)
		}
	}

	return nil
}
