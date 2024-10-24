#!/usr/bin/env bash
# This utitility script is used to clean up environment and tools for E2E tests
# This should be called after e2e_setup.sh

DEFAULT_CLUSTER_NAME=ephcont-e2e

BIN="$(pwd)/testbin"

main() {
    local cluster_name="$DEFAULT_CLUSTER_NAME"

    parse_cli_args

    echo "Deleting KinD cluster..."
    cleanup || true
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
            *)
                break
            ;;
        esac
        shift
    done
}

cleanup() {
    $BIN/kind delete cluster --name $cluster_name
    rm -rf "$BIN"
}

main "$@"
