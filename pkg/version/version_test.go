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

// Setting up tests and mocks
func BeforeEach() {
	version = "v1.0.0-dirty"
	gitCommitID = "9c474a"
}

func TestGetVersion(t *testing.T) {
	BeforeEach()

	v := GetVersion()
	if v != version {
		t.Fatalf("expected version %s but got %s", version, v)
	}
}

func TestGetGitCommitID(t *testing.T) {
	BeforeEach()

	id := GetGitCommitID()
	if id != gitCommitID {
		t.Fatalf("expected gitCommitID %s but got %s", version, id)
	}
}

func TestNewVersionInfo(t *testing.T) {
	BeforeEach()

	info := NewVersionInfo()
	expected := &VersionInfo{
		Version:     version,
		GitCommitID: gitCommitID,
	}
	if info == nil || info.Version != expected.Version || info.GitCommitID != expected.GitCommitID {
		t.Fatalf("expected versionInfo %+v but got %+v", expected, info)
	}
}
