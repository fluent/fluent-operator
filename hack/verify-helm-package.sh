#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

FLUENT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
FLUENT_HELM_DIR=${FLUENT_ROOT}/charts/fluent-operator
_tmpdir=/tmp/fluent-operator

function verify:package:helm:files {
    mkdir -p ${_tmpdir}

    pushd "$FLUENT_HELM_DIR" >/dev/null
    helm package . -d ${_tmpdir} > /dev/null
    mv ${_tmpdir}/*.tgz "${_tmpdir}/fluent-operator.tgz"
    helm_checksum=$(tar -xOzf "${FLUENT_HELM_DIR}/../fluent-operator.tgz" | sort | sha1sum | awk '{ print $1 }')
    temp_helm_checksum=$(tar -xOzf "${_tmpdir}/fluent-operator.tgz" | sort | sha1sum | awk '{ print $1 }')
    popd >/dev/null

    if [ "$helm_checksum" != "$temp_helm_checksum" ]; then
      echo "Helm package fluent-operator.tgz not updated or the helm chart is not expected."
      echo "Please run:  make update-helm-package"
      exit 1
    fi
}

function cleanup {
  #echo "Removing workspace: ${_tmpdir}"
  rm -rf "${_tmpdir}"
}

trap cleanup EXIT SIGINT

verify:package:helm:files
