package browser

import (
	"context"
	"io"

	"github.com/Khaym03/REG/internal/config"
	"github.com/go-rod/rod"
)

type PageFunc func(*rod.Page) error

type Reconfigurer interface {
	Reconfigure(ctx context.Context, cfg config.BrowserConfig) error
}

// Responsibilities:
// build browser,
// reconnect if needed,
// close browser,
// configuration,
// crash recovery.
type BrowserRuntime interface {
	Browser() *rod.Browser
	Reconfigurer
	io.Closer
}

type BrowserPool interface {
	Acquire(ctx context.Context) (BrowserRuntime, error)
	Release(runtime BrowserRuntime)
	io.Closer
}
