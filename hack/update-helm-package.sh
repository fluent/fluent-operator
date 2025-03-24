#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

FLUENT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
FLUENT_HELM_DIR=${FLUENT_ROOT}/charts/fluent-operator
_tmpdir=/tmp/fluent-operator

pushd "$FLUENT_HELM_DIR" >/dev/null
helm package . -d ${_tmpdir} > /dev/null
mv ${_tmpdir}/*.tgz "$FLUENT_HELM_DIR/../fluent-operator.tgz"
rm -rf "${_tmpdir}"
popd >/dev/null
