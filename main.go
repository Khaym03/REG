package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Khaym03/REG/app"
	"github.com/Khaym03/REG/container"
	"github.com/Khaym03/REG/domain"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer stop()

	mustLoadEnv()

	dateRange := getDateRangeFromFlags()
	user := loadCredential()
	browser := container.BuildBrowser()
	defer browser.MustClose()

	c := container.BuildContainer(browser)

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.Workflow.Run(
			ctx,
			app.WorkFlowInput{
				User: user,
				Date: dateRange,
			},
		)
	}()

	err := <-errCh
	if err != nil {
		log.Println(err)
	}

	<-ctx.Done()

}

func getDateRangeFromFlags() domain.DateRange {
	var fromStr, toStr string

	flag.StringVar(&fromStr, "from", "", "Start date in YYYY-MM-DD format")
	flag.StringVar(&toStr, "to", "", "End date in YYYY-MM-DD format")
	flag.Parse()

	// Default: Last year to present
	dateRange := domain.DateRange{
		From: time.Now().AddDate(-1, 0, 0),
		To:   time.Now(),
	}

	const layout = "2006-01-02"

	if fromStr != "" {
		parsedFrom, err := time.Parse(layout, fromStr)
		if err != nil {
			log.Fatalf("Invalid 'from' date format: %v", err)
		}
		dateRange.From = parsedFrom
	}

	if toStr != "" {
		parsedTo, err := time.Parse(layout, toStr)
		if err != nil {
			log.Fatalf("Invalid 'to' date format: %v", err)
		}
		dateRange.To = parsedTo
	}

	return dateRange
}

func mustLoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func loadCredential() domain.User {
	return domain.User{
		Username: os.Getenv("REG_TEST_USERNAME"),
		Password: os.Getenv("REG_TEST_PASSWORD"),
	}
}
