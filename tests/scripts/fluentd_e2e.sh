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
IMAGE_NAME="ghcr.io/fluent/fluent-operator/fluent-operator"
KIND_CLUSTER="${KIND_CLUSTER:-fluent-operator-test-e2e}"
BUILD_ARCH="${E2E_IMAGE_ARCH:-}"

if [ -z "$BUILD_ARCH" ]; then
  case "$(uname -m)" in
    arm64|aarch64) BUILD_ARCH="arm64" ;;
    x86_64|amd64) BUILD_ARCH="amd64" ;;
    *)
      echo "Unsupported architecture: $(uname -m). Set E2E_IMAGE_ARCH to amd64 or arm64."
      exit 1
      ;;
  esac
fi

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

  echo "Cleaning up…"
  pushd "$PROJECT_ROOT" >/dev/null || true
  make cleanup-test-e2e KIND_CLUSTER="$KIND_CLUSTER"
  popd >/dev/null || true
}

function prepare_cluster() {
  echo "wait the control-plane ready…"
  kubectl wait --for=condition=Ready "node/${KIND_CLUSTER}-control-plane" --timeout=60s
}

function build_image() {
  pushd "$PROJECT_ROOT" >/dev/null
  make "build-op-${BUILD_ARCH}" -e "FO_IMG=$IMAGE_NAME:$IMAGE_TAG"
  kind load docker-image "$IMAGE_NAME:$IMAGE_TAG" --name "$KIND_CLUSTER"

  # Build and load Fluentd image for e2e tests
  local fd_img="${FD_IMG:-ghcr.io/fluent/fluent-operator/fluentd:v1.19.2}"
  echo "Building Fluentd image for e2e tests…"
  make "build-fd-${BUILD_ARCH}" -e "FD_IMG=$fd_img"
  kind load docker-image "$fd_img" --name "$KIND_CLUSTER"

  popd >/dev/null
}

function start_fluent_operator() {
  pushd "$PROJECT_ROOT" >/dev/null
  sed -E \
    -e "s#^([[:space:]]*)image: ${IMAGE_NAME}:.*#\1image: ${IMAGE_NAME}:${IMAGE_TAG}\n\1imagePullPolicy: Never#" \
    < manifests/setup/setup.yaml | kubectl create -f -
  kubectl -n "$LOGGING_NAMESPACE" wait --for=condition=available deployment/fluent-operator --timeout=60s
  popd >/dev/null
}

function run_test() {
  export ACK_GINKGO_RC=true
  if ! "$GINKGO_BIN" -v "$E2E_DIR/e2e/fluentd/fluentd.test" -- "$debugflag"; then
    echo "Integration suite has failures, Please check !!"
    exit 1
  else
    echo "Integration suite successfully passed all the tests !!"
    exit 0
  fi
}

function main() {
  trap cleanup EXIT

  echo -e "\nBuilding testcases…"
  build_ginkgo_test

  echo -e "\nPreparing cluster…"
  prepare_cluster

  echo -e "\nBuilding image…"
  build_image

  echo -e "\nStart fluent operator…"
  start_fluent_operator

  echo -e "\nRunning test…"
  run_test
}

main "$@"
