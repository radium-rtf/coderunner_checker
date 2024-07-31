package app

import (
	"log/slog"
	"os"

	appgrpc "github.com/radium-rtf/coderunner_checker/internal/app/grpc"
	"github.com/radium-rtf/coderunner_checker/internal/domain"
	checkersrv "github.com/radium-rtf/coderunner_checker/internal/services/checker"
)

type App struct {
	Server *appgrpc.App
}

func New(cfg *domain.Config) (*App, error) {
	checkerSrv, err := checkersrv.NewCheckerSrv(cfg.Sandbox)
	if err != nil {
		return nil, err
	}

	logger := setupLogger(cfg.Env)

	grpcApp := appgrpc.New(logger, cfg.Server, checkerSrv)

	app := &App{
		Server: grpcApp,
	}

	return app, nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case domain.LocalEnv:
		log = slog.New(slog.Default().Handler())
	case domain.DevEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case domain.ProdEnv:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
