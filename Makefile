MAKEFLAGS = --warn-undefined-variables

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

VERSION ?= $(shell cat VERSION | tr -d " \t\n\r")
FB_VERSION?=$(shell grep -v '^#' cmd/fluent-watcher/fluentbit/VERSION | tr -d " \t\n\r")
# Image URL to use all building/pushing image targets
FB_IMG ?= ghcr.io/fluent/fluent-operator/fluent-bit:v${FB_VERSION}
FB_IMG_DEBUG ?= ghcr.io/fluent/fluent-operator/fluent-bit:v${FB_VERSION}-debug
FD_IMG ?= ghcr.io/fluent/fluent-operator/fluentd:v1.19.1
FO_IMG ?= kubesphere/fluent-operator:$(VERSION)

ARCH ?= arm64

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= crd:generateEmbeddedObjectMeta=true,allowDangerousTypes=true
OPERATOR_SDK_VERSION ?= v1.42.0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
	GOBIN = $(shell go env GOPATH)/bin
else
	GOBIN = $(shell go env GOBIN)
endif

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

shellcheck:
	@find . -type f -name *.sh -exec docker run --rm -v $(shell pwd):/mnt koalaman/shellcheck:stable {} +

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=charts/fluent-bit-crds/templates
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=charts/fluentd-crds/templates
	kubectl kustomize config/crd/bases/ | sed -e '/creationTimestamp/d' > manifests/setup/fluent-operator-crd.yaml
	kubectl kustomize manifests/setup/ | sed -e '/creationTimestamp/d' > manifests/setup/setup.yaml
	hack/mutate-crds.sh

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

test: manifests generate fmt vet setup-envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test $$(go list ./... | grep -v /e2e) -coverprofile cover.out

# Utilize Kind or modify the e2e tests to load the image locally, enabling compatibility with other vendors.
.PHONY: test-e2e  # Run the e2e tests against a Kind k8s instance that is spun up.
test-e2e: fluentd-e2e

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

.PHONY: lint-config
lint-config: golangci-lint ## Verify golangci-lint linter configuration
	$(GOLANGCI_LINT) config verify

##@ Build

binary:
	go build -o bin/fb-manager ./cmd/fluent-manager
	go build -o bin/fb-watcher ./cmd/fluent-watcher/fluentbit
	go build -o bin/fd-watcher ./cmd/fluent-watcher/fluentd

verify: verify-crds verify-codegen

verify-crds:
	./hack/verify-crds.sh

verify-codegen:
	./hack/verify-codegen.sh

build: generate fmt vet ## Build manager binary.
	go build -o bin/fluent-manager ./cmd/fluent-manager
	go build -o bin/fb-watcher ./cmd/fluent-watcher/fluentbit
	go build -o bin/fd-watcher ./cmd/fluent-watcher/fluentd

run: manifests generate fmt vet ## Run a controller from your host.
	go run cmd/fluent-manager/main.go

# Build amd64/arm64 Fluent Operator container image
.PHONY: build-op
build-op:
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}

# Build amd64/arm64 Fluent Bit container image
.PHONY: build-fb
build-fb: prepare-build
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

.PHONY: build-fb-debug
build-fb-debug: prepare-build
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-watcher/fluentbit/Dockerfile.debug . -t ${FB_IMG_DEBUG}

# Build all amd64 docker images
.PHONY: build-amd64
build-amd64: build-op-amd64 build-fb-amd64 build-fd-amd64

# Build all arm64 docker images
.PHONY: build-arm64
build-arm64: build-op-arm64 build-fb-arm64 build-fd-arm64

# Build amd64 Fluent Operator container image
.PHONY: build-op-amd64
build-op-amd64:
	docker build --platform=linux/amd64 -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}

# Build arm64 Fluent Operator container image
.PHONY: build-op-arm64
build-op-arm64:
	docker build --platform=linux/arm64 -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}

# Build amd64 Fluent Bit container image
.PHONY: build-fb-amd64
build-fb-amd64:
	docker build --platform=linux/amd64 -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

# Build arm64 Fluent Bit container image
.PHONY: build-fb-arm64
build-fb-arm64:
	docker build --platform=linux/arm64 -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

# Build amd64 Fluentd container image
.PHONY: build-fd-amd64
build-fd-amd64:
	docker build --platform=linux/amd64 -f cmd/fluent-watcher/fluentd/Dockerfile . -t ${FD_IMG}

# Build arm64 Fluentd container image
.PHONY: build-fd-arm64
build-fd-arm64:
	docker build --platform=linux/arm64 -f cmd/fluent-watcher/fluentd/Dockerfile . -t ${FD_IMG}

# Prepare for arm64 building
prepare-build:
	cmd/fluent-watcher/hooks/post-hook.sh

# Push the amd64 docker image
push-amd64:
	docker push ${FO_IMG}${AMD64}

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd/bases/ | kubectl create -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd/bases/ | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	kubectl create -f manifests/setup/setup.yaml

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	kubectl delete -f manifests/setup/setup.yaml

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GINKGO = $(LOCALBIN)/ginkgo
CODE_GENERATOR = $(LOCALBIN)/code-generator
KIND = $(LOCALBIN)/kind

## Tool Versions
KUSTOMIZE_VERSION ?= v5.6.0
CONTROLLER_TOOLS_VERSION ?= v0.18.0
#ENVTEST_VERSION is the version of controller-runtime release branch to fetch the envtest setup script (i.e. release-0.20)
ENVTEST_VERSION ?= $(shell go list -m -f "{{ .Version }}" sigs.k8s.io/controller-runtime | awk -F'[v.]' '{printf "release-%d.%d", $$2, $$3}')
#ENVTEST_K8S_VERSION is the version of Kubernetes to use for setting up ENVTEST binaries (i.e. 1.31)
ENVTEST_K8S_VERSION ?= $(shell go list -m -f "{{ .Version }}" k8s.io/api | awk -F'[v.]' '{printf "1.%d", $$3}')
GOLANGCI_LINT_VERSION ?= v2.1.0
GINKGO_VERSION ?= v2.27.2
CODE_GENERATOR_VERSION ?= v0.32.3
KIND_VERSION ?= v0.30.0

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: setup-envtest
setup-envtest: envtest ## Download the binaries required for ENVTEST in the local bin directory.
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	$(call go-install-tool,$(GINKGO),github.com/onsi/ginkgo/v2/ginkgo,$(GINKGO_VERSION))

.PHONY: code-generator
code-generator: $(CODE_GENERATOR) ## Download code-generator locally if necessary.
$(CODE_GENERATOR): $(LOCALBIN)
	$(call go-install-tool,$(CODE_GENERATOR),k8s.io/code-generator,$(CODE_GENERATOR_VERSION))

KIND_CLUSTER ?= fluent-operator-test-e2e

.PHONY: setup-test-e2e
setup-test-e2e: ## Set up a Kind cluster for e2e tests if it does not exist
	@command -v $(KIND) >/dev/null 2>&1 || { \
		echo "Kind is not installed. Please install Kind manually."; \
		exit 1; \
	}
	@case "$$($(KIND) get clusters)" in \
		*"$(KIND_CLUSTER)"*) \
		echo "Kind cluster '$(KIND_CLUSTER)' already exists. Skipping creation." ;; \
		*) \
		echo "Creating Kind cluster '$(KIND_CLUSTER)'..."; \
		$(KIND) create cluster --name $(KIND_CLUSTER) ;; \
	esac

.PHONY: cleanup-test-e2e
cleanup-test-e2e:
	$(KIND) delete cluster --name $(KIND_CLUSTER)

.PHONY: kind
kind: $(KIND) ## Download code-generator locally if necessary.
$(KIND): $(LOCALBIN)
	$(call go-install-tool,$(KIND),sigs.k8s.io/kind,$(KIND_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef

.PHONY: operator-sdk
OPERATOR_SDK ?= $(LOCALBIN)/operator-sdk
operator-sdk: ## Download operator-sdk locally if necessary.
ifeq (,$(wildcard $(OPERATOR_SDK)))
ifeq (, $(shell which operator-sdk 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPERATOR_SDK)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPERATOR_SDK) https://github.com/operator-framework/operator-sdk/releases/download/$(OPERATOR_SDK_VERSION)/operator-sdk_$${OS}_$${ARCH} ;\
	chmod +x $(OPERATOR_SDK) ;\
	}
else
OPERATOR_SDK = $(shell which operator-sdk)
endif
endif
.PHONY: bundle
bundle: manifests kustomize operator-sdk ## Generate bundle manifests and metadata, then validate generated files.
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle $(BUNDLE_GEN_FLAGS)
	$(OPERATOR_SDK) bundle validate ./bundle

.PHONY: bundle-build
bundle-build: ## Build the bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: bundle-push
bundle-push: ## Push the bundle image.
	$(MAKE) docker-push IMG=$(BUNDLE_IMG)

.PHONY: opm
OPM = $(LOCALBIN)/opm
opm: ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
ifeq (,$(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPM)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.55.0/$${OS}-$${ARCH}-opm ;\
	chmod +x $(OPM) ;\
	}
else
OPM = $(shell which opm)
endif
endif

# A comma-separated list of bundle images (e.g. make catalog-build BUNDLE_IMGS=example.com/operator-bundle:v0.1.0,example.com/operator-bundle:v0.2.0).
# These images MUST exist in a registry and be pull-able.
BUNDLE_IMGS ?= $(BUNDLE_IMG)

# The image tag given to the resulting catalog image (e.g. make catalog-build CATALOG_IMG=example.com/operator-catalog:v0.2.0).
CATALOG_IMG ?= $(IMAGE_TAG_BASE)-catalog:v$(VERSION)

# Set CATALOG_BASE_IMG to an existing catalog image tag to add $BUNDLE_IMGS to that image.
ifneq ($(origin CATALOG_BASE_IMG), undefined)
FROM_INDEX_OPT := --from-index $(CATALOG_BASE_IMG)
endif

# Build a catalog image by adding bundle images to an empty catalog using the operator package manager tool, 'opm'.
# This recipe invokes 'opm' in 'semver' bundle add mode. For more information on add modes, see:
# https://github.com/operator-framework/community-operators/blob/7f1438c/docs/packaging-operator.md#updating-your-existing-operator
.PHONY: catalog-build
catalog-build: opm ## Build a catalog image.
	$(OPM) index add --container-tool docker --mode semver --tag $(CATALOG_IMG) --bundles $(BUNDLE_IMGS) $(FROM_INDEX_OPT)

# Push the catalog image.
.PHONY: catalog-push
catalog-push: ## Push a catalog image.
	$(MAKE) docker-push IMG=$(CATALOG_IMG)

go-deps: # download go dependencies
	go mod download

docs-update: # update api docs
	go run ./cmd/doc-gen/main.go

fluentd-e2e: ginkgo # make e2e tests
	tests/scripts/fluentd_e2e.sh

fluentd-helm-e2e: ginkgo # make helm e2e tests
	tests/scripts/fluentd_helm_e2e.sh

update-helm-package: # update helm repo
	./hack/update-helm-package.sh

.PHONY: helm-docs
helm-docs:
	cd charts/fluentd-crds && helm-docs
	cd charts/fluent-bit-crds && helm-docs
