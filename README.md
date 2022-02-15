# Fluent Bit Operator

Fluent Bit Operator facilitates the deployment of Fluent Bit and provides great flexibility in building a logging layer based on Fluent Bit.

Once installed, the Fluent Bit Operator provides the following features:

- **Fluent Bit Management**: Deploy and destroy Fluent Bit DaemonSet automatically.
- **Custom Configuration**: Select input/filter/output plugins via labels.
- **Dynamic Reloading**: Update configuration without rebooting Fluent Bit pods.

## Table of contents

- [Fluent Bit Operator](#fluent-bit-operator)
  - [Table of contents](#table-of-contents)
  - [Overview](#overview)
  - [Get Started](#get-started)
    - [Prerequisites](#prerequisites)
    - [Install](#install)
        - [Deploy Fluent Bit Operator with YAML](#deploy-fluent-bit-operator-with-yaml)
        - [Deploy Fluent Bit Operator with Helm](#deploy-fluent-bit-operator-with-helm)
    - [Quick Start](#quick-start)
    - [Configure Custom Watch Namespaces](#configure-custom-watch-namespaces)
    - [Collect Kubernetes logs](#collect-kubernetes-logs)
        - [Deploy the Kubernetes logging stack with YAML](#deploy-the-kubernetes-logging-stack-with-yaml)
        - [Deploy the Kubernetes logging stack with Helm](#deploy-the-kubernetes-logging-stack-with-helm)
    - [Collect auditd logs](#collect-auditd-logs)
    - [Fluentd](#fluentd)
      - [Collect logs from Fluentbit](#collect-logs-from-fluentbit)
        - [Enable Fluentbit forward plugin](#enable-fluentbit-forward-plugin)
        - [Fluentd ClusterFluentdConfig example](#fluentd-clusterfluentdconfig-example)
        - [Fluentd FluentdConfig example](#fluentd-fluentdconfig-example)
        - [Fluentd mixed FluentdConfig and ClusterFluentdConfig example](#fluentd-mixed-fluentdconfig-and-clusterfluentdconfig-example)
        - [Fluentd ClusterOutput and Output example in multi-tenant scenario](#fluentd-clusteroutput-and-output-example-in-multi-tenant-scenario)
        - [Fluentd outputs to kafka/es example](#fluentd-outputs-to-kafkaes-example)
        - [Fluentd buffer example](#fluentd-buffer-example)
      - [Collect logs over HTTP](#collect-logs-over-http)
  - [Monitoring](#monitoring)
  - [API Doc](#api-doc)
    - [Fluent Bit](#fluent-bit)
    - [Fluentd](#fluentd-1)
  - [Best Practice](#best-practice)
    - [Plugin Grouping](#plugin-grouping)
    - [Path Convention](#path-convention)
  - [Custom Parser](#custom-parser)
  - [Roadmap](#roadmap)
  - [Development](#development)
    - [Requirements](#requirements)
    - [Running](#running)
  - [Contributing](#contributing)
    - [Documentation](#documentation)
    - [Manifests](#manifests)

## Overview

Fluent Bit Operator defines five custom resources using CustomResourceDefinition (CRD):

- **`FluentBit`**: Defines the Fluent Bit DaemonSet and its configs. A custom Fluent Bit image `kubesphere/fluent-bit` is required to work with FluentBit Operator for dynamic configuration reloading.
- **`FluentBitConfig`**: Select input/filter/output plugins and generates the final config into a Secret.
- **`Input`**: Defines input config sections.
- **`Parser`**: Defines parser config sections.
- **`Filter`**: Defines filter config sections.
- **`Output`**: Defines output config sections.

Each **`Input`**, **`Parser`**, **`Filter`**, **`Output`** represents a Fluent Bit config section, which are selected by **`FluentBitConfig`** via label selectors. The operator watches those objects, constructs the final config, and finally creates a Secret to store the config. This secret will be mounted into the Fluent Bit DaemonSet. The entire workflow looks like below:

![Fluent Bit workflow](docs/images/fluent-bit-operator-workflow.svg)

To enable fluent-bit to pick up and use the latest config whenever the fluent-bit config changes, a wrapper called fluent-bit watcher is added to restart the fluent-bit process as soon as fluent-bit config changes are detected. This way the fluent-bit pod needn't be restarted to reload the new config. The fluent-bit config is reloaded in this way because there is no reload interface in fluent-bit itself. Please refer to this [known issue](https://github.com/fluent/fluent-bit/issues/365) for more details.

![fluentbit-operator](docs/images/fluentbit-operator.svg)

Besides, we have finished [the pr #189](https://github.com/fluent/fluentbit-operator/pull/189) to integrate fluent-operator as a forward log layer, aims to collect logs from fluentbit or other apps. The whole workflow could be described as below.

![Fluent-operator](docs/images/fluent-operator.svg)

## Get Started

### Prerequisites

Kubernetes v1.16.13+ is necessary for running Fluent Bit Operator.

### Install

##### Deploy Fluent Bit Operator with YAML

Install the latest stable version

```shell
kubectl apply -f https://raw.githubusercontent.com/kubesphere/fluentbit-operator/release-0.12/manifests/setup/setup.yaml

# You can change the namespace in manifests/setup/kustomization.yaml in corresponding release branch 
# and then use command below to install to another namespace
# kubectl kustomize manifests/setup/ | kubectl apply -f -
```

Install the development version

```shell
kubectl apply -f https://raw.githubusercontent.com/kubesphere/fluentbit-operator/master/manifests/setup/setup.yaml

# You can change the namespace in manifests/setup/kustomization.yaml 
# and then use command below to install to another namespace
# kubectl kustomize manifests/setup/ | kubectl apply -f -
```

##### Deploy Fluent Bit Operator with Helm

> Note: For the Helm-based installation you need Helm v3.2.1 or later.

Fluent Bit Operator supports `docker` as well as `containerd` and `CRI-O`. `containerd` and `CRI-O` use the `CRI Log` format which is slightly different and requires additional parsing to parse JSON application logs. You should set different `containerRuntime` depending on  your container runtime.

If your container runtime is `docker`

```shell
helm install fluentbit-operator  --create-namespace -n kubesphere-logging-system charts/fluentbit-operator/ --set containerRuntime=docker
```

If your container runtime is `containerd`

```shell
helm install fluentbit-operator --create-namespace -n kubesphere-logging-system charts/fluentbit-operator/  --set containerRuntime=containerd
```

If your container runtime is `cri-o`

```shell
helm install fluentbit-operator --create-namespace -n kubesphere-logging-system charts/fluentbit-operator/  --set containerRuntime=crio
```

### Quick Start

The quick start instructs you to deploy fluent bit with `dummy` as input and `stdout` as output, which is equivalent to execute the binary with `fluent-bit -i dummy -o stdout`.

```shell
kubectl apply -f https://raw.githubusercontent.com/kubesphere/fluentbit-operator/master/manifests/quick-start/quick-start.yaml
```

Once everything is up, you'll observe messages in fluent bit pod logs like below:

```shell
[0] my_dummy: [1587991566.000091658, {"message"=>"dummy"}]
[1] my_dummy: [1587991567.000061572, {"message"=>"dummy"}]
[2] my_dummy: [1587991568.000056842, {"message"=>"dummy"}]
[3] my_dummy: [1587991569.000896217, {"message"=>"dummy"}]
[0] my_dummy: [1587991570.000172328, {"message"=>"dummy"}]
```

It means the FluentBit Operator works properly if you see the above messages, you can delete the quick start test by

```shell
kubectl delete -f https://raw.githubusercontent.com/kubesphere/fluentbit-operator/master/manifests/quick-start/quick-start.yaml
```

### Configure Custom Watch Namespaces

When starting the operator, you can pass an optional set of namespaces to watch for resources in with the `--watch-namespaces` flag. It takes a comma-separated list of namespaces to watch:

```
...
      spec:
        containers:
          image: kubesphere/fluentbit-operator
        - args:
          - --watch-namespaces=namespace1,namespace2
```

### Collect Kubernetes logs

This guide provisions a logging pipeline including the Fluent Bit DaemonSet and its log input/filter/output configurations to collect Kubernetes logs including container logs and kubelet logs.

![logging stack](docs/images/logging-stack.svg)

> Note that you need a running Elasticsearch v5+ cluster to receive log data before start. **Remember to adjust [output-elasticsearch.yaml](manifests/logging-stack/output-elasticsearch.yaml) to your own es setup**. Kafka and Fluentd outputs are optional and are turned off by default.

##### Deploy the Kubernetes logging stack with YAML

```shell
kubectl apply -f manifests/logging-stack

# You can change the namespace in manifests/logging-stack/kustomization.yaml 
# and then use command below to install to another namespace
# kubectl kustomize manifests/logging-stack/ | kubectl apply -f -
```

##### Deploy the Kubernetes logging stack with Helm

If your container runtime is `docker`

```shell
helm upgrade fluentbit-operator --create-namespace -n kubesphere-logging-system charts/fluentbit-operator/  --set Kubernetes=true,containerRuntime=docker
```

If your container runtime is `containerd`

```shell
helm upgrade fluentbit-operator --create-namespace -n kubesphere-logging-system charts/fluentbit-operator/  --set Kubernetes=true,containerRuntime=containerd
```

If your container runtime is `cri-o`

```shell
helm upgrade fluentbit-operator --create-namespace -n kubesphere-logging-system charts/fluentbit-operator/  --set Kubernetes=true,containerRuntime=crio
```

Within a couple of minutes, you should observe an index available:

```shell
$ curl localhost:9200/_cat/indices
green open ks-logstash-log-2020.04.26 uwQuoO90TwyigqYRW7MDYQ 1 1  99937 0  31.2mb  31.2mb
```

Success!

### Collect auditd logs

The Linux audit framework provides a CAPP-compliant (Controlled Access Protection Profile) auditing system that reliably collects information about any security-relevant (or non-security-relevant) event on a system. Refer to `manifests/logging-stack/auditd`, it supports a method for collecting audit logs from the Linux audit framework.

```shell
kubectl apply -f manifests/logging-stack/auditd

# You can change the namespace in manifests/logging-stack/auditd/kustomization.yaml 
# and then use command below to install to another namespace
# kubectl kustomize manifests/logging-stack/auditd/ | kubectl apply -f -
```

Within a couple of minutes, you should observe an index available:

```shell
$ curl localhost:9200/_cat/indices
green open ks-logstash-log-2021.04.06 QeI-k_LoQZ2h1z23F3XiHg  5 1 404879 0 298.4mb 149.2mb
```

### Fluentd 

Fluentd acts as a forward log layer to collect logs from fluentbit and other Apps. 

#### Collect logs from Fluentbit

##### Enable Fluentbit forward plugin

At first, we should enable the forward plugin in fluentbit to send logs to fluentd.

```shell
kubectl apply -f manifests/fluentd/fluentbit-output-forward.yaml
```

##### Fluentd ClusterFluentdConfig example

```shell
kubectl apply -f manifests/fluentd/fluentd-cluster-cfg-output-es.yaml
```

##### Fluentd FluentdConfig example

```shell
kubectl apply -f manifests/fluentd/fluentd-namespaced-cfg-output-es.yaml
```

##### Fluentd mixed FluentdConfig and ClusterFluentdConfig example

```shell
kubectl apply -f manifests/fluentd/fluentd-mixed-cfgs-output-es.yaml
```

##### Fluentd ClusterOutput and Output example in multi-tenant scenario

```shell
kubectl apply -f manifests/fluentd/fluentd-mixed-cfgs-multi-tenant-output.yaml
```

##### Fluentd outputs to kafka/es example

```shell
kubectl apply -f manifests/fluentd/fluentd-cluster-cfg-output-kafka.yaml
kubectl apply -f manifests/fluentd/fluentd-cluster-cfg-output-es.yaml
```

##### Fluentd buffer example

```shell
kubectl apply -f manifests/fluentd/fluentd-cluster-cfg-output-buffer-example.yaml
```

#### Collect logs over HTTP

```shell
kubectl apply -f manifests/quick-start/fluentd-http.yaml
```

## Monitoring

Fluent Bit comes with a built-in HTTP Server. According to the official [documentation](https://docs.fluentbit.io/manual/administration/monitoring) of fluentbit You can enable this by enabling the HTTP server from the fluent bit configuration file:

```conf
[SERVICE]
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_PORT    2020
```

When you use the fluent-operator, You can enable this from `FluentBitConfig` manifest. Example is below:

```yaml
apiVersion: fluentbit.fluent.io/v1alpha2
kind: FluentBitConfig
metadata:
  name: fluent-bit-config
  namespace: logging-system
spec:
  filterSelector:
    matchLabels:
      fluentbit.fluent.io/enabled: 'true'
  inputSelector:
    matchLabels:
      fluentbit.fluent.io/enabled: 'true'
  outputSelector:
    matchLabels:
      fluentbit.fluent.io/enabled: 'true'
  service:
    httpListen: 0.0.0.0
    httpPort: 2020
    httpServer: true
    parsersFile: parsers.conf

```

Once HTTP server is enabled, you should be able to get the information:

```bash
curl <podIP>:2020 | jq .

{
  "fluent-bit": {
    "version": "1.8.3",
    "edition": "Community",
    "flags": [
      "FLB_HAVE_PARSER",
      "FLB_HAVE_RECORD_ACCESSOR",
      "FLB_HAVE_STREAM_PROCESSOR",
      "FLB_HAVE_TLS",
      "FLB_HAVE_OPENSSL",
      "FLB_HAVE_AWS",
      "FLB_HAVE_SIGNV4",
      "FLB_HAVE_SQLDB",
      "FLB_HAVE_METRICS",
      "FLB_HAVE_HTTP_SERVER",
      "FLB_HAVE_SYSTEMD",
      "FLB_HAVE_FORK",
      "FLB_HAVE_TIMESPEC_GET",
      "FLB_HAVE_GMTOFF",
      "FLB_HAVE_UNIX_SOCKET",
      "FLB_HAVE_PROXY_GO",
      "FLB_HAVE_JEMALLOC",
      "FLB_HAVE_LIBBACKTRACE",
      "FLB_HAVE_REGEX",
      "FLB_HAVE_UTF8_ENCODER",
      "FLB_HAVE_LUAJIT",
      "FLB_HAVE_C_TLS",
      "FLB_HAVE_ACCEPT4",
      "FLB_HAVE_INOTIFY"
    ]
  }
}
```


## API Doc

### Fluent Bit

The list below shows supported plugins which are based on Fluent Bit v1.7.x+. For more information, please refer to the API docs of each plugin.

- [Input](docs/fluentbit.md#input)
  - [dummy](docs/plugins/fluentbit/input/dummy.md)
  - [tail](docs/plugins/fluentbit/input/tail.md)
  - [systemd](docs/plugins/fluentbit/input/systemd.md)
- [Parser](docs/fluentbit.md#parser)
  - [json](docs/plugins/fluentbit/parser/json.md)
  - [logfmt](docs/plugins/fluentbit/parser/logfmt.md)
  - [lstv](docs/plugins/fluentbit/parser/lstv.md)
  - [regex](docs/plugins/fluentbit/parser/regex.md)
- [Filter](docs/fluentbit.md#filter)
  - [kubernetes](docs/plugins/fluentbit/filter/kubernetes.md)
  - [modify](docs/plugins/fluentbit/filter/modify.md)
  - [nest](docs/plugins/fluentbit/filter/nest.md)
  - [parser](docs/plugins/fluentbit/filter/parser.md)
  - [grep](docs/plugins/fluentbit/filter/grep.md)
  - [record modifier](docs/plugins/fluentbit/filter/recordmodifier.md)
  - [lua](docs/plugins/fluentbit/filter/lua.md)
  - [throttle](docs/plugins/fluentbit/filter/throttle.md)
  - [aws](docs/plugins/fluentbit/filter/aws.md)
  - [multiline](docs/plugins/fluentbit/filter/multiline.md)
- [Output](docs/fluentbit.md#output)
  - [elasticsearch](docs/plugins/fluentbit/output/elasticsearch.md)
  - [file](docs/plugins/fluentbit/output/file.md)
  - [forward](docs/plugins/fluentbit/output/forward.md)
  - [http](docs/plugins/fluentbit/output/http.md)
  - [kafka](docs/plugins/fluentbit/output/kafka.md)
  - [null](docs/plugins/fluentbit/output/null.md)
  - [stdout](docs/plugins/fluentbit/output/stdout.md)
  - [tcp](docs/plugins/fluentbit/output/tcp.md)
  - [loki](docs/plugins/fluentbit/output/loki.md)
  - [syslog](docs/plugins/fluentbit/output/syslog.md)
  - [datadog](docs/plugins/fluentbit/output/datadog.md)

### Fluentd

The list below shows supported plugins which are based on Fluentd v1.14.4+. For more information, please refer to the API docs of each plugin.

- Common
  - [buffer](docs/plugins/fluentd/common/buffer.md)
  - [format](docs/plugins/fluentd/common/format.md)
  - [parse](docs/plugins/fluentd/common/parse.md)
  - [time](docs/plugins/fluentd/common/common.md#time)
  - [inject](docs/plugins/fluentd/common/common.md#inject)
  - [security](docs/plugins/fluentd/common/common.md#security)
  - [user](docs/plugins/fluentd/common/common.md#user)
  - [transport](docs/plugins/fluentd/common/common.md#transport)
  - [client](docs/plugins/fluentd/common/common.md#client)
  - [auth](docs/plugins/fluentd/common/common.md#auth)
  - [server](docs/plugins/fluentd/common/common.md#server)
  - [service_discovery](docs/plugins/fluentd/common/common.md#ServiceDiscovery)
- [Input](docs/fluentd/input/types.md)
  - [http](docs/plugins/fluentd/input/http.md)
  - [forward](docs/plugins/fluentd/input/forward.md)
- [Filter](docs/fluentd/filter/types.md)
  - [parser](docs/plugins/fluentd/filter/parser.md)
  - [grep](docs/plugins/fluentd/filter/grep.md)
  - [record modifier](docs/plugins/fluentd/filter/record_modifier.md)
  - [stdout](docs/plugins/fluentd/filter/stdout.md)
- [Output](docs/plugins/fluentd/output/types.md)
  - [elasticsearch](docs/plugins/fluentd/output/elasticsearch.md)
  - [forward](docs/plugins/fluentd/output/forward.md)
  - [http](docs/plugins/fluentd/output/http.md)
  - [kafka](docs/plugins/fluentd/output/kafka.md)
  - stdout


## Best Practice

### Plugin Grouping

Input, filter, and output plugins are connected by label selectors. For input and output plugins, always create `Input` or `Output` CRs for every plugin. Don't aggregate multiple inputs or outputs into one `Input` or `Output` object, except you have a good reason to do so. Take the demo `logging stack` for example, we have one yaml file for each output.

However, for filter plugins, if you want a filter chain, the order of filters matters. You need to organize multiple filters into an array as the demo [logging stack](manifests/logging-stack/filter-kubernetes.yaml) suggests.

### Path Convention

Path to file in Fluent Bit config should be well regulated. Fluent Bit Operator adopts the following convention internally.

|Dir Path|Description|
|---|---|
|/fluent-bit/tail|Stores tail related files, eg. file tracking db. Using [fluentbit.spec.positionDB](docs/fluentbit.md#fluentbitspec) will mount a file `pos.db` under this dir by default.|
|/fluent-bit/secrets/{secret_name}|Stores secrets, eg. TLS files. Specify secrets to mount in [fluentbit.spec.secrets](docs/fluentbit.md#fluentbitspec), then you have access.|
|/fluent-bit/config|Stores the main config file and user-defined parser config file.|

> Note that ServiceAccount files are mounted at `/var/run/secrets/kubernetes.io/serviceaccount`.

## Custom Parser

To enable parsers, you must set the value of `FluentBitConfig.Spec.Service.ParsersFile` to `parsers.conf`. Your custom parsers will be included into the built-in parser config via `@INCLUDE /fluent-bit/config/parsers.conf`. Note that the parsers.conf contains a few built-in parsers, for example, docker. Read [parsers.conf](https://github.com/kubesphere/fluentbit-operator/blob/master/conf/parsers.conf) for more information.

Check out the demo in the folder `/manifests/regex-parser` for how to use a custom regex parser.

## Roadmap

- [ ] Support containerd log format
- [ ] Add Fluentd CRDs as the log aggregation layer with group name `fluentd.fluent.io`
- [ ] Add FluentBit Cluster CRDs with new group name `fluentbit.fluent.io`
- [ ] Rename the entire project to Fluent Operator
- [ ] Support more Fluentd & FluentBit plugins

## Development

### Requirements
- golang v1.16+.requirement
- kubectl v1.16.13+.
- kubebuilder v2.3+ (the project is build with v2.3.2)
- Access to a Kubernetes cluster v1.16.13+

### Running

1. Install CRDs: `make install`
2. Run: `make run`

## Contributing

### Documentation

[API Doc](docs/fluentbit.md) is generated automatically. To modify it, edit the comment above struct fields, then run `go run cmd/doc-gen/main.go`.

### Manifests

Most files under the folder [manifests/setup](manifests/setup) are automatically generated from [config](config). Don't edit them directly, run `make manifests` instead, then replace these files accordingly.
