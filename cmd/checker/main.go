package main

import (
	"context"
	"log"

	"github.com/radium-rtf/coderunner_checker/internal/config"
	"github.com/radium-rtf/coderunner_checker/internal/checker/grpc"
	serverutils "github.com/radium-rtf/coderunner_checker/internal/utils/server"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()

	checker, err := checker.NewChecker(cfg.SandboxConfig)
	if err != nil {
		log.Fatalln(err)
	}

	checkergrpc.RegisterCheckerServer(server, checker)

	err = serverutils.Run(ctx, server, cfg.ServerConfig)
	if err != nil {
		log.Fatalln(err)
	}
}
