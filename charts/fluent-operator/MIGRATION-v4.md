# Migration Guide: Fluent Operator Helm Chart v3.x to v4.0

## Table of Contents

- [Overview](#overview)
- [Breaking Changes](#breaking-changes)
  - [1. Default Container Runtime Changed](#1-default-container-runtime-changed)
  - [2. CRD Dependencies Removed](#2-crd-dependencies-removed)
  - [3. initContainers Removed](#3-initcontainers-removed)
  - [4. Log Path Configuration Removed](#4-log-path-configuration-removed)
- [Default Paths by Runtime](#default-paths-by-runtime)
- [Migration Scenarios](#migration-scenarios)
  - [Scenario 1: Using Containerd (Default)](#scenario-1-using-containerd-default---no-changes-needed)
  - [Scenario 2: Using CRI-O](#scenario-2-using-cri-o---minimal-changes)
  - [Scenario 3: Using Docker](#scenario-3-using-docker)
- [CRD Management Changes](#crd-management-changes)
- [CRD Installation Methods](#crd-installation-methods)
  - [Method 1: Standard Installation (Recommended)](#method-1-standard-installation-recommended)
  - [Method 2: Helm-Managed CRDs (Advanced)](#method-2-helm-managed-crds-advanced)
- [Migration from v3.x to v4.0](#migration-from-v3x-to-v40)
  - [Upgrading Standard Installation](#upgrading-standard-installation)
  - [Migrating to Helm-Managed CRDs](#migrating-to-helm-managed-crds)
  - [Fresh v4.0 Installation](#fresh-v40-installation)
- [Legacy Chart Migration](#legacy-chart-migration)
- [Forward Looking: Planned Changes in v5.0](#forward-looking-planned-changes-in-v50)

---

## Overview

Major changes/themes for v4.0:

1. **Container Runtime Simplification**: Removes dynamic detection for the `docker` runtime via initContainers and adopts static, configuration-based paths. The `docker` runtime has not been used widely since Kubernetes v1.24 (2022) and modern Kubernetes distributions now use the `containerd` runtime.

2. **CRD Management Modernization**: CRDs are now included in the main `fluent-operator` chart's `crds/` directory following Helm v3 best practices. A new optional `fluent-operator-crds` chart provides full Helm lifecycle management for advanced users. The legacy `fluent-bit-crds` and `fluentd-crds` dependency sub-charts have been removed.

3. **Simplified Architecture**: The main chart no longer has dependencies, providing a cleaner, more maintainable structure. Users choose their preferred CRD management method (standard or Helm-managed) based on their operational needs.

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

### 2. CRD Dependencies Removed

**What Changed:**

- Chart no longer has dependencies on `fluent-bit-crds` and `fluentd-crds` sub-charts
- `fluentbit.crdsEnable` and `fluentd.crdsEnable` values removed
- CRDs now included in main chart's `crds/` directory or managed via separate `fluent-operator-crds` chart

**Impact:**

- Simpler chart with no dependencies
- Values referencing `crdsEnable` will be ignored
- Default behavior: CRDs installed from `crds/` directory

**Who Is Affected:**

- Users who explicitly set `fluentbit.crdsEnable=false` or `fluentd.crdsEnable=false`
- Users relying on the old sub-chart dependency structure

**Migration:**

If you were disabling CRD installation in v3.x:

```diff
  fluentbit:
-   crdsEnable: false
    enable: true
```

In v4.0, use `--skip-crds` flag instead:

```bash
helm install fluent-operator fluent/fluent-operator --skip-crds
```

### 3. initContainers Removed

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

### 4. Log Path Configuration Removed

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

## CRD Management Changes

**What Changed:**

- CRDs now included in the main `fluent-operator` chart's `crds/` directory (Helm v3 standard)
- New optional `fluent-operator-crds` chart available for Helm-managed CRD lifecycle
- Removed dependencies on legacy `fluent-bit-crds` and `fluentd-crds` sub-charts
- Two installation methods now available depending on your CRD management needs

**Impact:**

- **Standard Installation**: CRDs in `crds/` directory are installed automatically but NOT upgraded/deleted by Helm
- **Advanced Installation**: Optional `fluent-operator-crds` chart provides full Helm lifecycle management
- Simpler, cleaner chart structure with no dependencies

**Who Is Affected:**

- All users upgrading from v3.x
- Users who want automatic CRD upgrades
- Users managing CRDs separately from the operator

---

## CRD Installation Methods

v4.0 provides two methods for managing CRDs. Choose based on your operational requirements:

### Method 1: Standard Installation (Recommended)

The main `fluent-operator` chart now includes CRDs in the `crds/` directory following Helm v3 best practices.

```bash
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update

helm install fluent-operator fluent/fluent-operator
```

**Behavior:**
- ✅ CRDs installed automatically on first install
- ⚠️ Helm does NOT upgrade CRDs on `helm upgrade` (manual upgrade required)
- ⚠️ Helm does NOT delete CRDs on `helm uninstall`

**When to use:**
- Simple installations
- You're comfortable with manual CRD upgrades
- Most users (recommended default)

**Upgrading CRDs:**

When upgrading, manually apply CRD updates before upgrading the chart:

```bash
# Update repository
helm repo update

# Extract and apply CRD updates
helm pull fluent/fluent-operator --version 4.0.0 --untar
kubectl apply -f fluent-operator/crds/

# Then upgrade the chart
helm upgrade fluent-operator fluent/fluent-operator --version 4.0.0
```

**Skipping CRDs:**

If you manage CRDs separately:

```bash
helm install fluent-operator fluent/fluent-operator --skip-crds
```

---

### Method 2: Helm-Managed CRDs (Advanced)

The new `fluent-operator-crds` chart provides full Helm management of CRDs.

```bash
# Step 1: Install CRDs with Helm management
helm install fluent-operator-crds fluent/fluent-operator-crds

# Step 2: Install operator (skip CRDs since already installed)
helm install fluent-operator fluent/fluent-operator --skip-crds
```

**Behavior:**
- ✅ CRDs automatically upgrade with `helm upgrade fluent-operator-crds`
- ✅ Fine-grained control (enable/disable Fluent Bit or Fluentd CRDs)
- ⚠️ CRDs deleted on `helm uninstall` (unless protected with annotation)

**When to use:**
- Advanced users
- GitOps workflows requiring full automation
- Organizations needing complete CRD lifecycle management

**Protecting CRDs from Deletion:**

```bash
helm install fluent-operator-crds fluent/fluent-operator-crds \
  --set additionalAnnotations."helm\.sh/resource-policy"=keep
```

**Selective CRD Installation:**

```bash
# Only Fluent Bit CRDs
helm install fluent-operator-crds fluent/fluent-operator-crds \
  --set fluentd.enabled=false

# Only Fluentd CRDs
helm install fluent-operator-crds fluent/fluent-operator-crds \
  --set fluent-bit.enabled=false
```

---

## Migration from v3.x to v4.0

### Upgrading Standard Installation

If you were using v3.x with default settings:

```bash
# Update repository
helm repo update

# Upgrade to v4.0 (existing CRDs are preserved)
helm upgrade fluent-operator fluent/fluent-operator --version 4.0.0
```

**Note:** Existing CRDs from v3.x will continue to work. The new CRDs in the `crds/` directory will only be installed on fresh installations. To update CRDs, manually apply them (see Method 1 above).

### Migrating to Helm-Managed CRDs

If you want to switch to full Helm management of CRDs:

```bash
# Step 1: Install the new CRDs chart (existing CRDs will be adopted)
helm install fluent-operator-crds fluent/fluent-operator-crds \
  --set additionalAnnotations."helm\.sh/resource-policy"=keep

# Step 2: Upgrade operator to v4.0 with --skip-crds
helm upgrade fluent-operator fluent/fluent-operator --version 4.0.0 --skip-crds
```

### Fresh v4.0 Installation

For new installations, simply choose your preferred method:

**Standard:**
```bash
helm install fluent-operator fluent/fluent-operator
```

**Helm-Managed:**
```bash
helm install fluent-operator-crds fluent/fluent-operator-crds
helm install fluent-operator fluent/fluent-operator --skip-crds
```

---

## Legacy Chart Migration

The following charts have been **removed** in v4.0:
- `fluent-bit-crds` (standalone chart) - Replaced by CRDs in main chart or `fluent-operator-crds`
- `fluentd-crds` (standalone chart) - Replaced by CRDs in main chart or `fluent-operator-crds`

**If you were using these charts directly:**

```bash
# Uninstall legacy charts
helm uninstall fluent-bit-crds
helm uninstall fluentd-crds

# Choose new method:
# Option A: Use main chart (CRDs preserved)
helm install fluent-operator fluent/fluent-operator

# Option B: Use new unified CRD chart
helm install fluent-operator-crds fluent/fluent-operator-crds
helm install fluent-operator fluent/fluent-operator --skip-crds
```

## Forward Looking: Planned Changes in v5.0

**Future Change (v5.0):**

- The `fluent-operator-env` _ConfigMap_, which is used to provide backwards compatibility with fluent-operator =<3.5, will be completely removed
