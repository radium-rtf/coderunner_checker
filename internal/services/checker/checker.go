package checkersrv

import (
	"context"
	"fmt"

	coderunner "github.com/radium-rtf/coderunner_lib"
	libConfig "github.com/radium-rtf/coderunner_lib/config"
	"github.com/radium-rtf/coderunner_lib/file"
	"github.com/radium-rtf/coderunner_lib/info"
	"github.com/radium-rtf/coderunner_lib/limit"
	"github.com/radium-rtf/coderunner_lib/profile"

	"github.com/docker/docker/client"
	"github.com/radium-rtf/coderunner_checker/internal/config"
	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

type CheckerService struct {
	client *coderunner.Runner
	rules  config.Rules
}

func NewCheckerSrv(cfg config.SandboxConfig) (*CheckerService, error) {
	libCfg := libConfig.NewConfig(
		libConfig.WithUID(cfg.UUID),
		libConfig.WithUser(cfg.User),
		libConfig.WithWorkDir(cfg.WorkDir),
	)

	client, err := coderunner.NewRunner(libCfg, client.WithHost(cfg.Host))
	if err != nil {
		return nil, err
	}

	checker := &CheckerService{
		client: client,
		rules:  cfg.Rules,
	}
	return checker, nil
}

func (c *CheckerService) Close() error {
	return c.client.Close()
}

func (c *CheckerService) RunTests(ctx context.Context, req *checker.TestRequest, tests []*checker.ArrayTestsRequest_TestCase) <-chan *domain.TestResult {
	results := make(chan *domain.TestResult)

	go func(ctx context.Context, req *checker.TestRequest, tests []*checker.ArrayTestsRequest_TestCase, results chan *domain.TestResult) {
		defer close(results)

		sandboxInfo := c.getSandboxInfo(req)

		for i, test := range tests {
			testInfo, err := c.runTest(ctx, sandboxInfo, test)

			res := &domain.TestResult{
				Info:   testInfo,
				Error:  err,
				Number: int64(i + 1),
			}
			results <- res

			if err != nil || !testInfo.Success {
				return
			}
		}
	}(ctx, req, tests, results)

	return results
}

func (c *CheckerService) getSandboxInfo(req *checker.TestRequest) *domain.SandboxInfo {
	rule := c.rules[req.Lang]

	profile := profile.NewProfile(
		profile.Name(req.Lang),
		profile.Image(rule.Image),
	)

	limits := limit.NewLimits(
		limit.WithTimeout(req.Timeout.AsDuration()),
		limit.WithMemoryInBytes(req.MemoryLimitBytes),
	)

	return &domain.SandboxInfo{
		Profile: profile,
		Limits:  limits,
		Cmd:     rule.Launch,
		Rule:    rule,
		Code:    req.Code,
		Client:  c.client,
	}
}

func (c *CheckerService) runTest(ctx context.Context, sandboxInfo *domain.SandboxInfo, test *checker.ArrayTestsRequest_TestCase) (*domain.TestInfo, error) {
	files := []file.File{
		file.NewFile(sandboxInfo.Rule.Filename, file.StringContent(sandboxInfo.Code)),
		file.NewFile(domain.InputFile, file.StringContent(test.Stdin)),
	}

	sandbox, err := c.client.NewSandbox(ctx, sandboxInfo.Cmd, sandboxInfo.Profile, sandboxInfo.Limits, files)
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
		Info:    resInfo,
		Test:    test,
		Success: resInfo.Status == info.StatusOK && resInfo.Logs.StdOut == test.Stdout,
	}

	return testInfo, nil
}
