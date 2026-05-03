package scraper

import (
	"log"
	"time"

	"github.com/Khaym03/REG/domain"
	"github.com/go-rod/rod"
)

func WithRetry(
	attempts int,
	delay time.Duration,
) func(domain.PageFunc) domain.PageFunc {

	return func(next domain.PageFunc) domain.PageFunc {
		return func(p *rod.Page) error {
			var err error

			for i := 0; i < attempts; i++ {
				err = next(p)
				if err == nil {
					return nil
				}

				if i == attempts-1 {
					break
				}

				time.Sleep(delay)
				delay *= 2
				log.Println("Retrying")
			}

			return err
		}
	}
}
