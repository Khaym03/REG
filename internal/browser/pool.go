package browser

import (
	"context"
	"fmt"
)

const defaultBrowserPoolSize = 1

var _ BrowserPool = (*browserPool)(nil)

type browserPool struct {
	pool    chan BrowserRuntime
	factory BrowserRuntimeFactory
}

type BrowserRuntimeFactory func(ctx context.Context) (BrowserRuntime, error)

// DefaultBrowserRuntimeFactory creates a new browser runtime.
// It uses context.Background() for initialization to ensure the browser instance
// persists across multiple workflow executions, as the workflow context is canceled
// upon completion of each run.
func DefaultBrowserRuntimeFactory(ctx context.Context) (BrowserRuntime, error) {
	rt := NewBrowserRunTime()
	err := rt.Init(context.Background())
	if err != nil {
		return nil, fmt.Errorf("runtime.Init: %w", err)
	}

	return rt, nil
}

// Close implements [BrowserPool].
func (b *browserPool) Close() error {
	close(b.pool)
	for brt := range b.pool {
		err := brt.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Zero values have default config
type BrowserPoolConfig struct {
	size    int
	factory BrowserRuntimeFactory
}

func NewBrowserPool(cfg BrowserPoolConfig) *browserPool {
	if cfg.size == 0 {
		cfg.size = defaultBrowserPoolSize
	}

	if cfg.factory == nil {
		cfg.factory = DefaultBrowserRuntimeFactory
	}

	return &browserPool{
		pool:    make(chan BrowserRuntime, cfg.size),
		factory: cfg.factory,
	}
}

// Acquire implements [BrowserPool].
func (b *browserPool) Acquire(ctx context.Context) (BrowserRuntime, error) {
	select {
	case browser := <-b.pool:
		if browser != nil {
			return browser, nil
		}
	default:
	}

	return b.factory(ctx)
}

// Release implements [BrowserPool].
func (b *browserPool) Release(rt BrowserRuntime) {
	select {
	case b.pool <- rt:
	default:
		_ = rt.Close() // pool full, discard it
	}
}
