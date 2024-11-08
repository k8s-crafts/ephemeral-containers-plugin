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

package k8s_test

import (
	"context"
	"os"
	"path"

	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/k8s"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("K8s", func() {
	var t *test

	BeforeEach(func() {
		t = newTest()
	})

	JustBeforeEach(func() {
		for _, ns := range t.namespaces {

			testNs := t.newNamespace(ns)
			testPod := t.newPod("testpod", ns)

			_, err := t.clientset.CoreV1().Namespaces().Create(context.Background(), testNs, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())

			_, err = t.clientset.CoreV1().Pods(ns).Create(context.Background(), testPod, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
		}
	})

	When("creating Kubernetes clientset", func() {
		BeforeEach(func() {
			t.setEnvKubeConfig()
		})

		AfterEach(func() {
			t.UnsetEnvKubeConfig()
		})

		It("should create a non-nil kubeconfig", func() {
			kubeConfig := k8s.NewKubeConfig()
			t.expectKubeConfig(kubeConfig)
		})

		It("should create a clientset", func() {
			kubeConfig := k8s.NewKubeConfig()
			t.expectKubeConfig(kubeConfig)

			clientset, err := k8s.NewClientset(kubeConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(clientset).ToNot(BeNil())
		})
	})

	When("listing pods", func() {
		Context("in a namespace", func() {
			It("should return pods", func() {
				pods, err := t.clientset.ListPods(context.Background(), t.namespaces[0])
				Expect(err).ToNot(HaveOccurred())
				Expect(pods).To(HaveLen(1))
			})
		})
		Context("in all namespaces", func() {
			It("should return pods", func() {
				pods, err := t.clientset.ListPods(context.Background(), "")
				Expect(err).ToNot(HaveOccurred())
				Expect(pods).To(HaveLen(2))
			})
		})
	})

	When("getting a pod", func() {
		It("should return the pod", func() {
			pod, err := t.clientset.GetPod(context.Background(), t.namespaces[0], "testpod")
			Expect(err).ToNot(HaveOccurred())
			Expect(pod).ToNot(BeNil())

			Expect(pod.Kind).To(Equal("Pod"))
			Expect(pod.Name).To(Equal("testpod"))
			Expect(pod.Namespace).To(Equal(t.namespaces[0]))

		})
	})

	When("updating the ephemeralcontainer subresource for a pod", func() {
		It("should update", func() {
			pod, err := t.clientset.GetPod(context.Background(), t.namespaces[0], "testpod")
			Expect(err).ToNot(HaveOccurred())
			Expect(pod).ToNot(BeNil())

			newCont := corev1.EphemeralContainer{
				EphemeralContainerCommon: corev1.EphemeralContainerCommon{
					Name:  "debugger",
					Image: "busybox:1.27",
				},
			}

			pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, newCont)

			pod, err = t.clientset.UpdateEphemeralContainersForPod(context.Background(), pod)
			Expect(err).ToNot(HaveOccurred())
			Expect(pod).ToNot(BeNil())

			Expect(pod.Spec.EphemeralContainers).To(ContainElement(newCont))
		})
	})
})

type testInput struct {
	clientset  *k8s.KubeClientset
	namespaces []string
}

type test struct {
	*testInput
}

func newTest() *test {
	return &test{
		testInput: &testInput{
			namespaces: []string{"test-ns-0", "test-ns-1"},
			clientset: &k8s.KubeClientset{
				Interface: fake.NewClientset(),
			},
		},
	}
}

func (t *test) newNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (t *test) newPod(name, namespace string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			EphemeralContainers: []corev1.EphemeralContainer{
				{
					EphemeralContainerCommon: corev1.EphemeralContainerCommon{
						Name:  "debugger",
						Image: "busybox:1.28",
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "main",
					Image: "registry.k8s.io/pause:3.1",
				},
			},
		},
	}
}

func (t *test) expectKubeConfig(kubeConfig *k8s.KubeConfig) {
	Expect(kubeConfig).ToNot(BeNil())
	Expect(kubeConfig.ConfigFlags).ToNot(BeNil())
}

func (t *test) setEnvKubeConfig() {
	Expect(os.Setenv("KUBECONFIG", path.Join("testdata", "kubeconfig"))).ToNot(HaveOccurred())

}

func (t *test) UnsetEnvKubeConfig() {
	Expect(os.Unsetenv("KUBECONFIG")).ToNot(HaveOccurred())
}
