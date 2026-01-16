# fluent-operator

![Version: 4.0.0](https://img.shields.io/badge/Version-4.0.0-informational?style=flat-square) ![AppVersion: 3.6.0](https://img.shields.io/badge/AppVersion-3.6.0-informational?style=flat-square)

Fluent Operator provides great flexibility in building a logging layer based on Fluent Bit and Fluentd.

**Homepage:** <https://www.fluentd.org/>

## v4.0 Release

**Please see [MIGRATION-v4](./MIGRATION-v4.md) for important information regarding the upgrade to v4.0!**

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| wenchajun | <dehaocheng@kubesphere.io> |  |
| marcofranssen | <marco.franssen@gmail.com> | <https://marcofranssen.nl> |
| joshuabaird | <joshbaird@gmail.com> |  |

## Source Code

* <https://github.com/fluent/fluent-operator>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
|  | fluentd-crds | 0.2.1 |
| https://fluent.github.io/helm-charts | fluent-bit-crds | 0.2.3 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| Kubernetes | bool | `true` |  |
| containerRuntime | string | `"containerd"` |  |
| fluentbit.additionalVolumes | list | `[]` |  |
| fluentbit.additionalVolumesMounts | list | `[]` |  |
| fluentbit.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0].key | string | `"node-role.kubernetes.io/edge"` |  |
| fluentbit.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0].operator | string | `"DoesNotExist"` |  |
| fluentbit.annotations | object | `{}` |  |
| fluentbit.args | list | `[]` |  |
| fluentbit.command | list | `[]` |  |
| fluentbit.crdsEnable | bool | `true` |  |
| fluentbit.disableLogVolumes | bool | `false` |  |
| fluentbit.enable | bool | `true` |  |
| fluentbit.envVars | list | `[]` |  |
| fluentbit.filter.containerd.enable | bool | `true` |  |
| fluentbit.filter.kubernetes.annotations | bool | `false` |  |
| fluentbit.filter.kubernetes.enable | bool | `true` |  |
| fluentbit.filter.kubernetes.labels | bool | `false` |  |
| fluentbit.filter.multiline.buffer | bool | `false` |  |
| fluentbit.filter.multiline.emitterMemBufLimit | int | `120` |  |
| fluentbit.filter.multiline.emitterType | string | `"memory"` |  |
| fluentbit.filter.multiline.enable | bool | `false` |  |
| fluentbit.filter.multiline.flushMs | int | `2000` |  |
| fluentbit.filter.multiline.keyContent | string | `"log"` |  |
| fluentbit.filter.multiline.parsers[0] | string | `"go"` |  |
| fluentbit.filter.multiline.parsers[1] | string | `"python"` |  |
| fluentbit.filter.multiline.parsers[2] | string | `"java"` |  |
| fluentbit.filter.systemd.enable | bool | `true` |  |
| fluentbit.hostNetwork | bool | `false` |  |
| fluentbit.image.registry | string | `"ghcr.io"` |  |
| fluentbit.image.repository | string | `"fluent/fluent-operator/fluent-bit"` |  |
| fluentbit.image.tag | string | `"4.2.2"` |  |
| fluentbit.imagePullSecrets | list | `[]` |  |
| fluentbit.initContainers | list | `[]` |  |
| fluentbit.input.fluentBitMetrics | object | `{}` |  |
| fluentbit.input.nodeExporterMetrics | object | `{}` |  |
| fluentbit.input.systemd.enable | bool | `true` |  |
| fluentbit.input.systemd.includeKubelet | bool | `true` |  |
| fluentbit.input.systemd.path | string | `"/var/log/journal"` |  |
| fluentbit.input.systemd.pauseOnChunksOverlimit | string | `"off"` |  |
| fluentbit.input.systemd.storageType | string | `"memory"` |  |
| fluentbit.input.systemd.stripUnderscores | string | `"off"` |  |
| fluentbit.input.systemd.systemdFilter.enable | bool | `true` |  |
| fluentbit.input.systemd.systemdFilter.filters | list | `[]` |  |
| fluentbit.input.tail.bufferChunkSize | string | `""` |  |
| fluentbit.input.tail.bufferMaxSize | string | `""` |  |
| fluentbit.input.tail.enable | bool | `true` |  |
| fluentbit.input.tail.memBufLimit | string | `"100MB"` |  |
| fluentbit.input.tail.path | string | `"/var/log/containers/*.log"` |  |
| fluentbit.input.tail.pauseOnChunksOverlimit | string | `"off"` |  |
| fluentbit.input.tail.readFromHead | bool | `false` |  |
| fluentbit.input.tail.refreshIntervalSeconds | int | `10` |  |
| fluentbit.input.tail.skipEmptyLines | bool | `true` |  |
| fluentbit.input.tail.skipLongLines | bool | `true` |  |
| fluentbit.input.tail.storageType | string | `"memory"` |  |
| fluentbit.kubeedge.enable | bool | `false` |  |
| fluentbit.kubeedge.prometheusRemoteWrite.host | string | `"<cloud-prometheus-service-host>"` |  |
| fluentbit.kubeedge.prometheusRemoteWrite.port | string | `"<cloud-prometheus-service-port>"` |  |
| fluentbit.labels | object | `{}` |  |
| fluentbit.livenessProbe.enabled | bool | `true` |  |
| fluentbit.livenessProbe.failureThreshold | int | `8` |  |
| fluentbit.livenessProbe.httpGet.path | string | `"/"` |  |
| fluentbit.livenessProbe.httpGet.port | int | `2020` |  |
| fluentbit.livenessProbe.initialDelaySeconds | int | `10` |  |
| fluentbit.livenessProbe.periodSeconds | int | `10` |  |
| fluentbit.livenessProbe.successThreshold | int | `1` |  |
| fluentbit.livenessProbe.timeoutSeconds | int | `15` |  |
| fluentbit.logLevel | string | `""` |  |
| fluentbit.namespaceClusterFbCfg | string | `""` |  |
| fluentbit.namespaceFluentBitCfgSelector | object | `{}` |  |
| fluentbit.nodeSelector | object | `{}` |  |
| fluentbit.output.es.bufferSize | string | `"20MB"` |  |
| fluentbit.output.es.enable | bool | `false` |  |
| fluentbit.output.es.host | string | `"<Elasticsearch url like elasticsearch-logging-data.kubesphere-logging-system.svc>"` |  |
| fluentbit.output.es.logstashPrefix | string | `"ks-logstash-log"` |  |
| fluentbit.output.es.port | int | `9200` |  |
| fluentbit.output.es.traceError | bool | `true` |  |
| fluentbit.output.kafka.brokers | string | `"<kafka broker list like xxx.xxx.xxx.xxx:9092,yyy.yyy.yyy.yyy:9092>"` |  |
| fluentbit.output.kafka.enable | bool | `false` |  |
| fluentbit.output.kafka.logLevel | string | `"info"` |  |
| fluentbit.output.kafka.topics | string | `"ks-log"` |  |
| fluentbit.output.loki.enable | bool | `false` |  |
| fluentbit.output.loki.host | string | `"127.0.0.1"` |  |
| fluentbit.output.loki.httpPassword | string | `"mypass"` |  |
| fluentbit.output.loki.httpUser | string | `"myuser"` |  |
| fluentbit.output.loki.logLevel | string | `"info"` |  |
| fluentbit.output.loki.port | int | `3100` |  |
| fluentbit.output.loki.retryLimit | string | `"no_limits"` |  |
| fluentbit.output.loki.tenantID | string | `""` |  |
| fluentbit.output.opensearch | object | `{}` |  |
| fluentbit.output.opentelemetry | object | `{}` |  |
| fluentbit.output.prometheusMetricsExporter | object | `{}` |  |
| fluentbit.output.stackdriver | object | `{}` |  |
| fluentbit.output.stdout.enable | bool | `false` |  |
| fluentbit.parsers.javaMultiline.enable | bool | `false` |  |
| fluentbit.podSecurityContext | object | `{}` |  |
| fluentbit.ports | list | `[]` |  |
| fluentbit.positionDB.hostPath.path | string | `"/var/lib/fluent-bit/"` |  |
| fluentbit.priorityClassName | string | `""` |  |
| fluentbit.rbacRules | object | `{}` |  |
| fluentbit.resources.limits.cpu | string | `"500m"` |  |
| fluentbit.resources.limits.memory | string | `"200Mi"` |  |
| fluentbit.resources.requests.cpu | string | `"10m"` |  |
| fluentbit.resources.requests.memory | string | `"25Mi"` |  |
| fluentbit.schedulerName | string | `""` |  |
| fluentbit.secrets | list | `[]` |  |
| fluentbit.securityContext | object | `{}` |  |
| fluentbit.service.storage | object | `{}` |  |
| fluentbit.serviceAccountAnnotations | object | `{}` |  |
| fluentbit.serviceMonitor.enable | bool | `false` |  |
| fluentbit.serviceMonitor.interval | string | `"30s"` |  |
| fluentbit.serviceMonitor.metricRelabelings | list | `[]` |  |
| fluentbit.serviceMonitor.path | string | `"/api/v2/metrics/prometheus"` |  |
| fluentbit.serviceMonitor.relabelings | list | `[]` |  |
| fluentbit.serviceMonitor.scrapeTimeout | string | `"10s"` |  |
| fluentbit.serviceMonitor.secure | bool | `false` |  |
| fluentbit.serviceMonitor.tlsConfig | object | `{}` |  |
| fluentbit.tolerations[0].operator | string | `"Exists"` |  |
| fluentd.crdsEnable | bool | `true` |  |
| fluentd.enable | bool | `false` |  |
| fluentd.envVars | list | `[]` |  |
| fluentd.extras | object | `{}` |  |
| fluentd.forward.port | int | `24224` |  |
| fluentd.image.registry | string | `"ghcr.io"` |  |
| fluentd.image.repository | string | `"fluent/fluent-operator/fluentd"` |  |
| fluentd.image.tag | string | `"v1.19.1"` |  |
| fluentd.imagePullSecrets | list | `[]` |  |
| fluentd.logLevel | string | `""` |  |
| fluentd.mode | string | `"collector"` |  |
| fluentd.name | string | `"fluentd"` |  |
| fluentd.output.es.buffer.enable | bool | `false` |  |
| fluentd.output.es.buffer.path | string | `"/buffers/es"` |  |
| fluentd.output.es.buffer.type | string | `"file"` |  |
| fluentd.output.es.enable | bool | `false` |  |
| fluentd.output.es.host | string | `"elasticsearch-logging-data.kubesphere-logging-system.svc"` |  |
| fluentd.output.es.logstashPrefix | string | `"ks-logstash-log"` |  |
| fluentd.output.es.port | int | `9200` |  |
| fluentd.output.kafka.brokers | string | `"my-cluster-kafka-bootstrap.default.svc:9091,my-cluster-kafka-bootstrap.default.svc:9092,my-cluster-kafka-bootstrap.default.svc:9093"` |  |
| fluentd.output.kafka.buffer.enable | bool | `false` |  |
| fluentd.output.kafka.buffer.path | string | `"/buffers/kafka"` |  |
| fluentd.output.kafka.buffer.type | string | `"file"` |  |
| fluentd.output.kafka.enable | bool | `false` |  |
| fluentd.output.kafka.topicKey | string | `"kubernetes_ns"` |  |
| fluentd.output.opensearch | object | `{}` |  |
| fluentd.podSecurityContext | object | `{}` |  |
| fluentd.port | int | `24224` |  |
| fluentd.priorityClassName | string | `""` |  |
| fluentd.replicas | int | `1` |  |
| fluentd.resources.limits.cpu | string | `"500m"` |  |
| fluentd.resources.limits.memory | string | `"500Mi"` |  |
| fluentd.resources.requests.cpu | string | `"100m"` |  |
| fluentd.resources.requests.memory | string | `"128Mi"` |  |
| fluentd.schedulerName | string | `""` |  |
| fluentd.securityContext | object | `{}` |  |
| fluentd.watchedNamespaces[0] | string | `"kube-system"` |  |
| fluentd.watchedNamespaces[1] | string | `"default"` |  |
| fullnameOverride | string | `""` |  |
| nameOverride | string | `""` |  |
| namespaceOverride | string | `""` |  |
| operator.affinity | object | `{}` |  |
| operator.annotations | object | `{}` |  |
| operator.disableComponentControllers | string | `""` |  |
| operator.enable | bool | `true` |  |
| operator.extraArgs | list | `[]` |  |
| operator.image.registry | string | `"ghcr.io"` |  |
| operator.image.repository | string | `"fluent/fluent-operator/fluent-operator"` |  |
| operator.image.tag | string | `""` |  |
| operator.imagePullSecrets | list | `[]` |  |
| operator.labels | object | `{}` |  |
| operator.nodeSelector | object | `{}` |  |
| operator.podSecurityContext | object | `{}` |  |
| operator.priorityClassName | string | `""` |  |
| operator.rbac.additionalRules | list | `[]` |  |
| operator.rbac.clusterRole.name | string | `"fluent-operator"` |  |
| operator.rbac.clusterRoleBinding.name | string | `"fluent-operator"` |  |
| operator.rbac.create | bool | `true` | Specifies whether to create the ClusterRole and ClusterRoleBinding. |
| operator.resources.limits.cpu | string | `"100m"` |  |
| operator.resources.limits.memory | string | `"60Mi"` |  |
| operator.resources.requests.cpu | string | `"100m"` |  |
| operator.resources.requests.memory | string | `"20Mi"` |  |
| operator.securityContext | object | `{}` |  |
| operator.service.annotations | object | `{}` |  |
| operator.service.enable | bool | `true` |  |
| operator.service.labels | object | `{}` |  |
| operator.service.port | int | `8080` |  |
| operator.service.portName | string | `"metrics"` |  |
| operator.service.type | string | `"ClusterIP"` |  |
| operator.serviceAccount.name | string | `"fluent-operator"` |  |
| operator.serviceMonitor.enable | bool | `false` |  |
| operator.serviceMonitor.interval | string | `"30s"` |  |
| operator.serviceMonitor.metricRelabelings | list | `[]` |  |
| operator.serviceMonitor.path | string | `"/metrics"` |  |
| operator.serviceMonitor.relabelings | list | `[]` |  |
| operator.serviceMonitor.scrapeTimeout | string | `"10s"` |  |
| operator.serviceMonitor.secure | bool | `false` |  |
| operator.serviceMonitor.tlsConfig | object | `{}` |  |
| operator.tolerations | list | `[]` |  |

