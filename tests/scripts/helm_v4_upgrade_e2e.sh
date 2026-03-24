#!/usr/bin/env bash
# Upgrade path e2e tests: fluent-operator v3.x → v4.0
#
# Scenario A: Standard helm upgrade (same chart, newer version)
# Scenario B: Migrate from bundled-CRD install to fluent-operator-crds chart

set -Eeo pipefail

PROJECT_ROOT="$PWD"
LOGGING_NAMESPACE=fluent
IMAGE_TAG=$(date "+%Y-%m-%d-%H-%M-%S")
KIND_CLUSTER="${KIND_CLUSTER:-fluent-operator-test-e2e}"
FO_IMG="kubesphere/fluent-operator:${IMAGE_TAG}"
V3_VERSION="3.5.0"
V3_OPERATOR_IMG="ghcr.io/fluent/fluent-operator/fluent-operator:${V3_VERSION}"
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


wait_operator() {
  kubectl -n "$LOGGING_NAMESPACE" wait --for=condition=available \
    deployment/fluent-operator --timeout=120s \
    --context "kind-${KIND_CLUSTER}"
  pass "fluent-operator deployment is available"
}

# Scale down the operator to prevent it from re-adding finalizers, then strip them.
drain_cr_finalizers() {
  local ctx="kind-${KIND_CLUSTER}"

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
  sleep 2
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
  kubectl get crds --context "$ctx" -o name 2>/dev/null \
    | grep -E 'fluentbit\.fluent\.io|fluentd\.fluent\.io' \
    | xargs -r kubectl delete --context "$ctx" 2>/dev/null || true
}

# ─── setup ────────────────────────────────────────────────────────────────────

prepare_cluster() {
  log "Waiting for control-plane node…"
  kubectl wait --for=condition=Ready "node/${KIND_CLUSTER}-control-plane" \
    --context "kind-${KIND_CLUSTER}" --timeout=60s

  log "Pre-pulling v3 operator image and loading into Kind…"
  docker pull "${V3_OPERATOR_IMG}"
  kind load docker-image "${V3_OPERATOR_IMG}" --name "$KIND_CLUSTER"

  log "Building and loading v4 operator image…"
  pushd "$PROJECT_ROOT" >/dev/null
  case "$(uname -m)" in
    arm64|aarch64) make build-op-arm64 -e "FO_IMG=${FO_IMG}" ;;
    *)             make build-op-amd64 -e "FO_IMG=${FO_IMG}" ;;
  esac
  kind load docker-image "${FO_IMG}" --name "$KIND_CLUSTER"
  popd >/dev/null
}

# ─── scenario A: standard helm upgrade (v3 → v4) ─────────────────────────────

scenario_a_standard_upgrade() {
  log "Scenario A: Standard helm upgrade (v3.x → v4.0)"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  log "Installing fluent-operator v${V3_VERSION} from Helm registry…"
  helm install fluent-operator fluent/fluent-operator \
    --version "${V3_VERSION}" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --wait --timeout 120s

  wait_operator
  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  # The v3 chart installs CRDs via sub-chart crds/ dirs — they are NOT tracked
  # in the Helm release. Confirm there is exactly one release: fluent-operator.
  local release_count
  release_count=$(helm list -n "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" --short 2>/dev/null | wc -l | tr -d ' ')
  if [[ "$release_count" -eq 1 ]]; then
    pass "Exactly one Helm release present after v3 install (sub-charts are bundled, not separate)"
  else
    echo "  DEBUG: helm list output:"
    helm list -n "$LOGGING_NAMESPACE" --kube-context "kind-${KIND_CLUSTER}" 2>/dev/null || true
    fail "Expected 1 Helm release after v3 install, got ${release_count}"
  fi

  log "Manually applying v4 CRDs (required step — Helm never auto-upgrades crds/)…"
  # --server-side avoids the 262144-byte annotation limit that client-side apply hits
  # on large CRDs like fluentbits and fluentds.
  # --force-conflicts is required because Helm owns the existing CRD fields
  # (installed them via crds/ dir); we need kubectl to take over management.
  kubectl apply --server-side --force-conflicts \
    -f "${PROJECT_ROOT}/charts/fluent-operator/crds/" \
    --context "kind-${KIND_CLUSTER}"
  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  log "Upgrading fluent-operator to v4…"
  helm upgrade fluent-operator \
    "${PROJECT_ROOT}/charts/fluent-operator" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set "operator.image.tag=${IMAGE_TAG}" \
    --set "operator.image.registry=kubesphere" \
    --set "operator.image.repository=fluent-operator" \
    --wait --timeout 120s

  wait_operator
  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  # The chart creates a FluentBit CR. Verify the operator reconciles it.
  # (v3 chart also created a FluentBit CR, so the DaemonSet may already exist.)
  echo "  Waiting for operator to reconcile FluentBit CR into a DaemonSet…"
  local retries=24
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
    pass "FluentBit DaemonSet exists after upgrade (workload continuity confirmed)"
  else
    echo "  DEBUG: resources in namespace $LOGGING_NAMESPACE:"
    kubectl -n "$LOGGING_NAMESPACE" get all --context "kind-${KIND_CLUSTER}" 2>&1 || true
    echo "  DEBUG: operator logs:"
    kubectl -n "$LOGGING_NAMESPACE" logs deployment/fluent-operator \
      --context "kind-${KIND_CLUSTER}" --tail=30 2>&1 || true
    fail "FluentBit DaemonSet not found after upgrade (workload continuity broken)"
  fi

  cleanup_helm
}

# ─── scenario B: migrate to fluent-operator-crds chart ────────────────────────

scenario_b_migrate_to_crds_chart() {
  log "Scenario B: Migrate v3 install to fluent-operator-crds chart"

  kubectl create ns "$LOGGING_NAMESPACE" --context "kind-${KIND_CLUSTER}"

  log "Installing fluent-operator v${V3_VERSION} from Helm registry…"
  helm install fluent-operator fluent/fluent-operator \
    --version "${V3_VERSION}" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --wait --timeout 120s

  wait_operator
  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  log "Manually applying v4 CRDs before upgrade…"
  kubectl apply --server-side --force-conflicts \
    -f "${PROJECT_ROOT}/charts/fluent-operator/crds/" \
    --context "kind-${KIND_CLUSTER}"

  log "Upgrading fluent-operator to v4…"
  helm upgrade fluent-operator \
    "${PROJECT_ROOT}/charts/fluent-operator" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set "operator.image.tag=${IMAGE_TAG}" \
    --set "operator.image.registry=kubesphere" \
    --set "operator.image.repository=fluent-operator" \
    --wait --timeout 120s

  wait_operator

  # At this point: CRDs exist in the cluster with NO Helm release annotations
  # (they came from crds/ directories of v3 sub-charts, which Helm never tracks).
  # Helm cannot adopt pre-existing unowned CRDs via 'helm install' — it requires
  # the ownership label and annotations to be present first.
  log "Adding Helm ownership metadata to existing CRDs so fluent-operator-crds can adopt them…"
  local ctx="kind-${KIND_CLUSTER}"
  while IFS= read -r crd; do
    kubectl label "$crd" --context "$ctx" \
      app.kubernetes.io/managed-by=Helm --overwrite
    kubectl annotate "$crd" --context "$ctx" \
      "meta.helm.sh/release-name=fluent-operator-crds" \
      "meta.helm.sh/release-namespace=${LOGGING_NAMESPACE}" --overwrite
  done < <(kubectl get crds --context "$ctx" -o name 2>/dev/null \
             | grep -E 'fluentbit\.fluent\.io|fluentd\.fluent\.io')
  pass "Helm ownership metadata applied to all fluent.io CRDs"

  log "Installing fluent-operator-crds (CRDs now have ownership metadata)…"
  helm install fluent-operator-crds \
    "${PROJECT_ROOT}/charts/fluent-operator-crds" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --set 'additionalAnnotations.helm\.sh/resource-policy=keep' \
    --wait --timeout 60s

  pass "fluent-operator-crds installed successfully (existing CRDs adopted)"

  # Verify the CRD chart now tracks the CRDs
  local manifest_lines
  manifest_lines=$(helm get manifest fluent-operator-crds \
    -n "$LOGGING_NAMESPACE" --kube-context "kind-${KIND_CLUSTER}" 2>/dev/null \
    | grep -c "CustomResourceDefinition" || true)
  if [[ "$manifest_lines" -gt 0 ]]; then
    pass "fluent-operator-crds Helm release tracks ${manifest_lines} CRD definitions"
  else
    fail "fluent-operator-crds Helm release manifest contains no CRD definitions"
  fi

  # Upgrade the operator to use --skip-crds (CRDs now owned by fluent-operator-crds)
  log "Upgrading fluent-operator with --skip-crds (CRDs now owned by fluent-operator-crds)…"
  drain_cr_finalizers
  helm upgrade fluent-operator \
    "${PROJECT_ROOT}/charts/fluent-operator" \
    --namespace "$LOGGING_NAMESPACE" \
    --kube-context "kind-${KIND_CLUSTER}" \
    --skip-crds \
    --set "operator.image.tag=${IMAGE_TAG}" \
    --set "operator.image.registry=kubesphere" \
    --set "operator.image.repository=fluent-operator" \
    --wait --timeout 120s

  wait_operator
  assert_crds_present "fluentbit.fluent.io"
  assert_crds_present "fluentd.fluent.io"

  cleanup_helm
}

# ─── summary ──────────────────────────────────────────────────────────────────

print_summary() {
  echo ""
  echo "======================================================="
  echo "  Helm v4 Upgrade E2E Test Summary"
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

scenario_a_standard_upgrade
scenario_b_migrate_to_crds_chart

popd >/dev/null

print_summary
