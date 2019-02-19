# fluentbit-operator

FluentBit operator for Kubernetes based on Fluent-bit.

## What is this operator for?

This operator helps you to pack together logging information with your applications. With the help of Custom Resource Definition you can describe the behaviour of your application within its charts. The operator does the rest.


## Project status: Alpha

This is the first version of this operator to showcase our plans of handling logs on Kubernetes. This version includes only basic configuration that will expand quickly. Stay tuned.


## Installing the operator

```
helm repo add kubesphere http://kubernetes-charts.kubesphere.io/branch/master
helm install kubesphere/fluentbit-operator
```

## Example

The following steps set up an example configuration for sending nginx logs to ElasticSearch.

### Create FluentBit resource

Create a manifest that defines that you want to parse the kubernetes logs on the standard output of pods, and store them in the given ElasticSearch.

```
apiVersion: "logging.kubesphere.io/v1alpha1"
kind: "FluentBit"
metadata:
  name: "fluent-bit"
  namespace: "kubesphere-logging-system"
spec:
  service:
    - type: fluentbit_service
      name: fluentbit-service
      parameters:
        - name: Flush
          value: "1"
        - name: Daemon
          value: "Off"
        - name: Log_Level
          value: "info"
        - name: Parsers_File
          value: "parsers.conf"
  input:
    - type: fluentbit_input
      name: fluentbit-input
      parameters:
        - name: Name
          value: "tail"
        - name: Path
          value: "/var/log/containers/*.log"
        - name: Parser
          value: "docker"
        - name: Tag
          value: "kube.*"
        - name: Refresh_Interval
          value: "5"
        - name: Mem_Buf_Limit
          value: "5MB"
        - name: Skip_Long_Lines
          value: "On"
        - name: DB
          value: "/tail-db/tail-containers-state.db"
        - name: DB.Sync
          value: "Normal"
  filter:
    - type: fluentbit_filter
      name: fluentbit-filter
      parameters:
        - name: Name
          value: "kubernetes"
        - name: Match
          value: "kube.*"
        - name: Kube_URL
          value: "https://kubernetes.default.svc:443"
        - name: Kube_CA_File
          value: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
        - name: Kube_Token_File
          value: "/var/run/secrets/kubernetes.io/serviceaccount/token"
  output:
    - type: fluentbit_output
      name: fluentbit-output
      parameters:
        - name: Name
          value: "es"
        - name: Match
          value: "kube.*"
        - name: Host
          value: "elasticsearch-logging-data.kubesphere-logging-system.svc"
        - name: Port
          value: "9200"
        - name: Logstash_Format
          value: "On"
        - name: Replace_Dots
          value: "on"
        - name: Retry_Limit
          value: "False"
        - name: Type
          value: "flb_type"
        - name: Time_Key
          value: "@timestamp"
        - name: Logstash_Prefix
          value: "logstash"
  settings:
    - type: fluentbit_settings
      name: fluentbit-settings
      parameters:
        - name: Enable
          value: "false"
```

## Contributing

If you find this project useful here's how you can help:

- Send a pull request with your new features and bug fixes
- Help new users with issues they may encounter
- Support the development of this project and star this repo!
