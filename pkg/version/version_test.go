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

package version_test

import (
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	Context("when querying", func() {
		It("should return version information", func() {
			actual := version.NewVersionInfo()

			Expect(actual).ToNot(BeNil())
			// NOTE: Bump version for release
			Expect(actual.Version).To(Equal("v1.3.0-dev"))
		})
	})
})
