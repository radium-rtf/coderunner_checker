package grpc

import (
	"context"
	"github.com/radium-rtf/coderunner_checker/internal/api"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/radium-rtf/coderunner_checker/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddr = ":8081"
)

func getConfig() *config.Config {
	return &config.Config{
		ServerConfig: config.ServerConfig{
			Env:     config.TestEnv,
			Address: serverAddr,
		},
		SandboxConfig: config.SandboxConfig{
			User:    "root",
			UUID:    0,
			WorkDir: "/sandbox",
			Host:    "unix:///var/run/docker.sock",
			Rules: config.Rules{
				"python": config.Rule{
					Filename: "main.py",
					Image:    "python",
					Launch:   "python3",
				},
			},
		},
	}
}

func runServer(ctx context.Context, cfg *config.Config) (*server.Server, error) {
	checker, err := api.NewChecker(cfg.SandboxConfig)
	if err != nil {
		return nil, err
	}

	server, err := server.New(ctx, cfg.ServerConfig)
	if err != nil {
		return nil, err
	}

	checkergrpc.RegisterCheckerServer(server.GetRegistrar(), checker)

	if err = server.Start(); err != nil {
		return nil, err
	}

	return server, nil
}

func getConnection() (*grpc.ClientConn, error) {
	return grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
