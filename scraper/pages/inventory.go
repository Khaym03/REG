package pages

import (
	"fmt"
	"log"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

type InventoryPage struct {
	page *rod.Page
}

func NewInventoryPage(p *rod.Page) *InventoryPage {
	return &InventoryPage{page: p}
}

func (ip *InventoryPage) Open() error {
	return navigate(ip.page, c.InventoryURL)
}

func (ip *InventoryPage) SelectItem(newItem domain.Rubro) (err error) {
	// click the Select2 container to open the dropdown
	if err = click(ip.page, select2Container); err != nil {
		return
	}

	xpathOption := fmt.Sprintf(
		`//li[contains(translate(text(), "%s", "%s"), translate("%s", "%s", "%s"))]`,
		uppercase, lowercase, newItem.Name, uppercase, lowercase,
	)

	if err = click(ip.page, xpathOption); err != nil {
		return fmt.Errorf("item option '%s' not found in dropdown: %w", newItem.Name, err)
	}

	return
}

func (ip *InventoryPage) ExtractExistingRubros() ([]domain.Rubro, error) {
	rows, err := ip.page.Elements(inventoryTableRows)
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
	}

	return existingOnes, nil
}

func (ip *InventoryPage) Submit() error {
	return click(ip.page, uploadButton)
}

const (
	select2Container   = ".select2-selection"
	inventoryTableRows = "table tbody tr"
	uploadButton       = `//button[i[contains(@class, 'fa-upload')]]`
	uppercase          = "ABCDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÑ"
	lowercase          = "abcdefghijklmnopqrstuvwxyzáéíóúñ"
)
