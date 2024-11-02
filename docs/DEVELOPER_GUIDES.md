## Developer Guides

The guide describes workflow for developing the plugin. Feel free to suggest any improvement to ensure a seamless on-boarding for new contributors.

### Requirements

- Linux distros (e.g. Fedora, Ubuntu)
- [Go](https://go.dev/doc/install) `v1.22`
- Make
- [yq](https://mikefarah.gitbook.io/yq#install) `v4+`
- Kubernetes `>= v1.25.0-0` (e.g. [minikube](https://minikube.sigs.k8s.io/docs/) or [CRC](https://developers.redhat.com/products/openshift-local/getting-started))

### Fork the repository

Before you begin, please fork the repository at https://github.com/k8s-crafts/ephemeral-containers-plugin. You can follow [GitHub guide](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo) on how to accomplish this.

### Clone the repository

Once the fork is created, run the following (i.e. replace `<github-username>` with yours) to create a local copy of the repository:

```bash
git clone git@github.com:<github-username>/ephemeral-containers-plugin.git
```

### Build and install the plugin (development version)

To build the plugin, run:

```
make build
```

An executable will be created at `build/kubectl-ephemeral_containers`. To install the plugin into `PATH`, run:

```bash
make install
```

Then, the plugin will be registered with `kubectl`. You can run `kubectl ephemeral-containers help` to validate.

### Run tests

To execute unit tests, run:

```bash
make test
```

A coverage report will created at `cover.out`.

To execute E2E tests, run:

```bash
make test-e2e
```

### Run lints

The project uses [golangci-lint](https://golangci-lint.run/) for running linters. To run linters, run:

```bash
make lint-fix
```

This will apply lint fixes. If you'd like to only report lint issues, run:

```bash
make lint
```

### Add license headers

The source files are required to have license headers. To add license headers if missing, run:

```bash
make add-license
```

If you are running `make fmt`, the license headers are also added if missing.

### Add a Go package

To add a new package, for example, `github.com/google/go-cmp/cmp`, in the project, run:

```bash
go get -u github.com/google/go-cmp/cmp
```

This will fetch the package and modify `go.mod` and `go.sum`. Then, you import it in the source code.

**Notes**:

- Make sure to run `go mod tidy` regularly to remove any unused packages.
- Make sure to justify the additions of new packages to maintainers.
