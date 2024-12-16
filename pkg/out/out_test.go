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

package out_test

import (
	"bytes"

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/out"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Out", func() {
	var t *test

	BeforeEach(func() {
		t = newTest()
	})

	Context("when printing contents", func() {
		JustBeforeEach(func() {
			out.SetOutFile(t.f)
		})

		Context("to outFile", func() {
			It("should set outFile", func() {
				Expect(out.GetOutFile()).To(Equal(t.f))
			})
			It("should print content with newline", func() {
				out.Ln(t.format, t.subs...)
				t.expectContent(t.content + "\n")
			})
			It("should print content as is", func() {
				out.Stringf(t.format, t.subs...)
				t.expectContent(t.content)
			})
		})
		Context("to errFile", func() {
			JustBeforeEach(func() {
				out.SetErrFile(t.f)
			})

			It("should set errFile", func() {
				Expect(out.GetErrFile()).To(Equal(t.errFile))
			})
			It("should print content with newline", func() {
				out.ErrLn(t.format, t.subs...)
				t.expectContent(t.content + "\n")
			})
			It("should print content as is", func() {
				out.Errf(t.format, t.subs...)
				t.expectContent(t.content)
			})
		})
	})
})

// Input for test cases
type testInput struct {
	f       *bytes.Buffer
	errFile *bytes.Buffer

	format string
	subs   []interface{}

	content string
}

type test struct {
	*testInput
}

func (t *test) expectContent(content string) {
	Expect(t.f.String()).To(Equal(content))
}

func newTest() *test {
	return &test{
		testInput: &testInput{
			f:       new(bytes.Buffer),
			errFile: new(bytes.Buffer),

			format: "this is a format for content: %s",
			subs:   []interface{}{"my-content"},

			content: "this is a format for content: my-content",
		},
	}
}
