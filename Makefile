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
FD_IMG ?= ghcr.io/fluent/fluent-operator/fluentd:v1.17.1
FO_IMG ?= kubesphere/fluent-operator:$(VERSION)
FD_IMG_BASE ?= ghcr.io/fluent/fluent-operator/fluentd:v1.17.0-arm64-base

ARCH ?= arm64

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= crd:generateEmbeddedObjectMeta=true,allowDangerousTypes=true

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
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=charts/fluent-operator/charts/fluent-bit-crds/crds
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=charts/fluent-operator/charts/fluentd-crds/crds
	kubectl kustomize config/crd/bases/ | sed -e '/creationTimestamp/d' > manifests/setup/fluent-operator-crd.yaml
	kubectl kustomize manifests/setup/ | sed -e '/creationTimestamp/d' > manifests/setup/setup.yaml

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin

install-setup-envtest: ## Install the setup-envtest tool if it is not already installed
	if ! command -v setup-envtest &> /dev/null; then \
		go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest; \
	fi

setup-envtest: install-setup-envtest ## Download and set up the envtest binary
	source <(setup-envtest use -p env)

test: manifests generate fmt vet setup-envtest ## Run tests.
	go test ./apis/... -coverprofile cover.out

# Utilize Kind or modify the e2e tests to load the image locally, enabling compatibility with other vendors.
.PHONY: test-e2e  # Run the e2e tests against a Kind k8s instance that is spun up.
test-e2e:
	go test ./test/e2e/ -v -ginkgo.v

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

##@ Build

binary:
	go build -o bin/fb-manager cmd/fluent-manager/main.go
	go build -o bin/fb-watcher cmd/fluent-watcher/fluentbit/main.go
	go build -o bin/fd-watcher cmd/fluent-watcher/fluentd/main.go

verify: verify-crds verify-codegen

verify-crds:
	./hack/verify-crds.sh

verify-codegen:
	./hack/verify-codegen.sh

build: generate fmt vet ## Build manager binary.
	go build -o bin/fluent-manager cmd/fluent-manager/main.go
	go build -o bin/fb-watcher cmd/fluent-watcher/fluentbit/main.go
	go build -o bin/fd-watcher cmd/fluent-watcher/fluentd/main.go

run: manifests generate fmt vet ## Run a controller from your host.
	go run cmd/fluent-manager/main.go

# Build amd64/arm64 Fluent Operator container image
build-op:
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}

# Build amd64/arm64 Fluent Bit container image
build-fb: prepare-build
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

build-fb-debug: prepare-build
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-watcher/fluentbit/Dockerfile.debug . -t ${FB_IMG_DEBUG}

# Build all amd64 docker images
build-amd64: build-op-amd64 build-fb-amd64 build-fd-amd64

# Build amd64 Fluent Operator container image
build-op-amd64:
	docker build --platform=linux/amd64 -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}

# Build amd64 Fluent Bit container image
build-fb-amd64:
	docker build --platform=linux/amd64 -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

# Build amd64 Fluentd container image
build-fd-amd64:
	docker build --platform=linux/amd64 -f cmd/fluent-watcher/fluentd/Dockerfile.amd64 . -t ${FD_IMG}

build-fd-arm64-base: prepare-build
	docker buildx build --push --platform linux/arm64 -f cmd/fluent-watcher/fluentd/Dockerfile.arm64.base . -t ${FD_IMG_BASE}

# Use docker buildx to build arm64 Fluentd container image
build-fd-arm64: prepare-build
	docker buildx build --push --platform linux/arm64 -f cmd/fluent-watcher/fluentd/Dockerfile.arm64.quick . -t ${FD_IMG}${ARCH} --build-arg ${FD_IMG_BASE} --build-arg ${FD_IMG_BASE_TAG}

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
CONTROLLER_TOOLS_VERSION ?= v0.15.0
ENVTEST_VERSION ?= release-0.19
GOLANGCI_LINT_VERSION ?= v1.59.1
GINKGO_VERSION ?= v2.23.4
CODE_GENERATOR_VERSION ?= v0.32.3
KIND_VERSION ?= v0.17.0

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	$(call go-install-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5,$(KUSTOMIZE_VERSION))

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	$(call go-install-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen,$(CONTROLLER_TOOLS_VERSION))

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	$(call go-install-tool,$(GINKGO),github.com/onsi/ginkgo/v2/ginkgo,$(GINKGO_VERSION))

.PHONY: code-generator
code-generator: $(CODE_GENERATOR) ## Download code-generator locally if necessary.
$(CODE_GENERATOR): $(LOCALBIN)
	$(call go-install-tool,$(CODE_GENERATOR),k8s.io/code-generator,$(CODE_GENERATOR_VERSION))

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
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.23.0/$${OS}-$${ARCH}-opm ;\
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

e2e: ginkgo # make e2e tests
	tests/scripts/fluentd_e2e.sh

helm-e2e: ginkgo # make helm e2e tests
	tests/scripts/fluentd_helm_e2e.sh

update-helm-package: # update helm repo
	./hack/update-helm-package.sh
