VERSION?=$(shell cat VERSION | tr -d " \t\n\r")
# Image URL to use all building/pushing image targets
FB_IMG ?= kubesphere/fluent-bit:v1.8.11
FD_IMG ?= kubesphere/fluentd:v1.14.4
FO_IMG ?= kubesphere/fluent-operator:$(VERSION)

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

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
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentbit/..." output:crd:artifacts:config=charts/fluent-operator/crds
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/fluentd/..." output:crd:artifacts:config=charts/fluent-operator/crds
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

verify: verify-crds

verify-crds:
	sudo chmod a+x ./hack/verify-crds.sh && ./hack/verify-crds.sh

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
build-fb:
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

# Build amd64/arm64 Fluentd container image
build-fd:
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-watcher/fluentd/Dockerfile.complete . -t ${FD_IMG}

# Build all amd64 docker images
build-amd64: build-op-amd64 build-fb-amd64 build-fd-amd64

# Build amd64 Fluent Operator container image
build-op-amd64:
	docker build -f cmd/fluent-manager/Dockerfile . -t ${FO_IMG}

# Build amd64 Fluent Bit container image
build-fb-amd64:
	docker build -f cmd/fluent-watcher/fluentbit/Dockerfile . -t ${FB_IMG}

# Build amd64 Fluent Bit container image
build-fd-amd64:
	docker build -f cmd/fluent-watcher/fluentd/Dockerfile.complete . -t ${FD_IMG}

# Push the amd64 docker image
push-amd64:
	docker push ${FO_IMG}${AMD64}

##@ Deployment

install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd/bases/ | kubectl apply -f -

uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd/bases/ | kubectl delete -f -

deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	kubectl apply -k manifests/setup/setup.yaml

undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config.
	kubectl delete -k manifests/setup/setup.yaml

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
controller-gen: go-deps ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1)

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize: go-deps ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

go-deps: # download go dependencies
	go mod download

docs-update:
	go run ./cmd/doc-gen/main.go

e2e:
	chmod a+x tests/scripts/fluentd_e2e.sh && bash tests/scripts/fluentd_e2e.sh