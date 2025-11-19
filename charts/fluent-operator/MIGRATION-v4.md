# Migration Guide: Fluent Operator Helm Chart v3.x to v4.0

## Overview
Major changes/themes for v4.0:

1. **Container Runtime Simplification**: Removes dynamic detection for the `docker` runtime via initContainers and adopts static, configuration-based paths. The `docker` runtime has not been used widely since Kubernetes v1.24 (2022) and modern Kubernetes distributions now use the `containerd` runtime.

2. **Fluentd CRDs Separation**: Fluentd CRDs have been moved from a sub-chart to a separate, independently versioned `fluentd-crds` chart. This allows for independent lifecycle management of the CRDs and operator and simplifies the process of updating and managing CRDs.

## Breaking Changes

### 1. Default Container Runtime Changed

**What Changed:**

- Default `containerRuntime` changed from `docker` to `containerd`

**Impact:**

- Users who never explicitly set `containerRuntime` will now use containerd defaults
- Log parser will change from `docker` to `cri` format
- Systemd filter will target `containerd.service` instead of `docker.service`

**Who Is Affected:**

- Users still running Docker container runtime (uncommon - Docker support was removed from Kubernetes in v1.24, May 2022)
- Users who relied on the default value without explicit configuration

**Migration:**


The containerRuntime now defaults to `containerd`. In `v3.x` the implicit default was `docker`. Use `containerRuntime: docker` to maintain `v3.x` behavior.

```diff
+ containerRuntime: containerd
```

### 2. initContainers Removed

**What Changed:**

- Removed dynamic Docker root directory detection via initContainer
- Removed `operator.initcontainer.*` configuration options
- Container log paths are now statically configured

**Impact:**

- initContainer no longer runs before the operator starts
- Removes dependency on third party outdated images for initContainers
- Docker socket no longer mounted for path detection

**Who Is Affected:**

- Users who customized `operator.initcontainer.image` or `operator.initcontainer.resources`
- Users with Docker installations using non-standard root directories

**Migration:**

The operator initContainer has been removed.

```diff
  operator:
-   initcontainer:
-     image:
-       registry: docker.io
-       repository: docker
-       tag: "20.10"
    resources:
      limits:
        cpu: 100m
        memory: 64Mi
```

### 3. Log Path Configuration Removed

**What Changed:**

- Removed `operator.logPath.containerd` and `operator.logPath.crio`
- Removed ability to configure custom log paths
- Log paths are now determined automatically based on `containerRuntime`

**Impact:**

- Old `operator.logPath.*` configuration is ignored
- Each container runtime uses its standard default path

**Who Is Affected:**

- Users who set custom paths via `operator.logPath.containerd` or `operator.logPath.crio`
- Users with non-standard container log directory locations

**Migration:**

If you were using custom log paths, you must ensure your container runtime uses the standard default paths shown below, or adjust your container runtime configuration to use these standard paths.

## Default Paths by Runtime

v4.0 uses the following default paths based on the configured `containerRuntime`:

| Container Runtime | Default Path |
|-------------------|--------------|
| `containerd` | `/var/log/containers` |
| `crio` | `/var/log/containers` |
| `docker` | `/var/lib/docker/containers` |

## Migration Scenarios

### Scenario 1: Using Containerd (Default) - No Changes Needed ✅

```yaml
# v3.x
containerRuntime: containerd  # or not set
# ... rest of config

# v4.0 - No changes required!
# The new defaults work out of the box
```

### Scenario 2: Using CRI-O - Minimal Changes

```yaml
# v3.x
containerRuntime: crio
# ... rest of config

# v4.0 - Explicitly set runtime (same as before)
containerRuntime: crio
# Default path /var/log/containers works for most CRI-O installations
```

### Scenario 3: Using Docker

```yaml
# v3.x
containerRuntime: docker
# (relied on automatic detection)

# v4.0 - Must explicitly set runtime
containerRuntime: docker
# Uses default path: /var/lib/docker/containers
# This works for standard Docker installations
# If your Docker uses a custom root directory, you must reconfigure Docker
# to use the standard path
```

### 4. Fluentd CRDs Moved to Separate Chart

**What Changed:**

- Fluentd CRDs are no longer included as a sub-chart of `fluent-operator`
- Fluentd CRDs are now available as a separate top-level chart: `fluentd-crds`
- The `fluentd-crds` chart is versioned independently from the `fluent-operator` chart

**Impact:**

- Fluentd CRDs will not be installed automatically when installing/upgrading the `fluent-operator` chart
- Users must install the `fluentd-crds` chart separately if they use Fluentd
- Users can upgrade the `fluentd-crds` chart independently of the operator

**Who Is Affected:**

- Users who use Fluentd with the fluent-operator
- Users who relied on Fluentd CRDs being automatically installed with the operator chart

**Migration:**

If you use Fluentd, you must install the `fluentd-crds` chart separately before or alongside the `fluent-operator` chart.

```bash
# Add the Fluent Helm repository
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update

# Install fluentd-crds chart first
helm install fluentd-crds fluent/fluentd-crds --version 0.1.0

# Then install or upgrade fluent-operator
helm upgrade --install fluent-operator fluent/fluent-operator --version 4.0.0
```

## Forward Looking: Planned Changes in v5.0

**Future Change (v5.0):**

- The `fluent-operator-env` _ConfigMap_, which is used to provide backwards compatibility with fluent-operator =<3.5, will be completely removed
