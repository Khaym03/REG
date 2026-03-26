package inventory

import (
	"fmt"
	"log"
	"time"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
)

const (
	// selectInput  = `//select[@wire:model.defer='data.WHlvWDNuaFpuQ2lpN1lLOXV4OVgxUT09']`
	uploadButton = `//button[i[contains(@class, 'fa-upload')]]`
)

type Inventory struct {
	store map[string]struct{}
	page  *rod.Page
}

func NewInventory(page *rod.Page) *Inventory {
	m := make(map[string]struct{})

	page.MustNavigate(c.InventoryURL)
	page.MustWaitLoad()

	firstload(page, m)

	return &Inventory{page: page, store: m}
}

const (
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÑ"
	lowercase = "abcdefghijklmnopqrstuvwxyzáéíóúñ"
)

func (inv *Inventory) Insert(newItem string) {
	_, exist := inv.store[newItem]
	if !exist {
		time.Sleep(time.Second * 6)

		// click the Select2 container to open the dropdown
		inv.page.MustElement(".select2-selection").MustClick()

		xpathOption := fmt.Sprintf(
			`//li[contains(translate(text(), "%s", "%s"), translate("%s", "%s", "%s"))]`,
			uppercase, lowercase, newItem, uppercase, lowercase,
		)

		// Wait for the option to be visible before clicking
		inv.page.MustElementX(xpathOption).MustClick()

		inv.page.MustElementX(uploadButton).MustClick()

		log.Println("New item added to UI:", newItem)
		inv.store[newItem] = struct{}{}
		return
	}

}

func firstload(page *rod.Page, store map[string]struct{}) {
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

		store[rubro] = struct{}{}
		log.Println("New Item found:", rubro)
	}
}
