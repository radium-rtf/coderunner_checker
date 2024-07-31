package checkergrpc

import (
	"github.com/radium-rtf/coderunner_checker/internal/domain"
	checkerutils "github.com/radium-rtf/coderunner_checker/internal/utils/checker"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type CheckerAPI struct {
	checker.UnimplementedCheckerServer

	checkerSrv domain.CheckerSrv
}

func Register(server *grpc.Server, checkerSrv domain.CheckerSrv) {
	checker.RegisterCheckerServer(server, &CheckerAPI{checkerSrv: checkerSrv})
}

func (c *CheckerAPI) Check(in *checker.ArrayTestsRequest, stream checker.Checker_CheckServer) error {
	sandboxInfo := c.checkerSrv.GetSandbox(in.Request)

	for i, test := range in.Tests {
		testInfo, err := c.checkerSrv.RunTest(stream.Context(), sandboxInfo, test)
		if err != nil {
			return status.Errorf(codes.Internal, "test %d: %v", i+1, err)
		}

		res := &checker.TestResponse{
			Number:   int64(i + 1),
			Duration: durationpb.New(testInfo.Time.Diff()),
		}
		checkerutils.SetResponseInfo(res, testInfo, in.Request.FullInfoWa)

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Internal, "test %d: %v", i+1, err)
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
