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

package files

import (
	"context"
	"os"
	"os/exec"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

const (
	DEFAULT_EDITOR   string = "vi"
	ENV_EDITOR       string = "EDITOR"
	ENV_KUBE_EDITOR  string = "KUBE_EDITOR"
	TMP_FILE_PATTERN string = "ephemeral-containers-"
)

// Edit a k8s resource and return the updated one
func EditResource[r runtime.Object](ctx context.Context, editor string, obj r, result r) (r, error) {
	f, err := os.CreateTemp(os.TempDir(), TMP_FILE_PATTERN)
	if err != nil {
		return result, err
	}

	content, err := yaml.Marshal(obj)
	if err != nil {
		return result, err
	}

	if _, err = f.Write(content); err != nil {
		return result, err
	}

	if err = OpenEditorForFile(ctx, editor, f.Name()); err != nil {
		return result, err
	}

	content, err = os.ReadFile(f.Name())
	if err != nil {
		return result, err
	}

	err = yaml.Unmarshal(content, result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Execute editor command for a file path and await closing editor
func OpenEditorForFile(ctx context.Context, editor, path string, args ...string) error {
	cmd := exec.CommandContext(ctx, editor, append(args, path)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Get editor executable from sources if not set on cli flags
// Precedence:
// * --editor flag
// * KUBE_EDITOR env var
// * EDITOR env var
// * Default to vi (vim)
func GetEditorCmd(fromFlag string) string {
	return GetValueFromSources(
		func() string {
			return fromFlag
		},
		func() string {
			return os.Getenv(ENV_KUBE_EDITOR)
		},
		func() string {
			return os.Getenv(ENV_EDITOR)
		},
	)
}

// Get value from multiple sources until a non-empty value is returned
func GetValueFromSources(sources ...func() string) string {
	for _, source := range sources {
		if val := source(); len(val) > 0 {
			return val
		}
	}
	return ""
}
