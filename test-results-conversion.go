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

func optionalTestProblem(p *surefireProblem) *Issue {
	if p == nil {
		return nil
	}

	return &Issue{Message: p.Message, Detail: p.Data}
}

func toReRunIssues(runs []surefireRerun) []RerunIssue {
	if runs == nil || len(runs) == 0 {
		return nil
	}

	issues := make([]RerunIssue, len(runs))
	for i, r := range runs {
		issues[i] = RerunIssue{
			Message:     r.Message,
			Stacktrace:  r.Stacktrace,
			SystemOut:   r.SystemOut,
			SystemError: r.SystemError,
		}
	}

	return issues
}

func amountOf(runs []surefireRerun) int {
	if runs == nil {
		return 0
	}

	return len(runs)
}
