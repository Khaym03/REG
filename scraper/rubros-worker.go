package scraper

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type RodRubroWorker struct {
	browser *rod.Browser
	workers int
}

func (w *RodRubroWorker) Process(ctx context.Context, ids []string) ([]domain.Rubro, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	jobs := make(chan string)
	rubrosMap := make(map[string]domain.Rubro)

	// workers
	for i := 0; i < w.workers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			page := w.browser.MustPage()
			defer page.Close()

			for id := range jobs {

				select {
				case <-ctx.Done():
					return
				default:
				}

				url := fmt.Sprintf("%s/%s", constants.GuidesURL, id)

				page.MustNavigate(url)
				page.MustWaitLoad()

				rubros := extractRubrosFromGuide(page)

				mu.Lock()
				for _, r := range rubros {
					rubrosMap[r.Name] = r
				}
				mu.Unlock()
			}
		}(i)
	}

	// feed jobs
	go func() {
		defer close(jobs)
		for _, id := range ids {
			jobs <- id
		}
	}()

	wg.Wait()

	// map → slice
	var result []domain.Rubro
	for _, r := range rubrosMap {
		result = append(result, r)
	}

	return result, nil
}

func extractRubrosFromGuide(page *rod.Page) []domain.Rubro {
	table := page.MustElementR("h4", "RUBROS").MustParent().MustParent().MustNext()

	rows := table.MustElements("tbody tr")

	var listaRubros []domain.Rubro

	for _, row := range rows {
		cols := row.MustElements("td")

		// Ensure the row has the expected columns
		if len(cols) >= 5 {
			item := domain.Rubro{
				Name: strings.TrimSpace(cols[0].MustText()),
				// Cantidad:     strings.TrimSpace(cols[1].MustText()),
				// PrecioVenta:  strings.TrimSpace(cols[2].MustText()),
				// Presentacion: strings.TrimSpace(cols[3].MustText()),
				// Marca:        strings.TrimSpace(cols[4].MustText()),
			}
			listaRubros = append(listaRubros, item)
		}
	}

	for _, r := range listaRubros {
		fmt.Printf("Producto: %s\n", r.Name)
	}

	return listaRubros
}
