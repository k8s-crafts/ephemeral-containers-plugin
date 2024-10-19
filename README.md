# kubectl ephemeral-containers

![Release](https://img.shields.io/badge/Version-v1.1.0-informational?style=flat-square&label=Release)
![Go version](https://img.shields.io/github/go-mod/go-version/k8s-crafts/ephemeral-containers-plugin?style=flat-square)
[![Lint](https://img.shields.io/github/actions/workflow/status/k8s-crafts/ephemeral-containers-plugin/lint.yaml?style=flat-square&logo=github&label=Lint)](https://github.com/k8s-crafts/ephemeral-containers-plugin/actions/workflows/lint.yaml)
[![Test](https://img.shields.io/github/actions/workflow/status/k8s-crafts/ephemeral-containers-plugin/test.yaml?style=flat-square&logo=github&label=Test)
](https://github.com/k8s-crafts/ephemeral-containers-plugin/actions/workflows/test.yaml)

A `kubectl` plugin for directly modifying Pods' [`ephemeralContainers` spec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#ephemeralcontainer-v1-core).

An [ephemeral container](https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/) is a temporary container that is injected into existing Pods for some user-initiated actions, for example, troubleshooting. However, ephemeral container specs must be handled as a Pod's `ephemeralcontainers` subresource. Consequently, `kubectl edit` cannot be used for modifying `pod.spec.ephemeralcontainers`.

This project is a plugin (i.e. an extension) to `kubectl` that allows such direct editing, bringing back the experience of `kubectl edit`. For more information on how to extend `kubectl`, see [guides](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

## Requirements

- Kubernetes `>= v1.25.0-0`
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)

## Installation

### From released binaries

Download the binary from the [release page](https://github.com/k8s-crafts/ephemeral-containers-plugin/releases) and copy it to a directory on `PATH`.

### From source

In addition to released binaries, you can install the plugin from source.

#### With `go install`

```bash
go install github.com/k8s-crafts/ephemeral-containers-plugin@latest
```

**Note**: The binary will be installed under `GOBIN` (i.e. `go env GOBIN`) as `ephemeral-containers-plugin`. It must be renamed to `kubectl-ephemeral_containers` to register the plugin with `kubectl`.

#### With source repository

```bash
git clone git@github.com:k8s-crafts/ephemeral-containers-plugin.git
cd ephemeral-containers-plugin
# Build the plugin and install it on PATH
# Default to $HOME/bin
make build install
```

## User Guides

For details on how to use the plugin, see [User Guides](./docs/USER_GUIDES.md).

## Contributing

For details on how to contribute, see [Contributing Guides](./CONTRIBUTING.md) and [Developer Guides](./docs/DEVELOPER_GUIDES.md).

## License

The project follows [Apache License 2.0](./LICENSE) license.
