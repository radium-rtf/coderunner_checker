package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/radium-rtf/coderunner_checker/internal/app"
	"github.com/radium-rtf/coderunner_checker/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// Gracefull shutdown
	go func(cancel context.CancelFunc) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop
		cancel()
	}(cancel)

	err = app.Server.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
