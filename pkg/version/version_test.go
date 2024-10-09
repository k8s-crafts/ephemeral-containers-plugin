// Copyright 2024 k8s-crafts Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import "testing"

func beforeEach() {
	version = "v0.0.0-unset"
	gitCommitID = ""
}

func TestNewVersionInfo(t *testing.T) {
	tests := []struct {
		description string
		setup       func()
		expected    *VersionInfo
	}{
		{
			description: "should return default version info if unset",
			setup:       func() {},
			expected: &VersionInfo{
				Version:     "v0.0.0-unset",
				GitCommitID: "",
			},
		},
		{
			description: "should return correct version info if set during build",
			setup: func() {
				version = "v1.0.0"
				gitCommitID = "9c474a"
			},
			expected: &VersionInfo{
				Version:     "v1.0.0",
				GitCommitID: "9c474a",
			},
		},
	}

	for _, test := range tests {
		beforeEach()
		t.Run(test.description, func(t *testing.T) {
			test.setup()

			actual := NewVersionInfo()
			expected := test.expected
			if actual == nil || actual.Version != expected.Version || actual.GitCommitID != expected.GitCommitID {
				t.Fatalf("expected versionInfo %+v but received %+v", expected, actual)
			}
		})
	}

}
