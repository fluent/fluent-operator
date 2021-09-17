REV ?= $(shell git rev-parse --short HEAD)
PREFIX ?= do-
IMAGE_TAG ?= $(REV:$(PREFIX)%=%)

ifdef release
	REV = $(shell git rev-list --tags --max-count=1)
	GIT_TAG ?= $(shell git describe --tags $(REV))
	IMAGE_TAG = $(GIT_TAG:$(PREFIX)%=%)
endif

$(info using image tag: $(IMAGE_TAG))

.PHONY: image-operator
image-operator:
	docker build -f cmd/manager/Dockerfile -t digitaloceanapps/fluent-bit-operator:$(IMAGE_TAG) .
ifdef latest
	docker tag digitaloceanapps/fluent-bit-operator:$(IMAGE_TAG) digitaloceanapps/fluent-bit-operator:latest
endif

.PHONY: image-push-operator
image-push-operator: image-operator
	docker push digitaloceanapps/fluent-bit-operator:$(IMAGE_TAG)
ifdef latest
	docker push digitaloceanapps/fluent-bit-operator:latest
endif
