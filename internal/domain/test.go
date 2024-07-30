package domain

import (
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/radium-rtf/coderunner_lib/info"
)

type TestInfo struct {
	info.Info
	Test *checker.TestCase
}
