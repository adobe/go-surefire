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
	"regexp"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestConvertNonIntersectedTestCasesResult(t *testing.T) {
	suites := []surefireTestsuite{
		{
			Name: "Failure-Suite",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name: "Test-1",
					Time: 1.0,
					Failure: &surefireProblem{
						Message: "Failure-1",
						Data:    "Failure-1-Data",
					},
				},
				{
					Name: "Test-2",
					Time: 1.0,
					Error: &surefireProblem{
						Message: "Error-1",
						Data:    "Error-1-Data",
					},
				},
			},
		},
		{
			Name: "Success-Suite",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name: "Test-3",
					Time: 1.0,
				},
			},
		},
		{
			Name: "Skipped-Suite",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name:    "Test-4",
					Time:    1.0,
					Skipped: &surefireSkipped{Message: "Skipped"},
				},
			},
		},
	}
	assert := a.New(t)
	testResult := NewJUnitReportsReaderBuilder().Build().FromJUnitRepresentation(suites)
	assert.Equal(3, len(testResult.TestSuites()))
	assert.Equal(1, testResult.Errors())
	assert.Equal(1, testResult.Failures())
	assert.Equal(1, testResult.Skipped())
	assert.Equal(4, testResult.Tests())
	assert.Equal(0, testResult.Flakes())

	failureSuite := suiteByName("Failure-Suite", testResult.TestSuites())
	assert.NotNil(failureSuite)
	assert.Equal(2, len(failureSuite.NonSuccessfulTestCases()))
	assert.Equal(1.0, failureSuite.Time())

	testCase1 := caseByName("Test-1", failureSuite.NonSuccessfulTestCases())
	assert.NotNil(testCase1)
	assert.Equal("Failure-Suite", testCase1.Suite.Name())
	assert.NotNil(testCase1.Issue)
	assert.Equal("Failure-1", testCase1.Issue.Message)
	assert.Equal("Failure-1-Data", testCase1.Issue.Detail)

	testCase2 := caseByName("Test-2", failureSuite.NonSuccessfulTestCases())
	assert.NotNil(testCase2)
	assert.Equal("Failure-Suite", testCase2.Suite.Name())
	assert.NotNil(testCase2.Issue)
	assert.Equal("Error-1", testCase2.Issue.Message)
	assert.Equal("Error-1-Data", testCase2.Issue.Detail)

	successSuite := suiteByName("Success-Suite", testResult.TestSuites())
	assert.NotNil(successSuite)
	assert.Equal(1, len(successSuite.SuccessfulTestCases()))
	assert.Equal(1.0, successSuite.Time())
	testCase3 := caseByName("Test-3", successSuite.SuccessfulTestCases())
	assert.NotNil(testCase3)

	skippedSuite := suiteByName("Skipped-Suite", testResult.TestSuites())
	assert.NotNil(skippedSuite)
	assert.Equal(1, len(skippedSuite.SkippedTestCases()))
	assert.Equal(1.0, skippedSuite.Time())

	testCase4 := caseByName("Test-4", skippedSuite.SkippedTestCases())
	assert.NotNil(testCase4)
	assert.Equal("Skipped-Suite", testCase4.Suite.Name())
	assert.NotNil(testCase4.Skipped)
	assert.Equal("Skipped", testCase4.Skipped.Message)
}

func TestConvertIntersectedTestCasesResult(t *testing.T) {
	suites := []surefireTestsuite{
		{
			Name: "Intersected-Suite",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name: "Test-1",
					Time: 1.0,
					Failure: &surefireProblem{
						Message: "Failure-1",
						Data:    "Failure-1-Data",
					},
				},
				{
					Name: "Test-2",
					Time: 1.0,
					Error: &surefireProblem{
						Message: "Error-1",
						Data:    "Error-1-Data",
					},
				},
				{
					Name: "Test-3",
					Time: 1.0,
				},
				{
					Name:    "Test-4",
					Time:    1.0,
					Skipped: &surefireSkipped{Message: "Skipped"},
				},
			},
		},
	}
	assert := a.New(t)
	testResult := NewJUnitReportsReaderBuilder().Build().FromJUnitRepresentation(suites)
	assert.Equal(1, len(testResult.TestSuites()))
	assert.Equal(1, testResult.Errors())
	assert.Equal(1, testResult.Failures())
	assert.Equal(1, testResult.Skipped())
	assert.Equal(4, testResult.Tests())
	assert.Equal(0, testResult.Flakes())

	suite := suiteByName("Intersected-Suite", testResult.TestSuites())
	assert.NotNil(suite)
	assert.Equal(2, len(suite.NonSuccessfulTestCases()))
	assert.Equal(1, len(suite.SuccessfulTestCases()))
	assert.Equal(1, len(suite.SkippedTestCases()))

	testCase1 := caseByName("Test-1", suite.NonSuccessfulTestCases())
	assert.NotNil(testCase1)
	assert.Equal("Intersected-Suite", testCase1.Suite.Name())
	assert.NotNil(testCase1.Issue)
	assert.Equal("Failure-1", testCase1.Issue.Message)
	assert.Equal("Failure-1-Data", testCase1.Issue.Detail)

	testCase2 := caseByName("Test-2", suite.NonSuccessfulTestCases())
	assert.NotNil(testCase2)
	assert.Equal("Intersected-Suite", testCase2.Suite.Name())
	assert.NotNil(testCase2.Issue)
	assert.Equal("Error-1", testCase2.Issue.Message)
	assert.Equal("Error-1-Data", testCase2.Issue.Detail)

	assert.Equal(1, len(suite.SuccessfulTestCases()))

	testCase3 := caseByName("Test-3", suite.SuccessfulTestCases())
	assert.Equal("Intersected-Suite", testCase3.Suite.Name())
	assert.NotNil(testCase3)

	assert.Equal(1, len(suite.SkippedTestCases()))

	testCase4 := caseByName("Test-4", suite.SkippedTestCases())
	assert.NotNil(testCase4)
	assert.Equal("Intersected-Suite", testCase4.Suite.Name())
	assert.NotNil(testCase4.Skipped)
	assert.Equal("Skipped", testCase4.Skipped.Message)
}

func TestFlakyRuns(t *testing.T) {
	suites := []surefireTestsuite{
		{
			Name: "Flaky-Suite",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name: "FlakyTest-1",
					Time: 1.0,
					ReRunErrors: []surefireRerun{
						{
							Message:     "Rerun-Error-1",
							Stacktrace:  "Rerun-Error-1-Stacktrace",
							SystemError: "Rerun-Error-1-SystemError",
							SystemOut:   "Rerun-Error-1-SystemOut",
						},
						{
							Message:     "Rerun-Error-2",
							Stacktrace:  "Rerun-Error-2-Stacktrace",
							SystemError: "Rerun-Error-2-SystemError",
							SystemOut:   "Rerun-Error-2-SystemOut",
						},
					},
					ReRunFailures: []surefireRerun{
						{
							Message:     "Rerun-Failure-1",
							Stacktrace:  "Rerun-Failure-1-Stacktrace",
							SystemError: "Rerun-Failure-1-SystemError",
							SystemOut:   "Rerun-Failure-1-SystemOut",
						},
						{
							Message:     "Rerun-Failure-2",
							Stacktrace:  "Rerun-Failure-2-Stacktrace",
							SystemError: "Rerun-Failure-2-SystemError",
							SystemOut:   "Rerun-Failure-2-SystemOut",
						},
					},
					FlakyError: []surefireRerun{
						{
							Message:     "FlakyError-1",
							Stacktrace:  "FlakyError-1-Stacktrace",
							SystemError: "FlakyError-1-SystemError",
							SystemOut:   "FlakyError-1-SystemOut",
						},
						{
							Message:     "FlakyError-2",
							Stacktrace:  "FlakyError-2-Stacktrace",
							SystemError: "FlakyError-2-SystemError",
							SystemOut:   "FlakyError-2-SystemOut",
						},
					},
					FlakyFailure: []surefireRerun{
						{
							Message:     "FlakyFailure-1",
							Stacktrace:  "FlakyFailure-1-Stacktrace",
							SystemError: "FlakyFailure-1-SystemError",
							SystemOut:   "FlakyFailure-1-SystemOut",
						},
						{
							Message:     "FlakyFailure-2",
							Stacktrace:  "FlakyFailure-2-Stacktrace",
							SystemError: "FlakyFailure-2-SystemError",
							SystemOut:   "FlakyFailure-2-SystemOut",
						},
					},
				},
			},
		},
	}
	assert := a.New(t)
	testResult := NewJUnitReportsReaderBuilder().Build().FromJUnitRepresentation(suites)
	assert.Equal(1, len(testResult.TestSuites()))
	assert.Equal(1, testResult.Tests())
	assert.Equal(1, testResult.Flakes())

	failureSuite := suiteByName("Flaky-Suite", testResult.TestSuites())
	assert.NotNil(failureSuite)
	assert.Equal(1, len(failureSuite.SuccessfulTestCases()))

	flakyTest := caseByName("FlakyTest-1", failureSuite.SuccessfulTestCases())
	assert.NotNil(flakyTest)

	assert.NotNil(flakyTest.RerunFailures)
	assert.Equal(2, flakyTest.AmountRerunFailures)
	for _, r := range flakyTest.RerunFailures {
		assert.NotNil(r)
		assert.Regexp("Rerun-Failure-[1-2]", r.Message)
		assert.Regexp("Rerun-Failure-[1-2]-Stacktrace", r.Stacktrace)
		assert.Regexp("Rerun-Failure-[1-2]-SystemError", r.SystemError)
		assert.Regexp("Rerun-Failure-[1-2]-SystemOut", r.SystemOut)
	}

	assert.NotNil(flakyTest.RerunErrors)
	assert.Equal(2, flakyTest.AmountRerunErrors)
	for _, r := range flakyTest.RerunErrors {
		assert.NotNil(r)
		assert.Regexp("Rerun-Error-[1-2]", r.Message)
		assert.Regexp("Rerun-Error-[1-2]-Stacktrace", r.Stacktrace)
		assert.Regexp("Rerun-Error-[1-2]-SystemError", r.SystemError)
		assert.Regexp("Rerun-Error-[1-2]-SystemOut", r.SystemOut)
	}

	assert.NotNil(flakyTest.FlakyFailures)
	assert.Equal(2, flakyTest.AmountFlakyFailures)
	for _, r := range flakyTest.FlakyFailures {
		assert.NotNil(r)
		assert.Regexp("FlakyFailure-[1-2]", r.Message)
		assert.Regexp("FlakyFailure-[1-2]-Stacktrace", r.Stacktrace)
		assert.Regexp("FlakyFailure-[1-2]-SystemError", r.SystemError)
		assert.Regexp("FlakyFailure-[1-2]-SystemOut", r.SystemOut)
	}

	assert.NotNil(flakyTest.FlakyErrors)
	assert.Equal(2, flakyTest.AmountFlakyErrors)
	for _, r := range flakyTest.FlakyErrors {
		assert.NotNil(r)
		assert.Regexp("FlakyError-[1-2]", r.Message)
		assert.Regexp("FlakyError-[1-2]-Stacktrace", r.Stacktrace)
		assert.Regexp("FlakyError-[1-2]-SystemError", r.SystemError)
		assert.Regexp("FlakyError-[1-2]-SystemOut", r.SystemOut)
	}
}

func TestAssignStaticLabels(t *testing.T) {
	suites := []surefireTestsuite{
		{
			Name: "Success-Suite",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name: "Test-3",
					Time: 1.0,
				},
			},
		},
	}

	assert := a.New(t)
	testResult := NewJUnitReportsReaderBuilder().WithLabeler(assignStaticLabeler).Build().FromJUnitRepresentation(suites)

	assert.NotNil(testResult)

	assert.NotNil(suiteByName("Success-Suite", testResult.TestSuites()).Labels())
	assert.Contains(suiteByName("Success-Suite", testResult.TestSuites()).Labels(), "myCategory")
}

func TestAssignMultipleLabelsByRegex(t *testing.T) {
	suites := []surefireTestsuite{
		{
			Name: "Success-SuiteIT",
			Time: 1.0,
			Testcases: []surefireTestcase{
				{
					Name: "Test-3",
					Time: 1.0,
				},
			},
		},
	}

	assert := a.New(t)
	testResult := NewJUnitReportsReaderBuilder().WithLabeler(regexLabeler).Build().FromJUnitRepresentation(suites)

	assert.NotNil(testResult)

	assert.NotNil(suiteByName("Success-SuiteIT", testResult.TestSuites()).Labels())
	assert.Contains(suiteByName("Success-SuiteIT", testResult.TestSuites()).Labels(), "ATest")
	assert.Contains(suiteByName("Success-SuiteIT", testResult.TestSuites()).Labels(), "ATest")
	assert.Contains(suiteByName("Success-SuiteIT", testResult.TestSuites()).Labels(), "Integration-Test")
}

func regexLabeler(testcase TestSuite) []string {
	cats := make([]string, 2)
	r, _ := regexp.Compile(".+IT")

	if r.MatchString(testcase.Name()) {
		cats = append(cats, "Integration-Test")
	}
	cats = append(cats, "ATest")

	return cats
}

func assignStaticLabeler(dummy TestSuite) []string {
	cats := []string{"myCategory"}
	return cats
}

func suiteByName(suitename string, suites []TestSuite) TestSuite {
	for _, s := range suites {
		if s.Name() == suitename {
			return s
		}
	}
	return nil
}

func caseByName(casename string, cases []TestCase) *TestCase {
	for _, c := range cases {
		if c.Name == casename {
			return &c
		}
	}
	return nil
}
