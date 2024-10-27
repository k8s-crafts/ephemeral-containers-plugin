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

package testutils

import (
	"fmt"

	"github.com/blang/semver"
	"k8s.io/client-go/kubernetes"

	kbutils "sigs.k8s.io/kubebuilder/v4/test/e2e/utils"
)

var (
	MinK8sVersion semver.Version = semver.MustParse("v1.25.0")
)

type TestResource struct {
	// For invoking kubectl command
	Kubectl    *kbutils.Kubectl
	Client     *kubernetes.Clientset
	K8sVersion *kbutils.KubernetesVersion
	PluginName string
}

func NewTestResource() (*TestResource, error) {
	kubectl := &kbutils.Kubectl{
		CmdContext: &kbutils.CmdContext{},
		Namespace:  fmt.Sprintf("%s-%s", "e2e", generateRandom(4)),
	}

	k8sVersion, err := kubectl.Version()
	if err != nil {
		return nil, err
	}

	return &TestResource{
		Kubectl:    kubectl,
		K8sVersion: &k8sVersion,
		PluginName: "ephemeral-containers",
	}, nil
}

// ephemeralcontainer subresource is supported on Kubernetes >= v1.25
func (t *TestResource) IsKubeAPICompatible() bool {
	return semver.MustParse(t.K8sVersion.ServerVersion.GitVersion).GTE(MinK8sVersion)
}

func (t *TestResource) CreateNamespace() error {
	_, err := t.Kubectl.Command("create", "namespace", t.Kubectl.Namespace)
	return err
}

func (t *TestResource) DeleteNamespace() error {
	_, err := t.Kubectl.Delete(false, "namespace", t.Kubectl.Namespace)
	return err
}
