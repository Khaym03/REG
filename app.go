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
	"github.com/Khaym03/REG/internal/workflow/service"
	"github.com/mustafaturan/bus/v3"
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
	eventBus  *bus.Bus
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		mu:       sync.Mutex{},
		eventBus: event.NewBus(),
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

	a.registerEventHandlers()
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
		application, cleanup, err := service.NewApplication(
			ctx,
			browserConf,
			a.eventBus,
		)
		if err != nil {
			log.Error(err)
			return
		}
		defer cleanup()

		work := workflow.NewReceptionWorkflow(application)

		err = work.Run(
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

func (a *App) registerEventHandlers() {
	onStatResult := bus.Handler{
		Handle: func(ctx context.Context, e bus.Event) {
			d, ok := e.Data.(stats.Stats)
			if !ok {
				return
			}

			runtime.EventsEmit(a.ctx, event.StatsTopic, d)
		},
		Matcher: event.Matcher(event.StatsTopic),
	}
	a.eventBus.RegisterHandler(event.StatsTopic, onStatResult)

	activeEvents := []string{
		event.LogginTopic,
		event.LogoutTopic,
		event.GuidesGatherTopic,
		event.InventorySyncTopic,
		event.ReceptionTopic,
	}

	justEmitEvent := func(e string) bus.Handler {
		return bus.Handler{
			Handle: func(ctx context.Context, ev bus.Event) {
				log.Printf("HANDLER REGISTERED FOR: %s", e)
				log.Printf("EVENT RECEIVED: %+v", ev)
				runtime.EventsEmit(a.ctx, e, "")
			},
			Matcher: event.Matcher(e),
		}
	}

	for _, e := range activeEvents {
		a.eventBus.RegisterHandler(e, justEmitEvent(e))
	}

}

func (a *App) Topics() event.Topics {
	return event.StructTopics()
}

func (a *App) Ignore(_ stats.Stats) {}

func (a *App) GetUser() auth.User {
	return auth.LoadCredential()
}
