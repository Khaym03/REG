package scraper

import (
	"context"
	"log"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/go-rod/rod"
)

var _ domain.InventoryService = (*InventoryScraper)(nil)

type InventoryScraper struct {
}

func NewInventoryScraper() *InventoryScraper {
	return &InventoryScraper{}
}

func (i *InventoryScraper) Insert(
	ctx context.Context,
	session domain.Session,
	newItem domain.Rubro,
) (err error) {

	s, err := session.NewIsolated(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err = s.Close(); err != nil {
			log.Println(err)
		}
	}()

	return s.Do(ctx, func(p *rod.Page) error {
		inventoryPage := pages.NewInventoryPage(p)

		if err := inventoryPage.Open(); err != nil {
			return err
		}

		if err := inventoryPage.SelectItem(newItem); err != nil {
			return err
		}

		if err := inventoryPage.Submit(); err != nil {
			return err
		}

		log.Println("New item added to UI:", newItem.Name)

		return nil
	})
}

func (i InventoryScraper) Snapshot(
	ctx context.Context,
	session domain.Session,
) ([]domain.Rubro, error) {

	var rubros []domain.Rubro

	return rubros, session.Do(ctx, func(p *rod.Page) error {
		inventoryPage := pages.NewInventoryPage(p)

		if err := inventoryPage.Open(); err != nil {
			return err
		}

		r, err := inventoryPage.ExtractExistingRubros()
		rubros = r
		return err

	})
}
