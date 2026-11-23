package inventory

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/browser"
	c "github.com/Khaym03/REG/internal/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type InventoryPage struct {
	page *rod.Page
}

func NewInventoryPage(p *rod.Page) *InventoryPage {
	return &InventoryPage{page: p}
}

func (ip *InventoryPage) Open() error {
	return browser.Navigate(ip.page, c.InventoryURL)
}

func (ip *InventoryPage) SelectItem(newItem Rubro) (err error) {
	// click the Select2 container to open the dropdown
	if err = browser.Click(ip.page, select2Container); err != nil {
		return
	}

	waitLivewire := ip.page.Timeout(c.DefaultTimeout).
		WaitRequestIdle(500*time.Millisecond, []string{livewireEndpointPattern}, nil, nil)

	if err = ip.clickSelect2Option(newItem.Name); err != nil {
		return fmt.Errorf("item option '%s' not found in dropdown: %w", newItem.Name, err)
	}

	// Wait for Livewire's select2-events change handler (.set(id, val, false)) to finish
	waitLivewire()

	if err = ip.page.Timeout(c.DefaultTimeout).WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after selecting '%s': %w", newItem.Name, err)
	}

	return nil
}

func (ip *InventoryPage) clickSelect2Option(name string) error {
	exactXPath := fmt.Sprintf(
		`//li[contains(@class, "select2-results__option") and normalize-space(translate(text(), "%s", "%s")) = normalize-space(translate("%s", "%s", "%s"))]`,
		uppercase, lowercase, name, uppercase, lowercase,
	)

	if el, err := ip.page.Timeout(c.TimeoutShort).ElementX(exactXPath); err == nil && el != nil {
		return el.Timeout(c.DefaultTimeout).Click(proto.InputMouseButtonLeft, 1)
	}

	fallbackXPath := fmt.Sprintf(
		`//li[contains(@class, "select2-results__option") and contains(translate(text(), "%s", "%s"), translate("%s", "%s", "%s"))]`,
		uppercase, lowercase, name, uppercase, lowercase,
	)

	return browser.Click(ip.page, fallbackXPath)
}

func (ip *InventoryPage) ExtractExistingRubros() ([]Rubro, error) {
	rows, err := ip.page.Elements(inventoryTableRows)
	if err != nil {
		return nil, fmt.Errorf("failed to find inventory table rows: %w", err)
	}

	var existingOnes []Rubro
	for _, row := range rows {
		// Get columns using relative cells
		// td[2] is the Item (Rubro)
		// td[3] is the Balance (Saldo)
		cells, err := row.Elements("td")
		if err != nil || len(cells) < 3 {
			if len(cells) > 0 {
				log.Warn("skipping: ", cells.First().MustHTML())
			}
			continue
		}

		rubro, err := cells[1].Text()
		if err != nil {
			return nil, fmt.Errorf("failed to extract rubro text: %w", err)
		}

		existingOnes = append(existingOnes, Rubro{
			Name: strings.TrimSpace(rubro),
		})
	}

	return existingOnes, nil
}

func (ip *InventoryPage) Submit() error {
	btn, err := ip.page.Timeout(c.DefaultTimeout).ElementX(uploadButton)
	if err != nil {
		return fmt.Errorf("upload button not found: %w", err)
	}

	if err := btn.Timeout(c.DefaultTimeout).WaitEnabled(); err != nil {
		return fmt.Errorf("upload button remained disabled: %w", err)
	}

	waitLivewire := ip.page.Timeout(c.DefaultTimeout).
		WaitRequestIdle(time.Second, []string{livewireEndpointPattern}, nil, nil)

	if err := btn.Timeout(c.DefaultTimeout).Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click upload button: %w", err)
	}

	// Wait for Livewire wire:click POST request to complete
	waitLivewire()

	if err := ip.page.Timeout(c.DefaultTimeout).WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after inventory submit: %w", err)
	}

	return ip.checkErrorToast()
}

func (ip *InventoryPage) WaitForRubro(item Rubro) error {
	deadline := time.Now().Add(c.DefaultTimeout)
	target := strings.TrimSpace(item.Name)

	for time.Now().Before(deadline) {
		if err := ip.checkErrorToast(); err != nil {
			return err
		}

		rubros, err := ip.ExtractExistingRubros()
		if err == nil {
			for _, r := range rubros {
				if strings.EqualFold(strings.TrimSpace(r.Name), target) {
					return nil
				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("rubro %q was not persisted in inventory table after submit", item.Name)
}

func (ip *InventoryPage) checkErrorToast() error {
	toasts, err := ip.page.Elements(errorToastSelector)
	if err != nil || toasts.Empty() {
		return nil
	}

	for _, toast := range toasts {
		visible, err := toast.Visible()
		if err != nil || !visible {
			continue
		}

		text, err := toast.Text()
		if err != nil {
			return fmt.Errorf("inventory error toast detected")
		}

		msg := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), "×"))
		if msg != "" {
			return fmt.Errorf("inventory submit failed: %s", msg)
		}
	}

	return nil
}

const (
	select2Container        = ".select2-selection"
	inventoryTableRows      = "table tbody tr"
	uploadButton            = `//button[i[contains(@class, 'fa-upload')]]`
	errorToastSelector      = `.jq-toast-single.jq-icon-error`
	livewireEndpointPattern = `/livewire/message/`
	uppercase               = "ABCDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÑ"
	lowercase               = "abcdefghijklmnopqrstuvwxyzáéíóúñ"
)
