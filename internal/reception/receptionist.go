package reception

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/internal/domain"
)

type ReceptionistCommand struct {
	domain.DateRange
}

type ReceptionistHandler struct {
	repo    domain.ReceptionRepository
	scraper domain.ReceptionService
}

func NewReceptionistHandler(
	repo domain.ReceptionRepository,
	scraper domain.ReceptionService,
) *ReceptionistHandler {
	return &ReceptionistHandler{repo: repo, scraper: scraper}
}

func (r *ReceptionistHandler) Handle(
	ctx context.Context,
	session domain.Session,
	cmd ReceptionistCommand,
) error {
	dates := domain.MonthlyDateRanges(cmd.From, cmd.To, time.Now())

	for _, d := range dates {
		completed, err := r.repo.IsCompleted(ctx, d)
		if err != nil {
			return err
		}

		if completed {
			continue
		}

		result, err := r.scraper.Receive(ctx, session, d)
		r.repo.SaveProgress(ctx, d, result)

		if err != nil {
			return fmt.Errorf(
				"receptionist failed for range %s - %s: %w",
				d.From, d.To, err,
			)
		}

		if result.Completed {
			r.repo.MarkCompleted(ctx, d)
		}
	}

	return nil
}
