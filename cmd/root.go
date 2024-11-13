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
	"fmt"
	"os"

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/formatter"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	// KubeConfig reference
	kubeConfig *k8s.KubeConfig

	outputFormat    string
	outputFlagUsage string = fmt.Sprintf("Format for output. One of: default (%s for lists), %s, %s", formatter.Table, formatter.JSON, formatter.YAML)
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kubectl-ephemeral_containers",
		Short: "A kubectl plugin to directly modify pods.spec.ephemeralContainers",
		Long:  "A kubectl plugin to directly modify pods.spec.ephemeralContainers. It works by interacting the pod's ephemeralcontainers subresource",
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl ephemeral-containers",
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ns, _, err := kubeConfig.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				ExitError(err, 1)
			}

			// Use namespace "default" if none is set
			if len(ns) == 0 {
				ns = k8s.NAMESPACE_DEFAULT
			}
			kubeConfig.Namespace = &ns

			if err := kubeConfig.InitContext(kubeConfig.Timeout); err != nil {
				ExitError(err, 1)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			kubeConfig.ContextOptions.Cancel()
		},
	}

	// Define flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", formatter.Table, outputFlagUsage)

	// Define kube CLI generic flags to generate a KubeConfig
	kubeConfig.AddFlags(rootCmd.PersistentFlags())

	// Add subcommands
	rootCmd.AddCommand(NewEditCmd(), NewListCmd(), NewVersionCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		ExitError(err, 1)
	}
}

// Log errors and exit non-zero
func ExitError(err error, exitCode int) {
	out.ErrLn("%s", err.Error())
	os.Exit(exitCode)
}

func init() {
	kubeConfig = k8s.NewKubeConfig()
}
