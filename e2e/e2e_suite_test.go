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
