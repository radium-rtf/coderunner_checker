package domain

import (
	"context"

	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

type CheckerSrv interface {
	GetSandbox(*checker.TestRequest) *SandboxInfo
	RunTest(context.Context, *SandboxInfo, *checker.TestCase) (*TestInfo, error)
}
