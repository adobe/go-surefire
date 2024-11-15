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

// TestResults aggregates all TestSuites being read from the surefire reports and expose statistics
type TestResults interface {
	// All being read from the surefire reports. Except for those with an empty name attribute
	TestSuites() []TestSuite

	// The amount of a all tests
	Tests() int

	// The amount of tests that were successful
	Successes() int

	// The amount of failing tests
	Failures() int

	// The amount of tests in error
	Errors() int

	// The amount of skipped tests
	Skipped() int

	// The amount flaky tests
	Flakes() int
}

// Implementation of TestResults
type testResults struct {
	tests     int
	successes int
	failures  int
	errors    int
	skipped   int
	flakes    int
	suites    []TestSuite
}

// TestSuite represents a set of TestCase and exposes statistics
type TestSuite interface {
	// Returns all test cases
	TestCases() []TestCase

	// Returns non successful test cases, either failing or in error
	NonSuccessfulTestCases() []TestCase

	// Returns successful test cases
	SuccessfulTestCases() []TestCase

	// Returns skipped test cases
	SkippedTestCases() []TestCase

	// Returns flaky test cases
	FlakyTestCases() []TestCase

	// The amount of successful tests in this suite
	Success() int

	// The amount of failing tests in this suite
	Failure() int

	// The amount of tests in this suite that are in error
	Error() int

	// The amount of skipped tests in this suite
	Skipped() int

	// Name of the suite
	Name() string

	// Filename from which the result was coming from
	Filename() string

	// The time this suite needs to run
	Time() float64

	// Labels the suite is assigned to
	Labels() []string
}

// implementation of TestSuite
type testSuite struct {
	testcases []TestCase
	name      string
	filename  string
	time      float64
	successes int
	failures  int
	errors    int
	skipped   int

	// Labels the suite is assigned to
	labels []string
}

// TestCase represents a single test run
type TestCase struct {
	// Name of the test case
	Name string

	// Backreference to the suite this test case belongs to
	Suite TestSuite

	// Status of this test case
	Status Status

	// The time this test case needs to run
	Time float64

	// Full qualified name of the test case
	Fullname string

	// Classname of this test case
	Classname string

	// Issue this test case has. If nil there was no issue with it
	Issue *Issue

	// The failures which appeared that leads to a re-run of this failing test
	RerunFailures []RerunIssue

	// The amount of failure which appeared that leads to a re-run of this test
	AmountRerunFailures int

	// The errors which appeared that leads to a re-run of this test in error
	RerunErrors []RerunIssue

	// The amount of errors which appeared that leads to a re-run of this test
	AmountRerunErrors int

	// The failures which appeared that leads to a re-run of this flaky test
	FlakyFailures []RerunIssue

	// The amount of failures which appeared that leads to a re-run of this flaky test
	AmountFlakyFailures int

	// The errors which appeared that leads to a re-run of this flaky test
	FlakyErrors []RerunIssue

	// The amount of errors which appeared that leads to a re-run of this flaky test
	AmountFlakyErrors int

	// Set for a skipped test, nil otherwise
	Skipped *Skipped
}

// Issue encapsulates a failure or error
type Issue struct {
	// Message for that issue
	Message string
	// Details for that issue
	Detail string
}

// RerunIssue encapsulates a rerun failure or error
type RerunIssue struct {
	// Message for that RerunIssue
	Message string
	// Stacktrace for that RerunIssue
	Stacktrace string
	// SystemOut for that RerunIssue
	SystemOut string
	// SystemError for that RerunIssue
	SystemError string
}

// Issue encapsulates the message of a skipped test
type Skipped struct {
	Message string
}

type Labeler func(TestSuite) []string

// Status represents the status of a test case
type Status string

const (
	Success Status = "success"
	Skip    Status = "skipped"
	Failure Status = "failure"
	Error   Status = "error"
	Flaky   Status = "flaky"
)

func (r *testResults) TestSuites() []TestSuite {
	return r.suites
}

func (t *testSuite) NonSuccessfulTestCases() []TestCase {
	return t.filterTestCases(func(testCase TestCase) bool {
		return testCase.Issue != nil
	})
}

func (t *testSuite) SuccessfulTestCases() []TestCase {
	return t.filterTestCases(func(testCase TestCase) bool {
		return testCase.Issue == nil && testCase.Skipped == nil
	})
}

func (t *testSuite) SkippedTestCases() []TestCase {
	return t.filterTestCases(func(testCase TestCase) bool {
		return testCase.Skipped != nil
	})
}

func (t *testSuite) FlakyTestCases() []TestCase {
	return t.filterTestCases(func(testCase TestCase) bool {
		return testCase.AmountFlakyErrors != 0 ||
			testCase.AmountFlakyFailures != 0
	})
}

func (t *testSuite) filterTestCases(predicate func(testCase TestCase) bool) []TestCase {
	_cases := make([]TestCase, 0)

	for _, testCase := range t.testcases {
		if predicate(testCase) {
			_cases = append(_cases, testCase)
		}
	}

	return _cases
}

func (r *testSuite) TestCases() []TestCase {
	return r.testcases
}

func (r *testSuite) Success() int {
	return r.successes
}

func (r *testSuite) Failure() int {
	return r.failures
}

func (r *testSuite) Error() int {
	return r.errors
}

func (r *testSuite) Skipped() int {
	return r.skipped
}

func (r *testSuite) Name() string {
	return r.name
}

func (r *testSuite) Time() float64 {
	return r.time
}

func (r *testSuite) Labels() []string {
	return r.labels
}

func (r *testSuite) Filename() string {
	return r.filename
}

func (r *testResults) Successes() int {
	return r.successes
}

func (r *testResults) Failures() int {
	return r.failures
}

func (r *testResults) Errors() int {
	return r.errors
}

func (r *testResults) Skipped() int {
	return r.skipped
}

func (r *testResults) Tests() int {
	return r.tests
}

func (r *testResults) Flakes() int {
	return r.flakes
}

func (r *testResults) append(suite *testSuite) {
	r.tests += len(suite.TestCases())
	r.successes += suite.Success()
	r.errors += suite.Error()
	r.failures += suite.Failure()
	r.skipped += suite.Skipped()
	r.flakes += len(suite.FlakyTestCases())
	r.suites = append(r.suites, suite)
}
