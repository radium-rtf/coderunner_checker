package main

import (
	"context"
	"github.com/radium-rtf/coderunner_checker/internal/api"
	"github.com/radium-rtf/coderunner_checker/pkg/server"
	"log"
	"os/signal"
	"syscall"

	"github.com/radium-rtf/coderunner_checker/internal/config"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	checker, err := api.NewChecker(cfg.SandboxConfig)
	if err != nil {
		log.Fatalln(err)
	}

	server, err := server.New(ctx, cfg.ServerConfig)
	if err != nil {
		log.Fatalln(err)
	}

	checkergrpc.RegisterCheckerServer(server.GetRegistrar(), checker)

	if err = server.Start(); err != nil {
		log.Fatalln(err)
	}

	if err = server.Wait(); err != nil {
		log.Fatalln(err)
	}
}
