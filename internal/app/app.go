package app

import (
	"context"
	"log/slog"
	"os"

	appgrpc "github.com/radium-rtf/coderunner_checker/internal/app/grpc"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	checkersrv "github.com/radium-rtf/coderunner_checker/internal/services/checker"
)

// possible env values
const (
	localEnv = "local"
	devEnv   = "dev"
	prodEnv  = "prod"
)

type App struct {
	server *appgrpc.App
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	checkerSrv, err := checkersrv.NewCheckerSrv(cfg.Sandbox)
	if err != nil {
		return nil, err
	}

	logger := setupLogger(cfg.Env)

	grpcApp := appgrpc.New(ctx, logger, cfg.Server, checkerSrv)

	app := &App{
		server: grpcApp,
	}

	return app, nil
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) Wait() error {
	return a.server.Wait()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case localEnv:
		log = slog.New(slog.Default().Handler())
	case devEnv:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case prodEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
