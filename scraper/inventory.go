package scraper

import (
	"context"
	"log"

	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper/pages"
	"github.com/Khaym03/REG/session"
)

var _ domain.InventoryScraper = (*InventoryScraper)(nil)

type InventoryScraper struct {
}

func NewInventoryScraper() *InventoryScraper {
	return &InventoryScraper{}
}

func (i *InventoryScraper) Insert(ctx context.Context, newItem domain.Rubro) (err error) {
	page := session.FromContext(ctx).MainPage().Context(ctx)
	defer page.Close()

	inventoryPage := pages.NewInventoryPage(page)

	if err = inventoryPage.Open(); err != nil {
		return
	}

	if err = inventoryPage.SelectItem(newItem); err != nil {
		return
	}

	if err = inventoryPage.Submit(); err != nil {
		return
	}

	log.Println("New item added to UI:", newItem.Name)

	return nil

}

func (i InventoryScraper) RubrosSnapshot(ctx context.Context) ([]domain.Rubro, error) {
	page := session.FromContext(ctx).MainPage().Context(ctx)

	inventoryPage := pages.NewInventoryPage(page)

	if err := inventoryPage.Open(); err != nil {
		return nil, err
	}

	return inventoryPage.ExtractExistingRubros()
}
