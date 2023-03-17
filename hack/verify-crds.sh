#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/..
CRD_OPTIONS="crd:generateEmbeddedObjectMeta=true"

DIFFROOT="${SCRIPT_ROOT}/config/crd/bases/"
TMP_DIFFROOT="${SCRIPT_ROOT}/_tmp/config/crd/bases/"
_tmp="${SCRIPT_ROOT}/_tmp"

cleanup() {
  rm -rf "${_tmp}"
}
trap "cleanup" EXIT SIGINT

cleanup

mkdir -p "${TMP_DIFFROOT}"
cp -a "${DIFFROOT}"/* "${TMP_DIFFROOT}"

./bin/controller-gen $CRD_OPTIONS rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=config/crd/bases/
./bin/controller-gen $CRD_OPTIONS rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=config/crd/bases/
echo "diffing ${DIFFROOT} against freshly generated crds"
ret=0
diff -Naupr "${DIFFROOT}" "${TMP_DIFFROOT}" || ret=$?
cp -a "${TMP_DIFFROOT}"/* "${DIFFROOT}"
if [[ $ret -eq 0 ]]
then
  echo "${DIFFROOT} up to date."
else
  echo "${DIFFROOT} is out of date. Please rerun make manifests"
  exit 1
fi


