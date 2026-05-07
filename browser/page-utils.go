package browser

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	c "github.com/Khaym03/REG/constants"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const defaultTimeout = 10 * time.Second

func Navigate(page *rod.Page, url string) (err error) {
	log.Info("navigating to: ", url)
	if err = page.Timeout(defaultTimeout).Navigate(url); err != nil {
		return fmt.Errorf("failed to navigate: %w", err)
	}
	if err = WaitLoad(page); err != nil {
		return fmt.Errorf("wait load failed: %w", err)
	}

	return nil
}

func WaitLoad(page *rod.Page) error {
	return page.Timeout(defaultTimeout).WaitLoad()
}

func Click(page *rod.Page, query string) error {
	el, err := page.Timeout(defaultTimeout).Search(query)
	if err != nil {
		return err
	}

	return el.First.Click(proto.InputMouseButtonLeft, 1)
}

func FillInput(page *rod.Page, selector, value string) error {
	el, err := page.Timeout(defaultTimeout).Element(selector)
	if err != nil {
		return err
	}
	if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}
	return el.Input(value)
}

func FillInputTime(page *rod.Page, xpath string, t time.Time) error {
	log.Infof("filling input time with: %s", t.Format("2006-01-02"))
	el, err := page.Timeout(defaultTimeout).ElementX(xpath)
	if err != nil {
		return fmt.Errorf("date input element not found (%s): %w", xpath, err)
	}
	if err := el.InputTime(t); err != nil {
		return fmt.Errorf("failed to input time into %s: %w", xpath, err)
	}
	return nil
}

func SelectOption(page *rod.Page, parentSelector string, optionXPath string) (err error) {
	if err = Click(page, parentSelector); err != nil {
		return
	}

	if err = Click(page, optionXPath); err != nil {
		return err
	}

	return page.WaitDOMStable(c.TimeoutShort, 0.5)
}
