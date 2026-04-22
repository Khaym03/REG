package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/Khaym03/REG/session"
	"github.com/Khaym03/REG/utils"
)

var _ domain.Receptionist = (*ReceptionistScraper)(nil)

type ReceptionistScraper struct {
}

func NewReceptionistScraper() *ReceptionistScraper {
	return &ReceptionistScraper{}
}

func (r *ReceptionistScraper) Receive(ctx context.Context, date utils.DateRange) (domain.ReceptionResult, error) {
	page := pages.NewReceptionPage(session.FromContext(ctx).MainPage())

	var err error
	result := domain.ReceptionResult{}

	for {
		// Navigate to the receiver page at the start of every iteration
		// to ensure we have a fresh, non-stale DOM context.
		if err = page.Open(); err != nil {
			return result, err
		}

		err := page.ApplyFilters(date)
		if err != nil {
			return result, err
		}

		rows, err := page.Rows()
		if err != nil {
			return result, err
		}
		var processed bool

		// Process exactly one guide. If it returns true, we loop back
		// to re-navigate and re-filter for the remaining guides.
		for _, row := range rows {
			if row.IsExpired() {
				if err := row.TriggerReception(); err != nil {
					return result, fmt.Errorf("failed to trigger reception: %w", err)
				}

				if err := page.ConfirmReception(); err != nil {
					return result, fmt.Errorf("modal confirmation failed: %w", err)
				}

				result.Processed++
				processed = true
				break
			}
		}
		if !processed {
			result.Completed = true
			fmt.Println("No more expired guides found for this range.")
			break
		}

		fmt.Println("Guide successfully processed. Restarting sequence for the next one...")
		result.Processed++
		// Small buffer to allow the server to sync state changes
		time.Sleep(2 * time.Second)
	}

	return result, nil
}
