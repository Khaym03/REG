package scraper

import (
	"fmt"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func applyFiltersToGuideReceiver(page *rod.Page, date utils.DateRange) error {
	el, err := page.ElementX(filterAccordionSelector)

	if err != nil {
		return fmt.Errorf("filter accordion not found: %w", err)
	}
	if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click filter accordion: %w", err)
	}

	// Wait for the accordion animation to settle
	if err := page.WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after opening filters: %w", err)
	}

	if err := utils.SelectOption(page, selectStatusSelector, selectStatusOption); err != nil {
		return fmt.Errorf("failed to select status: %w", err)
	}
	if err := utils.SelectOption(page, selectReceptionStatus, selectReceptionOption); err != nil {
		return fmt.Errorf("failed to select reception status: %w", err)
	}

	if err := utils.FillInputTime(page, inputDateFromSelector, date.From); err != nil {
		return err
	}
	if err := utils.FillInputTime(page, inputDateToSelector, date.To); err != nil {
		return err
	}

	btn, err := page.ElementX(filterButtonSelector)
	if err != nil {
		return fmt.Errorf("filter submit button not found: %w", err)
	}
	if err := btn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click filter button: %w", err)
	}

	return page.WaitDOMStable(c.TimeoutShort, 0.5)
}
