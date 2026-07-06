package main

import (
	"context"
	"errors"
	"io"

	"sync"

	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/workflow"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	log "github.com/sirupsen/logrus"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	ErrWorkflowCanceled = errors.New("workflow canceled")
)

// App struct
type App struct {
	ctx       context.Context
	cancel    context.CancelCauseFunc
	mu        sync.Mutex
	loggerOut io.Writer
	eventBus  event.Bus
	app       *application.App

	sessionManager mediator.SessionMediator
}

// NewAppService creates a new App application struct
func NewAppService(
	app *application.App,
	sm mediator.SessionMediator,
	eventBus event.Bus,
) *App {
	return &App{
		app:            app,
		mu:             sync.Mutex{},
		eventBus:       eventBus,
		sessionManager: sm,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	a.ctx = ctx

	return nil
}

func (a *App) RunWorkflow(input workflow.WorkFlowInput, browserConf config.BrowserConfig) error {
	a.mu.Lock()
	if a.cancel != nil {
		a.cancel(ErrWorkflowCanceled)
	}

	ctx, cancel := context.WithCancelCause(a.ctx)
	a.cancel = cancel

	a.mu.Unlock()

	defer func() {
		a.mu.Lock()

		if a.cancel != nil {
			a.cancel = nil
		}

		a.mu.Unlock()

		cancel(ErrWorkflowCanceled)
	}()

	browserConf.LoggerOut = a.loggerOut

	input.BrowserConf = browserConf

	work, err := workflow.NewReceptionWorkflow(
		ctx,
		a.eventBus,
		a.sessionManager,
	)
	if err != nil {
		log.Error(err)
		return err
	}

	err = work.Run(
		ctx,
		input,
	)

	if err != nil {
		log.Printf("Run returned err=%v", err)
		log.Printf("ctx.Err()=%v", ctx.Err())
		log.Printf("context.Cause()=%v", context.Cause(ctx))

		return err
	}

	return nil
}

func (a *App) StopWorkflow() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancel != nil {
		log.Info("Workflow canceled")
		a.cancel(ErrWorkflowCanceled)
		a.cancel = nil
	}
}

func (a *App) Ignore(_ stats.Stats, _ event.Topic) {}
