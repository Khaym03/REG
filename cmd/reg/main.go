package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Khaym03/REG/internal/auth"
	"github.com/Khaym03/REG/internal/config"
	"github.com/Khaym03/REG/internal/domain"
	"github.com/Khaym03/REG/internal/event"
	"github.com/Khaym03/REG/internal/workflow"
	"github.com/joho/godotenv"
)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	eventBus := event.NewBus()

	sm := auth.NewSessionManager(eventBus)

	if err := sm.Init(ctx); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := sm.Close(); err != nil {
			log.Error(err)
		}
	}()

	work, err := workflow.NewReceptionWorkflow(ctx, eventBus, sm)
	if err != nil {
		log.Fatal(err)
	}

	err = work.Run(
		ctx,
		workflow.WorkFlowInput{
			User:        auth.LoadCredential(),
			Date:        getDateRangeFromFlags(),
			BrowserConf: config.BrowserConfFromENV(),
		},
	)

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
