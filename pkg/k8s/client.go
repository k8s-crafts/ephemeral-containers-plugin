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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Get a clientset to interact with Kubernetes API
// Config precedence:
// * --kubeconfig flag pointing at a file
// * KUBECONFIG environment variable pointing at a file
// * In-cluster config if running in cluster
// * $HOME/.kube/config if exists.
func NewClientset() (*kubernetes.Clientset, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

type PodFilterFn func(pod corev1.Pod) bool

// List pods by filters in the specified namespace
// If namespace is empty (i.e. ""), list in all namespaces
func ListPods(client *kubernetes.Clientset, namespace string, filters ...PodFilterFn) ([]corev1.Pod, error) {
	podList, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ApplyPodFilter(podList.Items, filters...), nil
}

// Apply filters (if any) to a list of pods
func ApplyPodFilter(pods []corev1.Pod, filters ...PodFilterFn) (result []corev1.Pod) {
	if len(filters) == 0 {
		return pods
	}

	for _, pod := range pods {
		for _, filter := range filters {
			if filter(pod) {
				result = append(result, pod)
			}
		}
	}

	return result
}
