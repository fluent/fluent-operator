# fluent-watcher

This directory contains `Dockerfiles` for the custom build of both `fluentd` and `fluent-bit` that `fluent-operator` uses.

These custom builds include a wrapper that watches for configuration changes and reloads/restarts the `fluentd` and `fluent-bit` processes.  These images are required for running the fluent-operator.

## Versioning

These images track upstream versions/tags.  For example, Fluent Operator's `ghcr.io/fluent/fluent-operator/fluent-bit:3.1.12` image is based off of the upstream `fluentd/fluent-bit:3.1.12` image.

Occasionally, changes to Fluent Operator's customizations need to be introduced into the images out of band of upstream version changes.  When this happens, we add a "patch version" to the tag.  For example, the `ghcr.io/fluent/fluent-operator/fluent-bit:3.1.12-1` image remains based off of the upstream `fluent/fluent-bit:3.1.2` image but contains a Fluent Operator-specific patch (`-1`).

We strive to never overwrite existing image tags (eg, `ghcr.io/fluent/fluent-operator/fluent-bit:3.1.12`) with customizations and/or patches.

## Building

As a maintainer, to build the `ghcr.io/fluent/fluent-operator/fluent-bit` and `ghcr.io/fluent/fluent-operator/fluentd` images, you can run the "Build Fluent Bit image" or "Publish Fluentd image" Github Action workflows in the "Actions" tab of this repository.

Always specify the upstream Fluent Bit version (eg, `3.1.2`, `3.1.3`, etc) when running this workflow.  If the CI workflow detects that an image tag already exists for the version specified, it will assume that a patch release needs to be built and will automatically add a patch version to the image tag (eg, `3.1.2-1`, `3.1.2-2`, etc).