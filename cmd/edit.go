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

package cmd

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

			pod, err := client.GetPod(kubeConfig.ContextOptions, *kubeConfig.Namespace, podName)
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
				if _, err = client.UpdateEphemeralContainersForPod(kubeConfig.ContextOptions, patch); err != nil {
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
