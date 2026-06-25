package main

import (
	"embed"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      false,
		DisableTimestamp: true,
	})
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

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

	app.RegisterService(
		application.NewService(
			NewAppService(app),
		),
	)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:         "REG",
		Width:         1024,
		Height:        600,
		Frameless:     true,
		DisableResize: true,
		BackgroundColour: application.RGBA{
			Red:   27,
			Green: 38,
			Blue:  54,
			Alpha: 1,
		},
	})

	if err := app.Run(); err != nil {
		println("Error:", err.Error())
	}
}
