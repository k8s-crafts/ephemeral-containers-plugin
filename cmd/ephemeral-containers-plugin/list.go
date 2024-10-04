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

package ephemeralcontainersplugin

import (
	"k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	"k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"os"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the Pods with ephemeral containers in the current namespace",
	Long:  "List the Pods with ephemeral containers in the current namespace",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := k8s.NewClientset()
		if err != nil {
			os.Exit(1)
		}

		pods, err := k8s.ListPods(client, namespace, func(pod corev1.Pod) bool {
			return len(pod.Spec.EphemeralContainers) > 0
		})
		if err != nil {
			os.Exit(1)
		}

		if len(pods) > 0 {
			out.Ln("Pods: %v", pods)
		} else {
			out.Ln("No pods with ephemeral containers in namespace %s", namespace)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
