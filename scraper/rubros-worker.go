package scraper

import (
	"context"
	"log"
	"sync"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/Khaym03/REG/session"
)

type RodRubroWorker struct {
	workers int
}

func NewRodRubroWorker(workers int) *RodRubroWorker {
	return &RodRubroWorker{workers: workers}
}

func (w *RodRubroWorker) Process(ctx context.Context, guides []domain.Guide) ([]domain.Rubro, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	jobs := make(chan domain.Guide)
	rubrosMap := make(map[string]domain.Rubro)

	tempBrowser := session.FromContext(ctx).Browser().MustIncognito()
	defer tempBrowser.Close()

	for i := 0; i < w.workers; i++ {
		wg.Go(func() {
			page := tempBrowser.MustPage()
			defer page.Close()

			// Wrap the page in our Page Object
			guidePage := pages.NewGuideDetailsPage(page)

			for guide := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if err := page.Navigate(guide.URL()); err != nil {
					log.Printf("Navigation error to %s: %v", guide.URL(), err)
					continue
				}

				if err := page.WaitLoad(); err != nil {
					log.Printf("Wait load error on %s: %v", guide.URL(), err)
					continue
				}

				rubros, err := guidePage.ExtractRubros()
				if err != nil {
					log.Printf("Extraction error on %s: %v", guide.URL(), err)
					continue
				}

				mu.Lock()
				for _, r := range rubros {
					rubrosMap[r.Name] = r
				}
				mu.Unlock()
			}
		})
	}

	// Jobs producer
	go func() {
		defer close(jobs)
		for _, guide := range guides {
			jobs <- guide
		}
	}()

	wg.Wait()

	var result []domain.Rubro
	for _, r := range rubrosMap {
		result = append(result, r)
	}

	return result, nil
}
