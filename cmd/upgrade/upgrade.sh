#!/bin/bash
namespace=kubesphere-logging-system
FluentbitOperator="fluentbit-operator"

function error_exit {
  echo "$1" 1>&2
  exit 1
}

function migrate(){
## Converting an existing configuration to a new one
local oldKind=$1
local newKind=$2
local list=$(kubectl get  $oldKind.logging.kubesphere.io -A -o json) || error_exit "Cannot get resource $oldKind"
local name=($(echo $list | jq -r '.items[].metadata.name | @json'))
local labels=($(echo $list | jq -r '.items[].metadata.labels | @json'))
local spec=($(echo $list | jq -r '.items[].spec | @json'))
local ns=($(echo $list | jq -r '.items[].metadata.namespace | @json'))
local size=${#spec[*]}
echo "Number of original $oldKind configuration files:$size"
for((i=0;i<${size};i++));do
if [[ "${kind}" = "fluentbits" ]]; then
cluster_resource_list[i]="{
\"apiVersion\": \"fluentbit.fluent.io/v1alpha2\",
\"kind\": \"${newKind}\",
\"metadata\": {
\"name\": ${name[i]},
\"labels\": ${labels[i]},
\"namespace\": \"${ns}\"
},
\"spec\": ${spec[i]}
}"
else
cluster_resource_list[i]="{
\"apiVersion\": \"fluentbit.fluent.io/v1alpha2\",
\"kind\": \"${kind}\",
\"metadata\": {
\"name\": ${name[i]},
\"labels\": ${labels[i]}
},
\"spec\": ${spec[i]}
}"
fi
done

## Uninstall the fluentbit-operator and the original configuration
for((i=0;i<${size};i++));do
echo "${name[i]}"
  temp=$(echo ${name[i]} | sed 's/"//g')
  echo "$temp"
   kubectl delete  $oldKind.logging.kubesphere.io $temp -n ${namespace}
done

for((i=0;i<${size};i++));do
echo ${cluster_resource_list[i]} | kubectl apply -f - || error_exit "Cannot apply resource $oldKind"
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
