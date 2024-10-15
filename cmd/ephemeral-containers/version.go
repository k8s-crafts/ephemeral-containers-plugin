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
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/formatter"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the plugin version",
	Long: `
Output the plugin version
`,
	Run: func(cmd *cobra.Command, args []string) {
		versionInfo := version.NewVersionInfo()
		output, err := formatter.FormatVersionOutput(outputFormat, versionInfo)
		if err != nil {
			ExitError(err, 1)
		}

		out.Ln("%s", output)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
