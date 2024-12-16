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

package edit

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

const (
	DEFAULT_EDITOR   string = "vi"
	ENV_EDITOR       string = "EDITOR"
	ENV_KUBE_EDITOR  string = "KUBE_EDITOR"
	TMP_FILE_PATTERN string = "ephemeral-containers-*.yaml"
)

// Edit a k8s resource and return the updated one
func EditResource[r runtime.Object](ctx context.Context, editor string, obj r, result r) (r, error) {
	f, err := os.CreateTemp(os.TempDir(), TMP_FILE_PATTERN)
	if err != nil {
		return result, err
	}
	defer func() {
		// Clean up
		err = errors.Join(err, f.Close(), os.Remove(f.Name()))
	}()

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
	sources := []string{
		fromFlag,
		os.Getenv(ENV_KUBE_EDITOR),
		os.Getenv(ENV_EDITOR),
	}

	for _, source := range sources {
		if len(source) > 0 {
			return source
		}
	}

	return DEFAULT_EDITOR
}
