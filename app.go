package main

import (
	"context"
	"io"
	"os"

	"sync"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/workflow"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	"github.com/mustafaturan/bus/v3"
	log "github.com/sirupsen/logrus"

	"github.com/wailsapp/wails/v3/pkg/application"
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
	cancel    context.CancelFunc
	mu        sync.Mutex
	loggerOut io.Writer
	eventBus  *bus.Bus
	app       *application.App
}

// NewAppService creates a new App application struct
func NewAppService(app *application.App) *App {
	return &App{
		app:      app,
		mu:       sync.Mutex{},
		eventBus: event.NewBus(),
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

		defer func() {
			done <- struct{}{}
		}()

		browserConf.LoggerOut = a.loggerOut

		input.BrowserConf = browserConf

		work, err := workflow.NewReceptionWorkflow(ctx, a.eventBus)
		if err != nil {
			log.Error(err)
			return
		}

		err = work.Run(
			ctx,
			input,
		)

		if err != nil {
			log.Println(err)
		}
	}()

	<-done
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

func (a *App) registerEventHandlers() {
	onStatResult := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			d, ok := e.Data.(stats.Stats)
			if !ok {
				return
			}

			a.app.Event.Emit(event.Stats, d)
		},
		Matcher: event.Matcher(event.Stats),
	}
	a.eventBus.RegisterHandler(event.Stats, onStatResult)

	activeEvents := []string{
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
		a.eventBus.RegisterHandler(e, justEmitEvent(e))
	}

}

func (a *App) Topics() event.Topics {
	return event.All()
}

func (a *App) Ignore(_ stats.Stats) {}

func (a *App) GetUser() auth.User {
	return auth.LoadCredential()
}
