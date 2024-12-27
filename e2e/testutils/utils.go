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

package testutils

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/util/wait"
	kbutils "sigs.k8s.io/kubebuilder/v4/test/e2e/utils"
)

type TestResource struct {
	// For invoking kubectl command
	*kbutils.Kubectl
	*kbutils.KubernetesVersion
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
		Kubectl:           kubectl,
		KubernetesVersion: &k8sVersion,
		PluginName:        "ephemeral-containers",
		AnotherNamespace:  fmt.Sprintf("%s-%s", "e2e1", generateRandom(4)),
	}, nil
}

func (t *TestResource) GetTestNamespaces() []string {
	return []string{t.Namespace, t.AnotherNamespace}
}

func (t *TestResource) SetEnv(key, value string) {
	// Last occurence takes precendence
	t.Env = append(t.Env, fmt.Sprintf("%s=%s", key, value))
}

func (t *TestResource) UnsetEnv(key string) {
	envs := t.Env
	for i, _env := range envs {
		if strings.Contains(_env, fmt.Sprintf("%s=", key)) {
			envs = append(envs[:i], envs[i+1:]...)
		}
	}
	t.Env = envs
}

// ephemeralcontainer subresource is supported on Kubernetes >= v1.25
func (t *TestResource) IsKubeAPICompatible() bool {
	return semver.Compare(t.ServerVersion.GitVersion, MinK8sVersion) >= 0
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
	_, err := t.Kubectl.Apply(false, "-n", namespace, "-f", path.Join(getTestdataDir(), "serviceaccount.yaml"))
	return err
}

func (t *TestResource) CreateTestPod(namespace string) error {
	_, err := t.Kubectl.Apply(false, "-n", namespace, "-f", path.Join(getTestdataDir(), "pod.yaml"))
	return err
}

func (t *TestResource) DeleteTestPod(namespace string) error {
	_, err := t.Kubectl.Delete(false, "-n", namespace, "--ignore-not-found=true", "-f", path.Join(getTestdataDir(), "pod.yaml"))
	return err
}

func (t *TestResource) ListEphemeralContainerNamesForTestPod(namespace string) (string, error) {
	return t.Kubectl.Command("get", "-n", namespace, fmt.Sprintf("pods/%s", TestPodName), "-o=jsonpath='{.spec.ephemeralContainers[*].name}'")
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

	_, err := t.Kubectl.Command(args...)
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
		_, err = t.Kubectl.Command("wait", "-n", namespace, "--for=condition=Ready", fmt.Sprintf("pods/%s", TestPodName))
		if err != nil {
			return false, err
		}
		return true, nil
	})
}
