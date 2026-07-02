package main

import (
	"context"
	"errors"
	"io"
	"os"

	"sync"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/workflow"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	ErrWorkflowCanceled = errors.New("workflow canceled")
)

type WailsLogWriter struct {
	ctx context.Context
	app *application.App
}

func (w *WailsLogWriter) Write(p []byte) (n int, err error) {
	w.app.Event.Emit("LOG", string(p))

	return len(p), nil
}

// App struct
type App struct {
	ctx       context.Context
	cancel    context.CancelCauseFunc
	mu        sync.Mutex
	loggerOut io.Writer
	eventBus  *bus.Bus
	app       *application.App

	sessionManager mediator.SessionMediator
}

// NewAppService creates a new App application struct
func NewAppService(app *application.App, sm mediator.SessionMediator) *App {
	return &App{
		app:            app,
		mu:             sync.Mutex{},
		eventBus:       event.NewBus(),
		sessionManager: sm,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	a.ctx = ctx

	writers := make([]io.Writer, 0, 2)
	writers = []io.Writer{&WailsLogWriter{ctx: a.ctx, app: a.app}}

	if config.IsDev() {
		writers = append(writers, os.Stdout)
	}

	a.loggerOut = io.MultiWriter(
		writers...,
	)

	log.SetOutput(a.loggerOut)

	a.registerEventHandlers()

	return nil
}

func (a *App) RunWorkflow(input workflow.WorkFlowInput, browserConf config.BrowserConfig) {
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

		return
	}

	err = work.Run(
		ctx,
		input,
	)

	if err != nil {
		log.Printf("Run returned err=%v", err)
		log.Printf("ctx.Err()=%v", ctx.Err())
		log.Printf("context.Cause()=%v", context.Cause(ctx))
	}

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

func (a *App) registerEventHandlers() {
	const ev = string(event.Stats)
	onStatResult := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			log.Printf("HANDLER REGISTERED FOR: %s", e)

			d, ok := e.Data.(stats.Stats)
			if !ok {
				return
			}
			log.Printf("EVENT RECEIVED: %+v", ev)
			a.app.Event.Emit(ev, d)
		},
		Matcher: event.Matcher(ev),
	}
	a.eventBus.RegisterHandler(ev, onStatResult)

	activeEvents := []event.Topic{
		event.WorkflowStarted,
		event.Login,
		event.BuildingBrowser,
		event.GuidesGather,
		event.InventorySync,
		event.Reception,
		event.Logout,
		event.DestroyingBrowser,
		event.WorkflowFinished,
	}

	justEmitEvent := func(e string) bus.Handler {
		return bus.Handler{
			Handle: func(ctx context.Context, ev bus.Event) {
				log.Printf("HANDLER REGISTERED FOR: %s", e)
				log.Printf("EVENT RECEIVED: %+v", ev)
				a.app.Event.Emit(e, "")
			},
			Matcher: event.Matcher(e),
		}
	}

	for _, e := range activeEvents {
		ev := string(e)
		a.eventBus.RegisterHandler(ev, justEmitEvent(ev))
	}

}

func (a *App) Ignore(_ stats.Stats, _ event.Topic) {}

func (a *App) GetUser() auth.User {
	return auth.LoadCredential()
}
