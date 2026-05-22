package main

import (
	"context"
	"io"
	"os"

	"sync"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/workflow"
	"github.com/Khaym03/REG/internal/workflow/service"
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

	writers := make([]io.Writer, 0, 2)
	writers = []io.Writer{&WailsLogWriter{ctx: a.ctx}}

	if config.IsDev() {
		writers = append(writers, os.Stdout)
	}

	a.loggerOut = io.MultiWriter(
		writers...,
	)

	log.SetOutput(a.loggerOut)
}

func (a *App) RunWorkflow(input workflow.WorkFlowInput, browserConf config.BrowserConfig) {
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
		application, cleanup := service.NewApplication(ctx, browserConf)
		defer cleanup()

		work := workflow.NewReceptionWorkflow(application)

		err := work.Run(
			ctx,
			input,
		)

		if err != nil {
			log.Println(err)
		}
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
	case <-done:
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
