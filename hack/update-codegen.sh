#!/usr/bin/env bash

# copied from: https://github.com/weaveworks/flagger/tree/master/hack

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(git rev-parse --show-toplevel)

# Grab code-generator version from go.sum.
CODEGEN_VERSION=$(grep 'k8s.io/code-generator' go.sum | awk '{print $2}' | head -1 |cut -b 1-7)
CODEGEN_PKG=$(echo `go env GOPATH`"/pkg/mod/k8s.io/code-generator@${CODEGEN_VERSION}")

# code-generator does work with go.mod but makes assumptions about
# the project living in `$GOPATH/src`. To work around this and support
# any location; create a temporary directory, use this as an output
# base, and copy everything back once generated.
TEMP_DIR=$(mktemp -d)
cleanup() {
    echo ">> Removing ${TEMP_DIR}"
    rm -rf ${TEMP_DIR}
}
trap "cleanup" EXIT SIGINT

echo ">> Temporary output directory ${TEMP_DIR}"

# Ensure we can execute.
chmod +x ${CODEGEN_PKG}/generate-groups.sh

${CODEGEN_PKG}/generate-groups.sh "client" \
    kubesphere.io/fluentbit-operator/apis/generated/kubesphere kubesphere.io/fluentbit-operator/apis \
    kubesphere.io:v1alpha2 \
    --output-base "${TEMP_DIR}" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt

${CODEGEN_PKG}/generate-groups.sh "client" \
    kubesphere.io/fluentbit-operator/apis/generated/fluentbit kubesphere.io/fluentbit-operator/apis \
    fluentbit.io:v1alpha2 \
    --output-base "${TEMP_DIR}" \
    --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt

# Copy everything back.
cp -r "${TEMP_DIR}/kubesphere.io/fluentbit-operator/." "${SCRIPT_ROOT}/"