package main

import (
	"embed"
	"os"

	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
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
	evBus := event.NewBus()
	manager := mediator.NewSessionMediator(evBus)

	app := application.New(application.Options{
		Name: "REG",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(manager),
			application.NewService(
				NewAccountsAPI(manager, evBus),
			),
		},
	})

	app.RegisterService(
		application.NewService(
			NewAppService(app, manager),
		),
	)

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
