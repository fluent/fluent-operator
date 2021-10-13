VERSION?=$(shell cat VERSION | tr -d " \t\n\r")
# Image URL to use all building/pushing image targets
FB_IMG ?= kubesphere/fluent-bit:v1.8.3
OP_IMG ?= kubesphere/fluentbit-operator:$(VERSION)
AMD64 ?= -amd64
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager cmd/manager/main.go

binary: 
	go build -o bin/manager cmd/manager/main.go
	go build -o bin/watcher cmd/fluent-bit-watcher/main.go

verify: verify-crds

verify-crds:
	sudo chmod a+x ./hack/verify-crds.sh && ./hack/verify-crds.sh

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	kubectl kustomize config/crd | sed -e '/creationTimestamp/d' > manifests/setup/fluentbit-operator-crd.yaml
	kubectl kustomize manifests/setup | sed -e '/creationTimestamp/d' > manifests/setup/setup.yaml

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."
	./hack/update-codegen.sh

# Build all amd64/arm64 docker images
build: build-op

# Build amd64/arm64 Fluent Bit container image
build-fb:
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/fluent-bit-watcher/Dockerfile . -t ${FB_IMG}

# Build amd64/arm64 Fluent Bit Operator container image
build-op: test
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/manager/Dockerfile . -t ${OP_IMG}

# Build all amd64 docker images
build-amd64: build-op-amd64 build-fb-amd64

# Build amd64 Fluent Bit container image
build-fb-amd64:
	docker build -f cmd/fluent-bit-watcher/Dockerfile . -t ${FB_IMG}${AMD64}

# Build amd64 Fluent Bit Operator container image
build-op-amd64: 
	docker build -f cmd/manager/Dockerfile . -t ${OP_IMG}${AMD64}

# Push the amd64 docker image
push-amd64:
	docker push ${OP_IMG}${AMD64}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go install -v sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
