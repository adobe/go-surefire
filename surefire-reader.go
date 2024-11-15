/*
Copyright 2023 Adobe. All rights reserved.
This file is licensed to you under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License. You may obtain a copy
of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
OF ANY KIND, either express or implied. See the License for the specific language
governing permissions and limitations under the License.
*/

package surefire

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

type JUnitReportsReader struct {
	labeler Labeler
}

func (b *JUnitReportsReader) FromReportFiles(surefireReportFiles []string) (TestResults, error) {
	surefireSuites, err := parseSurefireReports(surefireReportFiles)
	if err != nil {
		return nil, err
	}

	return b.FromJUnitRepresentation(surefireSuites), nil
}

func (b *JUnitReportsReader) FromJUnitRepresentation(surefireSuites []surefireTestsuite) TestResults {
	testResults := testResults{}

	for _, surefireSuite := range surefireSuites {
		testSuite := testSuite{
			name:      surefireSuite.Name,
			filename:  surefireSuite.Filename,
			time:      surefireSuite.Time,
			testcases: make([]TestCase, 0),
			successes: 0,
			failures:  0,
			errors:    0,
			skipped:   0,
			labels:    make([]string, 0),
		}

		for _, surefireTestCase := range surefireSuite.Testcases {
			var issue *Issue
			_failure := optionalTestProblem(surefireTestCase.Failure)
			_error := optionalTestProblem(surefireTestCase.Error)

			var status Status
			var skipped *Skipped

			_skipped := surefireTestCase.Skipped

			if _skipped != nil {
				testSuite.skipped++
				skipped = &Skipped{Message: _skipped.Message}
				status = Skip
			} else if _failure != nil {
				testSuite.failures++
				issue = _failure
				status = Failure
			} else if _error != nil {
				testSuite.errors++
				issue = _error
				status = Error
			} else {
				testSuite.successes++
				status = Success
			}

			if amountOf(surefireTestCase.FlakyError)+amountOf(surefireTestCase.FlakyFailure) > 0 {
				status = Flaky
			}

			testCase := TestCase{
				Name:      surefireTestCase.Name,
				Classname: surefireTestCase.Classname,
				Fullname:  surefireTestCase.Classname + "." + surefireTestCase.Name,
				Suite:     &testSuite,

				Time:              surefireSuite.Time,
				Issue:             issue,
				Skipped:           skipped,
				RerunErrors:       toReRunIssues(surefireTestCase.ReRunErrors),
				AmountRerunErrors: amountOf(surefireTestCase.ReRunErrors),

				RerunFailures:       toReRunIssues(surefireTestCase.ReRunFailures),
				AmountRerunFailures: amountOf(surefireTestCase.ReRunFailures),

				FlakyErrors:       toReRunIssues(surefireTestCase.FlakyError),
				AmountFlakyErrors: amountOf(surefireTestCase.FlakyError),

				FlakyFailures:       toReRunIssues(surefireTestCase.FlakyFailure),
				AmountFlakyFailures: amountOf(surefireTestCase.FlakyFailure),
				Status:              status,
			}

			testSuite.testcases = append(testSuite.testcases, testCase)
		}
		if b.labeler != nil {
			testSuite.labels = b.labeler(&testSuite)
		}
		testResults.append(&testSuite)
	}

	return &testResults
}

type JUnitReportsReaderBuilder struct {
	JUnitReportsReader JUnitReportsReader
}

func NewJUnitReportsReaderBuilder() *JUnitReportsReaderBuilder {
	return &JUnitReportsReaderBuilder{
		JUnitReportsReader: JUnitReportsReader{},
	}
}

func (b *JUnitReportsReaderBuilder) WithLabeler(labeler Labeler) *JUnitReportsReaderBuilder {
	b.JUnitReportsReader.labeler = labeler
	return b
}

func (b *JUnitReportsReaderBuilder) Build() *JUnitReportsReader {
	return &b.JUnitReportsReader
}

func parseSurefireReports(surefireReportFiles []string) ([]surefireTestsuite, error) {
	testsuites := make([]surefireTestsuite, 0)
	errors := make([]error, 0)
	var wg sync.WaitGroup
	reportMutex := sync.Mutex{}
	errorMutex := sync.Mutex{}
	wg.Add(len(surefireReportFiles))

	for _, file := range surefireReportFiles {
		go func(file string) {
			xmlFile, openFileError := os.Open(file)
			suite, readReportError := readReport(xmlFile)
			suite.Filename = file
			closeFileError := xmlFile.Close()
			if openFileError != nil || readReportError != nil || closeFileError != nil {
				errorMutex.Lock()
				if openFileError != nil {
					errors = append(errors, openFileError)
				}
				if readReportError != nil {
					errors = append(errors, readReportError)
				}
				if closeFileError != nil {
					errors = append(errors, closeFileError)
				}

				wg.Done()
				errorMutex.Unlock()
				return
			}
			reportMutex.Lock()
			if suite.Name != "" {
				testsuites = append(testsuites, suite)
			}
			wg.Done()
			reportMutex.Unlock()
		}(file)
	}
	wg.Wait()
	if len(errors) > 0 {
		slog.Error("error reading test report files", "errors", errors)
		return nil, fmt.Errorf("one or more error occured when reading test report files %s", errors)
	}
	return testsuites, nil
}

// readReport parses xml content from given reader and returns a Testsuite
func readReport(reader io.Reader) (surefireTestsuite, error) {
	var testsuite surefireTestsuite

	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(&testsuite)
	if err != nil {
		return surefireTestsuite{}, fmt.Errorf("error decoding XML: %s", err)

	}

	return testsuite, nil
}
