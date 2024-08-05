package grpc

import (
	"context"
	"github.com/radium-rtf/coderunner_checker/internal/app"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runServer(ctx context.Context, cfg *config.Config) error {
	app, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}
	return app.Run()
}

func getConnection(cfg config.ServerConfig) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
