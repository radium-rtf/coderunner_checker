package checker

import (
	"context"
	"fmt"
	"testing"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	serverutils "github.com/radium-rtf/coderunner_checker/internal/utils/server"
	checkergrpc "github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddr = ":8081"

	testsUrl = "http://localhost:8082/"

	code = `
print(int(input())*2)
`
)

func getConfig() *domain.Config {
	return &domain.Config{
		ServerConfig: domain.ServerConfig{
			Env:     domain.TestEnv,
			Address: serverAddr,
		},
		SandboxConfig: domain.SandboxConfig{
			User:    "root",
			UUID:    0,
			WorkDir: "/sandbox",
			Host:    "unix:///var/run/docker.sock",
			Rules: domain.Rules{
				"python": domain.Rule{
					Filename: "main.py",
					Image: "python",
					Launch: "python3",
				},
			},
		},
	}
}

func runServer(ctx context.Context, cfg *domain.Config) {
	server := grpc.NewServer()

	checker, err := NewChecker(cfg.SandboxConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	checkergrpc.RegisterCheckerServer(server, checker)

	err = serverutils.Run(ctx, server, cfg.ServerConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getConnection(t *testing.T) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestCheckSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := getConfig()

	go runServer(ctx, cfg)

	conn := getConnection(t)
	defer conn.Close()

	checker := checkergrpc.NewCheckerClient(conn)

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
	client, err := checker.Check(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(in.Tests); i++ {
		resp, err := client.Recv()
		if err != nil {
			t.Fatal(err)
		}

		if resp.Status != checkergrpc.Status_STATUS_SUCCESS {
			t.Errorf("want: %v, but was: %v", checkergrpc.Status_STATUS_SUCCESS, resp)
		}
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
