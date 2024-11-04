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

package cmd

import (
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/formatter"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	filterFn = func(pod corev1.Pod) bool {
		return len(pod.Spec.EphemeralContainers) > 0
	}

	allNamespace      bool
	allNamespaceUsage = "If true, list the pods in all namespaces"
)

func NewListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the Pods with ephemeral containers in the current namespace",
		Long: `
List the Pods with ephemeral containers in the current namespace
	`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := k8s.NewClientset(kubeConfig)
			if err != nil {
				ExitError(err, 1)
			}

			namespace := *kubeConfig.Namespace
			if allNamespace {
				namespace = ""
			}

			pods, err := k8s.ListPods(kubeConfig.ContextOptions, client, namespace, filterFn)
			if err != nil {
				ExitError(err, 1)
			}

			output, err := formatter.FormatListOutput(outputFormat, pods)
			if err != nil {
				ExitError(err, 1)
			}

			if len(output) > 0 {
				out.Ln("%v", output)
			} else if allNamespace {
				out.Ln("No pods with ephemeral containers found any namespaces")
			} else {
				out.Ln("No pods with ephemeral containers found in namespace %s", *kubeConfig.Namespace)
			}

		},
	}

	listCmd.Flags().BoolVarP(&allNamespace, "all-namespaces", "A", false, allNamespaceUsage)

	return listCmd
}
