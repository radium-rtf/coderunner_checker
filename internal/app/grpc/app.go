package appgrpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/radium-rtf/coderunner_checker/internal/api"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	"github.com/radium-rtf/coderunner_checker/internal/domain"
)

type App struct {
	log     *slog.Logger
	server  *grpc.Server
	cfg     config.ServerConfig
	toClose []domain.Closer
}

// New creates new gRPC server app.
func New(log *slog.Logger, cfg config.ServerConfig, checkerSrv domain.CheckerSrv) *App {
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

	server := grpc.NewServer(grpc.ChainStreamInterceptor(
		recovery.StreamServerInterceptor(recoveryOpts...),
		logging.StreamServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	api.RegisterChecker(server, checkerSrv)

	return &App{
		log:     log,
		server:  server,
		cfg:     cfg,
		toClose: []domain.Closer{checkerSrv},
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// Serve net.Listener and if context closed run GracefulStop()
func (a *App) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		a.server.GracefulStop()
		l.Close()

		for _, close := range a.toClose {
			err := close.Close()
			if err != nil {
				a.log.Error(err.Error())
			}
		}

		a.log.Info("grpc server stopped")
	}(ctx)

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	return a.server.Serve(l)
}
