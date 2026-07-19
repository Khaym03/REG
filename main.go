package main

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/mediator"
	"github.com/Khaym03/REG/internal/workflow/queries/stats"
	"github.com/Khaym03/REG/utils"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	loadEnvironmentFiles()

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

func loadEnvironmentFiles() {
	candidates := []string{".env", ".env.example"}

	if exePath, err := os.Executable(); err == nil {
		baseDirProd := utils.BaseDir()
		appDir := filepath.Dir(exePath)

		candidates = append(
			candidates,
			filepath.Join(appDir, ".env"),
			filepath.Join(appDir, ".env.example"),
			filepath.Join(baseDirProd, ".env"),
			filepath.Join(baseDirProd, ".env.example"),
		)
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		if _, err := os.Stat(candidate); err != nil {
			continue
		}

		if err := godotenv.Load(candidate); err != nil {
			log.Warnf("failed to load env file %s: %v", candidate, err)
		}
	}
}

func registerEvents() {
	application.RegisterEvent[stats.Stats](string(event.Stats))

	application.RegisterEvent[event.Empty](string(event.WorkflowStarted))
	application.RegisterEvent[event.Empty](string(event.BuildingBrowser))
	application.RegisterEvent[event.Empty](string(event.Login))
	application.RegisterEvent[event.Empty](string(event.GuidesGather))
	application.RegisterEvent[event.Empty](string(event.InventorySync))
	application.RegisterEvent[event.Empty](string(event.Reception))
	application.RegisterEvent[event.Empty](string(event.Logout))
	application.RegisterEvent[event.Empty](string(event.DestroyingBrowser))
	application.RegisterEvent[event.Empty](string(event.WorkflowFinished))

}
