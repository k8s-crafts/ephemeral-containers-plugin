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
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var (
	NAMESPACE_DEFAULT string = "default"
)

// Get a clientset to interact with Kubernetes API
// Config precedence:
// * --kubeconfig flag pointing at a file
// * KUBECONFIG environment variable pointing at a file
// * $HOME/.kube/config if exists.
func NewClientset(kubeConfig *genericclioptions.ConfigFlags) (*kubernetes.Clientset, error) {
	config, err := kubeConfig.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
