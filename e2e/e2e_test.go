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

package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("kubectl ephemeral-containers", func() {
	BeforeEach(func() {
		Expect(tr.CreateTestPod()).ToNot(HaveOccurred())
		Expect(tr.WaitForPodReady()).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(tr.DeleteTestPod()).ToNot(HaveOccurred())
	})

	Context("list", func() {
		When("there is no ephemeral container", func() {
			It("should return empty message", func() {
				actual, err := tr.RunPluginListCmd("")
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(fmt.Sprintf("No pods with ephemeral containers found in namespace %s\n", tr.Kubectl.Namespace)))
			})
		})
		When("there are ephemeral containers", func() {
			DescribeTable("should list in expected format", func(format string) {
				Expect(tr.RunDebugContainer(true)).ToNot(HaveOccurred())

				actual, err := tr.RunPluginListCmd(format)
				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(tr.NewListOutput(format)))
			},
				Entry("in Table", ""),
				Entry("in JSON", "json"),
				Entry("in YAML", "yaml"),
			)
		})

	})

	Context("help", func() {
		DescribeTable("should print help message for command", func(subCmd string) {
			actual, err := tr.RunPluginHelpCmd(subCmd)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).ToNot(BeEmpty())
			Expect(actual).To(ContainSubstring("Usage:"))
			Expect(actual).To(ContainSubstring("Flags:"))

			if len(subCmd) == 0 {
				Expect(actual).To(ContainSubstring("Available Commands:"))
			} else {
				Expect(actual).To(ContainSubstring("Global Flags:"))
			}
		},
			Entry("list", "list"),
			Entry("edit", "edit"),
			Entry("version", "version"),
			Entry("root", ""),
		)
	})

	Context("edit", func() {
		// TODO: Add e2e tests for edit command
	})
})
