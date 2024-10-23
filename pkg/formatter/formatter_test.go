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

package formatter_test

import (
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/formatter"
	"github.com/k8s-crafts/ephemeral-containers-plugin/pkg/version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Formatter", func() {
	var t *test

	Context("when formatting version", func() {
		BeforeEach(func() {
			t = newTest()
		})
		It("should return as table", func() {
			content, err := formatter.FormatVersionOutput(formatter.Table, t.version)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(t.versionTable))
		})

		It("should return as JSON", func() {
			content, err := formatter.FormatVersionOutput(formatter.JSON, t.version)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(t.versionJSON))
		})

		It("should return as YAML", func() {
			content, err := formatter.FormatVersionOutput(formatter.YAML, t.version)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(t.versionYAML))
		})
	})

	Context("when listing ephemeral containers for pod", func() {
		Context("with ephemeral containers", func() {
			BeforeEach(func() {
				t = newTestForPodWithEphemeralContainers()
			})

			It("should return names", func() {
				containers := formatter.ListEphemeralContainersForPod(t.pod)
				t.expectEphemeralContainers(containers)
			})
		})
		Context("without ephemeral containers", func() {
			BeforeEach(func() {
				t = newTestForPodWithoutEphemeralContainers()
			})

			It("should return empty", func() {
				containers := formatter.ListEphemeralContainersForPod(t.pod)
				t.expectEphemeralContainers(containers)
			})
		})
	})

	Context("when formatting pod list", func() {
		Context("with ephemeral containers", func() {
			BeforeEach(func() {
				t = newTestForPodWithEphemeralContainers()
			})

			It("should return as table", func() {
				content, err := formatter.FormatListOutput(formatter.Table, []corev1.Pod{t.pod})
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(t.listTable))
			})

			It("should return as JSON", func() {
				content, err := formatter.FormatListOutput(formatter.JSON, []corev1.Pod{t.pod})
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(t.listYAML))
			})

			It("should return as YAML", func() {
				content, err := formatter.FormatListOutput(formatter.YAML, []corev1.Pod{t.pod})
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(t.listJSON))
			})
		})
		Context("without ephemeral containers", func() {
			BeforeEach(func() {
				t = newTestForPodWithoutEphemeralContainers()
			})

			It("should return as table", func() {
				content, err := formatter.FormatListOutput(formatter.Table, make([]corev1.Pod, 0))
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(t.listTable))
			})

			It("should return as JSON", func() {
				content, err := formatter.FormatListOutput(formatter.JSON, make([]corev1.Pod, 0))
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(t.listYAML))
			})

			It("should return as YAML", func() {
				content, err := formatter.FormatListOutput(formatter.YAML, make([]corev1.Pod, 0))
				Expect(err).ToNot(HaveOccurred())
				Expect(content).To(Equal(t.listJSON))
			})
		})
	})
})

type testInput struct {
	pod        corev1.Pod
	containers []string

	listTable string
	listJSON  string
	listYAML  string

	version *version.VersionInfo

	versionTable string
	versionJSON  string
	versionYAML  string
}

type test struct {
	*testInput
}

func (t *test) expectEphemeralContainers(containers []string) {
	Expect(t.containers).To(ConsistOf(containers))
}

func newTest() *test {
	return &test{
		testInput: &testInput{
			version: &version.VersionInfo{
				Version: "v0.0.0-unset",
			},
			versionTable: "version: v0.0.0-unset",
			versionJSON: `{
  "version": "v0.0.0-unset"
}`,
			versionYAML: "version: v0.0.0-unset\n",
		},
	}
}

func newTestForPodWithEphemeralContainers() *test {
	t := newTest()
	t.containers = []string{"debug-container", "another-one"}
	t.pod = corev1.Pod{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			EphemeralContainers: []corev1.EphemeralContainer{
				{
					EphemeralContainerCommon: corev1.EphemeralContainerCommon{
						Name:  "debug-container",
						Image: "my-image:v1",
					},
				},
				{
					EphemeralContainerCommon: corev1.EphemeralContainerCommon{
						Name:  "another-one",
						Image: "my-image-1:v2",
					},
				},
			},
		},
	}
	t.listTable = `+--------+-----------+-----------------------------+
|  POD   | NAMESPACE |    EPHEMERAL CONTAINERS     |
+--------+-----------+-----------------------------+
| my-pod | default   | debug-container,another-one |
+--------+-----------+-----------------------------+
`
	t.listYAML = `[
  {
    "name": "my-pod",
    "namespace": "default",
    "ephemeralContainers": [
      "debug-container",
      "another-one"
    ]
  }
]`

	t.listJSON = `- ephemeralContainers:
  - debug-container
  - another-one
  name: my-pod
  namespace: default
`
	return t

}

func newTestForPodWithoutEphemeralContainers() *test {
	t := newTest()
	t.pod = corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{},
	}
	return t
}
