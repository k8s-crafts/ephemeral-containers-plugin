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
