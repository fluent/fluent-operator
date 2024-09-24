#!/usr/bin/env bash

# copied from: https://github.com/weaveworks/flagger/tree/main/hack

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
chmod +x ${CODEGEN_PKG}/kube_codegen.sh 

PACKAGE_PATH_BASE="github.com/fluent/fluent-operator/v3"

mkdir -p "${TEMP_DIR}/${PACKAGE_PATH_BASE}/apis/fluentbit" \
         "${TEMP_DIR}/${PACKAGE_PATH_BASE}/apis/fluentd" \
         "${TEMP_DIR}/${PACKAGE_PATH_BASE}/apis/generated"

source ${CODEGEN_PKG}/kube_codegen.sh kube::codegen::gen_client \
    --output-dir "${TEMP_DIR}" \
    --with-watch \
    --output-pkg "${PACKAGE_PATH_BASE}/apis/generated" \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
    ./apis

ls -lha $TEMP_DIR

# Copy everything back.
cp -r "${TEMP_DIR}/${PACKAGE_PATH_BASE}/." "${SCRIPT_ROOT}" 
