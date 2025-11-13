# Migration Guide: Fluent Operator Helm Chart v3.x to v4.0

## Overview

v4.0 simplifies container runtime configuration by removing dynamic detection for the `docker` runtime via initContainers and adopting static, configuration-based paths. The `docker` runtime has not been used widely since Kubernetes v1.24 (2022) and modern Kubernetes distributions now use the `containerd` runtime.

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

### 3. Log Path Configuration Simplified

**What Changed:**

- Removed `operator.logPath.containerd` and `operator.logPath.crio`
- Added new `operator.containerLogPath` for direct path specification

**Impact:**

- Old logPath configuration is ignored
- New configuration accepts full path to container logs (not just root directory)

**Who Is Affected:**

- Users who set custom paths via `operator.logPath.containerd` or `operator.logPath.crio`

**Migration:**


In v3 the root directory was specified for the log path. In v4 this will be the full path to the container logs folder.

```diff
  operator:
-   logPath:
-     containerd: /var/log
+   containerLogPath: "/var/log/containers"
```

## Default Paths by Runtime

v4.0 uses the following default paths when `operator.containerLogPath` is not explicitly set:

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

### Scenario 3: Using Docker with Standard Paths

```yaml
# v3.x
containerRuntime: docker
# (relied on automatic detection)

# v4.0 - Must explicitly set runtime
containerRuntime: docker
# Default /var/lib/docker/containers works for standard Docker installations
# No need to set operator.containerLogPath
```

### Scenario 4: Using Docker with Custom Paths

```yaml
# v3.x
containerRuntime: docker
# (used initContainer to detect custom docker root)

# v4.0 - Must explicitly configure path
containerRuntime: docker
operator:
  containerLogPath: "/custom/docker/root/containers"
```

### Scenario 5: Custom Log Paths (Any Runtime)

```yaml
# v3.x
containerRuntime: containerd
operator:
  logPath:
    containerd: /custom/log/path

# v4.0 - Use new configuration key
containerRuntime: containerd
operator:
  containerLogPath: "/custom/log/path/containers"
  # Note: append "/containers" to your old path
```

## Forward Looking: Planned Changes in v5.0

**Future Change (v5.0):**

- The `fluent-operator-env` _ConfigMap_, which is used to provide backwards compatibility with fluent-operator =<3.5, will be completely removed
