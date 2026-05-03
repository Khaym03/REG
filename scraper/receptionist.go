package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/go-rod/rod"
)

var _ domain.ReceptionService = (*ReceptionistScraper)(nil)

type ReceptionistScraper struct {
}

func NewReceptionistScraper() *ReceptionistScraper {
	return &ReceptionistScraper{}
}

func (r *ReceptionistScraper) Receive(
	ctx context.Context,
	session domain.Session,
	date domain.DateRange,
) (domain.ReceptionResult, error) {

	result := domain.ReceptionResult{}

	receive := func(p *rod.Page) error {
		page := pages.NewReceptionPage(p)

		for {

			processed, err := r.processNextExpiredGuide(page, date, &result)
			if err != nil {
				return err
			}

			if !processed {
				result.Completed = true
				fmt.Println("No more expired guides found for this range.", date)
				return nil
			}

			fmt.Println("Guide processed. Continuing...")
			time.Sleep(2 * time.Second)
		}
	}

	receive = WithRetry(3, time.Second*10)(receive)

	err := session.Do(ctx, receive)

	return result, err
}

func (r *ReceptionistScraper) processNextExpiredGuide(
	page *pages.ReceptionPage,
	date domain.DateRange,
	result *domain.ReceptionResult,
) (bool, error) {

	// Navigate to the receiver page at the start of every iteration
	// to ensure we have a fresh, non-stale DOM context.
	if err := page.Open(); err != nil {
		return false, err
	}

	if err := page.ApplyFilters(date); err != nil {
		return false, err
	}

	rows, err := page.Rows()
	if err != nil {
		return false, err
	}

	for _, row := range rows {
		if !row.IsExpired() {
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
