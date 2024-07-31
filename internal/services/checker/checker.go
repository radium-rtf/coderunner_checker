package checkersrv

import (
	"context"
	"fmt"
	"time"

	coderunner "github.com/radium-rtf/coderunner_lib"
	libConfig "github.com/radium-rtf/coderunner_lib/config"
	"github.com/radium-rtf/coderunner_lib/file"
	"github.com/radium-rtf/coderunner_lib/limit"
	"github.com/radium-rtf/coderunner_lib/profile"

	"github.com/docker/docker/client"
	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

type CheckerService struct {
	client *coderunner.Runner
	rules  domain.Rules
}

func NewCheckerSrv(cfg domain.SandboxConfig) (*CheckerService, error) {
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

func (c *CheckerService) GetSandbox(req *checker.TestRequest) *domain.SandboxInfo {
	rule := c.rules[req.Lang]

	profile := profile.NewProfile(
		profile.Name(req.Lang),
		profile.Image(rule.Image),
	)

	cmd := fmt.Sprintf(`cat %s | %s %s`, domain.InputFile, rule.Launch, rule.Filename)

	limits := limit.NewLimits(
		limit.WithTimeout(time.Duration(req.Timeout)*time.Second),
		limit.WithMemoryInBytes(req.MemoryLimit),
	)

	return &domain.SandboxInfo{
		Profile:  profile,
		Limits:   limits,
		Cmd:      cmd,
		Rule:     rule,
		UserCode: req.UserCode,
		Client:   c.client,
	}
}

func (c *CheckerService) RunTest(ctx context.Context, sandboxInfo *domain.SandboxInfo, test *checker.TestCase) (*domain.TestInfo, error) {
	files := []file.File{
		file.NewFile(sandboxInfo.Rule.Filename, file.StringContent(sandboxInfo.UserCode)),
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
		Info: resInfo,
		Test: test,
	}

	return testInfo, nil
}
