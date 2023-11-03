VERSION?=$(shell cat VERSION | tr -d " \t\n\r")
# Image URL to use all building/pushing image targets
FB_IMG ?= kubesphere/fluent-bit:v2.1.10
FB_IMG_DEBUG ?= kubesphere/fluent-bit:v2.1.10-debug
FD_IMG ?= kubesphere/fluentd:v1.15.3
FO_IMG ?= kubesphere/fluent-operator:$(VERSION)
FD_IMG_BASE ?= kubesphere/fluentd:v1.15.3-arm64-base

ARCH ?= arm64

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:generateEmbeddedObjectMeta=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

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

manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=config/crd/bases
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=charts/fluent-operator/charts/fluent-bit-crds/crds
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=charts/fluent-operator/charts/fluentd-crds/crds
	kubectl kustomize config/crd/bases/ | sed -e '/creationTimestamp/d' > manifests/setup/fluent-operator-crd.yaml
	kubectl kustomize manifests/setup/ | sed -e '/creationTimestamp/d' > manifests/setup/setup.yaml

generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."
	./hack/update-codegen.sh

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./apis/... -coverprofile cover.out

##@ Build

binary:
	go build -o bin/fb-manager cmd/fluent-manager/main.go
	go build -o bin/fb-watcher cmd/fluent-watcher/fluentbit/main.go
	go build -o bin/fd-watcher cmd/fluent-watcher/fluentd/main.go

verify: verify-crds verify-codegen

verify-crds:
	chmod a+x ./hack/verify-crds.sh && ./hack/verify-crds.sh

verify-codegen:
	chmod a+x ./hack/verify-codegen.sh && ./hack/verify-codegen.sh

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
	docker buildx build --push --platform linux/arm64 -f cmd/fluent-watcher/fluentd/Dockerfile.arm64.quick . -t ${FD_IMG}${ARCH}

# Prepare for arm64 building
prepare-build:
	chmod +x cmd/fluent-watcher/hooks/post-hook.sh && bash cmd/fluent-watcher/hooks/post-hook.sh

# Push the amd64 docker image
push-amd64:
	docker push ${FO_IMG}${AMD64}

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd/bases/ | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd/bases/ | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	kubectl apply -f manifests/setup/setup.yaml

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	kubectl delete -f manifests/setup/setup.yaml

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: go-deps ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.3)

GINKGO = $(shell pwd)/bin/ginkgo
ginkgo: go-deps ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(GINKGO),github.com/onsi/ginkgo/ginkgo@v1.16.5)


KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: go-deps ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v5@v5.0.0)

CODE_GENERATOR = $(shell go env GOPATH)/pkg/mod/k8s.io/code-generator@v0.26.1
code-generator: go-deps ## Download code-generator locally if necessary
	$(call go-get-tool,$(CODE_GENERATOR),k8s.io/code-generator@v0.26.1)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

go-deps: # download go dependencies
	go get k8s.io/code-generator@v0.26.1
	go mod download

docs-update: # update api docs
	go run ./cmd/doc-gen/main.go

e2e: ginkgo # make e2e tests
	chmod a+x tests/scripts/fluentd_e2e.sh && bash tests/scripts/fluentd_e2e.sh

helm-e2e: ginkgo # make helm e2e tests
	chmod a+x tests/scripts/fluentd_helm_e2e.sh && bash tests/scripts/fluentd_helm_e2e.sh

update-helm-package: # update helm repo
	chmod a+x ./hack/update-helm-package.sh && ./hack/update-helm-package.sh
