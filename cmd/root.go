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
