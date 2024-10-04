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

package flags

import (
	"flag"

	"github.com/spf13/pflag"
	klog "k8s.io/klog/v2"
)

// Override flags if applicable but take user's preferences if any
// Note: This func will parse flags
func SetFlagDefaults() {
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Resovle ErrHelp: pflag: help requested
	// ErrHelp is the error when the flag -help is invoked but no such flag is defined.
	// At this point, it is not yet defined. This works as a placerholder and will be correctly overriden later when cobra is initialized
	pflag.BoolP("help", "h", false, "")

	pflag.Parse()

	if !pflag.CommandLine.Changed("logtostderr") {
		if err := pflag.Set("logtostderr", "false"); err != nil {
			klog.Errorf("Failed to set default flag for logtostderr: %v", err)
		}
	}

	if !pflag.CommandLine.Changed("alsologtostderr") {
		if err := pflag.Set("alsologtostderr", "false"); err != nil {
			klog.Errorf("Failed to set default flag for alsologtostderr: %v", err)
		}
	}

}
