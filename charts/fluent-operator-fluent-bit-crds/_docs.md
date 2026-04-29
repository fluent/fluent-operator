## Advanced Installation

### Install Fluent Bit CRDs

```bash
helm install fluent-operator-fluent-bit-crds fluent/fluent-operator-fluent-bit-crds
```

## Protecting CRDs from Deletion

To prevent CRDs from being deleted on `helm uninstall`:

```bash
helm install fluent-operator-fluent-bit-crds fluent/fluent-operator-fluent-bit-crds \
  --set additionalAnnotations."helm\.sh/resource-policy"=keep
```

With this annotation, Helm will preserve the CRDs even if the chart is uninstalled.

## Using with fluent-operator

After installing the CRDs with this chart, install the operator with `--skip-crds`:

```bash
# Step 1: Install Fluent Bit CRDs
helm install fluent-operator-fluent-bit-crds fluent/fluent-operator-fluent-bit-crds

# Step 2: Install operator (skip CRDs since already installed)
helm install fluent-operator fluent/fluent-operator --skip-crds
```
