package scraper

import (
	"fmt"

	c "github.com/Khaym03/REG/constants"
	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
)

func applyFiltersToGuideReceiver(page *rod.Page, date utils.DateRange) (err error) {
	if err = click(page, filterAccordionSelector); err != nil {
		return err
	}

	// Wait for the accordion animation to settle
	if err = page.WaitDOMStable(c.TimeoutShort, 0.5); err != nil {
		return fmt.Errorf("page did not stabilize after opening filters: %w", err)
	}

	if err = selectOption(page, selectStatusSelector, selectStatusOption); err != nil {
		return fmt.Errorf("failed to select status: %w", err)
	}
	if err = selectOption(page, selectReceptionStatus, selectReceptionOption); err != nil {
		return fmt.Errorf("failed to select reception status: %w", err)
	}

	if err = fillInputTime(page, inputDateFromSelector, date.From); err != nil {
		return err
	}
	if err = fillInputTime(page, inputDateToSelector, date.To); err != nil {
		return err
	}

	if err = click(page, filterButtonSelector); err != nil {
		return err
	}

	return page.WaitDOMStable(c.TimeoutShort, 0.5)
}
