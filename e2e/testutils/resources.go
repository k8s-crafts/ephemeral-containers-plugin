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

package testutils

import (
	"fmt"
)

var (
	MinK8sVersion    string = "v1.25.0"
	TestPodName      string = "plugin-e2e"
	DebugImage       string = "docker.io/library/busybox:1.28"
	EphContainerName string = "debugger"
	PodName          string = "plugin-e2e"
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
`, PodName, t.Kubectl.Namespace, EphContainerName, PodName, t.AnotherNamespace, EphContainerName,
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

`, EphContainerName, PodName, t.Kubectl.Namespace, EphContainerName, PodName, t.AnotherNamespace)
	default: // table or empty
		return fmt.Sprintf(
			`+------------+-----------+----------------------+
|    POD     | NAMESPACE | EPHEMERAL CONTAINERS |
+------------+-----------+----------------------+
| %s | %s | %s             |
| %s | %s | %s             |
+------------+-----------+----------------------+

`, PodName, t.Kubectl.Namespace, EphContainerName, PodName, t.AnotherNamespace, EphContainerName)
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
`, PodName, namespace, EphContainerName,
		)
	case "yaml":
		return fmt.Sprintf(`- ephemeralContainers:
  - %s
  name: %s
  namespace: %s

`, EphContainerName, PodName, t.Kubectl.Namespace)
	default: // table or empty
		return fmt.Sprintf(
			`+------------+-----------+----------------------+
|    POD     | NAMESPACE | EPHEMERAL CONTAINERS |
+------------+-----------+----------------------+
| %s | %s | %s             |
+------------+-----------+----------------------+

`, PodName, namespace, EphContainerName)
	}
}

func (t *TestResource) NewListEmptyMessage(namespace string) string {
	if len(namespace) == 0 {
		return "No pods with ephemeral containers found any namespaces\n"
	}
	return fmt.Sprintf("No pods with ephemeral containers found in namespace %s\n", namespace)
}
