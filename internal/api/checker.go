package api

import (
	"context"
	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type (
	CheckerSrv interface {
		io.Closer
		RunTests(context.Context, *checker.TestRequest, []*checker.ArrayTestsRequest_TestCase) <-chan *domain.TestResult
	}

	CheckerAPI struct {
		checker.UnimplementedCheckerServer

		checkerSrv CheckerSrv
	}
)

func RegisterChecker(server grpc.ServiceRegistrar, checkerSrv CheckerSrv) {
	checker.RegisterCheckerServer(server, &CheckerAPI{checkerSrv: checkerSrv})
}

func (c *CheckerAPI) Check(in *checker.ArrayTestsRequest, stream checker.Checker_CheckServer) error {
	results := c.checkerSrv.RunTests(stream.Context(), in.Request, in.Tests)

	for result := range results {
		if result.Error != nil {
			return status.Errorf(codes.Internal, "test %d: %v", result.Number, result.Error)
		}

		res := GetResponse(result, in.Request.FullInfoWa)

		err := stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Internal, "test %d: %v", result.Number, err)
		}
	}

	return nil
}

func (c *CheckerAPI) CheckURL(in *checker.FileTestsRequest, stream checker.Checker_CheckURLServer) error {
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
