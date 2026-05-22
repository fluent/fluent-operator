# Fluent Operator Helm Chart Changelog

> [!NOTE]
> All notable changes to this project will be documented in this file; the format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!--
### Added - For new features.
### Changed - For changes in existing functionality.
### Deprecated - For soon-to-be removed features.
### Removed - For now removed features.
### Fixed - For any bug fixes.
### Security - In case of vulnerabilities.
-->

## [UNRELEASED]

## [v4.1.0] - 2026-05-21

### Added

- Operator deployment now includes configurable liveness and readiness probes (hitting `/healthz` and `/readyz`) ([#1956](https://github.com/fluent/fluent-operator/pull/1956))

### Changed

- Hardened default `podSecurityContext` and `securityContext` for the operator: `runAsNonRoot`, `runAsUser/Group 65532`, `readOnlyRootFilesystem`, drop `ALL` capabilities, `seccompProfile: RuntimeDefault` ([#1956](https://github.com/fluent/fluent-operator/pull/1956))
- Bumped default Fluent Bit image tag to `v5.0.5` ([#1968](https://github.com/fluent/fluent-operator/pull/1968))
- Bumped fluent-operator to v3.8.0

## [v4.0.0] - 2026-04-19

### Changed

- Initial release of v4 chart

<!--
RELEASE LINKS
-->
[UNRELEASED]: https://github.com/fluent/helm-charts/tree/main/charts/fluent-operator
[v4.1.0]: https://github.com/fluent/helm-charts/releases/tag/fluent-operator-4.1.0
[v4.0.0]: https://github.com/fluent/helm-charts/releases/tag/fluent-operator-4.0.0
