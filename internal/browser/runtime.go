package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/Khaym03/REG/internal/config"
	"github.com/go-rod/rod"
	log "github.com/sirupsen/logrus"
)

var _ BrowserRuntime = (*browserRuntime)(nil)

type browserRuntime struct {
	initOnce    sync.Once
	initErr     error
	browser     *rod.Browser
	browserConf config.BrowserConfig
	stateMu     sync.Mutex

	ready chan struct{}
}

func NewBrowserRunTime() *browserRuntime {
	return &browserRuntime{
		initOnce:    sync.Once{},
		browserConf: config.BrowserConfFromENV(),
		stateMu:     sync.Mutex{},
		ready:       make(chan struct{}),
	}
}

func (br *browserRuntime) Init(ctx context.Context) error {
	br.initOnce.Do(func() {
		b, err := BuildBrowser(ctx, config.BrowserConfFromENV())
		if err != nil {
			br.initErr = err
			close(br.ready)
			return
		}

		br.browser = b

		close(br.ready)
	})

	return br.initErr
}
func (br *browserRuntime) Browser() *rod.Browser {
	<-br.ready
	return br.browser
}

// Reconfigure implements [BrowserRuntime].
func (br *browserRuntime) Reconfigure(
	ctx context.Context,
	cfg config.BrowserConfig,
) error {
	br.stateMu.Lock()
	defer br.stateMu.Unlock()

	if br.browserConf.Equal(cfg) {
		return nil
	}

	if br.browser != nil {
		_ = br.browser.Close()
		br.browser = nil
	}

	// Use context.Background() to ensure the browser instance persists across
	// multiple workflow executions, as the provided ctx is canceled upon
	// completion of the current workflow run.
	browser, err := BuildBrowser(context.Background(), cfg)
	if err != nil {
		log.Printf("%T %#v", err, err)
		return fmt.Errorf("BuildBrowser: %w", err)
	}

	br.browser = browser
	br.browserConf = cfg

	return nil
}

func (br *browserRuntime) Close() error {
	<-br.ready
	return br.browser.Close()
}
