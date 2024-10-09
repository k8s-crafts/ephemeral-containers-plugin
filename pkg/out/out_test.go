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

package out

import (
	"bytes"
	"testing"
)

func beforeEach() {
	outFile = nil
	errFile = nil
}

func TestSetFileStreams(t *testing.T) {
	tests := []struct {
		outFile     *bytes.Buffer
		errFile     *bytes.Buffer
		description string
	}{
		{
			description: "should set correct outFile and errFile",
			outFile:     new(bytes.Buffer),
			errFile:     new(bytes.Buffer),
		},
	}

	for _, test := range tests {
		beforeEach()
		t.Run(test.description, func(t *testing.T) {
			SetOutFile(test.outFile)
			SetErrFile(test.errFile)

			if outFile != test.outFile {
				t.Errorf("failed to set outFile")
			}

			if errFile != test.errFile {
				t.Errorf("failed to set errFile")
			}
		})
	}
}

func TestOutputPrints(t *testing.T) {
	tests := []struct {
		description string

		isStderr bool

		outFile *bytes.Buffer
		errFile *bytes.Buffer

		fn     func(format string, a ...interface{})
		format string
		a      []interface{}

		expected string
	}{
		{
			description: "should correctly print content to stdout with Ln",
			outFile:     new(bytes.Buffer),
			errFile:     new(bytes.Buffer),
			fn:          Ln,
			format:      "this is a format for content: %s",
			a:           []interface{}{"my-content"},
			expected:    "this is a format for content: my-content\n",
		},
		{
			description: "should correctly print content to stdout with Stringf",
			outFile:     new(bytes.Buffer),
			errFile:     new(bytes.Buffer),
			fn:          Stringf,
			format:      "this is a format for content: %s",
			a:           []interface{}{"my-content"},
			expected:    "this is a format for content: my-content",
		},
		{
			description: "should correctly print content to stdout with ErrLn",
			isStderr:    true,
			outFile:     new(bytes.Buffer),
			errFile:     new(bytes.Buffer),
			fn:          ErrLn,
			format:      "this is a format for content: %s",
			a:           []interface{}{"my-content"},
			expected:    "this is a format for content: my-content\n",
		},
		{
			description: "should correctly print content to stdout with Errf",
			isStderr:    true,
			outFile:     new(bytes.Buffer),
			errFile:     new(bytes.Buffer),
			fn:          Errf,
			format:      "this is a format for content: %s",
			a:           []interface{}{"my-content"},
			expected:    "this is a format for content: my-content",
		},
	}

	for _, test := range tests {
		beforeEach()
		SetOutFile(test.outFile)
		SetErrFile(test.errFile)

		t.Run(test.description, func(t *testing.T) {
			// Expect output
			test.fn(test.format, test.a...)

			var actual string
			if test.isStderr {
				actual = test.errFile.String()
			} else {
				actual = test.outFile.String()
			}
			if test.expected != actual {
				t.Errorf("expected output %s but received %s", test.expected, actual)
			}
		})
	}
}
