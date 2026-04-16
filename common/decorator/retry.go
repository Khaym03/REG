package decorator

import (
	"context"
	"log"
	"time"

	"github.com/Khaym03/REG/common/decorator/command"
)

var DefaultRetryConfig = RetryDecoratorConfig{
	Attempts:       3,
	Delay:          10 * time.Second,
	AttemptTimeout: 10 * time.Second,
}

type RetryDecoratorConfig struct {
	Attempts       int
	Delay          time.Duration
	AttemptTimeout time.Duration
}

type RetryDecorator[C any] struct {
	base           command.CommandHandler[C]
	attempts       int
	delay          time.Duration
	attemptTimeout time.Duration
}

func NewRetryDecorator[C any](base command.CommandHandler[C], cfg RetryDecoratorConfig) RetryDecorator[C] {
	if cfg.Attempts < 1 {
		cfg.Attempts = 1
	}

	return RetryDecorator[C]{
		base:           base,
		attempts:       cfg.Attempts,
		delay:          cfg.Delay,
		attemptTimeout: cfg.AttemptTimeout,
	}
}

func (d RetryDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	for attempt := 1; attempt <= d.attempts; attempt++ {
		log.Printf("retry attempt %d/%d for %s", attempt, d.attempts, generateActionName(cmd))

		attemptCtx := ctx
		var cancel context.CancelFunc
		if d.attemptTimeout > 0 {
			attemptCtx, cancel = context.WithTimeout(ctx, d.attemptTimeout)
		}

		err = d.base.Handle(attemptCtx, cmd)
		if cancel != nil {
			cancel()
		}

		if err == nil {
			return nil
		}

		log.Printf("retry failed attempt %d/%d for %s: %v", attempt, d.attempts, generateActionName(cmd), err)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if attempt == d.attempts {
			return err
		}

		if d.delay > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(d.delay):
			}
		}
	}

	return err
}
