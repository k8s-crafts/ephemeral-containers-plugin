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
	"fmt"
	"k8s-crafts/ephemeral-containers-plugin/pkg/formatter"
	"k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	"k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-ephemeral_containers",
	Short: "A kubectl plugin to directly modify pods.spec.ephemeralContainers",
	Long:  "A kubectl plugin to directly modify pods.spec.ephemeralContainers. It works by interacting the pod's ephemeralcontainers subresource",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "kubectl ephemeral-containers",
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set namespace to "default" if unset
		if kubeConfig.Namespace == nil || *kubeConfig.Namespace == "" {
			kubeConfig.Namespace = &k8s.NAMESPACE_DEFAULT
		}

		kubeConfig.ContextOptions = k8s.NewContextOptions()
		if err := kubeConfig.ContextOptions.Init(kubeConfig.Timeout); err != nil {
			ExitError(err, 1)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		kubeConfig.ContextOptions.Cancel()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ExitError(err, 1)
	}
}

// Log errors and exit non-zero
func ExitError(err error, exitCode int) {
	out.Errf("%s", err.Error())
	os.Exit(exitCode)
}

var (
	kubeConfig *k8s.KubeConfig
	// Format for output
	outputFormat    string
	outputFlagUsage string = fmt.Sprintf("Format for output. One of: %s (default), %s, %s", formatter.Table, formatter.JSON, formatter.YAML)
)

func init() {
	// Initialize klog flag sets. These flags are added to pflags in main
	klog.InitFlags(nil)

	// Define flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", formatter.Table, outputFlagUsage)

	// Define kube CLI generic flags to generate a KubeConfig
	kubeConfig = k8s.NewKubeConfig(genericclioptions.NewConfigFlags(true))
	kubeConfig.AddFlags(rootCmd.PersistentFlags())
}
