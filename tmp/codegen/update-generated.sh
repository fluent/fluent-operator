#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
deepcopy \
kubesphere.io/fluentbit-operator/pkg/generated \
kubesphere.io/fluentbit-operator/pkg/apis \
fluentbit:v1alpha1 \
--go-header-file "./tmp/codegen/boilerplate.go.txt"
