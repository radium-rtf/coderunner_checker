package grpc

import (
	"context"
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

	codeWA = `
print(int(input()))
`

	codeTimeout = `
import time

time.sleep(30)
print(int(input())*2)
`

	codeOOMKilled = `
buf = [x for x in range(1024*1024*7)]
print(int(input())*2)
`

	codeUnknownError = `
0/0
`

	configPath = "../../config/config.yaml"
)

type StreamResponse struct {
	Status  checkergrpc.Status
	Number  int64
	Message any
}

type TestArray struct {
	Name     string
	Response []StreamResponse
	Request  *checkergrpc.ArrayTestsRequest
}

type TestFile struct {
	Name     string
	Response []StreamResponse
	Request  *checkergrpc.FileTestsRequest
}

func TestCheck(t *testing.T) {
	t.Parallel()

	testCases := []TestArray{
		{
			Name: "TestSuccess",
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
		{
			Name: "TestTimeout",
			Request: &checkergrpc.ArrayTestsRequest{
				Tests: []*checkergrpc.ArrayTestsRequest_TestCase{
					{
						Stdin:  "1\n",
						Stdout: "2\n",
					},
				},
				Request: &checkergrpc.TestRequest{
					Lang:             "python",
					Version:          "3.12",
					Code:             codeTimeout,
					Timeout:          durationpb.New(time.Second * 10),
					MemoryLimitBytes: 1024 * 1024 * 6,
					FullInfoWa:       false,
				},
			},
			Response: []StreamResponse{
				{
					Number: 1,
					Status: checkergrpc.Status_STATUS_TIMEOUT,
					Message: &checkergrpc.TestResponse_Text{
						Text: "time limit",
					},
				},
			},
		},
		{
			Name: "TestOOMKilled",
			Request: &checkergrpc.ArrayTestsRequest{
				Tests: []*checkergrpc.ArrayTestsRequest_TestCase{
					{
						Stdin:  "1\n",
						Stdout: "2\n",
					},
				},
				Request: &checkergrpc.TestRequest{
					Lang:             "python",
					Version:          "3.12",
					Code:             codeOOMKilled,
					Timeout:          durationpb.New(time.Second * 10),
					MemoryLimitBytes: 1024 * 1024 * 6,
					FullInfoWa:       false,
				},
			},
			Response: []StreamResponse{
				{
					Number: 1,
					Status: checkergrpc.Status_STATUS_OOM_KILLED,
					Message: &checkergrpc.TestResponse_Text{
						Text: "oom killed",
					},
				},
			},
		},
		{
			Name: "TestDivisionByZero",
			Request: &checkergrpc.ArrayTestsRequest{
				Tests: []*checkergrpc.ArrayTestsRequest_TestCase{
					{
						Stdin:  "1\n",
						Stdout: "2\n",
					},
				},
				Request: &checkergrpc.TestRequest{
					Lang:             "python",
					Version:          "3.12",
					Code:             codeUnknownError,
					Timeout:          durationpb.New(time.Second * 10),
					MemoryLimitBytes: 1024 * 1024 * 6,
					FullInfoWa:       false,
				},
			},
			Response: []StreamResponse{
				{
					Number: 1,
					Status: checkergrpc.Status_STATUS_UNKNOWN,
					Message: &checkergrpc.TestResponse_Text{
						Text: "Traceback (most recent call last):\n  File \"/sandbox/main.py\", line 2, in <module>\n    0/0\n    ~^~\nZeroDivisionError: division by zero\n",
					},
				},
			},
		},
		{
			Name: "TestWrongAnswer",
			Request: &checkergrpc.ArrayTestsRequest{
				Tests: []*checkergrpc.ArrayTestsRequest_TestCase{
					{
						Stdin:  "1\n",
						Stdout: "2\n",
					},
				},
				Request: &checkergrpc.TestRequest{
					Lang:             "python",
					Version:          "3.12",
					Code:             codeWA,
					Timeout:          durationpb.New(time.Second * 10),
					MemoryLimitBytes: 1024 * 1024 * 6,
					FullInfoWa:       false,
				},
			},
			Response: []StreamResponse{
				{
					Number: 1,
					Status: checkergrpc.Status_STATUS_WRONG_ANSWER,
					Message: &checkergrpc.TestResponse_Text{
						Text: "wrong answer on input:\n1\n",
					},
				},
			},
		},
		{
			Name: "TestWrongAnswerFullInfo",
			Request: &checkergrpc.ArrayTestsRequest{
				Tests: []*checkergrpc.ArrayTestsRequest_TestCase{
					{
						Stdin:  "1\n",
						Stdout: "2\n",
					},
				},
				Request: &checkergrpc.TestRequest{
					Lang:             "python",
					Version:          "3.12",
					Code:             codeWA,
					Timeout:          durationpb.New(time.Second * 10),
					MemoryLimitBytes: 1024 * 1024 * 6,
					FullInfoWa:       true,
				},
			},
			Response: []StreamResponse{
				{
					Number: 1,
					Status: checkergrpc.Status_STATUS_WRONG_ANSWER,
					Message: &checkergrpc.TestResponse_WrongAnswer_{
						WrongAnswer: &checkergrpc.TestResponse_WrongAnswer{
							Input:    "1\n",
							Actual:   "1\n",
							Expected: "2\n",
						},
					},
				},
			},
		},
	}

	cfg, err := config.LoadFromPath(configPath)
	require.NoError(t, err)

	conn, err := getConnection(cfg.Server)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
		conn.Close()
	})

	go runServer(ctx, cfg)

	checker := checkergrpc.NewCheckerClient(conn)

	for _, test := range testCases {
		test := test

		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			client, err := checker.Check(ctx, test.Request)
			require.NoError(t, err)

			for _, expResp := range test.Response {
				resp, err := client.Recv()
				require.NoError(t, err)

				assert.Equal(t, expResp.Status, resp.Status)
				assert.Equal(t, expResp.Message, resp.Message)
				assert.Equal(t, expResp.Number, resp.Number)
			}
		})
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
