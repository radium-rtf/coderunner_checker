package grpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"

	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

const (
	code = `
print(int(input())*2)
`
)

func TestCheckSuccess(t *testing.T) {
	in := &checkergrpc.ArrayTestsRequest{
		Tests: []*checkergrpc.TestCase{
			{
				Stdin:  "1",
				Stdout: "2\n",
			},
			{
				Stdin:  "2",
				Stdout: "4\n",
			},
		},
		Request: &checkergrpc.TestRequest{
			Lang:        "python",
			UserCode:    code,
			Timeout:     10,
			MemoryLimit: 1024 * 1024 * 6,
			FullInfoWa:  true,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cfg := getConfig()

	server, err := runServer(ctx, cfg)
	require.NoError(t, err)

	t.Cleanup(func() {
		cancel()
		server.Wait()
	})

	conn, err := getConnection()
	require.NoError(t, err)
	defer conn.Close()

	checker := checkergrpc.NewCheckerClient(conn)
	client, err := checker.Check(ctx, in)
	require.NoError(t, err)

	for i := 0; i < len(in.Tests); i++ {
		resp, err := client.Recv()
		require.NoError(t, err)
		require.Equal(t, checkergrpc.Status_STATUS_SUCCESS, resp.Status)
	}
}

// func TestCheckURL(t *testing.T) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	cfg := getConfig()

// 	go runServer(ctx, cfg)

// 	conn := getConnection(t)
// 	defer conn.Close()

// 	checker := checkergrpc.NewCheckerClient(conn)

// 	in := &checkergrpc.FileTestsRequest{
// 		Url: testsUrl,
// 		Request: &checkergrpc.TestRequest{
// 			Lang:        "python",
// 			UserCode:    code,
// 			Timeout:     1,
// 			MemoryLimit: 1024 * 1024 * 6,
// 		},
// 	}
// 	_, err := checker.CheckURL(ctx, in)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
