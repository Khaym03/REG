package main

import (
	"context"
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

	user := loadCredential()
	browser := container.BuildBrowser()
	defer browser.MustClose()

	c := container.BuildContainer(browser)

	lastYearToPresent := domain.DateRange{
		From: time.Now().AddDate(-1, 0, 0),
		To:   time.Now(),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.Workflow.Run(
			ctx,
			app.WorkFlowInput{
				User: user,
				Date: lastYearToPresent,
			},
		)
	}()

	err := <-errCh
	if err != nil {
		log.Println(err)
	}

	<-ctx.Done()

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
