#!/usr/bin/env bash

# inspired by: https://github.com/weaveworks/flagger/tree/main/hack

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(git rev-parse --show-toplevel)

DIFFROOT="${SCRIPT_ROOT}/apis"
TMP_DIFFROOT="${SCRIPT_ROOT}/_tmp/apis"

_tmp="${SCRIPT_ROOT}/_tmp"

cleanup() {
    echo ">> Removing ${_tmp}"
    rm -rf "${_tmp}"
}
trap "cleanup" EXIT SIGINT

cleanup

mkdir -p "${TMP_DIFFROOT}"
cp -a "${DIFFROOT}"/* "${TMP_DIFFROOT}"

"${SCRIPT_ROOT}/hack/update-codegen.sh"
echo "diffing ${DIFFROOT} against freshly generated clientset"
ret=0
diff -Naupr "${DIFFROOT}" "${TMP_DIFFROOT}" || ret=$?
cp -a "${TMP_DIFFROOT}"/* "${DIFFROOT}"
if [[ $ret -eq 0 ]]
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please rerun make generate"
  exit 1
fi
