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
)

// surefireTestsuite encapsulates the data from a single test suite
type surefireTestsuite struct {
	Suite     xml.Name           `xml:"testsuite"`
	Name      string             `xml:"name,attr"`
	Time      float64            `xml:"time,attr"`
	Tests     int                `xml:"tests,attr"`
	Errors    int                `xml:"errors,attr"`
	Skipped   int                `xml:"skipped,attr"`
	Failures  int                `xml:"failures,attr"`
	Testcases []surefireTestcase `xml:"testcase"`
	Filename  string
}

// surefireTestcase encapsulates the data from a single test case
type surefireTestcase struct {
	Name          string           `xml:"name,attr"`
	Classname     string           `xml:"classname,attr"`
	Time          float64          `xml:"time,attr"`
	Skipped       *surefireSkipped `xml:"skipped"`
	Failure       *surefireProblem `xml:"failure"`
	Error         *surefireProblem `xml:"error"`
	ReRunErrors   []surefireRerun  `xml:"rerunError"`
	ReRunFailures []surefireRerun  `xml:"rerunFailure"`
	FlakyError    []surefireRerun  `xml:"flakyError"`
	FlakyFailure  []surefireRerun  `xml:"flakyFailure"`
}

// surefireSkipped is present if the referencing test case was skipped
type surefireSkipped struct {
	Message string `xml:"message,attr"`
}

// surefireProblem is present if the referencing test case failed or errored
type surefireProblem struct {
	Message string `xml:"message,attr"`
	Data    string `xml:",chardata"`
}

// surefireRerun represents errors, failures and flakes from re-runs
type surefireRerun struct {
	Message     string `xml:"message,attr"`
	Stacktrace  string `xml:"stackTrace"`
	SystemOut   string `xml:"system-out"`
	SystemError string `xml:"system-err"`
}
