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
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodFilterFn func(pod corev1.Pod) bool

// List pods by filters in the specified namespace
// If namespace is empty (i.e. ""), list in all namespaces
func ListPods(ctx context.Context, client *kubernetes.Clientset, namespace string, filters ...PodFilterFn) ([]corev1.Pod, error) {
	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ApplyPodFilter(podList.Items, filters...), nil
}

// Get pod by name in a specific namespace
func GetPod(ctx context.Context, client *kubernetes.Clientset, namespace, name string) (*corev1.Pod, error) {
	return client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

// Update pod's ephemeralContainer subresource
func UpdateEphemeralContainersForPod(ctx context.Context, client *kubernetes.Clientset, pod *corev1.Pod) error {
	_, err := client.CoreV1().Pods(pod.Namespace).UpdateEphemeralContainers(ctx, pod.Name, pod, metav1.UpdateOptions{})
	return err
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

// Get pod name from CLI arguments
func GetPodNameFromArgs(args []string) (string, error) {
	switch len(args) {
	case 1:
		parts := strings.Split(args[0], "/")
		if len(parts) > 2 {
			return "", errors.New("single argument must have format: \"pod/pod-name\"")
		} else if len(parts) == 1 { // Assume pod-name
			return parts[0], nil
		}
		return parts[1], nil // pod/pod-name
	case 2:
		return args[1], nil // pod(s) pod-name
	default:
		return "", errors.New("invalid number of arguments. Expect 1 or 2 arguments: \"pod/pod-name\", or \"pod pod-name\"")
	}
}
