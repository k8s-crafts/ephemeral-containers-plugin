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

var (
	// When built, set flag -ldflag="-X k8s-crafts/ephemeral-containers-plugin/pkg/version.version=vX.Y.Z"
	version string = "v0.0.0-unset"

	// When built, set flag -ldflag="-X k8s-crafts/ephemeral-containers-plugin/pkg/version.gitCommitID=<commit-id>"
	gitCommitID string = ""
)

type VersionInfo struct {
	Version     string `json:"version,omitempty"`
	GitCommitID string `json:"gitCommitID,omitempty"`
}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{
		Version:     version,
		GitCommitID: gitCommitID,
	}
}
