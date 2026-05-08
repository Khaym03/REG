package inventory

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/browser"

	"github.com/go-rod/rod"
)

var _ InventoryService = (*InventoryScraper)(nil)

type InventoryScraper struct{}

func NewInventoryScraper() *InventoryScraper {
	return &InventoryScraper{}
}

func (i *InventoryScraper) Insert(
	ctx context.Context,
	session Session,
	newItem Rubro,
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

	insert := func(p *rod.Page) error {
		inventoryPage := NewInventoryPage(p)

		if err := inventoryPage.Open(); err != nil {
			return err
		}

		if err := inventoryPage.SelectItem(newItem); err != nil {
			return err
		}

		if err := inventoryPage.Submit(); err != nil {
			return err
		}

		log.Info("New item added to UI:", newItem.Name)

		return nil
	}

	insert = browser.WithRetry(3, time.Second*10)(insert)

	return s.Do(ctx, insert)
}

func (i InventoryScraper) Snapshot(
	ctx context.Context,
	session Session,
) ([]Rubro, error) {
	var rubros []Rubro

	snapshot := func(p *rod.Page) error {
		inventoryPage := NewInventoryPage(p)

		if err := inventoryPage.Open(); err != nil {
			return err
		}

		r, err := inventoryPage.ExtractExistingRubros()
		rubros = r
		return err
	}

	snapshot = browser.WithRetry(3, time.Second*10)(snapshot)

	return rubros, session.Do(ctx, snapshot)
}
