package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/radium-rtf/coderunner_checker/internal/app"
	"github.com/radium-rtf/coderunner_checker/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
