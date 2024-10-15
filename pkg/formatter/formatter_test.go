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

package formatter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/version"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestListEphemeralContainersForPod(t *testing.T) {
	tests := []struct {
		description string
		pod         *corev1.Pod
		expected    []string
	}{
		{
			description: "should return the name of ephemeral containers for Pod if available",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					EphemeralContainers: []corev1.EphemeralContainer{
						{
							EphemeralContainerCommon: corev1.EphemeralContainerCommon{
								Name:  "debug-container",
								Image: "my-image:v1",
							},
						},
					},
				},
			},
			expected: []string{"debug-container"},
		},
		{
			description: "should return empty if Pod has no ephemeral",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mypod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{},
			},
			expected: make([]string, 0),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := ListEphemeralContainersForPod(*test.pod)

			if len(actual) != len(test.expected) {
				t.Errorf("expected to get length %d, but received %d", len(test.expected), len(actual))
			}

			for index, name := range actual {
				if name != test.expected[index] {
					t.Errorf("expected to get name %s at index %d, but received %s", test.expected[index], index, name)
				}
			}
		})
	}
}

func TestConvertPodsToResourceData(t *testing.T) {
	tests := []struct {
		description string
		pods        []corev1.Pod
		expected    []ResourceData
	}{
		{
			description: "should return ResourceData that represents the Pod",
			pods: []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mypod-0",
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						EphemeralContainers: []corev1.EphemeralContainer{
							{
								EphemeralContainerCommon: corev1.EphemeralContainerCommon{
									Name:  "debug-container",
									Image: "my-image:v1",
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mypod-1",
						Namespace: "default",
					},
					Spec: corev1.PodSpec{},
				},
			},
			expected: []ResourceData{
				{
					Name:                "mypod-0",
					Namespace:           "default",
					EphemeralContainers: []string{"debug-container"},
				},
				{
					Name:                "mypod-1",
					Namespace:           "default",
					EphemeralContainers: make([]string, 0),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := ConvertPodsToResourceData(test.pods)

			if len(actual) != len(test.expected) {
				t.Errorf("expected to get length %d, but received %d", len(test.expected), len(actual))
			}

			isEqual := func(left, right ResourceData) bool {
				return left.Name == right.Name &&
					left.Namespace == right.Namespace &&
					len(left.EphemeralContainers) == len(right.EphemeralContainers) &&
					strings.Join(left.EphemeralContainers, ",") == strings.Join(right.EphemeralContainers, ",")
			}

			for index, data := range actual {
				if !isEqual(data, test.expected[index]) {
					t.Errorf("expected to get data %+v at index %d, but received %+v", test.expected[index], index, data)
				}
			}
		})
	}
}

func TestGetTableRow(t *testing.T) {
	tests := []struct {
		description string
		data        ResourceData
		expected    []string
	}{
		{
			description: "should return the table row for data with multiple ephemeral containers",
			data: ResourceData{
				Name:                "my-pod",
				Namespace:           "default",
				EphemeralContainers: []string{"my-container-0", "my-container-1"},
			},
			expected: []string{"my-pod", "default", "my-container-0,my-container-1"},
		},
		{
			description: "should return the table row for data with single ephemeral container",
			data: ResourceData{
				Name:                "my-pod",
				Namespace:           "default",
				EphemeralContainers: []string{"my-container-0"},
			},
			expected: []string{"my-pod", "default", "my-container-0"},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := GetTableRow(test.data)

			if len(actual) != len(test.expected) {
				t.Errorf("expected to get length %d, but received %d", len(test.expected), len(actual))
			}

			for index, name := range actual {
				if name != test.expected[index] {
					t.Errorf("expected to get name %s at index %d, but received %s", test.expected[index], index, name)
				}
			}
		})
	}
}

func TestFormatListOutput(t *testing.T) {
	tests := []struct {
		description string
		data        []ResourceData
		format      string
		expected    string
	}{
		{
			description: "should return the list output in table format",
			data: []ResourceData{
				{
					Name:                "my-pod",
					Namespace:           "default",
					EphemeralContainers: []string{"my-container-0", "my-container-1"},
				},
			},
			format: "table",
			expected: `+--------+-----------+-------------------------------+
|  POD   | NAMESPACE |     EPHEMERAL CONTAINERS      |
+--------+-----------+-------------------------------+
| my-pod | default   | my-container-0,my-container-1 |
+--------+-----------+-------------------------------+
`,
		},
		{
			description: "should return the list output in JSON format",
			data: []ResourceData{
				{
					Name:                "my-pod",
					Namespace:           "default",
					EphemeralContainers: []string{"my-container-0", "my-container-1"},
				},
			},
			format: "json",
			expected: `[
  {
    "name": "my-pod",
    "namespace": "default",
    "ephemeralContainers": [
      "my-container-0",
      "my-container-1"
    ]
  }
]`,
		},
		{
			description: "should return the list output in YAML format",
			data: []ResourceData{
				{
					Name:                "my-pod",
					Namespace:           "default",
					EphemeralContainers: []string{"my-container-0", "my-container-1"},
				},
			},
			format: "yaml",
			expected: `- ephemeralContainers:
  - my-container-0
  - my-container-1
  name: my-pod
  namespace: default
`,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual, err := FormatListOutput(test.format, test.data)
			if err != nil {
				t.Errorf("failed to output list: %v", err)
			}

			if actual != test.expected {
				fmt.Print(actual)
				t.Errorf("expected to get output %s, but received %s", test.expected, actual)
			}
		})
	}
}

func TestFormatVersionOutput(t *testing.T) {
	tests := []struct {
		description string
		version     *version.VersionInfo
		format      string
		expected    string
	}{
		{
			description: "should return the version output in table format",
			version: &version.VersionInfo{
				Version: "v0.0.0-unset",
			},
			format:   "table",
			expected: "version: v0.0.0-unset\n",
		},
		{
			description: "should return the version output in JSON format",
			version: &version.VersionInfo{
				Version: "v0.0.0-unset",
			},
			format: "json",
			expected: `{
  "version": "v0.0.0-unset"
}`,
		},
		{
			description: "should return the version output in YAML format",
			version: &version.VersionInfo{
				Version: "v0.0.0-unset",
			},
			format:   "yaml",
			expected: "version: v0.0.0-unset\n",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual, err := FormatVersionOutput(test.format, test.version)
			if err != nil {
				t.Errorf("failed to output version info: %v", err)
			}

			if actual != test.expected {
				t.Errorf("expected to get output %s, but received %s", test.expected, actual)
			}
		})
	}
}
