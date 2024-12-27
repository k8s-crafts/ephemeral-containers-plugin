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
	"github.com/spf13/pflag"
)

var _ = Describe("Cmd", func() {
	var t *test

	BeforeEach(func() {
		t = newTest()
	})

	Context("root command", func() {
		BeforeEach(func() {
			t.cmd = cmd.NewRootCmd()
			t.subCmds = []string{"edit", "list", "version"}
		})

		It("should have basic configurations", func() {
			t.expectCmdBasics()
		})

		It("should have CommandDisplayNameAnnotatio annotation", func() {
			t.expectCmdAnnotation(cobra.CommandDisplayNameAnnotation, "kubectl ephemeral-containers")
		})

		It("should have subcommands", func() {
			t.expectSubCommands()
		})

		It("should have persistent flags", func() {
			// output flag and some kubeconfig flags
			for _, flag := range []string{"output", "kubeconfig", "namespace"} {
				t.expectFlag(flag, true)
			}
		})
	})

	Context("edit command", func() {
		BeforeEach(func() {
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

		It("should have local flags", func() {
			for _, flag := range []string{"editor", "minify"} {
				t.expectFlag(flag, false)
			}
		})
	})

	Context("list command", func() {
		BeforeEach(func() {
			t.cmd = cmd.NewListCmd()
		})

		It("should have basic configurations", func() {
			t.expectCmdBasics()
		})

		It("should have local flags", func() {
			for _, flag := range []string{"all-namespaces"} {
				t.expectFlag(flag, false)
			}
		})
	})

	Context("version command", func() {
		BeforeEach(func() {
			t.cmd = cmd.NewVersionCmd()
		})

		It("should have basic configurations", func() {
			t.expectCmdBasics()
		})

		It("should have no local flags", func() {
			Expect(t.cmd.HasLocalFlags()).To(BeFalse())
		})
	})
})

type test struct {
	*testInput
}

type testInput struct {
	cmd     *cobra.Command
	subCmds []string // sub-commands' Use
}

func newTest() *test {
	return &test{
		testInput: &testInput{},
	}
}

func (t *test) expectCmdAnnotation(key, value string) {
	Expect(t.cmd.Annotations).To(HaveKeyWithValue(key, value))
}

func (t *test) expectCmdBasics() {
	Expect(t.cmd).ToNot(BeNil())
	Expect(t.cmd.Use).ToNot(BeEmpty())
	Expect(t.cmd.Short).ToNot(BeEmpty())
	Expect(t.cmd.Long).ToNot(BeEmpty())
}

func (t *test) expectFlag(name string, persistent bool) {
	var flags *pflag.FlagSet
	if persistent {
		flags = t.cmd.PersistentFlags()
	} else {
		flags = t.cmd.LocalFlags()
	}

	Expect(flags).ToNot(BeNil())
	Expect(flags.HasFlags()).To(BeTrue())

	names := make([]string, 0)
	flags.VisitAll(func(flag *pflag.Flag) {
		names = append(names, flag.Name)
	})

	Expect(names).To(ContainElement(name))
}

func (t *test) expectSubCommands() {
	subs := t.cmd.Commands()

	Expect(subs).ToNot(BeNil())
	Expect(subs).To(HaveLen(len(t.subCmds)))

	actual := make([]string, len(subs))
	for idx := range subs {
		actual[idx] = subs[idx].Use
	}
	Expect(actual).To(ConsistOf(t.subCmds))
}
