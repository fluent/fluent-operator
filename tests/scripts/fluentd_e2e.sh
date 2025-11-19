#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
set -o errtrace

debugflag=${1:-}

PROJECT_ROOT="$PWD"
# Portable way to get absolute path
E2E_DIR="$(cd "$(dirname "$0")/.." && pwd)"
LOGGING_NAMESPACE="fluent"
IMAGE_TAG="$(date "+%Y-%m-%d-%H-%M-%S")"
VERSION="$(tr -d " \t\n\r" < VERSION)"
LOG_FILE="$(mktemp)"
IMAGE_NAME="ghcr.io/fluent/fluent-operator/fluent-operator"

GINKGO_BIN="ginkgo"
if [ -f "$PROJECT_ROOT/bin/ginkgo" ]; then
  GINKGO_BIN="$PROJECT_ROOT/bin/ginkgo"
fi

function build_ginkgo_test() {
  pushd "$E2E_DIR" >/dev/null
  "$GINKGO_BIN" build -r e2e/fluentd/
  popd >/dev/null
}

# shellcheck disable=SC2329
function cleanup() {
  local exit_code=$?
  
  if [ "${SKIP_CLEANUP:-false}" == "true" ]; then
    echo "Skipping cleanup as requested."
    exit "$exit_code"
  fi

  echo "Cleaning up..."
  pushd "$PROJECT_ROOT" >/dev/null || true
  # kubectl delete -f manifests/setup/setup.yaml
  # kubectl delete ns $LOGGING_NAMESPACE
  kind delete cluster --name test || true
  popd >/dev/null || true
  rm -f "$LOG_FILE"
}

function prepare_cluster() {
  kind create cluster --name test
  kubectl create ns "$LOGGING_NAMESPACE"

  echo "wait the control-plane ready..."
  kubectl wait --for=condition=Ready node/test-control-plane --timeout=60s
}

function build_image() {
  pushd "$PROJECT_ROOT" >/dev/null
  make build-op-amd64 -e "FO_IMG=$IMAGE_NAME:$IMAGE_TAG"
  kind load docker-image "$IMAGE_NAME:$IMAGE_TAG" --name test
  popd >/dev/null
}

function start_fluent_operator() {
  pushd "$PROJECT_ROOT" >/dev/null
  sed "s#$IMAGE_NAME:${VERSION}#$IMAGE_NAME:$IMAGE_TAG#g" < manifests/setup/setup.yaml | kubectl create -f -
  kubectl -n "$LOGGING_NAMESPACE" wait --for=condition=available deployment/fluent-operator --timeout=60s
  popd >/dev/null
}

function run_test() {
  # inspired by github.com/kubeedge/kubeedge/tests/e2e/scripts/helm_keadm_e2e.sh
  echo "Logs will be written to $LOG_FILE"
  
  export ACK_GINKGO_RC=true
  "$GINKGO_BIN" -v "$E2E_DIR/e2e/fluentd/fluentd.test" -- "$debugflag"

  if [[ $? != 0 ]]; then
    echo "Integration suite has failures, Please check !!"
    exit 1
  else
    echo "Integration suite successfully passed all the tests !!"
    exit 0
  fi
}

function main() {
  trap cleanup EXIT

  echo -e "\nBuilding testcases..."
  build_ginkgo_test

  echo -e "\nPreparing cluster..."
  prepare_cluster

  echo -e "\nBuilding image..."
  build_image

  echo -e "\nStart fluent operator..."
  start_fluent_operator

  echo -e "\nRunning test..."
  run_test
}

main "$@"
