package grpc

import (
	"context"
	"fmt"
	"io"

	"testing"
	"time"

	"github.com/radium-rtf/coderunner_checker/internal/config"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	code = `
print(int(input())*2)
`

	configPath = "../../config/config.yaml"
)

type StreamResponse struct {
	Status  checkergrpc.Status
	Number  int64
	Message any
}

type TestCaseArray struct {
	Request  *checkergrpc.ArrayTestsRequest
	Response []StreamResponse
}

type TestCaseFile struct {
	Request  *checkergrpc.FileTestsRequest
	Response []StreamResponse
}

func TestCheckSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadFromPath(configPath)
	require.NoError(t, err)
	fmt.Println(cfg)

	go runServer(ctx, cfg)

	conn, err := getConnection(cfg.Server)
	require.NoError(t, err)
	defer conn.Close()

	checker := checkergrpc.NewCheckerClient(conn)

	testCases := []TestCaseArray{
		{
			Request: &checkergrpc.ArrayTestsRequest{
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
					Version:          "3.12",
					Code:             code,
					Timeout:          durationpb.New(time.Second * 10),
					MemoryLimitBytes: 1024 * 1024 * 6,
					FullInfoWa:       false,
				},
			},
			Response: []StreamResponse{
				{
					Number: 1,
					Status: checkergrpc.Status_STATUS_SUCCESS,
					Message: &checkergrpc.TestResponse_Text{
						Text: "success",
					},
				},
				{
					Number: 2,
					Status: checkergrpc.Status_STATUS_SUCCESS,
					Message: &checkergrpc.TestResponse_Text{
						Text: "success",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		client, err := checker.Check(ctx, test.Request)
		require.NoError(t, err)

		for _, expResp := range test.Response {
			resp, err := client.Recv()
			require.NoError(t, err)

			assert.Equal(t, expResp.Status, resp.Status)
			assert.Equal(t, expResp.Message, resp.Message)
			assert.Equal(t, expResp.Number, resp.Number)
		}
	}
}

func TestCheckURLSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadFromPath(configPath)
	require.NoError(t, err)

	go runServer(ctx, cfg)

	conn, err := getConnection(cfg.Server)
	require.NoError(t, err)
	defer conn.Close()

	checker := checkergrpc.NewCheckerClient(conn)

	in := &checkergrpc.FileTestsRequest{
		Url: "http://localhost:8082/",
		Request: &checkergrpc.TestRequest{
			Lang:             "python",
			Version:          "3.11",
			Code:             code,
			Timeout:          durationpb.New(time.Second * 10),
			MemoryLimitBytes: 1024 * 1024 * 6,
			FullInfoWa:       false,
		},
	}

	client, err := checker.CheckURL(ctx, in)
	require.NoError(t, err)

	for i := 1; ; i++ {
		resp, err := client.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		assert.Equal(t, checkergrpc.Status_STATUS_SUCCESS, resp.Status)
		assert.Equal(t, "succes", resp.Message)
		assert.Equal(t, i, resp.Number)
	}
}
