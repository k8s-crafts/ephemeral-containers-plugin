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
	"testing"

	. "github.com/k8s-crafts/ephemeral-containers-plugin/e2e/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var (
	tr *TestResource
)

var _ = BeforeSuite(func() {
	var err error

	// tr is defined globally
	tr, err = NewTestResource()
	Expect(err).ToNot(HaveOccurred())

	By("checking if kube API meets minimum required version")
	// Check if the kube API is supported
	Expect(tr.IsKubeAPICompatible()).To(BeTrue())

	// Create resources for tests
	for _, ns := range tr.GetTestNamespaces() {
		By(fmt.Sprintf("creating namespace %s and test resources", ns))

		Expect(tr.CreateNamespace(ns)).ToNot(HaveOccurred())
		Expect(tr.CreateServiceAccount(ns)).ToNot(HaveOccurred())
	}
})

var _ = AfterSuite(func() {
	for _, ns := range tr.GetTestNamespaces() {
		By(fmt.Sprintf("deleting namespace %s", ns))

		Expect(tr.DeleteNamespace(ns)).ToNot(HaveOccurred())
	}
})
