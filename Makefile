
# Image URL to use all building/pushing image targets
OP_IMG ?= kubespheredev/fluentbit-operator:v0.7.0
MIGRATOR_IMG ?= kubespheredev/fluentbit-operator:migrator
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
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

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
	kustomize build config/crd | sed -e '/creationTimestamp/d' > manifests/setup/fluentbit-operator-crd.yaml

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

# Build all docker images for amd64 and arm64
build: build-op

# Build the docker image for amd64 and arm64
build-op: test
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/manager/Dockerfile . -t ${OP_IMG}

# Build the docker image for amd64 and arm64
build-migtator: test
	docker buildx build --push --platform linux/amd64,linux/arm64 -f cmd/migrator/Dockerfile . -t ${MIGRATOR_IMG}

# Build all docker images for amd64
build-amd64: build-op-amd64 build-migtator-amd64

# Build the docker image for amd64
build-op-amd64: test
	docker build -f cmd/manager/Dockerfile . -t ${OP_IMG}${AMD64}

# Build the docker image for amd64
build-migtator-amd64: test
	docker build -f cmd/migrator/Dockerfile . -t ${MIGRATOR_IMG}${AMD64}

# Push the docker image
push-amd64:
	docker push ${OP_IMG}${AMD64}
	docker push ${MIGRATOR_IMG}${AMD64}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
