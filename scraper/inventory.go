package scraper

import (
	"context"
	"fmt"
	"log"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/session"
	"github.com/go-rod/rod/lib/proto"
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

	if err := page.Navigate(c.InventoryURL); err != nil {
		return fmt.Errorf("failed to navigate to inventory: %w", err)
	}
	if err = page.WaitLoad(); err != nil {
		return fmt.Errorf("wait load failed: %w", err)
	}

	// click the Select2 container to open the dropdown
	selection, err := page.Element(".select2-selection")
	if err != nil {
		return fmt.Errorf("select2 container not found: %w", err)
	}
	if err := selection.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click select2: %w", err)
	}

	xpathOption := fmt.Sprintf(
		`//li[contains(translate(text(), "%s", "%s"), translate("%s", "%s", "%s"))]`,
		uppercase, lowercase, newItem.Name, uppercase, lowercase,
	)

	// Wait for the option to be visible before clicking
	option, err := page.ElementX(xpathOption)
	if err != nil {
		return fmt.Errorf("item option '%s' not found in dropdown: %w", newItem.Name, err)
	}
	if err := option.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click item option: %w", err)
	}

	btn, err := page.ElementX(uploadButton)
	if err != nil {
		return fmt.Errorf("upload button not found: %w", err)
	}
	if err := btn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click upload button: %w", err)
	}

	log.Println("New item added to UI:", newItem.Name)

	return nil

}

func (i InventoryScraper) RubrosSnapshot(ctx context.Context) ([]domain.Rubro, error) {
	page := session.FromContext(ctx).MainPage().Context(ctx)
	var err error

	if err = page.Navigate(c.InventoryURL); err != nil {
		return nil, fmt.Errorf("failed to navigate to inventory: %w", err)
	}
	if err = page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("wait load failed: %w", err)
	}

	rows, err := page.Elements("table tbody tr")
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory table rows: %w", err)
	}

	var existingOnes []domain.Rubro
	for _, row := range rows {
		// Get columns using relative XPath
		// td[2] is the Item (Rubro)
		// td[3] is the Balance (Saldo)
		cells, err := row.Elements("td")
		if err != nil || len(cells) < 3 {
			log.Println("skipping: ", cells.First().MustHTML())
		}

		rubro := cells[1].MustText()

		existingOnes = append(existingOnes, domain.Rubro{
			Name: rubro,
		})

		log.Println("New Item found:", rubro)
	}

	return existingOnes, nil
}

const (
	uploadButton = `//button[i[contains(@class, 'fa-upload')]]`
	uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÑ"
	lowercase    = "abcdefghijklmnopqrstuvwxyzáéíóúñ"
)
