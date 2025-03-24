#!/bin/bash
namespace=kubesphere-logging-system
FluentbitOperator="fluentbit-operator"

function error_exit {
  echo "$1" 1>&2
  exit 1
}

function migrate(){
  ## Converting an existing configuration to a new one
  local list name labels spec ns size
  local oldKind=$1
  local newKind=$2
  list=$(kubectl get "$oldKind.logging.kubesphere.io" -A -o json) || error_exit "Cannot get resource $oldKind"
  mapfile -t name < <(echo "$list" | jq -r '.items[].metadata.name | @json')
  mapfile -t labels < <(echo "$list" | jq -r '.items[].metadata.labels | @json')
  mapfile -t spec < <(echo "$list" | jq -r '.items[].spec | @json')
  mapfile -t ns < <(echo "$list" | jq -r '.items[].metadata.namespace | @json')
  size=${#spec[*]}
  echo "Number of original $oldKind configuration files:$size"
  for ((i=0; i < size; i++)); do
  if [[ "$newKind" = "fluentbits" ]]; then
    cluster_resource_list[i]="{
    \"apiVersion\": \"fluentbit.fluent.io/v1alpha2\",
    \"kind\": \"${newKind}\",
    \"metadata\": {
    \"name\": ${name[i]},
    \"labels\": ${labels[i]},
    \"namespace\": \"${ns[i]}\"
    },
    \"spec\": ${spec[i]}
    }"
  else
    cluster_resource_list[i]="{
    \"apiVersion\": \"fluentbit.fluent.io/v1alpha2\",
    \"kind\": \"${oldKind}\",
    \"metadata\": {
    \"name\": ${name[i]},
    \"labels\": ${labels[i]}
    },
    \"spec\": ${spec[i]}
    }"
  fi
  done

  ## Uninstall the fluentbit-operator and the original configuration
  for ((i=0; i < size; i++)); do
  echo "${name[i]}"
    temp="${name[i]//\"/}"
    echo "$temp"
    kubectl delete "$oldKind.logging.kubesphere.io" "$temp" -n ${namespace}
  done

  for ((i=0; i<size; i++)); do
    kubectl apply -f "${cluster_resource_list[i]}" || error_exit "Cannot apply resource $oldKind"
  done
}

migrate "Input" "ClusterInput"
migrate "Parser" "ClusterParser"
migrate "Filter" "ClusterFilter"
migrate "Output" "ClusterOutput"
migrate "FluentbitConfig" "ClusterFluentBitConfig"
migrate "FluentBit" "FluentBit"

# Determine if Deployment exists
if kubectl get deployment -n $namespace $FluentbitOperator >/dev/null 2>&1; then
  # Delete Deployment if it exists
  kubectl delete deployment -n $namespace $FluentbitOperator
  kubectl delete clusterrolebinding kubesphere:operator:fluentbit-operator
  kubectl delete clusterrole kubesphere:operator:fluentbit-operator
  kubectl delete serviceaccount fluentbit-operator -n $namespace
  echo "Deployment $FluentbitOperator deleted"
else
  # If it does not exist, output the message
  echo "Deployment $FluentbitOperator does not exist"
fi

## Delete the old crd
kubectl get crd -o=jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | grep "logging.kubesphere.io" | xargs -I crd_name kubectl delete crd crd_name
