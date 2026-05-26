package reception

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/internal/common/decorator"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/repo"
	"github.com/mustafaturan/bus/v3"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type ReceptionistCommand struct {
	Date                   DateRange
	ReceiveGuidesInTransit bool
}

type ReceptionistHandler decorator.CommandHandler[ReceptionistCommand]

type receptionistHandler struct {
	repo    repo.ReceptionRepository
	scraper ReceptionService

	eventBus *bus.Bus
}

func NewReceptionistHandler(
	repo repo.ReceptionRepository,
	scraper ReceptionService,
	logger *logrus.Entry,
	eventBus *bus.Bus,
) ReceptionistHandler {
	return decorator.ApplyCommandDecorators(
		&receptionistHandler{
			repo:     repo,
			scraper:  scraper,
			eventBus: eventBus,
		},
		logger,
	)
}

func (r *receptionistHandler) Handle(
	ctx context.Context,
	session Session,
	cmd ReceptionistCommand,
) error {

	if err := r.eventBus.Emit(ctx, event.ReceptionTopic, struct{}{}); err != nil {
		log.Error(err)
	}

	dates := domain.MonthlyDateRanges(cmd.Date.From, cmd.Date.To, time.Now())

	for _, d := range dates {
		completed, err := r.repo.IsCompleted(ctx, d)
		if err != nil {
			return err
		}

		isNotCurrentMonth := d.From.Month() != time.Now().Month()
		if completed && isNotCurrentMonth {
			continue
		}

		result, err := r.scraper.Receive(ctx, session, ReceptionOptions{
			Date:                   cmd.Date,
			ReceiveGuidesInTransit: cmd.ReceiveGuidesInTransit,
		})
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
