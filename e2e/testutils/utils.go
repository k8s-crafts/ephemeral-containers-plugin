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
	"strings"
	"time"

	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	kbutils "sigs.k8s.io/kubebuilder/v4/test/e2e/utils"
)

type TestResource struct {
	// For invoking kubectl command
	Kubectl    *kbutils.Kubectl
	Client     *kubernetes.Clientset
	K8sVersion *kbutils.KubernetesVersion
	// Another namespace besides namespace in context
	AnotherNamespace string
	PluginName       string
}

func NewTestResource() (*TestResource, error) {
	kubectl := &kbutils.Kubectl{
		CmdContext: &kbutils.CmdContext{
			Env: []string{
				"EDITOR=vi",
				"KUBE_EDITOR=",
			},
		},
		Namespace: fmt.Sprintf("%s-%s", "e2e0", generateRandom(4)),
	}

	k8sVersion, err := kubectl.Version()
	if err != nil {
		return nil, err
	}

	return &TestResource{
		Kubectl:          kubectl,
		K8sVersion:       &k8sVersion,
		PluginName:       "ephemeral-containers",
		AnotherNamespace: fmt.Sprintf("%s-%s", "e2e1", generateRandom(4)),
	}, nil
}

func (t *TestResource) GetTestNamespaces() []string {
	return []string{t.Kubectl.Namespace, t.AnotherNamespace}
}

func (t *TestResource) SetEnv(key, value string) {
	// Last occurence takes precendence
	t.Kubectl.Env = append(t.Kubectl.Env, fmt.Sprintf("%s=%s", key, value))
}

func (t *TestResource) UnsetEnv(key string) {
	envs := t.Kubectl.Env
	for i, _env := range envs {
		if strings.Contains(_env, fmt.Sprintf("%s=", key)) {
			envs = append(envs[:i], envs[i+1:]...)
		}
	}
	t.Kubectl.Env = envs
}

// ephemeralcontainer subresource is supported on Kubernetes >= v1.25
func (t *TestResource) IsKubeAPICompatible() bool {
	return semver.Compare(t.K8sVersion.ServerVersion.GitVersion, MinK8sVersion) >= 0
}

func (t *TestResource) CreateNamespace(namespace string) error {
	_, err := t.Kubectl.Command("create", "namespace", namespace)
	return err
}

func (t *TestResource) DeleteNamespace(namespace string) error {
	_, err := t.Kubectl.Delete(false, "namespace", namespace)
	return err
}

func (t *TestResource) CreateServiceAccount(namespace string) error {
	_, err := t.Kubectl.Apply(true, "-n", namespace, "-f", path.Join(getTestdataDir(), "serviceaccount.yaml"))
	return err
}

func (t *TestResource) CreateTestPod(namespace string) error {
	_, err := t.Kubectl.Apply(true, "-n", namespace, "-f", path.Join(getTestdataDir(), "pod.yaml"))
	return err
}

func (t *TestResource) DeleteTestPod(namespace string) error {
	_, err := t.Kubectl.Delete(false, "-n", namespace, "--ignore-not-found=true", "-f", path.Join(getTestdataDir(), "pod.yaml"))
	return err
}

func (t *TestResource) ListEphemeralContainerNamesForTestPod(namespace string) (string, error) {
	return t.Kubectl.CommandInNamespace("get", "-n", namespace, fmt.Sprintf("pods/%s", TestPodName), "-o=jsonpath='{.spec.ephemeralContainers[*].name}'")
}

// Add an ephemeral containers by kubectl debug
// If interactive, it is not attached
func (t *TestResource) RunDebugContainerForTestPod(namespace, containerName string, interactive bool) error {
	args := []string{
		"debug",
		"-n", namespace,
		fmt.Sprintf("pods/%s", TestPodName),
		fmt.Sprintf("--image=%s", DebugImage),
		fmt.Sprintf("--container=%s", containerName),
	}

	if interactive {
		args = append(args, "-it", "--attach=false")
	}

	_, err := t.Kubectl.CommandInNamespace(args...)
	return err
}

func (t *TestResource) RunPluginHelpCmd(subCmd string) (string, error) {
	return t.Kubectl.Command(t.PluginName, "help", subCmd)
}

func (t *TestResource) RunPluginListCmd(format string, namespace string) (string, error) {
	// Setting namespace manually as flags cannot be set before plugin name
	args := []string{t.PluginName, "list"}

	if len(namespace) > 0 {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "--all-namespaces")
	}

	if len(format) > 0 {
		args = append(args, "-o", format)
	}

	return t.Kubectl.Command(args...)
}

func (t *TestResource) RunPluginEditCmd(namespace string, podName string) (string, error) {
	// Setting namespace manually as flags cannot be set before plugin name
	return t.Kubectl.Command(t.PluginName, "edit", "-n", namespace, podName)
}

func (t *TestResource) WaitForTestPodReady(namespace string) error {
	return wait.PollUntilContextCancel(context.TODO(), time.Second, true, func(ctx context.Context) (done bool, err error) {
		_, err = t.Kubectl.CommandInNamespace("wait", "-n", namespace, "--for=condition=Ready", fmt.Sprintf("pods/%s", TestPodName))
		if err != nil {
			return false, err
		}
		return true, nil
	})
}
