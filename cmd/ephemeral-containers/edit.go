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

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/edit"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
)

var (
	editor      string
	editorUsage string = "Editor to use. If unset, the plugin will look into environment variable KUBE_EDITOR, EDITOR or fall back to vim"

	minify      bool
	minifyUsage string = "If true, remove information not necessary for editting ephemeral containers. Default to false"
)

func NewEditCmd() *cobra.Command {
	editCmd := &cobra.Command{
		Use:   "edit",
		Short: "Command to edit the ephemeralContainers spec for a Pod",
		Long: `
	This command is a convenient wrapper that, in turn, uses the pod's ephemeralcontainers subresource.
	
	Note: The command only consider changes to "pod.spec.ephemeralContainers". Other changes are ignored.
	`,
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

			if minify {
				pod = k8s.MinifyPod(pod)
			}

			editedPod, err := edit.EditResource(kubeConfig.ContextOptions, edit.GetEditorCmd(editor), pod, &corev1.Pod{})
			if err != nil {
				ExitError(errors.Join(fmt.Errorf("failed to edit pod/%s", podName), err), 1)
			}

			patch, err := k8s.SanitizeEditedPod(pod, editedPod)
			if err != nil {
				ExitError(err, 1)
			}

			if patch != nil {
				if _, err = k8s.UpdateEphemeralContainersForPod(kubeConfig.ContextOptions, client, patch); err != nil {
					ExitError(err, 1)
				}
				out.Ln("pod/%s successfully edited", podName)
			} else {
				out.Ln("Edit cancelled, no changes made for pod/%s", podName)
			}
		},
	}

	// Set default to empty to allow search in env vars
	editCmd.Flags().StringVarP(&editor, "editor", "e", "", editorUsage)
	editCmd.Flags().BoolVarP(&minify, "minify", "", false, minifyUsage)

	return editCmd
}
