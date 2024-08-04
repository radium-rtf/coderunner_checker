package grpc

import (
	"context"

	"github.com/radium-rtf/coderunner_checker/internal/config"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"testing"
	"time"
)

const (
	code = `
print(int(input())*2)
`

	configPath = "../../config/config.yaml"
)

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
			Timeout:          durationpb.New(time.Second * 10),
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
