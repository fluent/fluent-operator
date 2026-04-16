# Copilot Instructions for fluent-operator

## Project Summary

Fluent Operator is a Kubernetes operator (written in Go) that manages Fluent Bit and Fluentd deployments via Custom Resource Definitions (CRDs). It is built with **kubebuilder v3** and **controller-runtime**. The Go module path is `github.com/fluent/fluent-operator/v3`.

- **Language**: Go (version specified in `go.mod`, currently 1.26.2)
- **Framework**: kubebuilder / controller-runtime
- **License**: Apache 2.0
- **Default branch**: `master`

## Build, Test, and Lint Commands

Always run `go mod download` before building or testing.

| Task | Command | Notes |
|---|---|---|
| **Download deps** | `go mod download` | Run first |
| **Generate deepcopy** | `make generate` | Regenerates `zz_generated.deepcopy.go` files |
| **Generate CRDs & manifests** | `make manifests` | Requires `kubectl` on PATH for kustomize steps |
| **Build binaries** | `make build` | Runs generate + fmt + vet, outputs to `bin/` |
| **Build only (no codegen)** | `make binary` | Builds `bin/fb-manager`, `bin/fb-watcher`, `bin/fd-watcher` |
| **Run unit tests** | `make test` | Runs manifests + generate + fmt + vet + envtest setup, then `go test` with race detector. Excludes `/e2e` tests. |
| **Run linter** | `make lint` | Uses golangci-lint v2 with config in `.golangci.yml` |
| **Lint fix** | `make lint-fix` | Auto-fixes lint issues |
| **Update API docs** | `make docs-update` | Runs `go run ./cmd/doc-gen/main.go` |
| **Verify CRDs match** | `make verify` | Runs `hack/verify-crds.sh` and `hack/verify-codegen.sh` |
| **Format** | `go fmt ./...` | |
| **Vet** | `go vet ./...` | |

### Important: After changing any API type in `apis/`

Always run in this order:
1. `make generate` — regenerates DeepCopy methods
2. `make manifests` — regenerates CRD YAML files in `config/crd/bases/`, `charts/fluent-operator/crds/`, `charts/fluent-operator-crds/templates/`, and `manifests/setup/`
3. `make docs-update` — regenerates API documentation

CI will fail if generated files are out of date (`git diff --exit-code` is checked).

### Tool versions (auto-installed to `./bin/`)

- controller-gen: v0.18.0
- kustomize: v5.6.0
- golangci-lint: v2.6.2
- ginkgo: v2.27.2
- kind: v0.30.0

## Repository Layout

```
apis/                        # CRD type definitions (the most commonly edited area)
  fluentbit/v1alpha2/        # Fluent Bit CRDs — types, plugins, webhook
    plugins/                 # Plugin type definitions (input, filter, output, parser)
  fluentd/v1alpha1/          # Fluentd CRDs — types, plugins
    plugins/                 # Fluentd plugin types (input, filter, output)
  generated/                 # Generated client code
cmd/
  fluent-manager/            # Main operator binary entrypoint (main.go)
  fluent-watcher/fluentbit/  # Fluent Bit watcher binary + Dockerfile
  fluent-watcher/fluentd/    # Fluentd watcher binary + Dockerfile
  doc-gen/                   # API doc generation tool
controllers/                 # Reconciler implementations
  fluentbit_controller.go
  fluentbitconfig_controller.go
  fluentd_controller.go
  fluentdconfig_controller.go
  collector_controller.go
pkg/
  operator/                  # Kubernetes resource builders for DaemonSet/StatefulSet
  fluentd/                   # Fluentd config rendering
  utils/                     # Shared utilities
config/                      # kubebuilder/kustomize config (RBAC, CRD bases, webhooks)
charts/                      # Helm charts (fluent-operator, fluent-operator-crds)
manifests/                   # YAML manifests for kubectl-based deployment
hack/                        # Shell scripts for verification and CRD mutation
tests/                       # E2E test scripts and specs
```

## CI Checks (run on every PR)

1. **Lint** (`.github/workflows/lint.yml`): `golangci-lint run` — runs on all pushes/PRs.
2. **Tests** (`.github/workflows/test.yml`): `go mod tidy && make test` — runs on all pushes/PRs.
3. **Main CI** (`.github/workflows/main.yaml`): Triggered on changes to Go/API/controller files:
   - `make generate manifests docs-update` then `git diff --exit-code` (generated code must be committed)
   - `make shellcheck` (shell script linting via Docker)
   - `make test`
   - `make verify` (CRDs and codegen match)
   - `make binary` (binary build)
4. **Helm CI** (`.github/workflows/helm-ci.yaml`): Helm chart linting and testing.

## Coding Conventions

- All API types use kubebuilder markers (`+kubebuilder:...`) for validation, defaulting, and CRD generation.
- New Go files must include the Apache 2.0 license header (see `hack/boilerplate.go.txt`).
- Commits require DCO sign-off (`Signed-off-by` line).
- Linter config is in `.golangci.yml` — enabled linters include: errcheck, govet, staticcheck, revive, misspell, unused, ginkgolinter, and others.
- Tests use Ginkgo/Gomega and envtest for controller tests.

## Tips

- Trust these instructions first. Only search the repo if information here is incomplete or incorrect.
- The `make test` target handles envtest binary setup automatically.
- When adding a new Fluent Bit or Fluentd plugin, add the type in `apis/fluentbit/v1alpha2/plugins/` or `apis/fluentd/v1alpha1/plugins/`, then wire it into the relevant config type and run all three generation steps.
- The `manifests/setup/setup.yaml` file is auto-generated — never edit it directly.
- CRD files in `config/crd/bases/` and `charts/` are auto-generated — never edit them directly.