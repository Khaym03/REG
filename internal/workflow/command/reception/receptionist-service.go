package reception

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/go-rod/rod"
)

var _ ReceptionService = (*ReceptionistScraper)(nil)

type ReceptionistScraper struct{}

func NewReceptionistScraper() *ReceptionistScraper {
	return &ReceptionistScraper{}
}

func (r *ReceptionistScraper) Receive(
	ctx context.Context,
	session Session,
	opt ReceptionOptions,
) (ReceptionResult, error) {
	result := ReceptionResult{}

	receive := func(p *rod.Page) error {
		var page Page = NewReceptionPage(p)

		for {

			processed, err := ProcessNextExpiredGuide(page, opt, &result)
			if err != nil {
				return err
			}

			if !processed {
				result.Completed = true
				log.Info("No more expired guides found for this range.", opt.Date)
				return nil
			}

			log.Info("Guide processed. Continuing...")
			time.Sleep(2 * time.Second)
		}
	}

	receive = browser.WithRetry(ctx, 3, time.Second*10)(receive)

	err := session.Do(ctx, receive)

	return result, err
}
func ProcessNextExpiredGuide(
	page Page,
	opt ReceptionOptions,
	result *ReceptionResult,
) (bool, error) {
	// Navigate to the receiver page at the start of every iteration
	// to ensure we have a fresh, non-stale DOM context.
	if err := page.Open(); err != nil {
		return false, err
	}

	if err := page.ApplyFilters(opt.Date); err != nil {
		return false, err
	}

	rows, err := page.Rows()
	if err != nil {
		return false, err
	}

	for _, row := range rows {
		if !row.IsExpired() && !opt.ReceiveGuidesInTransit {
			continue
		}

		if err := row.TriggerReception(); err != nil {
			return false, fmt.Errorf("trigger reception: %w", err)
		}

		if err := page.ConfirmReception(); err != nil {
			return false, fmt.Errorf("confirm reception: %w", err)
		}

		result.Processed++
		return true, nil
	}

	return false, nil
}
