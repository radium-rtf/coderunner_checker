package domain

import (
	"github.com/radium-rtf/coderunner_lib/info"
)

type Test struct {
	Stdin, Stdout string
}

type TestInfo struct {
	info.Info
	Test    *Test
	Success bool
}

type TestResult struct {
	Info   *TestInfo
	Error  error
	Number int64
}
