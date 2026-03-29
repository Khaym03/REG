package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/Khaym03/REG/adapters"
	"github.com/Khaym03/REG/app"
	"github.com/Khaym03/REG/app/command"
	"github.com/Khaym03/REG/common/decorator"
	"github.com/Khaym03/REG/domain"
	"github.com/Khaym03/REG/scraper"
	"github.com/Khaym03/REG/session"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	user := loadCredential()

	browser := buildBrowser()
	defer browser.MustClose()

	repo := adapters.NewJSONGuideRepository("state.json")

	var authService domain.AuthService = scraper.NewLoginScraper()

	sessionProvider := session.NewProvider(
		browser,
		decorator.NewLoggingDecorator(command.NewLoginHandler(authService)),
		decorator.NewLoggingDecorator(command.NewLogoutHandler(authService)),
	)

	// --- workflow ---
	workflow := app.NewReceptionWorkflow(
		sessionProvider,
		decorator.NewLoggingDecorator(
			command.NewGatherGuidesHandler(
				repo,
				scraper.NewGuidesScraper(sessionProvider),
				scraper.NewRodRubroWorker(sessionProvider, 1),
			),
		),
		decorator.NewLoggingDecorator(
			command.NewInventoryHandler(
				repo,
				scraper.NewInventoryScraper(sessionProvider),
			),
		),
		decorator.NewLoggingDecorator(
			command.NewReceptionistHandler(
				repo,
				scraper.NewReceptionistScraper(sessionProvider),
			),
		),
	)

	if err := workflow.Run(context.Background(), app.WorkFlowInput{
		User: user,
	}); err != nil {
		log.Fatal(err)
	}
}

func loadCredential() domain.User {
	return domain.User{
		Username: os.Getenv("REG_TEST_USERNAME"),
		Password: os.Getenv("REG_TEST_PASSWORD"),
	}
}

func buildBrowser() *rod.Browser {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	l := launcher.New().
		Headless(os.Getenv("REG_HEADLESS") == "1").
		Devtools(false).
		Leakless(false).
		UserDataDir(filepath.Join(rootDir, "rod_data"))

	return rod.New().
		ControlURL(l.MustLaunch()).
		Trace(os.Getenv("REG_ROD_VERBOSE") == "1").
		MustConnect()
}
