package scraper

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/go-rod/rod"
)

var _ domain.GuideCollector = (*GuidesScraper)(nil)

type GuidesScraper struct {
}

func NewGuidesScraper() *GuidesScraper {
	return &GuidesScraper{}
}

func (g GuidesScraper) Collect(
	ctx context.Context,
	session domain.Session,
	date domain.DateRange,
) (guides []domain.Guide, err error) {

	err = session.Do(ctx, func(p *rod.Page) error {
		ReceptionPage := pages.NewReceptionPage(p)

		ReceptionPage.Open()
		ReceptionPage.ApplyFilters(date)

		rows, _ := ReceptionPage.Rows()

		for _, row := range rows {
			id, err := row.ID()
			if err != nil {
				return err
			}

			guide, err := domain.NewGuide(id)
			if err != nil {
				return err
			}

			guides = append(guides, guide)

		}

		return nil
	})

	return
}
