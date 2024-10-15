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

package edit

import (
	"os"
	"testing"
)

func TestGetEditorCmd(t *testing.T) {
	type Env struct {
		Key   string
		Value string
	}

	setupEnv := func(t *testing.T, envs []Env) {
		for _, env := range envs {
			t.Setenv(env.Key, env.Value)
		}
	}

	tests := []struct {
		desciptions string
		fromFlag    string
		envs        []Env
		expected    string
	}{

		{
			desciptions: "should return editor cmd from cli flags",
			fromFlag:    "nano",
			envs: []Env{
				{
					Key:   ENV_EDITOR,
					Value: "code",
				},
				{
					Key:   ENV_KUBE_EDITOR,
					Value: "vi",
				},
			},
			expected: "nano",
		},
		{
			desciptions: "should return editor cmd from env var KUBE_EDITOR",
			fromFlag:    "",
			envs: []Env{
				{
					Key:   ENV_EDITOR,
					Value: "code",
				},
				{
					Key:   ENV_KUBE_EDITOR,
					Value: "vi",
				},
			},
			expected: "vi",
		},
		{
			desciptions: "should return editor cmd from env var EDITOR",
			fromFlag:    "",
			envs: []Env{
				{
					Key:   ENV_EDITOR,
					Value: "code",
				},
			},
			expected: "code",
		},
		{
			desciptions: "should return vim as default editor cmd",
			fromFlag:    "",
			// Env vars are set to empty
			// since the system might have already set them
			envs: []Env{
				{
					Key:   ENV_EDITOR,
					Value: "",
				},
				{
					Key:   ENV_KUBE_EDITOR,
					Value: "",
				},
			},
			expected: "vi",
		},
	}

	for _, test := range tests {
		t.Run(test.desciptions, func(t *testing.T) {
			setupEnv(t, test.envs)

			actual := GetEditorCmd(test.fromFlag)
			if actual != test.expected {
				t.Errorf("%v", os.Getenv("KUBE_EDITOR"))
				t.Fatalf("expect to get %v, but received %v", test.expected, actual)

			}
		})
	}
}
