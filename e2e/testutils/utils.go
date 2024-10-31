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
	"context"
	"fmt"
	"path"
	"time"

	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	kbutils "sigs.k8s.io/kubebuilder/v4/test/e2e/utils"
)

var (
	MinK8sVersion string = "v1.25.0"
	TestPodName   string = "plugin-e2e"
	DebugImage    string = "docker.io/library/busybox:1.28"
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
	return semver.Compare(t.K8sVersion.ServerVersion.GitVersion, MinK8sVersion) >= 0
}

func (t *TestResource) CreateNamespace() error {
	_, err := t.Kubectl.Command("create", "namespace", t.Kubectl.Namespace)
	return err
}

func (t *TestResource) DeleteNamespace() error {
	_, err := t.Kubectl.Delete(false, "namespace", t.Kubectl.Namespace)
	return err
}

func (t *TestResource) CreateServiceAccount() error {
	_, err := t.Kubectl.Apply(true, "-f", path.Join(getTestdataDir(), "serviceaccount.yaml"))
	return err
}

func (t *TestResource) CreateTestPod() error {
	_, err := t.Kubectl.Apply(true, "-f", path.Join(getTestdataDir(), "pod.yaml"))
	return err
}

func (t *TestResource) DeleteTestPod() error {
	_, err := t.Kubectl.Delete(false, "--ignore-not-found=true", "-f", path.Join(getTestdataDir(), "pod.yaml"))
	return err
}

// Add an ephemeral containers by kubectl debug
// If interactive, it is not attached
func (t *TestResource) RunDebugContainer(interactive bool) error {
	args := []string{
		"debug",
		fmt.Sprintf("pods/%s", TestPodName),
		fmt.Sprintf("--image=%s", DebugImage),
		"--container=debugger",
	}

	if interactive {
		args = append(args, "-it", "--attach=false")
	}

	_, err := t.Kubectl.CommandInNamespace(args...)
	return err
}

func (t *TestResource) RunPluginListCmd() (string, error) {
	// Setting namespace manually as flags cannot be set before plugin name
	return t.Kubectl.Command(t.PluginName, "list", "-n", t.Kubectl.Namespace)
}

func (t *TestResource) WaitForPodReady() error {
	return wait.PollUntilContextCancel(context.TODO(), time.Second, true, func(ctx context.Context) (done bool, err error) {
		// kubectl wait --for=condition=Ready=false pod/busybox1
		_, err = t.Kubectl.CommandInNamespace("wait", "--for=condition=Ready", fmt.Sprintf("pods/%s", TestPodName))
		if err != nil {
			return false, err
		}
		return true, nil
	})
}
