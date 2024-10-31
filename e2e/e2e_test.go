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

package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("kubectl ephemeral-containers", func() {
	BeforeEach(func() {
		Expect(tr.CreateTestPod()).ToNot(HaveOccurred())
		Expect(tr.WaitForPodReady()).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(tr.DeleteTestPod()).ToNot(HaveOccurred())
	})

	Context("list", func() {
		It("should return empty message if none", func() {
			By("checking checking pod ephemeralContainers spec")

			actual, err := tr.RunPluginListCmd()
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(fmt.Sprintf("No pods with ephemeral containers found in namespace %s\n", tr.Kubectl.Namespace)))
		})
		It("should return the list if any", func() {
			Expect(tr.RunDebugContainer(true)).ToNot(HaveOccurred())

			By("checking checking pod ephemeralContainers spec")
			expected := fmt.Sprintf(
				`+------------+-----------+----------------------+
|    POD     | NAMESPACE | EPHEMERAL CONTAINERS |
+------------+-----------+----------------------+
| plugin-e2e | %s  | debugger             |
+------------+-----------+----------------------+

`, tr.Kubectl.Namespace,
			)

			actual, err := tr.RunPluginListCmd()
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})
	})
})
