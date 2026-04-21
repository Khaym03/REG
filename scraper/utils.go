package scraper

import (
	"fmt"

	"github.com/go-rod/rod"
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
