package receiver

import (
	"log"

	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
)

type monthlyGuideIDs map[string][]string

type GuidesIDGatherer struct {
	monthlyGuideIDs `json:"guides-id"`
}

func (g *GuidesIDGatherer) ApplyFiltersToGuideReceiver(page *rod.Page) error {
	datesets := utils.MonthlyDateRangesCurrentToLastYear()

	page.MustElementX(filterAccordionSelector).MustClick()
	page.MustWaitDOMStable()

	utils.SelectOption(page, selectStatusSelector, selectStatusOption)
	utils.SelectOption(page, selectReceptionStatus, selectReceptionOption)

	for _, date := range datesets {
		rangeID := date.From.Format("2006-01")

		if _, exist := g.monthlyGuideIDs[rangeID]; exist {
			log.Printf("Skipping month %s: already processed", rangeID)
			continue
		}

		page.MustElementX(inputDateFromSelector).MustInputTime(date.From)
		page.MustElementX(inputDateToSelector).MustInputTime(date.To)
		page.MustElementX(filterButtonSelector).MustClick()
		page.MustWaitDOMStable()

		g.monthlyGuideIDs[rangeID] = CollectGuidesIDs(page)
	}

	return nil
}

func CollectGuidesIDs(page *rod.Page) []string {
	var ids []string

	rows, err := page.ElementsX(tableRowSelector)
	if err != nil || len(rows) == 0 {
		log.Println("No rows found in the table. Continuing...")
		return ids
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
		ids = append(ids, *idValue)
	}

	return ids
}
