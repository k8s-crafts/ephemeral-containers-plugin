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

package cmd_test

import (
	"github.com/k8s-crafts/ephemeral-containers-plugin/cmd"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Cmd", func() {
	var t *test

	Context("when creating", func() {
		BeforeEach(func() {
			t = newTest()
		})

		Context("root command", func() {
			JustBeforeEach(func() {
				t.cmd = cmd.NewRootCmd()
			})

			It("should have basic configurations", func() {
				t.expectCmdBasics()
			})

			It("should have CommandDisplayNameAnnotatio annotation", func() {
				t.expectCmdAnnotation(cobra.CommandDisplayNameAnnotation, "kubectl ephemeral-containers")
			})

			It("should have subcommands", func() {
				subs := []string{"edit", "list", "version"}
				t.expectSubCommands(subs)
			})
		})

		Context("edit command", func() {
			JustBeforeEach(func() {
				t.cmd = cmd.NewEditCmd()
			})

			It("should have basic configurations", func() {
				t.expectCmdBasics()
			})

			Context("when given arguments", func() {
				It("should accept 1 argument", func() {
					err := t.cmd.Args(t.cmd, []string{"pod/name"})
					Expect(err).ToNot(HaveOccurred())
				})
				It("should accept 2 arguments", func() {
					err := t.cmd.Args(t.cmd, []string{"pods", "pod-name"})
					Expect(err).ToNot(HaveOccurred())
				})
				It("should fail otherwise", func() {
					err := t.cmd.Args(t.cmd, []string{})
					Expect(err).To(HaveOccurred())

					err = t.cmd.Args(t.cmd, []string{"pods", "pod-name", "another-one"})
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("list command", func() {
			JustBeforeEach(func() {
				t.cmd = cmd.NewListCmd()
			})

			It("should have basic configurations", func() {
				t.expectCmdBasics()
			})
		})

		Context("version command", func() {
			JustBeforeEach(func() {
				t.cmd = cmd.NewVersionCmd()
			})

			It("should have basic configurations", func() {
				t.expectCmdBasics()
			})
		})
	})
})

type test struct {
	*testInput
}

type testInput struct {
	cmd *cobra.Command
}

func newTest() *test {
	return &test{
		testInput: &testInput{},
	}
}

func (t *test) expectCmdAnnotation(key, value string) {
	val, ok := t.cmd.Annotations[key]

	Expect(ok).To(BeTrue())
	Expect(val).To(Equal(value))
}

func (t *test) expectCmdBasics() {
	Expect(t.cmd).ToNot(BeNil())
	Expect(t.cmd.Use).ToNot(BeEmpty())
	Expect(t.cmd.Short).ToNot(BeEmpty())
	Expect(t.cmd.Long).ToNot(BeEmpty())
}

func (t *test) expectSubCommands(subUses []string) {
	subs := t.cmd.Commands()

	Expect(subs).ToNot(BeNil())
	Expect(subs).To(HaveLen(len(subUses)))

	for idx, sub := range subs {
		Expect(sub).ToNot(BeNil())
		Expect(sub.Use).To(Equal(subUses[idx]))
	}
}
