package scraper

import (
	"context"
	"fmt"
	"log"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

var _ domain.InventoryScraper = (*InventoryScraper)(nil)

type InventoryScraper struct {
	browser *rod.Browser
}

func NewInventoryScraper(
	browser *rod.Browser,
) *InventoryScraper {
	return &InventoryScraper{browser: browser}
}

func (i *InventoryScraper) Insert(ctx context.Context, newItem domain.Rubro) error {
	page := i.browser.MustPage(c.InventoryURL)
	defer page.Close()
	page.MustWaitLoad()

	// click the Select2 container to open the dropdown
	page.MustElement(".select2-selection").MustClick()

	xpathOption := fmt.Sprintf(
		`//li[contains(translate(text(), "%s", "%s"), translate("%s", "%s", "%s"))]`,
		uppercase, lowercase, newItem.Name, uppercase, lowercase,
	)

	// Wait for the option to be visible before clicking
	page.MustElementX(xpathOption).MustClick()

	page.MustElementX(uploadButton).MustClick()

	log.Println("New item added to UI:", newItem.Name)

	return nil

}

func (i InventoryScraper) RubrosSnapshot(ctx context.Context) ([]domain.Rubro, error) {
	var existingOnes []domain.Rubro

	page := i.browser.MustPage(c.InventoryURL)
	defer page.Close()
	page.MustWaitLoad()

	rows := page.MustElements("table tbody tr")

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
