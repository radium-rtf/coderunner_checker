package appgrpc

import (
	"context"
	checkersrv "github.com/radium-rtf/coderunner_checker/internal/services/checker"
	server2 "github.com/radium-rtf/coderunner_checker/pkg/server"
	"io"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/radium-rtf/coderunner_checker/internal/api"
	"github.com/radium-rtf/coderunner_checker/internal/config"
)

type App struct {
	log     *slog.Logger
	server  *server2.Server
	cfg     config.ServerConfig
	toClose []io.Closer
}

// New creates new gRPC server app.
func New(ctx context.Context, log *slog.Logger, cfg config.ServerConfig, checkerSrv *checkersrv.CheckerService) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadSent, logging.PayloadReceived,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcServer := server2.New(ctx, cfg, grpc.ChainStreamInterceptor(
		recovery.StreamServerInterceptor(recoveryOpts...),
		logging.StreamServerInterceptor(interceptorLogger(log), loggingOpts...),
	))

	api.RegisterChecker(grpcServer.GetRegistrar(), checkerSrv)

	return &App{
		log:     log,
		server:  grpcServer,
		cfg:     cfg,
		toClose: []io.Closer{checkerSrv},
	}
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a *App) Run() error {
	err := a.server.Start()
	if err != nil {
		return err
	}
	a.log.Info("grpc server started")
	return nil
}

func (a *App) Wait() error {
	err := a.server.Wait()
	if err != nil {
		return err
	}

	for _, close := range a.toClose {
		err := close.Close()
		if err != nil {
			return err
		}
	}

	a.log.Info("grpc server stopped")
	return nil
}
