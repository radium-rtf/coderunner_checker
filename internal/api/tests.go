package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/radium-rtf/coderunner_checker/internal/domain"
	"github.com/radium-rtf/coderunner_checker/pkg/api/checker/v1"
)

func getTestsDTO(testCases []*checker.ArrayTestsRequest_TestCase) []*domain.Test {
	tests := make([]*domain.Test, 0, len(testCases))

	for _, testCase := range testCases {
		test := &domain.Test{
			Stdin:  testCase.Stdin,
			Stdout: testCase.Stdout,
		}
		tests = append(tests, test)
	}

	return tests
}

func getTestsFromFile(url string) ([]*domain.Test, error) {
	buf, err := getDataFromResponse(url)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return nil, err
	}

	testsRaw, err := getRawTests(zipReader.File)
	if err != nil {
		return nil, err
	}

	tests, err := getValidTests(testsRaw)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

func getDataFromResponse(url string) (*bytes.Buffer, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func getRawTests(files []*zip.File) (map[int]*domain.Test, error) {
	filesCount := len(files) / 2
	testsRaw := make(map[int]*domain.Test, filesCount)

	for _, file := range files {
		isInput, isOutput := strings.HasSuffix(file.Name, "input.txt"), strings.HasSuffix(file.Name, "output.txt")
		if !isInput && !isOutput {
			continue
		}

		dir := filepath.Dir(file.Name)
		testNumber, err := strconv.Atoi(dir)
		if err != nil {
			continue
		}

		content, err := getFileContent(file)
		if err != nil {
			return nil, err
		}

		if testsRaw[testNumber] == nil {
			testsRaw[testNumber] = &domain.Test{}
		}

		test := testsRaw[testNumber]
		if isInput {
			test.Stdin = string(content)
			continue
		}
		test.Stdout = string(content)
	}

	return testsRaw, nil
}

func getFileContent(file *zip.File) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	err = f.Close()
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getValidTests(testsRaw map[int]*domain.Test) ([]*domain.Test, error) {
	tests := make([]*domain.Test, 0, len(testsRaw))
	for i := 0; i < len(testsRaw); i++ {
		test, ok := testsRaw[i+1]
		if !ok {
			return nil, fmt.Errorf("no test with number %d", i+1)
		}
		if test.Stdin == "" || test.Stdout == "" {
			return nil, fmt.Errorf("no stdin or stdout in test")
		}
		tests = append(tests, test)
	}
	return tests, nil
}
