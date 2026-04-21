package scraper

import (
	"context"
	"fmt"
	"log"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/session"
	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
)

var _ domain.GuideScraper = (*GuidesScraper)(nil)

type GuidesScraper struct {
}

func NewGuidesScraper() *GuidesScraper {
	return &GuidesScraper{}
}

func (g GuidesScraper) CollectGuides(ctx context.Context, date utils.DateRange) ([]domain.Guide, error) {
	var err error

	page := session.FromContext(ctx).MainPage().Context(ctx)

	if err = navigate(page, c.ReceptionURL); err != nil {
		return nil, fmt.Errorf("nav to reception failed: %w", err)
	}

	err = applyFiltersToGuideReceiver(page, date)
	if err != nil {
		return nil, err
	}

	return g.collectIDs(page)
}

func (g GuidesScraper) collectIDs(page *rod.Page) ([]domain.Guide, error) {
	var guides []domain.Guide

	rows, err := page.ElementsX(tableRowSelector)
	if err != nil || len(rows) == 0 {
		log.Println("No rows found in the table. Continuing...")
		return guides, err
	}
	log.Printf("Found %d rows. Processing...", len(rows))

	for i, row := range rows {
		// Check if the row actually has the data-id_ column.
		// We use ElementX here too to avoid panicking on a weirdly formatted row.
		column, err := row.ElementX(dataIDColumnSelector)
		if err != nil {
			log.Printf("Row %d does not contain a data-id_ attribute. Skipping.", i)
			continue
		}

		// Extract the attribute
		idValue, err := column.Attribute("data-id_")
		if err != nil || idValue == nil {
			continue
		}

		log.Printf("Found ID: %s\n", *idValue)
		guides = append(guides, domain.Guide{ID: *idValue})
	}

	return guides, nil
}

const (
	filterAccordionSelector = `//*[@id="accordion-filtros"]/div/div[1]/a`
	selectStatusSelector    = `//*[@id="select2-estatus_-container"]/..`
	selectReceptionStatus   = `//*[@id="select2-recepcion-container"]/..`
	selectStatusOption      = `//li[contains(text(), "APROBADA")]`
	selectReceptionOption   = `//li[contains(@id, "SIN_RECEPCIONAR")]`
	inputDateFromSelector   = `//*[@id="desde"]`
	inputDateToSelector     = `//*[@id="hasta"]`
	filterButtonSelector    = `//*[@id="collapse-filtro"]/div/form/div[3]/button`
	tableRowSelector        = `//table[@id="tabla-component"]/tbody/tr`
	dataIDColumnSelector    = `./td[@data-id_]`
)
