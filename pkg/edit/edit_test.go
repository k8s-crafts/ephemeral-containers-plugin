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
