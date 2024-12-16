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

package edit_test

import (
	"os"

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/edit"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Edit", func() {
	var t *test

	Context("when editting", func() {

		JustBeforeEach(func() {
			for _, env := range t.envs {
				os.Setenv(env.key, env.value)
			}
		})

		JustAfterEach(func() {
			for _, env := range t.envs {
				os.Unsetenv(env.key)
			}
		})

		Context("with editor from cmd flags", func() {
			BeforeEach(func() {
				t = newTestForEditorFromFlags()
			})

			It("should get editor cmd", func() {
				t.expectEditor(edit.GetEditorCmd("nano"))
			})
		})
		Context("with editor from env var KUBE_EDITOR", func() {
			BeforeEach(func() {
				t = newTestForEditorFromVarKubeEditor()
			})
			It("should get editor cmd", func() {
				t.expectEditor(edit.GetEditorCmd(""))
			})
		})
		Context("with editor from env var EDITOR", func() {
			BeforeEach(func() {
				t = newTestForEditorFromVarEditor()
			})
			It("should get editor cmd", func() {
				t.expectEditor(edit.GetEditorCmd(""))
			})
		})
		Context("with editor vim as default", func() {
			BeforeEach(func() {
				t = newTestForDefaultEditor()
			})
			It("should get editor cmd", func() {
				t.expectEditor(edit.GetEditorCmd(""))
			})
		})
	})
})

// Struct type representing an environment variable
// as a key-value pair
type env struct {
	key   string
	value string
}

// Input for tests
type testInput struct {
	envs   []env
	editor string
}

type test struct {
	*testInput
}

func (t *test) expectEditor(editor string) {
	Expect(t.editor).To(Equal(editor))
}

func newTestForEditorFromFlags() *test {
	return &test{
		testInput: &testInput{
			envs: []env{
				{
					key:   edit.ENV_EDITOR,
					value: "code",
				},
				{
					key:   edit.ENV_KUBE_EDITOR,
					value: "vi",
				},
			},
			editor: "nano",
		},
	}
}

func newTestForEditorFromVarKubeEditor() *test {
	return &test{
		testInput: &testInput{
			envs: []env{
				{
					key:   edit.ENV_EDITOR,
					value: "code",
				},
				{
					key:   edit.ENV_KUBE_EDITOR,
					value: "vi",
				},
			},
			editor: "vi",
		},
	}
}

func newTestForEditorFromVarEditor() *test {
	return &test{
		testInput: &testInput{
			envs: []env{
				{
					key:   edit.ENV_EDITOR,
					value: "code",
				},
			},
			editor: "code",
		},
	}
}

func newTestForDefaultEditor() *test {
	return &test{
		testInput: &testInput{
			envs:   make([]env, 0),
			editor: "vi",
		},
	}
}
