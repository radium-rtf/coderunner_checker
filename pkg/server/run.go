package server

import (
	"context"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

type Server struct {
	ctx  context.Context
	cfg  config.ServerConfig
	grpc *grpc.Server

	wg  sync.WaitGroup
	err error

	started atomic.Bool
}

// Serve gRPC server and when context closed stops the gRPC server gracefully.
func New(ctx context.Context, cfg config.ServerConfig, interceptor ...grpc.ServerOption) *Server {
	grpc := grpc.NewServer(interceptor...)

	server := &Server{grpc: grpc, cfg: cfg, ctx: ctx}
	server.wg.Add(1)

	return server
}

func (s *Server) Wait() error {
	s.wg.Wait()
	return s.err
}

func (s *Server) Start() error {
	if s.started.Swap(true) {
		return ErrAlreadyStarted
	}

	lis, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	errorChan := make(chan error)
	go func() {
		err := s.grpc.Serve(lis)
		if err != nil {
			errorChan <- err
		}
	}()

	go func(ctx context.Context) {
		defer s.wg.Done()
		defer lis.Close()
		defer s.grpc.GracefulStop()

		select {
		case <-ctx.Done():
			return
		case err := <-errorChan:
			s.err = err
			return
		}
	}(s.ctx)

	return nil
}

func (s *Server) GetRegistrar() grpc.ServiceRegistrar {
	return s.grpc
}
