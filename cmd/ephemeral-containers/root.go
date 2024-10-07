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
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-ephemeral_containers",
	Short: "A kubectl plugin to directly modify pods.spec.ephemeralContainers",
	Long:  "A kubectl plugin to directly modify pods.spec.ephemeralContainers. It works by interacting the pod's ephemeralcontainers subresource",
	Annotations: map[string]string{
		cobra.CommandDisplayNameAnnotation: "kubectl ephemeral-containers",
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	// Namespace to act on. If unset (i.e. ""), it means all namespaces
	namespace string

	// Format for output
	outputFormat string
)

func init() {
	// Initialize klog flag sets. These flags are added to pflags in main
	klog.InitFlags(nil)

	// Define flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", formatter.Table, fmt.Sprintf("Format for output. One of: %s (default), %s, %s", formatter.Table, formatter.JSON, formatter.YAML))
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "The namespace to operate on. If unset (i.e. \"\"), all namespaces are considered")
}
