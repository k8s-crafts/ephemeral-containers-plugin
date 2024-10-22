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
