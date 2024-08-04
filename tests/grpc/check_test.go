package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/radium-rtf/coderunner_checker/internal/app"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	code = `
print(int(input())*2)
`

	configPath = "../../config/config.yaml"
)

func runServer(ctx context.Context, cfg *config.Config) error {
	app, err := app.New(cfg)
	if err != nil {
		return err
	}
	return app.Server.Run(ctx)
}

func getConnection(cfg config.ServerConfig) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TestCheckSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadFromPath(configPath)
	require.NoError(t, err)

	go runServer(ctx, cfg)

	conn, err := getConnection(cfg.Server)
	require.NoError(t, err)
	defer conn.Close()

	checker := checkergrpc.NewCheckerClient(conn)

	in := &checkergrpc.ArrayTestsRequest{
		Tests: []*checkergrpc.ArrayTestsRequest_TestCase{
			{
				Stdin:  "1\n",
				Stdout: "2\n",
			},
			{
				Stdin:  "2\n",
				Stdout: "4\n",
			},
		},
		Request: &checkergrpc.TestRequest{
			Lang:             "python",
			Code:             code,
			Timeout:          durationpb.New(time.Second * 1),
			MemoryLimitBytes: 1024 * 1024 * 6,
			FullInfoWa:       false,
		},
	}

	client, err := checker.Check(ctx, in)
	require.NoError(t, err)

	for i := 0; i < len(in.Tests); i++ {
		resp, err := client.Recv()
		require.NoError(t, err)
		assert.Equal(t, checkergrpc.Status_STATUS_SUCCESS, resp.Status)
	}
}
