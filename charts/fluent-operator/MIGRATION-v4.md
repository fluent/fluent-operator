# Migration Guide: Fluent Operator Helm Chart v3.x to v4.0

## Overview

Major changes/themes for v4.0:

1. **Container Runtime Simplification**: Removes dynamic detection for the `docker` runtime via initContainers and adopts static, configuration-based paths. The `docker` runtime has not been used widely since Kubernetes v1.24 (2022) and modern Kubernetes distributions now use the `containerd` runtime.

2. **CRD lifecycle changes**: Both Fluent Bit and Fluentd CRDs have been moved from embedded sub-charts to separate, independently versioned `fluent-bit-crds` and `fluentd-crds` charts hosted in the Fluent Helm repository. CRDs are now in the `templates/` directory, enabling automatic upgrades via `helm upgrade`. **Important:** CRDs will now be deleted on `helm uninstall` unless protected with the `helm.sh/resource-policy: keep` annotation (see below for details).

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

### Scenario 1: Using Containerd (Default) - No Changes Needed

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

## CRDs Moved to Separate Charts with Automatic Upgrades

**What Changed:**

- Fluent Bit CRDs moved from embedded sub-chart to separate `fluent-bit-crds` chart in the Fluent Helm repository
- Fluentd CRDs moved from embedded sub-chart to separate `fluentd-crds` chart in the Fluent Helm repository
- CRDs relocated from `crds/` directory to `templates/` directory, changing their lifecycle behavior

**Impact:**

- **CRDs now upgrade automatically** with `helm upgrade` (previously required manual upgrade)
- **CRDs will be deleted on `helm uninstall`** (previously they were preserved).

**Who Is Affected:**

- All users who use Fluent Bit with the fluent-operator
- All users who use Fluentd with the fluent-operator
- Users who rely on CRDs and custom resources persisting after chart uninstall

**Migration:**

Before installing or upgrading, ensure the Fluent Helm repository is added:

```bash
# Add the Fluent Helm repository (required for the dependencies)
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update

# Install or upgrade fluent-operator (CRDs installed automatically via dependencies)
helm upgrade --install fluent-operator fluent/fluent-operator \
  --version 4.0.0 \
  --set fluentbit.enable=true \
  --set fluentd.enable=true
```

**Note:** Helm will automatically manage the CRD chart dependencies:
- `fluent-bit-crds` when `fluentbit.enable=true` and `fluentbit.crdsEnable=true` (default)
- `fluentd-crds` when `fluentd.enable=true` and `fluentd.crdsEnable=true` (default)

**Protecting CRDs from Deletion:**

To prevent CRDs from being deleted when the CRD charts are uninstalled, add the `helm.sh/resource-policy: keep` annotation:

```bash
helm upgrade --install fluent-operator fluent/fluent-operator \
  --version 4.0.0 \
  --set fluentbit.enable=true \
  --set fluentd.enable=true \
  --set fluent-bit-crds.additionalAnnotations."helm\.sh/resource-policy"=keep \
  --set fluentd-crds.additionalAnnotations."helm\.sh/resource-policy"=keep
```

With this annotation, Helm will preserve the CRDs and all custom resources even if the charts are uninstalled.

## Forward Looking: Planned Changes in v5.0

**Future Change (v5.0):**

- The `fluent-operator-env` _ConfigMap_, which is used to provide backwards compatibility with fluent-operator =<3.5, will be completely removed
