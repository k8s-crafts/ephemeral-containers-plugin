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

package testutils

import (
	"fmt"
)

var (
	MinK8sVersion    string = "v1.25.0"
	TestPodName      string = "plugin-e2e"
	DebugImage       string = "docker.io/library/busybox:1.28"
	EphContainerName string = "debugger"
)

func (t *TestResource) NewListOutput(format string, namespace string) string {
	if len(namespace) == 0 {
		return t.newListInAllNamespaces(format)
	}
	return t.newListInNamespace(format, namespace)
}

func (t *TestResource) newListInAllNamespaces(format string) string {
	switch format {
	case "json":
		return fmt.Sprintf(`[
  {
    "name": "%s",
    "namespace": "%s",
    "ephemeralContainers": [
      "%s"
    ]
  },
  {
    "name": "%s",
    "namespace": "%s",
    "ephemeralContainers": [
      "%s"
    ]
  }
]
`, TestPodName, t.Namespace, EphContainerName, TestPodName, t.AnotherNamespace, EphContainerName,
		)
	case "yaml":
		return fmt.Sprintf(`- ephemeralContainers:
  - %s
  name: %s
  namespace: %s
- ephemeralContainers:
  - %s
  name: %s
  namespace: %s

`, EphContainerName, TestPodName, t.Namespace, EphContainerName, TestPodName, t.AnotherNamespace)
	default: // table or empty
		return fmt.Sprintf(
			`+------------+-----------+----------------------+
|    POD     | NAMESPACE | EPHEMERAL CONTAINERS |
+------------+-----------+----------------------+
| %s | %s | %s             |
| %s | %s | %s             |
+------------+-----------+----------------------+

`, TestPodName, t.Namespace, EphContainerName, TestPodName, t.AnotherNamespace, EphContainerName)
	}
}

func (t *TestResource) newListInNamespace(format string, namespace string) string {
	switch format {
	case "json":
		return fmt.Sprintf(`[
  {
    "name": "%s",
    "namespace": "%s",
    "ephemeralContainers": [
      "%s"
    ]
  }
]
`, TestPodName, namespace, EphContainerName,
		)
	case "yaml":
		return fmt.Sprintf(`- ephemeralContainers:
  - %s
  name: %s
  namespace: %s

`, EphContainerName, TestPodName, t.Namespace)
	default: // table or empty
		return fmt.Sprintf(
			`+------------+-----------+----------------------+
|    POD     | NAMESPACE | EPHEMERAL CONTAINERS |
+------------+-----------+----------------------+
| %s | %s | %s             |
+------------+-----------+----------------------+

`, TestPodName, namespace, EphContainerName)
	}
}

func (t *TestResource) NewListEmptyMessage(namespace string) string {
	if len(namespace) == 0 {
		return "No pods with ephemeral containers found any namespaces\n"
	}
	return fmt.Sprintf("No pods with ephemeral containers found in namespace %s\n", namespace)
}
