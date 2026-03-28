package scraper

import (
	"github.com/Khaym03/REG/utils"
	"github.com/go-rod/rod"
)

func applyFiltersToGuideReceiver(page *rod.Page, date utils.DateRange) error {
	page.MustElementX(filterAccordionSelector).MustClick()
	page.MustWaitDOMStable()

	utils.SelectOption(page, selectStatusSelector, selectStatusOption)
	utils.SelectOption(page, selectReceptionStatus, selectReceptionOption)

	page.MustElementX(inputDateFromSelector).MustInputTime(date.From)
	page.MustElementX(inputDateToSelector).MustInputTime(date.To)
	page.MustElementX(filterButtonSelector).MustClick()
	page.MustWaitDOMStable()

	return nil
}
