package scraper

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/Khaym03/REG/session"
)

var _ domain.GuideCollector = (*GuidesScraper)(nil)

type GuidesScraper struct {
}

func NewGuidesScraper() *GuidesScraper {
	return &GuidesScraper{}
}

func (g GuidesScraper) Collect(ctx context.Context, date domain.DateRange) ([]domain.Guide, error) {
	ReceptionPage := pages.NewReceptionPage(session.FromContext(ctx).MainPage())

	ReceptionPage.Open()
	ReceptionPage.ApplyFilters(date)

	rows, _ := ReceptionPage.Rows()
	var guides []domain.Guide

	for _, row := range rows {
		id, err := row.ID()
		if err != nil {
			return guides, err
		}

		guide, err := domain.NewGuide(id)
		if err != nil {
			return guides, err
		}

		guides = append(guides, guide)
	}

	return guides, nil
}
