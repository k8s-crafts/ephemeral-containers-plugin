#!/usr/bin/env bash
# This utitility script is used to set up environment and tools for E2E tests

# Defaults
DEFAULT_KUBECTL_VERSION=v1.30.0
DEFAULT_KIND_VERSION=v0.24.0
DEFAULT_CLUSTER_NAME=ephcont-e2e

BIN="$(pwd)/testbin"

# Entry point
main() {
    local kubectl_version="$DEFAULT_KUBECTL_VERSION"
    local kind_version="$DEFAULT_KIND_VERSION"
    local cluster_name="$DEFAULT_CLUSTER_NAME"
    local config=
    local arch=

    case "$(uname -m)" in
        x86_64)
            arch=amd64
        ;;
        aarch64)
            arch=arm64
        ;;
        *)
            echo "ERROR: $(uname -m) is not supported"
            exit 2
        ;;
    esac

    parse_cli_args "$@"

    echo "Installing kind..."
    install_kind

    echo "Installing kubectl..."
    install_kubectl

    echo "Creating KinD cluster..."
    create_cluster
}

# Parse command line flags
parse_cli_args() {
    while :; do
        case "${1:-}" in
            --cluster-name)
                if [ -n "${2:-}" ]; then
                    cluster_name="$2"
                    shift
                else
                    echo "ERROR: missing cluster name for '--cluster-name'"
                    exit 2
                fi
            ;;
            --kind-version)
                if [ -n "${2:-}" ]; then
                    kind_version="$2"
                    shift
                else
                    echo "ERROR: missing kind version for '--kind-version'"
                    exit 2
                fi
            ;;          
            --kubectl-version)
                if [ -n "${2:-}" ]; then
                    kubectl_version="$2"
                    shift
                else
                    echo "ERROR: missing kubectl version for '--kubectl-version'"
                    exit 2
                fi
            ;;
            --config)
                if [ -n "${2:-}" ]; then
                    config="$2"
                    shift
                else
                    echo "ERROR: missing config file path for '--kubectl-version'"
                    exit 2
                fi
            ;;
            *)
                break
            ;;
        esac
        shift
    done
}

install_kind() {
    if check_tool_available "$BIN/kind" $kind_version; then
        echo "kind $kind_version is available."
    else
        mkdir -p "$BIN"

        curl -L -o "$BIN/kind" "https://kind.sigs.k8s.io/dl/$kind_version/kind-linux-$arch"
        chmod +x "$BIN/kind"
    fi
}

install_kubectl() {
    if check_tool_available "$BIN/kubectl" $kubectl_version; then
        echo "kubectl $kubectl_version is available."
    else
        mkdir -p "$BIN"

        curl -L -o "$BIN/kubectl" https://dl.k8s.io/release/$kubectl_version/bin/linux/$arch/kubectl
        chmod +x "$BIN/kubectl"
    fi
}

# Create a KinD cluster
# and write client configurations to testbin/kubeconfig
create_cluster() {
    KUBECONFIG="$BIN/kubeconfig" "$BIN/kind" create cluster --config "$config" --name "$cluster_name"
}

check_tool_available() {
    local tool="$1"
    local version="$2"

    # If the tool path exists (and executable) and the version matches
    if test -x "$tool" && grep "$version" <<<"$($tool version 2>/dev/null)" >/dev/null 2>&1; then
        return 0
    fi
    return 1
}

main "$@"
