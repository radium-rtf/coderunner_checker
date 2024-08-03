package tester

import (
	"context"
	"fmt"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	"time"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"

	coderunner "github.com/radium-rtf/coderunner_lib"
	"github.com/radium-rtf/coderunner_lib/file"
	"github.com/radium-rtf/coderunner_lib/limit"
	"github.com/radium-rtf/coderunner_lib/profile"
)

type Tester struct {
	limits   *limit.Limits
	cmd      string
	userCode string
	profile  profile.Profile
	rule     config.Rule
	client   *coderunner.Runner
}

func NewTester(client *coderunner.Runner, req *checker.TestRequest, rule config.Rule) *Tester {
	profile := profile.NewProfile(
		profile.Name(req.Lang),
		profile.Image(rule.Image),
	)

	cmd := fmt.Sprintf(`cat %s | %s %s`, domain.InputFile, rule.Launch, rule.Filename)

	limits := limit.NewLimits(
		limit.WithTimeout(time.Duration(req.Timeout)*time.Second),
		limit.WithMemoryInBytes(req.MemoryLimit),
	)

	return &Tester{
		profile:  profile,
		limits:   limits,
		cmd:      cmd,
		rule:     rule,
		userCode: req.UserCode,
		client:   client,
	}
}

func (s *Tester) RunTest(ctx context.Context, test *checker.TestCase) (*domain.TestInfo, error) {
	files := []file.File{
		file.NewFile(s.rule.Filename, file.StringContent(s.userCode)),
		file.NewFile(domain.InputFile, file.StringContent(test.Stdin)),
	}

	sandbox, err := s.client.NewSandbox(ctx, s.cmd, s.profile, s.limits, files)
	if err != nil {
		return nil, fmt.Errorf("cant create container: %v", err)
	}
	defer sandbox.Close()

	err = sandbox.Start()
	if err != nil {
		return nil, fmt.Errorf("cant start container: %v", err)
	}

	resInfo, err := sandbox.Info()
	if err != nil {
		return nil, fmt.Errorf("cant get result: %v", err)
	}

	testInfo := &domain.TestInfo{
		Info: resInfo,
		Test: test,
	}

	return testInfo, nil
}
