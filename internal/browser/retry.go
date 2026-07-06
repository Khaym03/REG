package browser

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-rod/rod"
)

func WithRetry(
	ctx context.Context,
	attempts int,
	delay time.Duration,
) func(PageFunc) PageFunc {
	return func(next PageFunc) PageFunc {
		return func(p *rod.Page) error {
			var err error

			for i := range attempts {
				err = next(p)
				if err == nil {
					return nil
				}

				// if !errors.Is(err, &rod.NavigationError{}) {
				// 	return err
				// }

				if i == attempts-1 {
					break
				}

				select {
				case <-ctx.Done():
					return context.Cause(ctx)
				case <-time.After(delay):
				}
				delay *= 2
				log.Warn("Retrying")
			}

			return err
		}
	}
}
