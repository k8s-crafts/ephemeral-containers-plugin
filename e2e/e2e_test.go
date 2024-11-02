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

	"github.com/k8s-crafts/ephemeral-containers-plugin/e2e/testutils"
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
		JustBeforeEach(func() {
			Expect(tr.RunDebugContainer(true)).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			tr.UnsetEnv("E2E_EDIT_ACTION")
			tr.UnsetEnv("container_name")
		})

		When("making no changes", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "")
			})

			It("should print a message", func() {
				actual, err := tr.RunPluginEditCmd(testutils.PodName)

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(fmt.Sprintf("Edit cancelled, no changes made for pod/%s\n", testutils.PodName)))
			})
		})

		When("modifying an existing ephemeral container", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "modify")
			})

			It("should fail with Forbidden error", func() {
				actual, err := tr.RunPluginEditCmd(testutils.PodName)

				Expect(err).To(HaveOccurred())
				Expect(actual).To(ContainSubstring(fmt.Sprintf("Pod \"%s\" is invalid: spec.ephemeralContainers: Forbidden: existing ephemeral containers \"%s\" may not be changed", testutils.PodName, testutils.EphContainerName)))
			})
		})

		When("deleting an existing ephemeral container", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "delete")
			})

			It("should fail with Forbidden error", func() {
				actual, err := tr.RunPluginEditCmd(testutils.PodName)

				Expect(err).To(HaveOccurred())
				Expect(actual).To(ContainSubstring(fmt.Sprintf("Pod \"%s\" is invalid: spec.ephemeralContainers: Forbidden: existing ephemeral containers \"%s\" may not be removed", testutils.PodName, testutils.EphContainerName)))
			})
		})

		When("adding a new ephemeral container", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "add")
				tr.SetEnv("container_name", "another-debugger")
			})

			It("should succeed with a message", func() {
				actual, err := tr.RunPluginEditCmd(testutils.PodName)

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(fmt.Sprintf("pod/%s successfully edited\n", testutils.TestPodName)))

				names, err := tr.ListEphemeralContainerNamesForTestPod()
				Expect(err).ToNot(HaveOccurred())

				Expect(names).To(ContainSubstring(testutils.EphContainerName))
				Expect(names).To(ContainSubstring("another-debugger"))
			})
		})
	})
})
