package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	user := loadCredential()

	browser := buildBrowser()
	defer browser.MustClose()

	repo := adapters.NewJSONGuideRepository("state.json")

	var authService domain.AuthService = scraper.NewLoginScraper()

	workflow := app.NewReceptionWorkflow(
		session.NewProvider(
			browser,
			decorator.NewLoggingDecorator(
				decorator.NewRetryDecorator(
					command.NewLoginHandler(authService),
					decorator.DefaultRetryConfig,
				),
			),
			decorator.NewLoggingDecorator(
				decorator.NewRetryDecorator(
					command.NewLogoutHandler(authService),
					decorator.DefaultRetryConfig,
				),
			),
		),
		decorator.NewLoggingDecorator(
			command.NewGatherGuidesHandler(
				repo,
				scraper.NewGuidesScraper(),
				scraper.NewRodRubroWorker(1),
			),
		),
		decorator.NewLoggingDecorator(
			command.NewInventoryHandler(
				repo,
				scraper.NewInventoryScraper(),
			),
		),
		decorator.NewLoggingDecorator(
			command.NewReceptionistHandler(
				repo,
				scraper.NewReceptionistScraper(),
			),
		),
	)

	go func() {
		if err := workflow.Run(context.Background(), app.WorkFlowInput{
			User: user,
		}); err != nil {
			log.Println(err)
		}

		cancel()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Interruption signal received (%v). Shutting down...", sig)
	case <-ctx.Done():
		log.Printf("The scraping process completed its normal execution.")
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
