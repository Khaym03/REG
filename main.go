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

	var authService domain.AuthService = scraper.NewLoginScraper(browser)
	var guideScraper domain.GuideScraper = scraper.NewGuidesScraper(browser)
	var rubroWorker domain.RubroWorker = scraper.NewRodRubroWorker(browser, 1)

	// --- handlers ---
	loginHandler := command.NewLoginHandler(authService)
	logoutHandler := command.NewLogoutHandler(authService)
	gatherHandler := command.NewGatherGuidesHandler(repo, guideScraper, rubroWorker)

	// --- decorators ---
	loginH := decorator.NewLoggingDecorator(loginHandler)
	logoutH := decorator.NewLoggingDecorator(logoutHandler)
	gatherH := decorator.NewLoggingDecorator(gatherHandler)

	// --- workflow ---
	workflow := app.NewReceptionWorkflow(
		loginH,
		logoutH,
		gatherH,
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
