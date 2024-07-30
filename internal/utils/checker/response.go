package checkerutils

import (
	"fmt"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/radium-rtf/coderunner_lib/info"
)

// This func sets to res status and message with testInfo.
func SetResponseInfo(res *checker.TestResponse, testInfo *domain.TestInfo, fullInfoWA bool) {
	if testInfo.Status == info.StatusOK && testInfo.Logs.StdOut == testInfo.Test.Stdout {
		setResponseSuccess(res)
		return
	}
	switch testInfo.Status {
	case info.StatusOK:
		setResponseWA(res, testInfo, fullInfoWA)
	case info.StatusUnknown:
		setResponseUnknown(res, testInfo)
	case info.StatusTimeout:
		setResponseTimeout(res)
	case info.StatusOOMKilled:
		setResponseOOMKilled(res)
	}
}

func setResponseSuccess(res *checker.TestResponse) {
	res.Status = checker.Status_STATUS_SUCCESS
	res.Message = &checker.TestResponse_Text{
		Text: "success",
	}
}

func setResponseWA(res *checker.TestResponse, testInfo *domain.TestInfo, fullInfoWA bool) {
	res.Status = checker.Status_STATUS_WRONG_ANSWER
	if fullInfoWA {
		res.Message = &checker.TestResponse_WrongAnswer_{
			WrongAnswer: &checker.TestResponse_WrongAnswer{
				Input:    testInfo.Test.Stdin,
				Actual:   testInfo.Logs.StdOut,
				Expected: testInfo.Test.Stdout,
			},
		}
		return
	}
	res.Message = &checker.TestResponse_Text{
		Text: fmt.Sprintf("wrong answer on input:\n%s", testInfo.Test.Stdin),
	}
}

func setResponseUnknown(res *checker.TestResponse, testInfo *domain.TestInfo) {
	res.Status = checker.Status_STATUS_ERROR
	res.Message = &checker.TestResponse_Text{
		Text: testInfo.Logs.StdErr,
	}
}

func setResponseTimeout(res *checker.TestResponse) {
	res.Status = checker.Status_STATUS_TIMEOUT
	res.Message = &checker.TestResponse_Text{
		Text: "time limit",
	}
}

func setResponseOOMKilled(res *checker.TestResponse) {
	res.Status = checker.Status_STATUS_OOM_KILLED
	res.Message = &checker.TestResponse_Text{
		Text: "memory limit",
	}
}
