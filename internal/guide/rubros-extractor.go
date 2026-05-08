package guide

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/go-rod/rod"
)

var _ domain.RubroExtractor = (*RodRubroWorker)(nil)

type RodRubroWorker struct {
	workers int
}

func NewRodRubroWorker(workers int) *RodRubroWorker {
	return &RodRubroWorker{workers: workers}
}

func (w *RodRubroWorker) FromGuides(
	ctx context.Context,
	session domain.Session,
	guides []domain.Guide,
) ([]domain.Rubro, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	jobs := make(chan domain.Guide)
	rubrosMap := make(map[string]domain.Rubro)

	tempSession, err := session.NewIsolated(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tempSession.Close(); err != nil {
			log.Println(err)
		}
	}()

	extractTheRubros := func(p *rod.Page) error {
		guidePage := NewGuideDetailsPage(p)

		for guide := range jobs {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			if err := p.Navigate(guide.URL); err != nil {
				log.Errorf("Navigation error to %s: %v", guide.URL, err)
				continue
			}

			if err := p.WaitLoad(); err != nil {
				log.Errorf("Wait load error on %s: %v", guide.URL, err)
				continue
			}

			rubros, err := guidePage.ExtractRubros()
			if err != nil {
				log.Errorf("Extraction error on %s: %v", guide.URL, err)
				continue
			}

			mu.Lock()
			for _, r := range rubros {
				rubrosMap[r.Name] = r
			}
			mu.Unlock()
		}

		return nil
	}

	extractTheRubros = browser.WithRetry(3, time.Second*10)(extractTheRubros)
	for i := 0; i < w.workers; i++ {
		wg.Go(func() {
			tempSession.Do(ctx, extractTheRubros)
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
