#!/usr/bin/env bash

# copied from: https://github.com/weaveworks/flagger/tree/master/hack

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(git rev-parse --show-toplevel)
CODEGEN_PKG=$(go list -m -f '{{.Dir}}' k8s.io/code-generator)

# code-generator does work with go.mod but makes assumptions about
# the project living in `$GOPATH/src`. To work around this and support
# any location; create a temporary directory, use this as an output
# base, and copy everything back once generated.
TEMP_DIR=$(mktemp -d)
cleanup() {
    echo ">> Removing ${TEMP_DIR}"
    rm -rf "${TEMP_DIR}"
}
trap "cleanup" EXIT SIGINT

echo ">> Temporary output directory ${TEMP_DIR}"

# Ensure we can execute.
chmod +x "${CODEGEN_PKG}"/generate-groups.sh

"${CODEGEN_PKG}"/generate-groups.sh "client" \
    github.com/fluent/fluent-operator/apis/generated github.com/fluent/fluent-operator/apis \
    "fluentbit:v1alpha2 fluentd:v1alpha1" \
    --output-base "${TEMP_DIR}" \
    --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

# Copy everything back.
cp -r "${TEMP_DIR}/github.com/fluent/fluent-operator/apis/." "${SCRIPT_ROOT}/apis/"
