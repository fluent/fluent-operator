#!/usr/bin/env bash

set -Eeo pipefail

PROJECT_ROOT="$PWD"
LOGGING_NAMESPACE=fluent
IMAGE_TAG=$(date "+%Y-%m-%d-%H-%M-%S")
KIND_CLUSTER="${KIND_CLUSTER:-fluent-operator-test-e2e}"
FO_IMG="kubesphere/fluent-operator:${IMAGE_TAG}"
PASS=0
FAIL=0

# ─── helpers ──────────────────────────────────────────────────────────────────

log()  { echo -e "\n==> $*"; }
pass() { echo "  PASS: $*"; PASS=$((PASS + 1)); }
fail() { echo "  FAIL: $*"; FAIL=$((FAIL + 1)); }

assert_crds_present() {
  local pattern="$1"
  local retries=20
  for ((i=0; i<retries; i++)); do
    if kubectl get crds --context "kind-${KIND_CLUSTER}" 2>&1 | grep -q "$pattern"; then
      pass "CRDs matching '$pattern' are present"
      return
    fi
    sleep 5
  done
  echo "  DEBUG: kubectl get crds output:"
  kubectl get crds --context "kind-${KIND_CLUSTER}" 2>&1 || true
  fail "CRDs matching '$pattern' not found after $((retries * 5))s"
}

assert_crds_absent() {
  local pattern="$1"
  local retries=20
  for ((i=0; i<retries; i++)); do
    if ! kubectl get crds --context "kind-${KIND_CLUSTER}" 2>&1 | grep -q "$pattern"; then
      pass "CRDs matching '$pattern' are absent (as expected)"
      return
    fi
    sleep 5
  done
  echo "  DEBUG: kubectl get crds output:"
  kubectl get crds --context "kind-${KIND_CLUSTER}" 2>&1 | grep "$pattern" || true
  fail "CRDs matching '$pattern' still present after $((retries * 5))s"
}

assert_annotation() {
  local crd="$1" key="$2" expected="$3"
  local actual
  actual=$(kubectl get crd "$crd" \
    --context "kind-${KIND_CLUSTER}" \
    -o jsonpath="{.metadata.annotations.$key}" 2>/dev/null || true)
  if [[ "$actual" == "$expected" ]]; then
    pass "CRD '$crd' has annotation '$key=$expected'"
  else
    fail "CRD '$crd' annotation '$key': expected '$expected', got '$actual'"
  fi
}

wait_operator() {
  kubectl -n "$LOGGING_NAMESPACE" wait --for=condition=available \
    deployment/fluent-operator --timeout=120s \
    --context "kind-${KIND_CLUSTER}"
  pass "fluent-operator deployment is available"
}

# Remove finalizers from all fluent.io CRs in the logging namespace.
# Must be called BEFORE helm uninstall of the operator, because the operator
# is responsible for removing its finalizers on deletion — if it's gone first,
# the CRs (and any CRDs they belong to) get stuck in Terminating forever.
drain_cr_finalizers() {
  local ctx="kind-${KIND_CLUSTER}"

  # Scale down the operator first so it cannot re-add finalizers after we strip them
  kubectl scale deployment/fluent-operator -n "$LOGGING_NAMESPACE" \
    --context "$ctx" --replicas=0 2>/dev/null || true
  kubectl -n "$LOGGING_NAMESPACE" wait pod \
    -l app.kubernetes.io/name=fluent-operator \
    --for=delete --context "$ctx" --timeout=30s 2>/dev/null || true

  for kind in fluentbits fluentdconfigs fluentbitconfigs collectors; do
    kubectl get "$kind" -n "$LOGGING_NAMESPACE" --context "$ctx" \
      -o name 2>/dev/null \
      | xargs -r kubectl patch --type=merge -n "$LOGGING_NAMESPACE" \
          --context "$ctx" -p '{"metadata":{"finalizers":[]}}' 2>/dev/null || true
  done
  # Give the API server a moment to process the patch before we proceed
  sleep 3
}

cleanup_helm() {
  local ctx="kind-${KIND_CLUSTER}"

  drain_cr_finalizers

  helm uninstall fluent-operator -n "$LOGGING_NAMESPACE" \
    --kube-context "$ctx" --wait --timeout 120s 2>/dev/null || true
  helm uninstall fluent-operator-crds -n "$LOGGING_NAMESPACE" \
    --kube-context "$ctx" --wait --timeout 120s 2>/dev/null || true
  kubectl delete ns "$LOGGING_NAMESPACE" --ignore-not-found \
    --context "$ctx" --timeout=60s
  # Remove any leftover CRDs so each scenario starts clean
  kubectl get crds --context "$ctx" -o name 2>/dev/null \
    | grep -E 'fluentbit\.fluent\.io|fluentd\.fluent\.io' \
    | xargs -r kubectl delete --context "$ctx" 2>/dev/null || true
}

# ─── setup ────────────────────────────────────────────────────────────────────

prepare_cluster() {
  log "Waiting for control-plane node…"
  kubectl wait --for=condition=Ready "node/${KIND_CLUSTER}-control-plane" \
    --context "kind-${KIND_CLUSTER}" --timeout=60s

  log "Building and loading operator image…"
  pushd "$PROJECT_ROOT" >/dev/null
  case "$(uname -m)" in
    arm64|aarch64) make build-op-arm64 -e "FO_IMG=${FO_IMG}" ;;
    *)             make build-op-amd64 -e "FO_IMG=${FO_IMG}" ;;
  esac
  kind load docker-image "${FO_IMG}" --name "$KIND_CLUSTER"
  popd >/dev/null
}

# ─── scenario 1: standard install (CRDs bundled in crds/) ─────────────────────

scenario_1_standard_install() {
  log "Scenario 1: Standard Install (CRDs bundled in crds/)"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  helm install fluent-operator \
    "${PROJECT_ROOT}/charts/fluent-operator" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set "operator.image.tag=${IMAGE_TAG}" \
    --set "operator.image.registry=kubesphere" \
    --set "operator.image.repository=fluent-operator" \
    --wait --timeout 60s

  wait_operator
  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  # The chart creates a FluentBit CR by default (fluentbit.enable=true).
  # Verify the operator reconciles it into a DaemonSet.
  echo "  Waiting for operator to reconcile FluentBit CR into a DaemonSet…"
  local retries=40
  local ds_found=false
  for ((i=0; i<retries; i++)); do
    if kubectl -n "$LOGGING_NAMESPACE" get daemonset fluent-bit \
         --context "kind-${KIND_CLUSTER}" &>/dev/null; then
      ds_found=true
      break
    fi
    sleep 5
  done
  if [[ "$ds_found" == "true" ]]; then
    pass "Operator reconciled FluentBit CR into a DaemonSet"
  else
    echo "  DEBUG: resources in namespace $LOGGING_NAMESPACE:"
    kubectl -n "$LOGGING_NAMESPACE" get all --context "kind-${KIND_CLUSTER}" 2>&1 || true
    echo "  DEBUG: operator logs:"
    kubectl -n "$LOGGING_NAMESPACE" logs deployment/fluent-operator \
      --context "kind-${KIND_CLUSTER}" --tail=50 2>&1 || true
    echo "  DEBUG: operator events:"
    kubectl -n "$LOGGING_NAMESPACE" get events --sort-by='.lastTimestamp' \
      --context "kind-${KIND_CLUSTER}" 2>&1 | tail -20 || true
    echo "  DEBUG: FluentBit CRs:"
    kubectl get fluentbits -A --context "kind-${KIND_CLUSTER}" 2>&1 || true
    fail "DaemonSet fluent-bit not created after $((retries * 5))s"
  fi

  # Drain finalizers before uninstalling so operator-owned CRs don't get stuck
  drain_cr_finalizers
  # Uninstall — CRDs in crds/ must survive
  helm uninstall fluent-operator -n "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}"
  assert_crds_present "fluentbit.fluent.io"
  pass "CRDs survived helm uninstall (crds/ directory behavior confirmed)"

  cleanup_helm
}

# ─── scenario 2: helm-managed CRDs (fluent-operator-crds chart) ───────────────

scenario_2_helm_managed_crds() {
  log "Scenario 2: Helm-Managed CRDs (fluent-operator-crds chart)"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  helm install fluent-operator-crds \
    "${PROJECT_ROOT}/charts/fluent-operator-crds" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --create-namespace \
    --wait --timeout 60s

  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  helm install fluent-operator \
    "${PROJECT_ROOT}/charts/fluent-operator" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --skip-crds \
    --set "operator.image.tag=${IMAGE_TAG}" \
    --set "operator.image.registry=kubesphere" \
    --set "operator.image.repository=fluent-operator" \
    --wait --timeout 60s

  wait_operator

  # Drain finalizers before uninstalling so operator-owned CRs don't get stuck.
  # Use --wait so all CR instances (ClusterFluentBitConfig, ClusterInput, etc.)
  # are fully deleted before we attempt to remove the CRDs.
  drain_cr_finalizers
  helm uninstall fluent-operator -n "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --wait --timeout 120s
  assert_crds_present "fluentbit.fluent.io"
  pass "CRDs survived operator uninstall (owned by fluent-operator-crds)"

  # Uninstall CRD chart — CRDs must be removed
  helm uninstall fluent-operator-crds -n "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}"
  assert_crds_absent "fluentbit.fluent.io"
  assert_crds_absent "fluentd.fluent.io"
  pass "CRDs removed after helm uninstall of fluent-operator-crds"

  cleanup_helm
}

# ─── scenario 3: selective CRD installation ───────────────────────────────────

scenario_3_selective_crds() {
  log "Scenario 3: Selective CRD Installation (Fluent Bit CRDs only)"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  helm install fluent-operator-crds \
    "${PROJECT_ROOT}/charts/fluent-operator-crds" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set "fluentd.enabled=false" \
    --wait --timeout 60s

  assert_crds_present "fluentbit.fluent.io"
  assert_crds_absent "fluentd.fluent.io"

  cleanup_helm
}

# ─── scenario 4: CRD protection annotation ────────────────────────────────────

scenario_4_crd_protection() {
  log "Scenario 4: CRD Protection (helm.sh/resource-policy=keep)"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  helm install fluent-operator-crds \
    "${PROJECT_ROOT}/charts/fluent-operator-crds" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set 'additionalAnnotations.helm\.sh/resource-policy=keep' \
    --wait --timeout 60s

  assert_annotation \
    "clusterfilters.fluentbit.fluent.io" \
    "helm\.sh/resource-policy" \
    "keep"

  # Uninstall — CRDs must survive because of the annotation
  helm uninstall fluent-operator-crds -n "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}"
  assert_crds_present "fluentbit.fluent.io"
  pass "CRDs survived uninstall due to helm.sh/resource-policy=keep"

  cleanup_helm
}

# ─── scenario 5: docker runtime backward compatibility ────────────────────────

scenario_5_docker_runtime_compat() {
  log "Scenario 5: Docker Runtime Backward Compatibility (containerRuntime=docker)"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  helm install fluent-operator \
    "${PROJECT_ROOT}/charts/fluent-operator" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set "operator.image.tag=${IMAGE_TAG}" \
    --set "operator.image.registry=kubesphere" \
    --set "operator.image.repository=fluent-operator" \
    --set "containerRuntime=docker" \
    --wait --timeout 60s

  wait_operator

  # Verify ClusterInput/tail uses parser: docker (not cri)
  local tail_parser
  tail_parser=$(kubectl get clusterinput tail \
    --context "kind-${KIND_CLUSTER}" \
    -o jsonpath='{.spec.tail.parser}' 2>/dev/null || true)
  if [[ "$tail_parser" == "docker" ]]; then
    pass "ClusterInput/tail uses parser=docker (correct for docker runtime)"
  else
    fail "ClusterInput/tail parser: expected 'docker', got '${tail_parser}'"
  fi

  # Verify the systemd ClusterInput is named 'docker' (runtime-derived name)
  if kubectl get clusterinput docker \
      --context "kind-${KIND_CLUSTER}" &>/dev/null; then
    pass "ClusterInput/docker (systemd) created with runtime-derived name"
  else
    echo "  DEBUG: existing ClusterInputs:"
    kubectl get clusterinput --context "kind-${KIND_CLUSTER}" 2>&1 || true
    fail "ClusterInput named 'docker' not found (systemd input should be named after runtime)"
  fi

  # Verify the systemd filter targets docker.service
  local systemd_filter
  systemd_filter=$(kubectl get clusterinput docker \
    --context "kind-${KIND_CLUSTER}" \
    -o jsonpath='{.spec.systemd.systemdFilter[0]}' 2>/dev/null || true)
  if [[ "$systemd_filter" == *"docker.service"* ]]; then
    pass "ClusterInput/docker systemdFilter targets docker.service"
  else
    fail "ClusterInput/docker systemdFilter: expected 'docker.service', got '${systemd_filter}'"
  fi

  cleanup_helm
}

# ─── summary ──────────────────────────────────────────────────────────────────

print_summary() {
  echo ""
  echo "======================================================="
  echo "  Helm v4 E2E Test Summary"
  echo "======================================================="
  echo "  PASSED: ${PASS}"
  echo "  FAILED: ${FAIL}"
  echo "======================================================="
  if [[ "$FAIL" -ne 0 ]]; then
    echo "  Some tests FAILED. Please check output above."
    exit 1
  else
    echo "  All tests PASSED."
    exit 0
  fi
}

# ─── main ─────────────────────────────────────────────────────────────────────

pushd "$PROJECT_ROOT" >/dev/null

prepare_cluster

scenario_1_standard_install
scenario_2_helm_managed_crds
scenario_3_selective_crds
scenario_4_crd_protection
scenario_5_docker_runtime_compat

popd >/dev/null

print_summary
