package domain

import (
	"github.com/radium-rtf/coderunner_checker/internal/config"
	coderunner "github.com/radium-rtf/coderunner_lib"
	"github.com/radium-rtf/coderunner_lib/limit"
	"github.com/radium-rtf/coderunner_lib/profile"
)

const (
	InputFile = "input.txt"
)

type SandboxInfo struct {
	Limits  *limit.Limits
	Cmd     string
	Code    string
	Profile profile.Profile
	Rule    config.Rule
	Client  *coderunner.Runner
}
