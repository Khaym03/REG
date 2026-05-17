package main

import (
	"context"
	"io"
	"os"

	"sync"

	"github.com/Khaym03/REG/app"
	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/container"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WailsLogWriter struct {
	ctx context.Context
}

func (w *WailsLogWriter) Write(p []byte) (n int, err error) {
	runtime.EventsEmit(w.ctx, "LOG", string(p))

	return len(p), nil
}

// App struct
type App struct {
	ctx       context.Context
	cancel    context.CancelFunc
	mu        sync.Mutex
	loggerOut io.Writer
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		mu: sync.Mutex{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.loggerOut = io.MultiWriter(
		os.Stdout,
		&WailsLogWriter{ctx: a.ctx},
	)

	log.SetOutput(a.loggerOut)
}

func (a *App) RunWorkflow(input app.WorkFlowInput, browserConf browser.BrowserConfig) {
	log.Info(input)
	log.Info(browserConf)

	a.mu.Lock()

	if a.cancel != nil {
		a.cancel()
	}

	ctx, cancel := context.WithCancel(a.ctx)

	a.cancel = cancel

	a.mu.Unlock()

	done := make(chan struct{}, 1)

	go func() {
		defer func() {
			a.mu.Lock()

			if a.cancel != nil {
				a.cancel = nil
			}

			a.mu.Unlock()

			cancel()
		}()

		browserConf.LoggerOut = a.loggerOut
		browser, err := browser.BuildBrowser(ctx, browserConf)
		if err != nil {
			log.Error(err)
			return
		}
		defer browser.MustClose()

		c := container.BuildContainer(browser)

		err = c.Workflow.Run(
			ctx,
			app.WorkFlowInput{
				User: input.User,
				Date: input.Date,
			},
		)

		if err != nil {
			log.Println(err)
		}
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return
	case <-done:
		return
	}
}

func (a *App) StopWorkflow() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancel != nil {
		log.Info("Workflow canceled")
		a.cancel()
		a.cancel = nil
	}
}

func (a *App) GetUser() auth.User {
	return auth.LoadCredential()
}
