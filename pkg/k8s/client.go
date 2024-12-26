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
	noopFunc          func() = func() {}
)

// Represent context with a cancel func
type ContextOptions struct {
	context.Context
	Cancel  context.CancelFunc
	SigChan chan os.Signal
}

// Set up the options with the following steps:
// * Create a Context with timeout if any. Otherwise, no timeout is set (i.e. context.Background())
// * Create a chan os.Signal to handle SIGTERM, SIGINT (Ctrl + C),SIGHUP (terminal is closed)
func (opts *ContextOptions) InitContext(timeout *string) error {
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
		ctx, cancel = context.Background(), noopFunc
	}

	opts.SigChan = make(chan os.Signal, 1)
	signal.Notify(opts.SigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	opts.Context = ctx
	opts.Cancel = func() {
		cancel()
		signal.Stop(opts.SigChan)
	}

	// Signal handler
	go func() {
		select {
		case <-opts.SigChan:
			// Receive a signal
			opts.Cancel()
		case <-opts.Context.Done():
		}
	}()

	return nil
}

// Cancel execution context
func (opts *ContextOptions) CancelContext() {
	opts.Cancel()
}

// Represent kube client configurations
type KubeConfig struct {
	*genericclioptions.ConfigFlags
	*ContextOptions
}

// Get a new KubeConfig struct
func NewKubeConfig() *KubeConfig {
	return &KubeConfig{
		ConfigFlags:    genericclioptions.NewConfigFlags(true),
		ContextOptions: &ContextOptions{},
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
		Interface: _clientset,
	}, nil
}
