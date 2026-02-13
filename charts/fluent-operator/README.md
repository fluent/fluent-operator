# fluent-operator

![Version: 4.0.0](https://img.shields.io/badge/Version-4.0.0-informational?style=flat-square) ![AppVersion: 3.6.0](https://img.shields.io/badge/AppVersion-3.6.0-informational?style=flat-square)

## Overview

A Helm chart for Fluent Operator, a tool that provides flexibility in building a logging pipeline based on Fluent-Bit and Fluentd.

## Installation

### Method 1: Standard Installation

The simplest way to install Fluent Operator with CRDs included:

```bash
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update

helm install fluent-operator fluent/fluent-operator
```

**Behavior:**
- CRDs automatically installed from the `crds/` directory
- Helm does NOT upgrade CRDs on `helm upgrade` (manual upgrade required)
- Helm does NOT delete CRDs on `helm uninstall`

**Upgrading CRDs:**

When upgrading, manually apply CRD updates before upgrading the chart. You can obtain the CRDs from the chart package or the repository:

```bash
# Option 1: Extract from the chart
helm pull fluent/fluent-operator --untar
kubectl apply -f fluent-operator/crds/

# Option 2: Clone the repository
git clone https://github.com/fluent/fluent-operator.git
cd fluent-operator
kubectl apply -f charts/fluent-operator/crds/

# Then upgrade the chart
helm upgrade fluent-operator fluent/fluent-operator
```

**Skipping CRDs:**

If you manage CRDs separately:

```bash
helm install fluent-operator fluent/fluent-operator --skip-crds
```

### Method 2: Helm-Managed CRDs (Advanced)

For full Helm lifecycle management of CRDs (automatic upgrades and deletions):

```bash
# Step 1: Install CRDs with Helm management
helm install fluent-operator-crds fluent/fluent-operator-crds

# Step 2: Install operator (skip CRDs since already installed)
helm install fluent-operator fluent/fluent-operator --skip-crds
```

**Behavior:**
- CRDs automatically upgrade with `helm upgrade fluent-operator-crds`
- Fine-grained control (enable/disable Fluent Bit or Fluentd CRDs)
- CRDs deleted on `helm uninstall` (unless protected with annotation)

**Protecting CRDs:**

```bash
helm install fluent-operator-crds fluent/fluent-operator-crds \
  --set additionalAnnotations."helm\.sh/resource-policy"=keep
```

See [fluent-operator-crds chart](https://github.com/fluent/fluent-operator/tree/master/charts/fluent-operator-crds) for more details.

## Upgrading

See [MIGRATION-v4.md](MIGRATION-v4.md) for detailed upgrade instructions from v3.x to v4.0.

### Upgrading from v3.x to v4.0

**Major Changes:**
- CRDs now in `crds/` directory (Helm v3 standard)
- New `fluent-operator-crds` chart available for Helm-managed CRDs
- Default container runtime changed to `containerd`
- Removed dependency on legacy CRD sub-charts

**Upgrade Steps:**

```bash
# Update repository
helm repo update

# Manually update CRDs first (Helm doesn't upgrade CRDs in crds/ directory)
helm pull fluent/fluent-operator --version 4.0.0 --untar
kubectl apply -f fluent-operator/crds/

# Then upgrade the chart
helm upgrade fluent-operator fluent/fluent-operator --version 4.0.0
```

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| wenchajun | <dehaocheng@kubesphere.io> |  |
| marcofranssen | <marco.franssen@gmail.com> | <https://marcofranssen.nl> |
| joshuabaird | <joshbaird@gmail.com> |  |

## Source Code

* <https://github.com/fluent/fluent-operator>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| Kubernetes | bool | `true` | Deploy default Fluent Bit pipeline to collect Kubernetes logs. See https://github.com/fluent/fluent-operator/tree/master/manifests/logging-stack |
| containerRuntime | string | `"containerd"` | Container runtime used by your Kubernetes cluster. Supported values: containerd, crio, docker |
| fluentbit | object | `{"additionalVolumes":[],"additionalVolumesMounts":[],"affinity":{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"node-role.kubernetes.io/edge","operator":"DoesNotExist"}]}]}}},"annotations":{},"args":[],"command":[],"disableLogVolumes":false,"enable":true,"envVars":[],"filter":{"containerd":{"enable":true},"kubernetes":{"annotations":false,"enable":true,"labels":false},"multiline":{"buffer":false,"emitterMemBufLimit":120,"emitterType":"memory","enable":false,"flushMs":2000,"keyContent":"log","parsers":["go","python","java"]},"systemd":{"enable":true}},"hostNetwork":false,"image":{"registry":"ghcr.io","repository":"fluent/fluent-operator/fluent-bit","tag":"4.2.2"},"imagePullSecrets":[],"initContainers":[],"input":{"fluentBitMetrics":{},"nodeExporterMetrics":{},"systemd":{"enable":true,"includeKubelet":true,"path":"/var/log/journal","pauseOnChunksOverlimit":"off","storageType":"memory","stripUnderscores":"off","systemdFilter":{"enable":true,"filters":[]}},"tail":{"bufferChunkSize":"","bufferMaxSize":"","enable":true,"memBufLimit":"100MB","path":"/var/log/containers/*.log","pauseOnChunksOverlimit":"off","readFromHead":false,"refreshIntervalSeconds":10,"skipEmptyLines":true,"skipLongLines":true,"storageType":"memory"}},"kubeedge":{"enable":false,"prometheusRemoteWrite":{"host":"<cloud-prometheus-service-host>","port":"<cloud-prometheus-service-port>"}},"labels":{},"livenessProbe":{"enabled":true,"failureThreshold":8,"httpGet":{"path":"/","port":2020},"initialDelaySeconds":10,"periodSeconds":10,"successThreshold":1,"timeoutSeconds":15},"logLevel":"","namespaceClusterFbCfg":"","namespaceFluentBitCfgSelector":{},"nodeSelector":{},"output":{"es":{"bufferSize":"20MB","enable":false,"host":"<Elasticsearch url like elasticsearch-logging-data.kubesphere-logging-system.svc>","logstashPrefix":"ks-logstash-log","port":9200,"traceError":true},"kafka":{"brokers":"<kafka broker list like xxx.xxx.xxx.xxx:9092,yyy.yyy.yyy.yyy:9092>","enable":false,"logLevel":"info","topics":"ks-log"},"loki":{"enable":false,"host":"127.0.0.1","httpPassword":"mypass","httpUser":"myuser","logLevel":"info","port":3100,"retryLimit":"no_limits","tenantID":""},"opensearch":{},"opentelemetry":{},"prometheusMetricsExporter":{},"stackdriver":{},"stdout":{"enable":false}},"parsers":{"javaMultiline":{"enable":false}},"podSecurityContext":{},"ports":[],"positionDB":{"hostPath":{"path":"/var/lib/fluent-bit/"}},"priorityClassName":"","rbacRules":{},"resources":{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"10m","memory":"25Mi"}},"schedulerName":"","secrets":[],"securityContext":{},"service":{"storage":{}},"serviceAccountAnnotations":{},"serviceMonitor":{"enable":false,"interval":"30s","metricRelabelings":[],"path":"/api/v2/metrics/prometheus","relabelings":[],"scrapeTimeout":"10s","secure":false,"tlsConfig":{}},"tolerations":[{"operator":"Exists"}]}` | Fluent Bit configuration |
| fluentbit.additionalVolumes | list | `[]` | Additional volumes that can be mounted by containers belonging to the pod |
| fluentbit.additionalVolumesMounts | list | `[]` | Additional volume mounts to mount into the container's filesystem |
| fluentbit.affinity | object | `{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"node-role.kubernetes.io/edge","operator":"DoesNotExist"}]}]}}}` | Affinity configuration for Fluent Bit pods Ref: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity |
| fluentbit.annotations | object | `{}` | Annotations to add to each Fluent Bit pod Request to Fluent Bit to exclude or not the logs generated by the Pod: `fluentbit.io/exclude: "true"` Prometheus can use this tag to automatically discover the Pod and collect monitoring data: `prometheus.io/scrape: "true"` |
| fluentbit.args | list | `[]` | Custom command line arguments for Fluent Bit containers |
| fluentbit.command | list | `[]` | Custom command for Fluent Bit containers |
| fluentbit.disableLogVolumes | bool | `false` | Removes the hostPath mounts for varlibcontainers, varlogs and systemd |
| fluentbit.enable | bool | `true` | Enable Fluent Bit deployment |
| fluentbit.envVars | list | `[]` | Environment variables that can be passed to Fluent Bit pods |
| fluentbit.filter.containerd.enable | bool | `true` | Enable containerd log format converter filter |
| fluentbit.filter.kubernetes.annotations | bool | `false` | Include Kubernetes annotations in logs |
| fluentbit.filter.kubernetes.enable | bool | `true` | Enable Kubernetes metadata filter |
| fluentbit.filter.kubernetes.labels | bool | `false` | Include Kubernetes labels in logs |
| fluentbit.filter.multiline.buffer | bool | `false` | Buffer for multiline filter |
| fluentbit.filter.multiline.emitterMemBufLimit | int | `120` | Emitter memory buffer limit in MB |
| fluentbit.filter.multiline.emitterType | string | `"memory"` | Emitter type for multiline filter |
| fluentbit.filter.multiline.enable | bool | `false` | Enable multiline filter |
| fluentbit.filter.multiline.flushMs | int | `2000` | Flush interval in milliseconds |
| fluentbit.filter.multiline.keyContent | string | `"log"` | Key content field for multiline filter |
| fluentbit.filter.multiline.parsers | list | `["go","python","java"]` | Multiline parsers to use |
| fluentbit.filter.systemd.enable | bool | `true` | Enable systemd filter |
| fluentbit.hostNetwork | bool | `false` | Use host network for Fluent Bit DaemonSet |
| fluentbit.image.registry | string | `"ghcr.io"` | Fluent Bit image registry |
| fluentbit.image.repository | string | `"fluent/fluent-operator/fluent-bit"` | Fluent Bit image repository |
| fluentbit.image.tag | string | `"4.2.2"` | Fluent Bit image tag |
| fluentbit.imagePullSecrets | list | `[]` | Image pull secrets for Fluent Bit Ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/ |
| fluentbit.initContainers | list | `[]` | Init containers for Fluent Bit pods Ref: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/ |
| fluentbit.input | object | `{"fluentBitMetrics":{},"nodeExporterMetrics":{},"systemd":{"enable":true,"includeKubelet":true,"path":"/var/log/journal","pauseOnChunksOverlimit":"off","storageType":"memory","stripUnderscores":"off","systemdFilter":{"enable":true,"filters":[]}},"tail":{"bufferChunkSize":"","bufferMaxSize":"","enable":true,"memBufLimit":"100MB","path":"/var/log/containers/*.log","pauseOnChunksOverlimit":"off","readFromHead":false,"refreshIntervalSeconds":10,"skipEmptyLines":true,"skipLongLines":true,"storageType":"memory"}}` | Configure the format of the config file to either classic or yaml. Default to classic when this value is not set configFileFormat: yaml Set a limit of memory that Tail plugin can use when appending data to the Engine. You can find more details here: https://docs.fluentbit.io/manual/pipeline/inputs/tail#config If the limit is reach, it will be paused; when the data is flushed it resumes. if the inbound traffic is less than 2.4Mbps, setting memBufLimit to 5MB is enough if the inbound traffic is less than 4.0Mbps, setting memBufLimit to 10MB is enough if the inbound traffic is less than 13.64Mbps, setting memBufLimit to 50MB is enough |
| fluentbit.input.fluentBitMetrics | object | `{}` | Fluent Bit metrics input configuration |
| fluentbit.input.nodeExporterMetrics | object | `{}` | Node exporter metrics input configuration |
| fluentbit.input.systemd.enable | bool | `true` | Enable systemd input |
| fluentbit.input.systemd.includeKubelet | bool | `true` | Include kubelet logs from systemd |
| fluentbit.input.systemd.path | string | `"/var/log/journal"` | Path to systemd journal |
| fluentbit.input.systemd.pauseOnChunksOverlimit | string | `"off"` | Pause input when chunks overlimit |
| fluentbit.input.systemd.storageType | string | `"memory"` | Storage type for systemd input buffering. Use "filesystem" for persistent buffering. |
| fluentbit.input.systemd.stripUnderscores | string | `"off"` | Strip underscores from systemd field names |
| fluentbit.input.systemd.systemdFilter.enable | bool | `true` | Enable systemd filter |
| fluentbit.input.systemd.systemdFilter.filters | list | `[]` | Systemd unit filters |
| fluentbit.input.tail.bufferChunkSize | string | `""` | Buffer chunk size for tail input |
| fluentbit.input.tail.bufferMaxSize | string | `""` | Buffer max size for tail input |
| fluentbit.input.tail.enable | bool | `true` | Enable tail input |
| fluentbit.input.tail.memBufLimit | string | `"100MB"` | Memory buffer limit for tail input |
| fluentbit.input.tail.path | string | `"/var/log/containers/*.log"` | Path to container logs |
| fluentbit.input.tail.pauseOnChunksOverlimit | string | `"off"` | Pause input when chunks overlimit |
| fluentbit.input.tail.readFromHead | bool | `false` | Read from head of log file |
| fluentbit.input.tail.refreshIntervalSeconds | int | `10` | Refresh interval for tail input in seconds |
| fluentbit.input.tail.skipEmptyLines | bool | `true` | Skip empty lines in logs |
| fluentbit.input.tail.skipLongLines | bool | `true` | Skip long lines in logs |
| fluentbit.input.tail.storageType | string | `"memory"` | Storage type for tail input buffering. Use "filesystem" for persistent buffering. |
| fluentbit.kubeedge.enable | bool | `false` | Enable KubeEdge integration for Fluent Bit |
| fluentbit.kubeedge.prometheusRemoteWrite.host | string | `"<cloud-prometheus-service-host>"` | Host of a cloud-side Prometheus-compatible server that can receive Prometheus remote write data |
| fluentbit.kubeedge.prometheusRemoteWrite.port | string | `"<cloud-prometheus-service-port>"` | Port of a cloud-side Prometheus-compatible server that can receive Prometheus remote write data |
| fluentbit.labels | object | `{}` | Additional custom labels for Fluent Bit pods |
| fluentbit.livenessProbe.enabled | bool | `true` | Enable liveness probe |
| fluentbit.livenessProbe.failureThreshold | int | `8` | Failure threshold for liveness probe |
| fluentbit.livenessProbe.initialDelaySeconds | int | `10` | Initial delay before liveness probe starts |
| fluentbit.livenessProbe.periodSeconds | int | `10` | Period between liveness probes |
| fluentbit.livenessProbe.successThreshold | int | `1` | Success threshold for liveness probe |
| fluentbit.livenessProbe.timeoutSeconds | int | `15` | Timeout for liveness probe |
| fluentbit.logLevel | string | `""` | Log level for Fluent Bit |
| fluentbit.namespaceClusterFbCfg | string | `""` | Using namespaceClusterFbCfg, deploy fluent-bit configmap and secret in this namespace. If it is not defined, it is in the namespace of the fluent-operator. |
| fluentbit.namespaceFluentBitCfgSelector | object | `{}` | Namespace selector for Fluent Bit configuration |
| fluentbit.nodeSelector | object | `{}` | Node selector for Fluent Bit pods Ref: https://kubernetes.io/docs/user-guide/node-selection/ |
| fluentbit.output.es.bufferSize | string | `"20MB"` | Buffer size for Elasticsearch output |
| fluentbit.output.es.enable | bool | `false` | Enable Elasticsearch output |
| fluentbit.output.es.host | string | `"<Elasticsearch url like elasticsearch-logging-data.kubesphere-logging-system.svc>"` | Elasticsearch host |
| fluentbit.output.es.logstashPrefix | string | `"ks-logstash-log"` | Logstash prefix for Elasticsearch indices |
| fluentbit.output.es.port | int | `9200` | Elasticsearch port |
| fluentbit.output.es.traceError | bool | `true` | Trace errors in Elasticsearch output |
| fluentbit.output.kafka.brokers | string | `"<kafka broker list like xxx.xxx.xxx.xxx:9092,yyy.yyy.yyy.yyy:9092>"` | Kafka broker list |
| fluentbit.output.kafka.enable | bool | `false` | Enable Kafka output |
| fluentbit.output.kafka.logLevel | string | `"info"` | Log level for Kafka output |
| fluentbit.output.kafka.topics | string | `"ks-log"` | Kafka topics |
| fluentbit.output.loki.enable | bool | `false` | Enable Loki output |
| fluentbit.output.loki.host | string | `"127.0.0.1"` | Loki host |
| fluentbit.output.loki.httpPassword | string | `"mypass"` | HTTP basic auth password for Loki |
| fluentbit.output.loki.httpUser | string | `"myuser"` | HTTP basic auth username for Loki |
| fluentbit.output.loki.logLevel | string | `"info"` | Log level for Loki output |
| fluentbit.output.loki.port | int | `3100` | Loki port |
| fluentbit.output.loki.retryLimit | string | `"no_limits"` | Retry limit for Loki output |
| fluentbit.output.loki.tenantID | string | `""` | Tenant ID for Loki |
| fluentbit.output.opensearch | object | `{}` | OpenSearch output configuration |
| fluentbit.output.opentelemetry | object | `{}` | OpenTelemetry output configuration |
| fluentbit.output.prometheusMetricsExporter | object | `{}` | Prometheus metrics exporter configuration Uncomment the following section to enable Prometheus metrics exporter. |
| fluentbit.output.stackdriver | object | `{}` | Stackdriver output configuration |
| fluentbit.output.stdout.enable | bool | `false` | Enable stdout output |
| fluentbit.parsers.javaMultiline.enable | bool | `false` | Enable Java multiline parser for generic springboot multiline log format |
| fluentbit.podSecurityContext | object | `{}` | Pod security context for Fluent Bit pods Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| fluentbit.ports | list | `[]` | Additional ports for Fluent Bit |
| fluentbit.positionDB | object | `{"hostPath":{"path":"/var/lib/fluent-bit/"}}` | Specify storage for position db. You will use it if tail input is enabled. Ref: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volumesource-v1-core |
| fluentbit.priorityClassName | string | `""` | Priority class for Fluent Bit pods Ref: https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#priorityclass |
| fluentbit.rbacRules | object | `{}` | Additional RBAC rules which will be applied to the Fluent Bit ClusterRole NOTE: As fluent-bit is managed by the fluent-operator, fluent-bit can only be granted permissions the operator also has Ref: https://kubernetes.io/docs/reference/access-authn-authz/rbac/#rolebinding-and-clusterrolebinding |
| fluentbit.resources | object | `{"limits":{"cpu":"500m","memory":"200Mi"},"requests":{"cpu":"10m","memory":"25Mi"}}` | Fluent Bit resource requests and limits. You can adjust it based on the log volume. If you do want to specify resources, adjust them as necessary |
| fluentbit.schedulerName | string | `""` | Scheduler name for Fluent Bit pods |
| fluentbit.secrets | list | `[]` | Secrets to mount in Fluent Bit pods |
| fluentbit.securityContext | object | `{}` | Security context for Fluent Bit container Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| fluentbit.service.storage | object | `{}` | Storage configuration for Fluent Bit buffering |
| fluentbit.serviceAccountAnnotations | object | `{}` | Additional custom annotations for Fluent Bit ServiceAccount |
| fluentbit.serviceMonitor.enable | bool | `false` | Enable Prometheus ServiceMonitor for Fluent Bit |
| fluentbit.serviceMonitor.interval | string | `"30s"` | Scrape interval |
| fluentbit.serviceMonitor.metricRelabelings | list | `[]` | Metric relabeling configs for ServiceMonitor |
| fluentbit.serviceMonitor.path | string | `"/api/v2/metrics/prometheus"` | Metrics path |
| fluentbit.serviceMonitor.relabelings | list | `[]` | Relabeling configs for ServiceMonitor |
| fluentbit.serviceMonitor.scrapeTimeout | string | `"10s"` | Scrape timeout |
| fluentbit.serviceMonitor.secure | bool | `false` | Use HTTPS for scraping |
| fluentbit.serviceMonitor.tlsConfig | object | `{}` | TLS configuration for ServiceMonitor |
| fluentbit.tolerations | list | `[{"operator":"Exists"}]` | Node tolerations for Fluent Bit pods Ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ |
| fluentd | object | `{"enable":false,"envVars":[],"extras":{},"forward":{"port":24224},"image":{"registry":"ghcr.io","repository":"fluent/fluent-operator/fluentd","tag":"v1.19.1"},"imagePullSecrets":[],"logLevel":"","mode":"collector","name":"fluentd","output":{"es":{"buffer":{"enable":false,"path":"/buffers/es","type":"file"},"enable":false,"host":"elasticsearch-logging-data.kubesphere-logging-system.svc","logstashPrefix":"ks-logstash-log","port":9200},"kafka":{"brokers":"my-cluster-kafka-bootstrap.default.svc:9091,my-cluster-kafka-bootstrap.default.svc:9092,my-cluster-kafka-bootstrap.default.svc:9093","buffer":{"enable":false,"path":"/buffers/kafka","type":"file"},"enable":false,"topicKey":"kubernetes_ns"},"opensearch":{}},"podSecurityContext":{},"port":24224,"priorityClassName":"","replicas":1,"resources":{"limits":{"cpu":"500m","memory":"500Mi"},"requests":{"cpu":"100m","memory":"128Mi"}},"schedulerName":"","securityContext":{},"watchedNamespaces":["kube-system","default"]}` | Fluentd configuration |
| fluentd.enable | bool | `false` | Enable Fluentd deployment |
| fluentd.envVars | list | `[]` | Environment variables that can be passed to Fluentd pods |
| fluentd.extras | object | `{}` | Extra configuration for Fluentd |
| fluentd.forward.port | int | `24224` | Forward input port |
| fluentd.image.registry | string | `"ghcr.io"` | Fluentd image registry |
| fluentd.image.repository | string | `"fluent/fluent-operator/fluentd"` | Fluentd image repository |
| fluentd.image.tag | string | `"v1.19.1"` | Fluentd image tag |
| fluentd.imagePullSecrets | list | `[]` | Image pull secrets for Fluentd |
| fluentd.logLevel | string | `""` | Log level for Fluentd |
| fluentd.mode | string | `"collector"` | Fluentd deployment mode. Valid values: "collector" (StatefulSet) or "agent" (DaemonSet) |
| fluentd.name | string | `"fluentd"` | Fluentd name |
| fluentd.output.es.buffer.enable | bool | `false` | Enable buffer for Elasticsearch output |
| fluentd.output.es.buffer.path | string | `"/buffers/es"` | Buffer path |
| fluentd.output.es.buffer.type | string | `"file"` | Buffer type |
| fluentd.output.es.enable | bool | `false` | Enable Elasticsearch output for Fluentd |
| fluentd.output.es.host | string | `"elasticsearch-logging-data.kubesphere-logging-system.svc"` | Elasticsearch host |
| fluentd.output.es.logstashPrefix | string | `"ks-logstash-log"` | Logstash prefix for Elasticsearch indices |
| fluentd.output.es.port | int | `9200` | Elasticsearch port |
| fluentd.output.kafka.brokers | string | `"my-cluster-kafka-bootstrap.default.svc:9091,my-cluster-kafka-bootstrap.default.svc:9092,my-cluster-kafka-bootstrap.default.svc:9093"` | Kafka broker list |
| fluentd.output.kafka.buffer.enable | bool | `false` | Enable buffer for Kafka output |
| fluentd.output.kafka.buffer.path | string | `"/buffers/kafka"` | Buffer path |
| fluentd.output.kafka.buffer.type | string | `"file"` | Buffer type |
| fluentd.output.kafka.enable | bool | `false` | Enable Kafka output for Fluentd |
| fluentd.output.kafka.topicKey | string | `"kubernetes_ns"` | Kafka topic key |
| fluentd.output.opensearch | object | `{}` | OpenSearch output configuration for Fluentd |
| fluentd.podSecurityContext | object | `{}` | Pod security context for Fluentd pod Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| fluentd.port | int | `24224` | Fluentd port |
| fluentd.priorityClassName | string | `""` | Priority class name for Fluentd pods |
| fluentd.replicas | int | `1` | Number of Fluentd replicas Applicable when the mode is "collector", and will be ignored when the mode is "agent" |
| fluentd.resources | object | `{"limits":{"cpu":"500m","memory":"500Mi"},"requests":{"cpu":"100m","memory":"128Mi"}}` | Fluentd resource requests and limits |
| fluentd.schedulerName | string | `""` | Scheduler name for Fluentd pods |
| fluentd.securityContext | object | `{}` | Container security context for Fluentd container Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| fluentd.watchedNamespaces | list | `["kube-system","default"]` | Namespaces to watch for Fluentd resources |
| fullnameOverride | string | `""` | Override full name of the chart |
| nameOverride | string | `""` | Override name of the chart |
| namespaceOverride | string | `""` | Override namespace where resources are deployed |
| operator | object | `{"affinity":{},"annotations":{},"disableComponentControllers":"","enable":true,"extraArgs":[],"image":{"registry":"ghcr.io","repository":"fluent/fluent-operator/fluent-operator","tag":""},"imagePullSecrets":[],"labels":{},"nodeSelector":{},"podSecurityContext":{},"priorityClassName":"","rbac":{"additionalRules":[],"clusterRole":{"name":"fluent-operator"},"clusterRoleBinding":{"name":"fluent-operator"},"create":true},"resources":{"limits":{"cpu":"100m","memory":"60Mi"},"requests":{"cpu":"100m","memory":"20Mi"}},"securityContext":{},"service":{"annotations":{},"enable":true,"labels":{},"port":8080,"portName":"metrics","type":"ClusterIP"},"serviceAccount":{"name":"fluent-operator"},"serviceMonitor":{"enable":false,"interval":"30s","metricRelabelings":[],"path":"/metrics","relabelings":[],"scrapeTimeout":"10s","secure":false,"tlsConfig":{}},"tolerations":[]}` | Fluent Operator configuration |
| operator.affinity | object | `{}` | Affinity for Fluent Operator pods Ref: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity |
| operator.annotations | object | `{}` | Annotations to add to Fluent Operator pods |
| operator.disableComponentControllers | string | `""` | Disable specific component controllers. Value can be "fluent-bit" or "fluentd" to disable that controller |
| operator.enable | bool | `true` | Enable Fluent Operator deployment. Set to false to disable creation of ClusterRole, ClusterRoleBinding, Deployment, and ServiceAccount |
| operator.extraArgs | list | `[]` | Extra arguments for the Fluent Operator controller |
| operator.image.registry | string | `"ghcr.io"` | Fluent Operator image registry |
| operator.image.repository | string | `"fluent/fluent-operator/fluent-operator"` | Fluent Operator image repository |
| operator.image.tag | string | Chart appVersion | Fluent Operator image tag (immutable tags are recommended). Overrides the image tag whose default is the chart appVersion |
| operator.imagePullSecrets | list | `[]` | Image pull secrets for Fluent Operator Ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/ |
| operator.labels | object | `{}` | Labels to add to Fluent Operator pods |
| operator.nodeSelector | object | `{}` | Node selector for Fluent Operator pods Ref: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector |
| operator.podSecurityContext | object | `{}` | Pod security context for Fluent Operator Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| operator.priorityClassName | string | `""` | Priority class name for Fluent Operator pods Ref: https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#priorityclass |
| operator.rbac.additionalRules | list | `[]` | Additional RBAC rules for Fluent Operator ClusterRole. Operator cannot give permissions it does not have |
| operator.rbac.clusterRole.name | string | `"fluent-operator"` | ClusterRole name |
| operator.rbac.clusterRoleBinding.name | string | `"fluent-operator"` | ClusterRoleBinding name |
| operator.rbac.create | bool | `true` | Specifies whether to create the ClusterRole and ClusterRoleBinding |
| operator.resources | object | `{"limits":{"cpu":"100m","memory":"60Mi"},"requests":{"cpu":"100m","memory":"20Mi"}}` | Fluent Operator resource requests and limits |
| operator.securityContext | object | `{}` | Container security context for Fluent Operator container Ref: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| operator.service.annotations | object | `{}` | Annotations for Fluent Operator service |
| operator.service.enable | bool | `true` | Enable Fluent Operator service |
| operator.service.labels | object | `{}` | Labels for Fluent Operator service |
| operator.service.port | int | `8080` | Service port |
| operator.service.portName | string | `"metrics"` | Service port name |
| operator.service.type | string | `"ClusterIP"` | Service type |
| operator.serviceAccount.name | string | `"fluent-operator"` | ServiceAccount name for Fluent Operator |
| operator.serviceMonitor.enable | bool | `false` | Enable Prometheus ServiceMonitor for Fluent Operator |
| operator.serviceMonitor.interval | string | `"30s"` | Scrape interval |
| operator.serviceMonitor.metricRelabelings | list | `[]` | Metric relabeling configs for ServiceMonitor |
| operator.serviceMonitor.path | string | `"/metrics"` | Metrics path |
| operator.serviceMonitor.relabelings | list | `[]` | Relabeling configs for ServiceMonitor |
| operator.serviceMonitor.scrapeTimeout | string | `"10s"` | Scrape timeout |
| operator.serviceMonitor.secure | bool | `false` | Use HTTPS for scraping |
| operator.serviceMonitor.tlsConfig | object | `{}` | TLS configuration for ServiceMonitor |
| operator.tolerations | list | `[]` | Tolerations for Fluent Operator pods Ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/ |

## Contributing

Contributions are welcome! Please read our [Contributing Guide](https://github.com/fluent/fluent-operator/blob/master/CONTRIBUTING.md).

## License

Apache License 2.0, see [LICENSE](https://github.com/fluent/fluent-operator/blob/master/LICENSE).

----------------------------------------------

Autogenerated from chart metadata using [helm-docs](https://github.com/norwoodj/helm-docs/).
