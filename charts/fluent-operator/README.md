# Fluent Operator Helm chart

[Fluent Operator](https://github.com/fluent/fluent-operator/) provides a Kubernetes-native logging pipeline based on Fluent-Bit and Fluentd.

## Installation

To install or upgrade Fluent Operator using Helm:

```shell
export FLUENT_OPERATOR_CONTAINER_RUNTIME="containerd" # or "cri-o", "docker" depending on the container runtime being used (see `values.yaml`)

helm repo add fluent https://fluent.github.io/helm-charts
helm upgrade --install fluent-operator fluent/fluent-operator \
  --create-namespace \
  --set containerRuntime=${FLUENT_OPERATOR_CONTAINER_RUNTIME}
```

By default, all CRDs required for Fluent Operator will be installed.  To prevent `helm install` from installing CRDs, you can set `fluent-bit.crdsEnable` or `fluentd.crdsEnable` to `false`.

## Upgrading

### Upgrading to v4.0

⚠️ **v4.0 contains breaking changes.** Please review the [Migration Guide](MIGRATION-v4.md) before upgrading.

**Key Changes:**

- Default `containerRuntime` changed from `docker` to `containerd`
- Removed initContainers for dynamic path detection for the `docker` runtime
- Simplified configuration with `operator.containerLogPath`

**Quick Migration:**

```yaml
# If using Docker, explicitly set in your values:
containerRuntime: docker

# If using a non-default for `docker` logs, use new configuration:
operator:
  containerLogPath: "/var/log/containers"
```

See [MIGRATION-v4.md](MIGRATION-v4.md) for complete migration instructions.

### Upgrading

Helm [does not manage the lifecycle of CRDs](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/), so if the Fluent Operator CRDs already exist, subsequent
chart upgrades will not add or remove CRDs even if they have changed.  During upgrades, users should manually update CRDs:

```shell
wget https://github.com/fluent/fluent-operator/releases/download/<version>/fluent-operator.tgz
tar -xf fluent-operator.tgz
kubectl replace -f fluent-operator/crds
```

## Chart Values

```shell
helm show values fluent/fluent-operator
```
