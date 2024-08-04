package domain

import (
	"context"

	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

type CheckerSrv interface {
	Closer
	RunTests(context.Context, *checker.TestRequest, []*checker.ArrayTestsRequest_TestCase) <-chan *TestResult
}
