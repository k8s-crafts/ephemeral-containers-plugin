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

package main

import (
	"flag"
	"os"

	"github.com/k8s-crafts/ephemeral-containers-plugin/cmd"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/out"
	"github.com/spf13/pflag"
	klog "k8s.io/klog/v2"
)

func main() {
	// Initialize klog flag sets
	klog.InitFlags(nil)

	// Add go FlagSet (i.e. from klog) to pflag
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Resovle ErrHelp: pflag: help requested
	// ErrHelp is the error when the flag -help is invoked but no such flag is defined.
	// At this point, it is not yet defined. This works as a placerholder and will be correctly overriden later when cobra is initialized
	pflag.BoolP("help", "h", false, "")

	pflag.Parse()

	// Override flags if applicable but take user's preferences if any
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

	// Initialize the out and err destinations
	out.SetOutFile(os.Stdout)
	out.SetErrFile(os.Stderr)

	// Execute the command
	cmd.Execute()
}
