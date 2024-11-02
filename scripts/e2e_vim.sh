#!/usr/bin/env bash
# This is a custom editor used for e2e tests
# Requirements
# - yq (v4)
# Note: The patches are predefined

if ! command -v yq >/dev/null 2>&1; then
    echo "yq is required, but missing. Install guide: https://mikefarah.gitbook.io/yq#install"
    exit 1
fi

# Path to temporary file containing Pod's manifest
pod_manifest="$1"

# Constants
export container_image=${container_image:-"docker.io/library/busybox:1.28"}
export container_name=${container_name:-"debugger-$RANDOM"}
export other_container_image=${other_container_image:-"docker.io/library/busybox:1.27"}

# Path to patch to apply to Pod manifest
case "${E2E_EDIT_ACTION:-""}" in
    add)
        # Add a new container
        container_spec="{\"image\": \"$container_image\", \"name\": \"$container_name\", \"imagePullPolicy\": \"IfNotPresent\", \"stdin\": true, \"tty\": true}"
        yq -i ".spec.ephemeralContainers |= . + [$container_spec]" "$pod_manifest"
    ;;
    delete)
        # Delete a container (assuming one is available)
        yq -i 'del(.spec.ephemeralContainers[0])' "$pod_manifest"
    ;;
    modify)
        # Modify container image (assuming one is available)
        yq -i '.spec.ephemeralContainers[0].image |= env(other_container_image)' "$pod_manifest"
    ;;
    *)
        exit 0
    ;;
esac
