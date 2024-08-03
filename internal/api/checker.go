package api

import (
	"github.com/docker/docker/client"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	coderunner "github.com/radium-rtf/coderunner_lib"
	libConfig "github.com/radium-rtf/coderunner_lib/config"

	"github.com/radium-rtf/coderunner_checker/internal/service/tester"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Checker struct {
	checker.UnimplementedCheckerServer

	client *coderunner.Runner
	rules  config.Rules
}

func NewChecker(cfg config.SandboxConfig) (*Checker, error) {
	libCfg := libConfig.NewConfig(
		libConfig.WithUID(cfg.UUID),
		libConfig.WithUser(cfg.User),
		libConfig.WithWorkDir(cfg.WorkDir),
	)

	client, err := coderunner.NewRunner(libCfg, client.WithHost(cfg.Host))
	if err != nil {
		return nil, err
	}

	checker := &Checker{
		client: client,
		rules:  cfg.Rules,
	}
	return checker, nil
}

func (c *Checker) Check(in *checker.ArrayTestsRequest, stream checker.Checker_CheckServer) error {
	tester := tester.NewTester(c.client, in.Request, c.rules[in.Request.Lang])

	for i, test := range in.Tests {
		testInfo, err := tester.RunTest(stream.Context(), test)
		if err != nil {
			return status.Errorf(codes.Internal, "test %d: %v", i+1, err)
		}

		res := &checker.TestResponse{
			Number:   int64(i + 1),
			Duration: durationpb.New(testInfo.Time.Diff()),
		}
		setResponseInfo(res, testInfo, in.Request.FullInfoWa)

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Internal, "test %d: %v", i+1, err)
		}
	}

	return nil
}

func (c *Checker) CheckURL(in *checker.FileTestsRequest, stream checker.Checker_CheckURLServer) error {
	// req, err := http.NewRequest("GET", in.Url, nil)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "cant create request")
	// }

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, status.Errorf(codes.InvalidArgument, "can't send request to %v", in.Url)
	// }
	// defer resp.Body.Close()

	// _, err = io.Copy(os.Stdout, resp.Body)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "can't read response body")
	// }

	return status.Errorf(codes.Unimplemented, "method CheckURL not implemented")
}
