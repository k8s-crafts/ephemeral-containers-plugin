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
