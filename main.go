package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Khaym03/REG/app"
	"github.com/Khaym03/REG/container"
	"github.com/Khaym03/REG/domain"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mustLoadEnv()

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
			},
		)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Println(err)
		}
		stop()
	case <-ctx.Done():
	}
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
