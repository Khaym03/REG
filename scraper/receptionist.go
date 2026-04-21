package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/session"
	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
)

var _ domain.Receptionist = (*ReceptionistScraper)(nil)

type ReceptionistScraper struct {
}

func NewReceptionistScraper() *ReceptionistScraper {
	return &ReceptionistScraper{}
}

func (r *ReceptionistScraper) Receive(ctx context.Context, date utils.DateRange) (domain.ReceptionResult, error) {
	page := session.FromContext(ctx).MainPage()

	var err error
	result := domain.ReceptionResult{}

	for {
		// Navigate to the receiver page at the start of every iteration
		// to ensure we have a fresh, non-stale DOM context.
		if err = navigate(page, c.ReceptionURL); err != nil {
			return result, err
		}

		page.MustWaitIdle()

		err := applyFiltersToGuideReceiver(page, date)
		if err != nil {
			return result, err
		}

		// Process exactly one guide. If it returns true, we loop back
		// to re-navigate and re-filter for the remaining guides.
		processed, err := r.processSingleExpiredGuide(page)
		if err != nil {
			log.Printf("Error during guide processing: %v", err)
			return result, err
		}

		if !processed {
			fmt.Println("No more expired guides found for this range.")
			result.Completed = true
			return result, nil
		}

		fmt.Println("Guide successfully processed. Restarting sequence for the next one...")
		result.Processed++
		// Small buffer to allow the server to sync state changes
		time.Sleep(2 * time.Second)
	}
}

func (r *ReceptionistScraper) processSingleExpiredGuide(page *rod.Page) (bool, error) {
	page.MustWaitDOMStable()

	rows, err := page.Elements("#tabla-component tbody tr")
	if err != nil {
		// Meaning no row to process / done
		return false, err
	}

	for _, row := range rows {
		// Check for the "VENCIDA" badge
		expiredBadge, err := row.Element(".badge-danger")
		if err != nil || expiredBadge == nil {
			continue
		}

		text, _ := expiredBadge.Text()
		if text != "VENCIDA" {
			continue
		}

		btn, err := row.Element(".recepcionar")
		if err != nil {
			continue
		}

		fmt.Println("Found expired guide. Triggering reception modal...")
		btn.MustScrollIntoView().MustClick()

		// This function now handles the navigation wait
		r.clickRecepcionarButtonInModal(page)

		// Return true to indicate we handled one guide and need to re-filter
		return true, nil
	}

	return false, nil
}

func (r *ReceptionistScraper) clickRecepcionarButtonInModal(page *rod.Page) {
	modalSelector := "#recepcionarGuia"
	buttonSelector := "#recepcionarGuia button.btn-success"

	// Wait for the modal to be fully visible
	page.MustElement(modalSelector).MustWaitVisible()

	// Initialize a navigation waiter BEFORE the click.
	// This forces Rod to wait until the browser finishes the redirect/reload.
	wait := page.MustWaitNavigation()

	fmt.Println("Confirming reception in modal...")
	page.MustElement(buttonSelector).MustClick()

	// Block execution here until navigation is complete
	wait()
}
