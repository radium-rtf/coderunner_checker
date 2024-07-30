package serverutils

import (
	"bufio"
	"context"
	"net"
	"os"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"google.golang.org/grpc"
)

// Serve gRPC server and when context closed stops the gRPC server gracefully.
func Run(ctx context.Context, server *grpc.Server, cfg domain.ServerConfig) error {
	if cfg.Env == domain.DevEnv {
		ctx, cancel := context.WithCancel(ctx)

		go func(ctx context.Context) {
			<-ctx.Done()
			server.GracefulStop()
		}(ctx)

		go handShutdown(cancel)
	}

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	return server.Serve(lis)
}

// Closed context with cancel func when bufio.Scanner scan something from os.Stdin
func handShutdown(cancel context.CancelFunc) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		cancel()
	}
}
