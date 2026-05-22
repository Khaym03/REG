package guide

import (
	"context"
	"time"

	"github.com/Khaym03/REG/internal/app/command/reception"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/domain"

	"github.com/go-rod/rod"
)

var _ GuideCollector = (*GuidesScraper)(nil)

type GuidesScraper struct{}

func NewGuidesScraper() *GuidesScraper {
	return &GuidesScraper{}
}

func (g GuidesScraper) Collect(
	ctx context.Context,
	session Session,
	date DateRange,
) (guides []Guide, err error) {
	collect := func(p *rod.Page) error {
		ReceptionPage := reception.NewReceptionPage(p)

		ReceptionPage.Open()
		ReceptionPage.ApplyFilters(date)

		rows, err := ReceptionPage.Rows()
		if err != nil {
			return err
		}

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
	}

	collect = browser.WithRetry(ctx, 3, time.Second*10)(collect)

	err = session.Do(ctx, collect)

	return
}
