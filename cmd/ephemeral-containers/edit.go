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

package ephemeralcontainers

import (
	"errors"
	"fmt"
	"k8s-crafts/ephemeral-containers-plugin/pkg/edit"
	"k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	"k8s-crafts/ephemeral-containers-plugin/pkg/out"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Command to edit the ephemeralContainers spec for a Pod",
	Long:  "This command is a convenient wrapper that, in turn, uses the pod's ephemeralcontainers subresource",
	// Format: "pod/pod-name", "pod pod-name", "pod-name"
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		podName, err := k8s.GetPodNameFromArgs(args)
		if err != nil {
			ExitError(err, 1)
		}

		client, err := k8s.NewClientset(kubeConfig)
		if err != nil {
			ExitError(err, 1)
		}

		pod, err := k8s.GetPod(kubeConfig.ContextOptions, client, *kubeConfig.Namespace, podName)
		if err != nil {
			ExitError(err, 1)
		}

		editorCmd := edit.GetEditorCmd(editor)
		obj, err := edit.EditResource(kubeConfig.ContextOptions, editorCmd, pod, &corev1.Pod{})
		if err != nil {
			ExitError(errors.Join(fmt.Errorf("failed to edit pod/%s", podName), err), 1)
		}

		if err = k8s.UpdateEphemeralContainersForPod(kubeConfig.ContextOptions, client, obj); err != nil {
			ExitError(err, 1)
		}

		out.Ln("pod/%s successfully edited", podName)
	},
}

var (
	editor      string
	editorUsage string = "Editor to use. If unset, the plugin will look into environment variable KUBE_EDITOR, EDITOR or fall back to vim"
)

func init() {
	// Set default to empty to allow search in env vars
	editCmd.Flags().StringVarP(&editor, "editor", "e", "", editorUsage)

	rootCmd.AddCommand(editCmd)
}
