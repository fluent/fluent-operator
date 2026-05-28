#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

VERIFY_PATHS=(
  "config/crd/bases"
  "charts/fluent-operator/crds"
  "charts/fluent-operator-fluent-bit-crds/templates"
  "charts/fluent-operator-fluentd-crds/templates"
  "manifests/setup/setup.yaml"
)

cd "${SCRIPT_ROOT}"

make manifests

echo "diffing checked-in manifests against freshly generated manifests"
if git diff --exit-code -- "${VERIFY_PATHS[@]}"
then
  echo "CRDs are up to date."
else
  echo "CRDs are out of date. Please rerun make manifests"
  exit 1
fi
