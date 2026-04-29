## Advanced Installation

### Method 1: Standard Installation

The simplest way to install Fluent Operator with CRDs included:

```bash
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update

helm install fluent-operator fluent/fluent-operator
```

**Behavior:**
- CRDs automatically installed from the `crds/` directory
- Helm does NOT upgrade CRDs on `helm upgrade` (manual upgrade required)
- Helm does NOT delete CRDs on `helm uninstall`

**Upgrading CRDs:**

When upgrading, manually apply CRD updates before upgrading the chart. You can obtain the CRDs from the chart package or the repository:

```bash
# Option 1: Extract from the chart
helm pull fluent/fluent-operator --untar
kubectl apply --server-side --force-conflicts -f fluent-operator/crds/

# Option 2: Clone the repository
git clone https://github.com/fluent/fluent-operator.git
cd fluent-operator
kubectl apply --server-side --force-conflicts -f charts/fluent-operator/crds/

# Then upgrade the chart
helm upgrade fluent-operator fluent/fluent-operator
```

**Skipping CRDs:**

If you manage CRDs separately:

```bash
helm install fluent-operator fluent/fluent-operator --skip-crds
```

### Method 2: Helm-Managed CRDs (Advanced)

For full Helm lifecycle management of CRDs (automatic upgrades and deletions), install the
`fluent-operator-fluent-bit-crds` and/or `fluent-operator-fluentd-crds` charts separately:

```bash
# Step 1: Install CRDs with Helm management
helm install fluent-operator-fluent-bit-crds fluent/fluent-operator-fluent-bit-crds
helm install fluent-operator-fluentd-crds fluent/fluent-operator-fluentd-crds

# Step 2: Install operator (skip CRDs since already installed)
helm install fluent-operator fluent/fluent-operator --skip-crds
```

You can install only the CRD chart(s) you need — omitting a chart means those CRDs are simply not installed.

**Behavior:**
- CRDs automatically upgrade with `helm upgrade fluent-operator-fluent-bit-crds` / `helm upgrade fluent-operator-fluentd-crds`
- Fine-grained control — install only Fluent Bit CRDs, only Fluentd CRDs, or both
- CRDs deleted on `helm uninstall` (unless protected with annotation)

**Protecting CRDs:**

```bash
helm install fluent-operator-fluent-bit-crds fluent/fluent-operator-fluent-bit-crds \
  --set additionalAnnotations."helm\.sh/resource-policy"=keep

helm install fluent-operator-fluentd-crds fluent/fluent-operator-fluentd-crds \
  --set additionalAnnotations."helm\.sh/resource-policy"=keep
```

See the [fluent-operator-fluent-bit-crds](https://github.com/fluent/fluent-operator/tree/master/charts/fluent-operator-fluent-bit-crds) and [fluent-operator-fluentd-crds](https://github.com/fluent/fluent-operator/tree/master/charts/fluent-operator-fluentd-crds) charts for more details.

## Upgrading

See [MIGRATION-v4.md](https://github.com/fluent/fluent-operator/blob/master/charts/fluent-operator/MIGRATION-v4.md) for detailed upgrade instructions from v3.x to v4.0.

### Upgrading from v3.x to v4.0

**Major Changes:**
- CRDs now in `crds/` directory (Helm v3 standard)
- New `fluent-operator-fluent-bit-crds` and `fluent-operator-fluentd-crds` charts available for Helm-managed CRDs
- Default container runtime changed to `containerd`
- Removed dependency on legacy CRD sub-charts

**Upgrade Steps:**

```bash
# Update repository
helm repo update

# Manually update CRDs first (Helm doesn't upgrade CRDs in crds/ directory)
helm pull fluent/fluent-operator --version 4.0.0 --untar
kubectl apply -f fluent-operator/crds/

# Then upgrade the chart
helm upgrade fluent-operator fluent/fluent-operator --version 4.0.0
```
