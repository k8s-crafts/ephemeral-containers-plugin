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
	"strings"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
)

type PodFilterFn func(pod corev1.Pod) bool

type KubeClientset struct {
	kubernetes.Interface
}

// List pods by filters in the specified namespace
// If namespace is empty (i.e. ""), list in all namespaces
func (client *KubeClientset) ListPods(ctx context.Context, namespace string, filters ...PodFilterFn) ([]corev1.Pod, error) {
	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ApplyPodFilter(podList.Items, filters...), nil
}

// Get pod by name in a specific namespace
func (client *KubeClientset) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return setGVK(pod), nil
}

// Update pod's ephemeralContainer subresource
func (client *KubeClientset) UpdateEphemeralContainersForPod(ctx context.Context, pod *corev1.Pod) (*corev1.Pod, error) {
	return client.CoreV1().Pods(pod.Namespace).UpdateEphemeralContainers(ctx, pod.Name, pod, metav1.UpdateOptions{})
}

// Explicitly set GVK for Pod
// See: https://github.com/kubernetes/kubernetes/issues/80609
func setGVK(pod *corev1.Pod) *corev1.Pod {
	gvk := schema.GroupVersionKind{
		Kind:    "Pod",
		Version: "v1",
		Group:   "",
	}
	pod.SetGroupVersionKind(gvk)
	return pod
}

// Validate the pod struct from edited manifests
func SanitizeEditedPod(original, edited *corev1.Pod) (*corev1.Pod, error) {
	if !cmp.Equal(original.Name, edited.Name) {
		return nil, fmt.Errorf("pod's name cannot be changed. Expected %s but got %s", original.Name, edited.Name)
	}

	if !cmp.Equal(original.Namespace, edited.Namespace) {
		return nil, fmt.Errorf("pod's namespace cannot be changed. Expected %s but got %s", original.Namespace, edited.Namespace)
	}

	// Nothing changes in spec.ephemeralContainers
	if cmp.Equal(original.Spec.EphemeralContainers, edited.Spec.EphemeralContainers) {
		return nil, nil
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      original.Name,
			Namespace: original.Namespace,
		},
		Spec: corev1.PodSpec{
			EphemeralContainers: edited.Spec.DeepCopy().EphemeralContainers,
		},
	}, nil
}

func MinifyPod(pod *corev1.Pod) *corev1.Pod {
	result := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		Spec: corev1.PodSpec{
			EphemeralContainers: pod.Spec.DeepCopy().EphemeralContainers,
		},
	}

	return setGVK(result)
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
