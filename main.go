package main

import (
	"embed"
	"os"

	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:      os.Getenv("APP_ENV") == "dev",
		DisableTimestamp: true,
	})

	log.SetReportCaller(true)
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name: "REG",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	evBus := event.NewBus(app.Event)
	registerEvents()

	manager := mediator.NewSessionMediator(evBus)

	services := [...]application.Service{
		application.NewService(manager),
		application.NewService(NewAccountsAPI(manager, evBus)),
		application.NewService(NewAppService(app, manager, evBus)),
	}

	for i := range services {
		app.RegisterService(services[i])
	}

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "REG",
		Width:     1024,
		Height:    600,
		Frameless: true,
		// DisableResize: true,
		BackgroundColour: application.RGBA{
			Red:   27,
			Green: 38,
			Blue:  54,
			Alpha: 1,
		},
	})

	app.OnShutdown(func() {
		manager.Close()
	})

	if err := app.Run(); err != nil {
		println("Error:", err.Error())
	}
}

func registerEvents() {
	for _, e := range event.AviableTopis {
		switch e {
		case string(event.Stats):
			application.RegisterEvent[stats.Stats](e)
		default:
			application.RegisterEvent[event.Empty](e)
		}
	}
}
