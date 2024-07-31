package api

import (
	"fmt"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
	"github.com/radium-rtf/coderunner_lib/info"
	"google.golang.org/protobuf/types/known/durationpb"
)

func GetResponse(result *domain.TestResult, fullInfoWA bool) *checker.TestResponse {
	if result.Info.Success {
		return getResponseSuccess(result)
	}

	switch result.Info.Status {
	case info.StatusOK:
		return getResponseWA(result, fullInfoWA)
	case info.StatusUnknown:
		return getResponseUnknown(result)
	case info.StatusTimeout:
		return getResponseTimeout(result)
	case info.StatusOOMKilled:
		return getResponseOOMKilled(result)
	default:
		return nil
	}
}

func getResponseSuccess(result *domain.TestResult) *checker.TestResponse {
	return &checker.TestResponse{
		Status: checker.Status_STATUS_SUCCESS,
		Message: &checker.TestResponse_Text{
			Text: "success",
		},
		Duration: durationpb.New(result.Info.Time.Diff()),
		Number: result.Number,
	}
}

func getResponseWA(result *domain.TestResult, fullInfoWA bool) *checker.TestResponse {
	res := &checker.TestResponse{
		Status: checker.Status_STATUS_WRONG_ANSWER,
		Duration: durationpb.New(result.Info.Time.Diff()),
		Number: result.Number,
	}

	if fullInfoWA {
		res.Message = &checker.TestResponse_WrongAnswer_{
			WrongAnswer: &checker.TestResponse_WrongAnswer{
				Input:    result.Info.Test.Stdin,
				Actual:   result.Info.Logs.StdOut,
				Expected: result.Info.Test.Stdout,
			},
		}
	} else {
		res.Message = &checker.TestResponse_Text{
			Text: fmt.Sprintf("wrong answer on input:\n%s", result.Info.Test.Stdin),
		}
	}

	return res
}

func getResponseUnknown(result *domain.TestResult) *checker.TestResponse {
	return &checker.TestResponse{
		Status: checker.Status_STATUS_UNKNOWN,
		Message: &checker.TestResponse_Text{
			Text: result.Info.Logs.StdErr,
		},
		Duration: durationpb.New(result.Info.Time.Diff()),
		Number: result.Number,
	}
}

func getResponseTimeout(result *domain.TestResult) *checker.TestResponse {
	return &checker.TestResponse{
		Status: checker.Status_STATUS_TIMEOUT,
		Message: &checker.TestResponse_Text{
			Text: "time limit",
		},
		Duration: durationpb.New(result.Info.Time.Diff()),
		Number: result.Number,
	}
}

func getResponseOOMKilled(result *domain.TestResult) *checker.TestResponse {
	return &checker.TestResponse{
		Status: checker.Status_STATUS_OOM_KILLED,
		Message: &checker.TestResponse_Text{
			Text: "memory limit",
		},
		Duration: durationpb.New(result.Info.Time.Diff()),
		Number: result.Number,
	}
}
