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
