// MIT License

// Copyright (c) 2024 k8s-crafts Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package e2e_test

import (
	"fmt"

	"github.com/k8s-crafts/ephemeral-containers-plugin/e2e/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("kubectl ephemeral-containers", func() {
	BeforeEach(func() {
		for _, ns := range tr.GetTestNamespaces() {
			By(fmt.Sprintf("creating and waiting for test pod in namespace %s", ns))

			Expect(tr.CreateTestPod(ns)).ToNot(HaveOccurred())
			Expect(tr.WaitForTestPodReady(ns)).ToNot(HaveOccurred())
		}
	})

	AfterEach(func() {
		for _, ns := range tr.GetTestNamespaces() {
			By(fmt.Sprintf("deleting test pod in namespace %s", ns))

			Expect(tr.DeleteTestPod(ns)).ToNot(HaveOccurred())
		}
	})

	Context("list", func() {
		commonTests := func(allNamespace bool) {
			var namespace string

			BeforeEach(func() {
				// Leave namespace empty if in all-namespace mode
				if !allNamespace {
					namespace = tr.Kubectl.Namespace
				}
			})

			When("there is no ephemeral container", func() {
				It("should return empty message", func() {
					By("running kubectl ephemeral-containers list")

					actual, err := tr.RunPluginListCmd("", namespace)
					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(tr.NewListEmptyMessage(namespace)))
				})
			})

			When("there are ephemeral containers", func() {
				JustBeforeEach(func() {
					for _, ns := range tr.GetTestNamespaces() {
						By(fmt.Sprintf("adding an ephemeral container  to test pod in namespace %s", ns))

						Expect(tr.RunDebugContainerForTestPod(ns, testutils.EphContainerName, true)).ToNot(HaveOccurred())
					}
				})

				DescribeTable("should list in expected format", func(format string) {
					By("running kubectl ephemeral-containers list")

					actual, err := tr.RunPluginListCmd(format, namespace)
					Expect(err).ToNot(HaveOccurred())
					Expect(actual).To(Equal(tr.NewListOutput(format, namespace)))
				},
					Entry("in Table", ""),
					Entry("in JSON", "json"),
					Entry("in YAML", "yaml"),
				)
			})
		}

		Context("in a namespace", func() {
			commonTests(false)
		})

		Context("in all namespaces", func() {
			commonTests(true)
		})
	})

	Context("help", func() {
		DescribeTable("should print help message for command", func(subCmd string) {
			By("running kubectl ephemeral-containers help")

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
			for _, ns := range tr.GetTestNamespaces() {
				By(fmt.Sprintf("adding an ephemeral container  to test pod in namespace %s", ns))

				Expect(tr.RunDebugContainerForTestPod(ns, testutils.EphContainerName, true)).ToNot(HaveOccurred())
			}
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
				By("running kubectl ephemeral-containers edit")

				actual, err := tr.RunPluginEditCmd(tr.Kubectl.Namespace, testutils.TestPodName)

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(fmt.Sprintf("Edit cancelled, no changes made for pod/%s\n", testutils.TestPodName)))
			})
		})

		When("modifying an existing ephemeral container", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "modify")
			})

			It("should fail with Forbidden error", func() {
				By("running kubectl ephemeral-containers edit")

				actual, err := tr.RunPluginEditCmd(tr.Kubectl.Namespace, testutils.TestPodName)

				Expect(err).To(HaveOccurred())
				Expect(actual).To(ContainSubstring(fmt.Sprintf("Pod \"%s\" is invalid: spec.ephemeralContainers: Forbidden: existing ephemeral containers \"%s\" may not be changed", testutils.TestPodName, testutils.EphContainerName)))
			})
		})

		When("deleting an existing ephemeral container", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "delete")
			})

			It("should fail with Forbidden error", func() {
				By("running kubectl ephemeral-containers edit")

				actual, err := tr.RunPluginEditCmd(tr.Kubectl.Namespace, testutils.TestPodName)

				Expect(err).To(HaveOccurred())
				Expect(actual).To(ContainSubstring(fmt.Sprintf("Pod \"%s\" is invalid: spec.ephemeralContainers: Forbidden: existing ephemeral containers \"%s\" may not be removed", testutils.TestPodName, testutils.EphContainerName)))
			})
		})

		When("adding a new ephemeral container", func() {
			BeforeEach(func() {
				tr.SetEnv("E2E_EDIT_ACTION", "add")
				tr.SetEnv("container_name", "another-debugger")
			})

			It("should succeed with a message", func() {
				By("running kubectl ephemeral-containers edit")

				actual, err := tr.RunPluginEditCmd(tr.Kubectl.Namespace, testutils.TestPodName)

				Expect(err).ToNot(HaveOccurred())
				Expect(actual).To(Equal(fmt.Sprintf("pod/%s successfully edited\n", testutils.TestPodName)))

				names, err := tr.ListEphemeralContainerNamesForTestPod(tr.Kubectl.Namespace)
				Expect(err).ToNot(HaveOccurred())

				Expect(names).To(ContainSubstring(testutils.EphContainerName))
				Expect(names).To(ContainSubstring("another-debugger"))
			})
		})
	})
})
