package scraper

import (
	"fmt"
	"time"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func navigate(page *rod.Page, url string) (err error) {
	if err = page.Navigate(url); err != nil {
		return fmt.Errorf("failed to navigate to login: %w", err)
	}
	if err = page.WaitLoad(); err != nil {
		return fmt.Errorf("wait load failed: %w", err)
	}

	return nil
}

func click(page *rod.Page, query string) error {
	el, err := page.Search(query)
	if err != nil {
		return err
	}

	return el.First.Click(proto.InputMouseButtonLeft, 1)
}

func fillInput(page *rod.Page, selector, value string) error {
	el, err := page.Element(selector)
	if err != nil {
		return err
	}
	if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	return el.Input(value)
}

func fillInputTime(page *rod.Page, xpath string, t time.Time) error {
	el, err := page.ElementX(xpath)
	if err != nil {
		return fmt.Errorf("date input element not found (%s): %w", xpath, err)
	}
	if err := el.InputTime(t); err != nil {
		return fmt.Errorf("failed to input time into %s: %w", xpath, err)
	}
	return nil
}

func selectOption(page *rod.Page, parentSelector string, optionXPath string) (err error) {
	if err = click(page, parentSelector); err != nil {
		return
	}

	if err = click(page, optionXPath); err != nil {
		return err
	}

	return page.WaitDOMStable(c.TimeoutShort, 0.5)
}
