package scraper

import (
	"context"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/Khaym03/REG/session"
	"github.com/Khaym03/REG/utils"
)

var _ domain.GuideScraper = (*GuidesScraper)(nil)

type GuidesScraper struct {
}

func NewGuidesScraper() *GuidesScraper {
	return &GuidesScraper{}
}

func (g GuidesScraper) CollectGuides(ctx context.Context, date utils.DateRange) ([]domain.Guide, error) {
	ReceptionPage := pages.NewReceptionPage(session.FromContext(ctx).MainPage())

	ReceptionPage.Open()
	ReceptionPage.ApplyFilters(date)

	rows, _ := ReceptionPage.Rows()
	var guides []domain.Guide

	for _, row := range rows {
		guides = append(guides, domain.Guide{ID: row.ID()})
	}

	return guides, nil
}
