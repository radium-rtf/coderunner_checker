package api

import (
	"context"
	"fmt"
	"io"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	CheckerSrv interface {
		io.Closer
		RunTests(context.Context, *checker.TestRequest, []*domain.Test) (<-chan *domain.TestResult, error)
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
	err := in.ValidateAll()
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	tests := getTestsDTO(in.Tests)

	results, err := c.checkerSrv.RunTests(stream.Context(), in.Request, tests)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

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
	err := in.ValidateAll()
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	tests, err := getTestsFromFile(in.Url)
	if err != nil {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("bad file url: %v", err))
	}

	results, err := c.checkerSrv.RunTests(stream.Context(), in.Request, tests)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

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
