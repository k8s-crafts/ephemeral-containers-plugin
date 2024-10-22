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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/version"
	"github.com/olekukonko/tablewriter"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

const (
	JSON  string = "json"
	YAML  string = "yaml"
	Table string = "table" // Default format
)

var (
	TableHeaders []string = []string{"Pod", "Namespace", "Ephemeral Containers"}
)

type ResourceData struct {
	Name                string   `json:"name,omitempty"`
	Namespace           string   `json:"namespace,omitempty"`
	EphemeralContainers []string `json:"ephemeralContainers"`
}

// List the name of ehemeral containers for a Pod
func ListEphemeralContainersForPod(pod corev1.Pod) (containers []string) {
	for _, container := range pod.Spec.EphemeralContainers {
		containers = append(containers, container.Name)
	}
	return containers
}

// Convert Pod data to simplified version
func ConvertPodsToResourceData(pods []corev1.Pod) (data []ResourceData) {
	for _, pod := range pods {
		data = append(data, ResourceData{
			Name:                pod.Name,
			Namespace:           pod.Namespace,
			EphemeralContainers: ListEphemeralContainersForPod(pod),
		})
	}
	return data
}

// Get a table row from resource data
func GetTableRow(data ResourceData) []string {
	return []string{data.Name, data.Namespace, strings.Join(data.EphemeralContainers, ",")}
}

// Formatter for list output
func FormatListOutput(format string, pods []corev1.Pod) (string, error) {
	data := ConvertPodsToResourceData(pods)
	if len(data) == 0 {
		return "", nil
	}

	switch format {
	case JSON:
		jsonOut, err := json.MarshalIndent(data, "", "  ")
		return string(jsonOut), err
	case YAML:
		yamlOut, err := yaml.Marshal(data)
		return string(yamlOut), err
	default:
		var buffer bytes.Buffer
		table := tablewriter.NewWriter(&buffer)

		// Add header
		table.SetHeader(TableHeaders)

		for _, d := range data {
			table.Append(GetTableRow(d))
		}

		table.Render()

		return buffer.String(), nil
	}
}

// Formatter for version output
func FormatVersionOutput(format string, version *version.VersionInfo) (string, error) {
	if version == nil {
		return "", nil
	}

	switch format {
	case JSON:
		jsonOut, err := json.MarshalIndent(version, "", "  ")
		return string(jsonOut), err
	case YAML:
		yamlOut, err := yaml.Marshal(version)
		return string(yamlOut), err
	default:
		return fmt.Sprintf("version: %v", version.Version), nil
	}
}
