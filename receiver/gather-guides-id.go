package receiver

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
)

type monthlyGuideIDs map[string][]string

type GuidesIDGatherer struct {
	monthlyGuideIDs `json:"guides-id"`
	ExtractedRubros map[string]bool `json:"extracted-rubros"`
}

func (g *GuidesIDGatherer) ApplyFiltersToGuideReceiver(page *rod.Page) error {
	page.MustElementX(filterAccordionSelector).MustClick()
	page.MustWaitDOMStable()
	utils.SelectOption(page, selectStatusSelector, selectStatusOption)
	utils.SelectOption(page, selectReceptionStatus, selectReceptionOption)

	idsChan := make(chan string, 100)

	// Run the worker pool in the background so we can stream IDs to it
	var finalRubros map[string]bool
	done := make(chan struct{})
	go func() {
		finalRubros = g.processGuidesWorker(idsChan, page)
		close(done)
	}()

	datesets := utils.MonthlyDateRangesCurrentToLastYear()
	for _, date := range datesets {
		rangeID := date.From.Format("2006-01")

		if _, exist := g.monthlyGuideIDs[rangeID]; exist {
			continue
		}

		page.MustElementX(inputDateFromSelector).MustInputTime(date.From)
		page.MustElementX(inputDateToSelector).MustInputTime(date.To)
		page.MustElementX(filterButtonSelector).MustClick()
		page.MustWaitDOMStable()

		ids := CollectGuidesIDs(page)
		g.monthlyGuideIDs[rangeID] = ids

		// Stream IDs to workers immediately
		for _, id := range ids {
			idsChan <- id
		}
	}

	close(idsChan)
	<-done
	g.ExtractedRubros = finalRubros

	log.Println("Final Rubros: ", finalRubros)

	return nil
}

func (g *GuidesIDGatherer) processGuidesWorker(idsChan <-chan string, page *rod.Page) map[string]bool {
	var wg sync.WaitGroup
	var mu sync.Mutex
	uniqueRubros := make(map[string]bool)
	numWorkers := 1

	for i := range numWorkers {
		wg.Go(func() {
			tempBrowser := page.Browser()
			tempPage := tempBrowser.MustIncognito().MustPage()

			log.Printf("[Worker %d] Started", i)
			for id := range idsChan {
				time.Sleep(time.Second * 5)
				url := fmt.Sprintf("https://sica.sunagro.gob.ve/guias/%s", id)
				tempPage.MustNavigate(url)
				tempPage.MustWaitLoad()
				rubros := ExtractRubrosFromGuide(tempPage)

				mu.Lock()
				for _, r := range rubros {
					uniqueRubros[r.Nombre] = true
				}
				mu.Unlock()
			}

			tempBrowser.Close()
			tempPage.Close()

			log.Printf("[Worker %d] Finished", i)
		})
	}

	wg.Wait()
	return uniqueRubros
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
