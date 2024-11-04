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

package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var (
	NAMESPACE_DEFAULT string = "default"
)

// Represent context with a cancel func
type ContextOptions struct {
	context.Context
	Cancel  context.CancelFunc
	SigChan chan os.Signal
}

// Get a new ContextOptions struct
func NewContextOptions() *ContextOptions {
	return &ContextOptions{}
}

// Set up the options with the following steps:
// * Create a Context with timeout if any. Otherwise, no timeout is set (i.e. context.Background())
// * Create a chan os.Signal to handle SIGTERM, SIGINT (Ctrl + C),SIGHUP (terminal is closed)
func (opts *ContextOptions) Init(timeout *string) error {
	// Global context
	var ctx context.Context
	var cancel context.CancelFunc
	// Initialize the context with timeout
	// "0" means no time-out was requested
	if timeout != nil && *timeout != "0" {
		duration, err := time.ParseDuration(*timeout)
		if err != nil {
			return errors.Join(fmt.Errorf("failed to parse duration: %s", *timeout), err)
		}
		ctx, cancel = context.WithTimeout(context.Background(), duration)
	} else {
		ctx, cancel = context.Background(), func() {
			// no-op
		}
	}

	// Signal handler
	opts.SigChan = make(chan os.Signal, 1)

	signal.Notify(opts.SigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		select {
		case <-opts.SigChan:
			// Receive a signal
			opts.Cancel()
		case <-opts.Context.Done():
		}
	}()

	opts.Context = ctx
	opts.Cancel = func() {
		signal.Stop(opts.SigChan)
		cancel()
	}
	return nil
}

// Represent kube client configurations
type KubeConfig struct {
	*genericclioptions.ConfigFlags
	ContextOptions *ContextOptions
}

// Get a new KubeConfig struct
func NewKubeConfig(configFlags *genericclioptions.ConfigFlags) *KubeConfig {
	return &KubeConfig{
		ConfigFlags: configFlags,
	}
}

// Get a clientset to interact with Kubernetes API
// Config precedence:
// * --kubeconfig flag pointing at a file
// * KUBECONFIG environment variable pointing at a file
// * $HOME/.kube/config if exists.
func NewClientset(kubeConfig *KubeConfig) (*KubeClientset, error) {
	config, err := kubeConfig.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	_clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeClientset{
		Clientset: _clientset,
	}, nil
}
