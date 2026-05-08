package browser

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-rod/rod"
)

func WithRetry(
	attempts int,
	delay time.Duration,
) func(PageFunc) PageFunc {
	return func(next PageFunc) PageFunc {
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
				log.Warn("Retrying")
			}

			return err
		}
	}
}
