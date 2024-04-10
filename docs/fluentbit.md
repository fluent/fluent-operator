# API Docs
This Document documents the types introduced by the fluentbit Operator.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.
## Table of Contents
* [ClusterFilter](#clusterfilter)
* [ClusterFilterList](#clusterfilterlist)
* [ClusterFluentBitConfig](#clusterfluentbitconfig)
* [ClusterFluentBitConfigList](#clusterfluentbitconfiglist)
* [ClusterInput](#clusterinput)
* [ClusterInputList](#clusterinputlist)
* [ClusterMultilineParser](#clustermultilineparser)
* [ClusterMultilineParserList](#clustermultilineparserlist)
* [ClusterOutput](#clusteroutput)
* [ClusterOutputList](#clusteroutputlist)
* [ClusterParser](#clusterparser)
* [ClusterParserList](#clusterparserlist)
* [Collector](#collector)
* [CollectorList](#collectorlist)
* [CollectorService](#collectorservice)
* [CollectorSpec](#collectorspec)
* [Decorder](#decorder)
* [Filter](#filter)
* [FilterItem](#filteritem)
* [FilterList](#filterlist)
* [FilterSpec](#filterspec)
* [FluentBit](#fluentbit)
* [FluentBitConfig](#fluentbitconfig)
* [FluentBitConfigList](#fluentbitconfiglist)
* [FluentBitConfigSpec](#fluentbitconfigspec)
* [FluentBitList](#fluentbitlist)
* [FluentBitService](#fluentbitservice)
* [FluentBitSpec](#fluentbitspec)
* [InputSpec](#inputspec)
* [MultilineParser](#multilineparser)
* [MultilineParserList](#multilineparserlist)
* [NamespacedFluentBitCfgSpec](#namespacedfluentbitcfgspec)
* [Output](#output)
* [OutputList](#outputlist)
* [OutputSpec](#outputspec)
* [Parser](#parser)
* [ParserList](#parserlist)
* [ParserSpec](#parserspec)
* [Script](#script)
* [Service](#service)
* [Storage](#storage)
# ClusterFilter

ClusterFilter defines a cluster-level Filter configuration.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec | Specification of desired Filter configuration. | FilterSpec |

[Back to TOC](#table-of-contents)
# ClusterFilterList

ClusterFilterList contains a list of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterFilter |

[Back to TOC](#table-of-contents)
# ClusterFluentBitConfig

ClusterFluentBitConfig is the Schema for the cluster-level fluentbitconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FluentBitConfigSpec |

[Back to TOC](#table-of-contents)
# ClusterFluentBitConfigList

ClusterFluentBitConfigList contains a list of ClusterFluentBitConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterFluentBitConfig |

[Back to TOC](#table-of-contents)
# ClusterInput

ClusterInput is the Schema for the inputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | InputSpec |

[Back to TOC](#table-of-contents)
# ClusterInputList

ClusterInputList contains a list of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterInput |

[Back to TOC](#table-of-contents)
# ClusterMultilineParser

ClusterMultilineParser is the Schema for the cluster-level multiline parser API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | MultilineParserSpec |

[Back to TOC](#table-of-contents)
# ClusterMultilineParserList

ClusterMultilineParserList contains a list of ClusterMultilineParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterMultilineParser |

[Back to TOC](#table-of-contents)
# ClusterOutput

ClusterOutput is the Schema for the cluster-level outputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | OutputSpec |

[Back to TOC](#table-of-contents)
# ClusterOutputList

ClusterOutputList contains a list of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterOutput |

[Back to TOC](#table-of-contents)
# ClusterParser

ClusterParser is the Schema for the cluster-level parsers API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | ParserSpec |

[Back to TOC](#table-of-contents)
# ClusterParserList

ClusterParserList contains a list of ClusterParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterParser |

[Back to TOC](#table-of-contents)
# Collector

Collector is the Schema for the fluentbits API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | CollectorSpec |
| status |  | CollectorStatus |

[Back to TOC](#table-of-contents)
# CollectorList

CollectorList contains a list of Collector


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Collector |

[Back to TOC](#table-of-contents)
# CollectorService

CollectorService defines the service of the FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| name | Name is the name of the FluentBit service. | string |
| annotations | Annotations to add to each Fluentbit service. | map[string]string |
| labels | Labels to add to each FluentBit service | map[string]string |

[Back to TOC](#table-of-contents)
# CollectorSpec

CollectorSpec defines the desired state of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| image | Fluent Bit image. | string |
| args | Fluent Bit Watcher command line arguments. | []string |
| imagePullPolicy | Fluent Bit image pull policy. | corev1.PullPolicy |
| imagePullSecrets | Fluent Bit image pull secret | []corev1.LocalObjectReference |
| resources | Compute Resources required by container. | corev1.ResourceRequirements |
| nodeSelector | NodeSelector | map[string]string |
| affinity | Pod's scheduling constraints. | *corev1.Affinity |
| tolerations | Tolerations | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#toleration-v1-core) |
| fluentBitConfigName | Fluentbitconfig object associated with this Fluentbit | string |
| secrets | The Secrets are mounted into /fluent-bit/secrets/<secret-name>. | []string |
| runtimeClassName | RuntimeClassName represents the container runtime configuration. | string |
| priorityClassName | PriorityClassName represents the pod's priority class. | string |
| volumes | List of volumes that can be mounted by containers belonging to the pod. | []corev1.Volume |
| volumesMounts | Pod volumes to mount into the container's filesystem. | []corev1.VolumeMount |
| annotations | Annotations to add to each Fluentbit pod. | map[string]string |
| serviceAccountAnnotations | Annotations to add to the Fluentbit service account | map[string]string |
| securityContext | SecurityContext holds pod-level security attributes and common container settings. | *corev1.PodSecurityContext |
| hostNetwork | Host networking is requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false. | bool |
| pvc | PVC definition | *corev1.PersistentVolumeClaim |
| rbacRules | RBACRules represents additional rbac rules which will be applied to the fluent-bit clusterrole. | []rbacv1.PolicyRule |
| disableService | By default will build the related service according to the globalinputs definition. | bool |
| bufferPath | The path where buffer chunks are stored. | *string |
| ports | Ports represents the pod's ports. | []corev1.ContainerPort |
| service | Service represents configurations on the fluent-bit service. | CollectorService |
| schedulerName | SchedulerName represents the desired scheduler for the Fluentbit collector pods | string |

[Back to TOC](#table-of-contents)
# Decorder




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| decodeField | If the content can be decoded in a structured message, append that structure message (keys and values) to the original log message. | string |
| decodeFieldAs | Any content decoded (unstructured or structured) will be replaced in the same key/value, no extra keys are added. | string |

[Back to TOC](#table-of-contents)
# Filter

Filter is the Schema for namespace level filter API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FilterSpec |

[Back to TOC](#table-of-contents)
# FilterItem




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| grep | Grep defines Grep Filter configuration. | *[filter.Grep](plugins/filter/grep.md) |
| recordModifier | RecordModifier defines Record Modifier Filter configuration. | *[filter.RecordModifier](plugins/filter/recordmodifier.md) |
| kubernetes | Kubernetes defines Kubernetes Filter configuration. | *[filter.Kubernetes](plugins/filter/kubernetes.md) |
| modify | Modify defines Modify Filter configuration. | *[filter.Modify](plugins/filter/modify.md) |
| nest | Nest defines Nest Filter configuration. | *[filter.Nest](plugins/filter/nest.md) |
| parser | Parser defines Parser Filter configuration. | *[filter.Parser](plugins/filter/parser.md) |
| lua | Lua defines Lua Filter configuration. | *[filter.Lua](plugins/filter/lua.md) |
| throttle | Throttle defines a Throttle configuration. | *[filter.Throttle](plugins/filter/throttle.md) |
| rewriteTag | RewriteTag defines a RewriteTag configuration. | *[filter.RewriteTag](plugins/filter/rewritetag.md) |
| aws | Aws defines a Aws configuration. | *[filter.AWS](plugins/filter/aws.md) |
| multiline | Multiline defines a Multiline configuration. | *[filter.Multiline](plugins/filter/multiline.md) |
| customPlugin | CustomPlugin defines a Custom plugin configuration. | *custom.CustomPlugin |

[Back to TOC](#table-of-contents)
# FilterList

FilterList contains a list of Filters


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Filter |

[Back to TOC](#table-of-contents)
# FilterSpec

FilterSpec defines the desired state of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | A pattern to match against the tags of incoming records. It's case-sensitive and support the star (*) character as a wildcard. | string |
| matchRegex | A regular expression to match against the tags of incoming records. Use this option if you want to use the full regex syntax. | string |
| logLevel |  | string |
| filters | A set of filter plugins in order. | []FilterItem |

[Back to TOC](#table-of-contents)
# FluentBit

FluentBit is the Schema for the fluentbits API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FluentBitSpec |
| status |  | FluentBitStatus |

[Back to TOC](#table-of-contents)
# FluentBitConfig

FluentBitConfig is the Schema for the API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | NamespacedFluentBitCfgSpec |

[Back to TOC](#table-of-contents)
# FluentBitConfigList

FluentBitConfigList contains a list of Collector


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []FluentBitConfig |

[Back to TOC](#table-of-contents)
# FluentBitConfigSpec

FluentBitConfigSpec defines the desired state of ClusterFluentBitConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| service | Service defines the global behaviour of the Fluent Bit engine. | *Service |
| inputSelector | Select input plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| filterSelector | Select filter plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| outputSelector | Select output plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| parserSelector | Select parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| multilineParserSelector | Select multiline parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| namespace | If namespace is defined, then the configmap and secret for fluent-bit is in this namespace. If it is not defined, it is in the namespace of the fluentd-operator | *string |

[Back to TOC](#table-of-contents)
# FluentBitList

FluentBitList contains a list of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []FluentBit |

[Back to TOC](#table-of-contents)
# FluentBitService

FluentBitService defines the service of the FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| name | Name is the name of the FluentBit service. | string |
| annotations | Annotations to add to each Fluentbit service. | map[string]string |
| labels | Labels to add to each FluentBit service | map[string]string |

[Back to TOC](#table-of-contents)
# FluentBitSpec

FluentBitSpec defines the desired state of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| disableService | DisableService tells if the fluentbit service should be deployed. | bool |
| image | Fluent Bit image. | string |
| args | Fluent Bit Watcher command line arguments. | []string |
| command | Fluent Bit Watcher command. | []string |
| imagePullPolicy | Fluent Bit image pull policy. | corev1.PullPolicy |
| imagePullSecrets | Fluent Bit image pull secret | []corev1.LocalObjectReference |
| internalMountPropagation | MountPropagation option for internal mounts | *corev1.MountPropagationMode |
| positionDB | Storage for position db. You will use it if tail input is enabled. | [corev1.VolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volume-v1-core) |
| containerLogRealPath | Container log path | string |
| resources | Compute Resources required by container. | corev1.ResourceRequirements |
| nodeSelector | NodeSelector | map[string]string |
| affinity | Pod's scheduling constraints. | *corev1.Affinity |
| tolerations | Tolerations | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#toleration-v1-core) |
| fluentBitConfigName | Fluentbitconfig object associated with this Fluentbit | string |
| namespaceFluentBitCfgSelector | NamespacedFluentBitCfgSelector selects the namespace FluentBitConfig associated with this FluentBit | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| secrets | The Secrets are mounted into /fluent-bit/secrets/<secret-name>. | []string |
| runtimeClassName | RuntimeClassName represents the container runtime configuration. | string |
| priorityClassName | PriorityClassName represents the pod's priority class. | string |
| volumes | List of volumes that can be mounted by containers belonging to the pod. | []corev1.Volume |
| volumesMounts | Pod volumes to mount into the container's filesystem. | []corev1.VolumeMount |
| disableLogVolumes | DisableLogVolumes removes the hostPath mounts for varlibcontainers, varlogs and systemd. | bool |
| annotations | Annotations to add to each Fluentbit pod. | map[string]string |
| serviceAccountAnnotations | Annotations to add to the Fluentbit service account | map[string]string |
| labels | Labels to add to each FluentBit pod | map[string]string |
| securityContext | SecurityContext holds pod-level security attributes and common container settings. | *corev1.PodSecurityContext |
| containerSecurityContext | ContainerSecurityContext holds container-level security attributes. | *corev1.SecurityContext |
| hostNetwork | Host networking is requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false. | bool |
| envVars | EnvVars represent environment variables that can be passed to fluentbit pods. | []corev1.EnvVar |
| livenessProbe | LivenessProbe represents the pod's liveness probe. | *corev1.Probe |
| readinessProbe | ReadinessProbe represents the pod's readiness probe. | *corev1.Probe |
| initContainers | InitContainers represents the pod's init containers. | []corev1.Container |
| ports | Ports represents the pod's ports. | []corev1.ContainerPort |
| rbacRules | RBACRules represents additional rbac rules which will be applied to the fluent-bit clusterrole. | []rbacv1.PolicyRule |
| dnsPolicy | Set DNS policy for the pod. Defaults to \"ClusterFirst\". Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. | corev1.DNSPolicy |
| metricsPort | MetricsPort is the port used by the metrics server. If this option is set, HttpPort from ClusterFluentBitConfig needs to match this value. Default is 2020. | int32 |
| service | Service represents configurations on the fluent-bit service. | FluentBitService |
| schedulerName | SchedulerName represents the desired scheduler for fluent-bit pods. | string |

[Back to TOC](#table-of-contents)
# InputSpec

InputSpec defines the desired state of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| alias | A user friendly alias name for this input plugin. Used in metrics for distinction of each configured input. | string |
| logLevel |  | string |
| dummy | Dummy defines Dummy Input configuration. | *[input.Dummy](plugins/input/dummy.md) |
| tail | Tail defines Tail Input configuration. | *[input.Tail](plugins/input/tail.md) |
| systemd | Systemd defines Systemd Input configuration. | *[input.Systemd](plugins/input/systemd.md) |
| nodeExporterMetrics | NodeExporterMetrics defines Node Exporter Metrics Input configuration. | *[input.NodeExporterMetrics](plugins/input/nodeexportermetrics.md) |
| prometheusScrapeMetrics | PrometheusScrapeMetrics  defines Prometheus Scrape Metrics Input configuration. | *[input.PrometheusScrapeMetrics](plugins/input/prometheusscrapemetrics.md) |
| fluentBitMetrics | FluentBitMetrics defines Fluent Bit Metrics Input configuration. | *[input.FluentbitMetrics](plugins/input/fluentbitmetrics.md) |
| customPlugin | CustomPlugin defines Custom Input configuration. | *custom.CustomPlugin |
| forward | Forward defines forward  input plugin configuration | *[input.Forward](plugins/input/forward.md) |
| openTelemetry | OpenTelemetry defines the OpenTelemetry input plugin configuration | *[input.OpenTelemetry](plugins/input/opentelemetry.md) |
| http | HTTP defines the HTTP input plugin configuration | *[input.HTTP](plugins/input/http.md) |
| mqtt | MQTT defines the MQTT input plugin configuration | *[input.MQTT](plugins/input/mqtt.md) |
| collectd | Collectd defines the Collectd input plugin configuration | *[input.Collectd](plugins/input/collectd.md) |
| statsd | StatsD defines the StatsD input plugin configuration | *[input.StatsD](plugins/input/statsd.md) |
| nginx | Nginx defines the Nginx input plugin configuration | *[input.Nginx](plugins/input/nginx.md) |
| syslog | Syslog defines the Syslog input plugin configuration | *[input.Syslog](plugins/input/syslog.md) |
| tcp | TCP defines the TCP input plugin configuration | *[input.TCP](plugins/input/tcp.md) |

[Back to TOC](#table-of-contents)
# MultilineParser

MultilineParser is the Schema of namespace-level multiline parser API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | MultilineParserSpec |

[Back to TOC](#table-of-contents)
# MultilineParserList

MultilineParserList contains a list of MultilineParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []MultilineParser |

[Back to TOC](#table-of-contents)
# NamespacedFluentBitCfgSpec

NamespacedFluentBitCfgSpec defines the desired state of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| filterSelector | Select filter plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| outputSelector | Select output plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| parserSelector | Select parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| service | Service defines the global behaviour of the Fluent Bit engine. | *Service |
| clusterParserSelector | Select cluster level parser config | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| multilineParserSelector | Select multiline parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| clusterMultilineParserSelector | Select cluster level multiline parser config | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |

[Back to TOC](#table-of-contents)
# Output

Output is the schema for namespace level output API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | OutputSpec |

[Back to TOC](#table-of-contents)
# OutputList

OutputList contains a list of Outputs


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Output |

[Back to TOC](#table-of-contents)
# OutputSpec

OutputSpec defines the desired state of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | A pattern to match against the tags of incoming records. It's case sensitive and support the star (*) character as a wildcard. | string |
| matchRegex | A regular expression to match against the tags of incoming records. Use this option if you want to use the full regex syntax. | string |
| alias | A user friendly alias name for this output plugin. Used in metrics for distinction of each configured output. | string |
| logLevel | Set the plugin's logging verbosity level. Allowed values are: off, error, warn, info, debug and trace, Defaults to the SERVICE section's Log_Level | string |
| azureBlob | AzureBlob defines AzureBlob Output Configuration | *[output.AzureBlob](plugins/output/azureblob.md) |
| azureLogAnalytics | AzureLogAnalytics defines AzureLogAnalytics Output Configuration | *[output.AzureLogAnalytics](plugins/output/azureloganalytics.md) |
| cloudWatch | CloudWatch defines CloudWatch Output Configuration | *[output.CloudWatch](plugins/output/cloudwatch.md) |
| retry_limit | RetryLimit represents configuration for the scheduler which can be set independently on each output section. This option allows to disable retries or impose a limit to try N times and then discard the data after reaching that limit. | string |
| es | Elasticsearch defines Elasticsearch Output configuration. | *[output.Elasticsearch](plugins/output/elasticsearch.md) |
| file | File defines File Output configuration. | *[output.File](plugins/output/file.md) |
| forward | Forward defines Forward Output configuration. | *[output.Forward](plugins/output/forward.md) |
| http | HTTP defines HTTP Output configuration. | *[output.HTTP](plugins/output/http.md) |
| kafka | Kafka defines Kafka Output configuration. | *[output.Kafka](plugins/output/kafka.md) |
| null | Null defines Null Output configuration. | *[output.Null](plugins/output/null.md) |
| stdout | Stdout defines Stdout Output configuration. | *[output.Stdout](plugins/output/stdout.md) |
| tcp | TCP defines TCP Output configuration. | *[output.TCP](plugins/output/tcp.md) |
| loki | Loki defines Loki Output configuration. | *[output.Loki](plugins/output/loki.md) |
| syslog | Syslog defines Syslog Output configuration. | *[output.Syslog](plugins/output/syslog.md) |
| influxDB | InfluxDB defines InfluxDB Output configuration. | *[output.InfluxDB](plugins/output/influxdb.md) |
| datadog | DataDog defines DataDog Output configuration. | *[output.DataDog](plugins/output/datadog.md) |
| firehose | Firehose defines Firehose Output configuration. | *[output.Firehose](plugins/output/firehose.md) |
| kinesis | Kinesis defines Kinesis Output configuration. | *[output.Kinesis](plugins/output/kinesis.md) |
| stackdriver | Stackdriver defines Stackdriver Output Configuration | *[output.Stackdriver](plugins/output/stackdriver.md) |
| splunk | Splunk defines Splunk Output Configuration | *[output.Splunk](plugins/output/splunk.md) |
| opensearch | OpenSearch defines OpenSearch Output configuration. | *[output.OpenSearch](plugins/output/opensearch.md) |
| opentelemetry | OpenTelemetry defines OpenTelemetry Output configuration. | *[output.OpenTelemetry](plugins/output/opentelemetry.md) |
| prometheusExporter | PrometheusExporter_types defines Prometheus exporter configuration to expose metrics from Fluent Bit. | *[output.PrometheusExporter](plugins/output/prometheusexporter.md) |
| prometheusRemoteWrite | PrometheusRemoteWrite_types defines Prometheus Remote Write configuration. | *[output.PrometheusRemoteWrite](plugins/output/prometheusremotewrite.md) |
| s3 | S3 defines S3 Output configuration. | *[output.S3](plugins/output/s3.md) |
| gelf | Gelf defines GELF Output configuration. | *[output.Gelf](plugins/output/gelf.md) |
| customPlugin | CustomPlugin defines Custom Output configuration. | *custom.CustomPlugin |

[Back to TOC](#table-of-contents)
# Parser

Parser is the Schema for namespace level parser API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | ParserSpec |

[Back to TOC](#table-of-contents)
# ParserList

ParserList contains a list of Parsers


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Parser |

[Back to TOC](#table-of-contents)
# ParserSpec

ParserSpec defines the desired state of ClusterParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| json | JSON defines json parser configuration. | *[parser.JSON](plugins/parser/json.md) |
| regex | Regex defines regex parser configuration. | *[parser.Regex](plugins/parser/regex.md) |
| ltsv | LTSV defines ltsv parser configuration. | *[parser.LSTV](plugins/parser/lstv.md) |
| logfmt | Logfmt defines logfmt parser configuration. | *[parser.Logfmt](plugins/parser/logfmt.md) |
| decoders | Decoders are a built-in feature available through the Parsers file, each Parser definition can optionally set one or multiple decoders. There are two type of decoders type: Decode_Field and Decode_Field_As. | []Decorder |

[Back to TOC](#table-of-contents)
# Script




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| Name |  | string |
| Content |  | string |

[Back to TOC](#table-of-contents)
# Service




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| daemon | If true go to background on start | *bool |
| flushSeconds | Interval to flush output | *int64 |
| graceSeconds | Wait time on exit | *int64 |
| hcErrorsCount | the error count to meet the unhealthy requirement, this is a sum for all output plugins in a defined HC_Period, example for output error: [2022/02/16 10:44:10] [ warn] [engine] failed to flush chunk '1-1645008245.491540684.flb', retry in 7 seconds: task_id=0, input=forward.1 > output=cloudwatch_logs.3 (out_id=3) | *int64 |
| hcRetryFailureCount | the retry failure count to meet the unhealthy requirement, this is a sum for all output plugins in a defined HC_Period, example for retry failure: [2022/02/16 20:11:36] [ warn] [engine] chunk '1-1645042288.260516436.flb' cannot be retried: task_id=0, input=tcp.3 > output=cloudwatch_logs.1 | *int64 |
| hcPeriod | The time period by second to count the error and retry failure data point | *int64 |
| healthCheck | enable Health check feature at http://127.0.0.1:2020/api/v1/health Note: Enabling this will not automatically configure kubernetes to use fluentbit's healthcheck endpoint | *bool |
| httpListen | Address to listen | string |
| httpPort | Port to listen | *int32 |
| httpServer | If true enable statistics HTTP server | *bool |
| logFile | File to log diagnostic output | string |
| logLevel | Diagnostic level (error/warning/info/debug/trace) | string |
| parsersFile | Optional 'parsers' config file (can be multiple) | string |
| parsersFiles | backward compatible | []string |
| storage | Configure a global environment for the storage layer in Service. It is recommended to configure the volume and volumeMount separately for this storage. The hostPath type should be used for that Volume in Fluentbit daemon set. | *Storage |
| emitterName | Per-namespace re-emitter configuration | string |
| emitterMemBufLimit |  | string |
| emitterStorageType |  | string |

[Back to TOC](#table-of-contents)
# Storage




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| path | Select an optional location in the file system to store streams and chunks of data/ | string |
| sync | Configure the synchronization mode used to store the data into the file system | string |
| checksum | Enable the data integrity check when writing and reading data from the filesystem | string |
| backlogMemLimit | This option configure a hint of maximum value of memory to use when processing these records | string |
| maxChunksUp | If the input plugin has enabled filesystem storage type, this property sets the maximum number of Chunks that can be up in memory | *int64 |
| metrics | If http_server option has been enabled in the Service section, this option registers a new endpoint where internal metrics of the storage layer can be consumed | string |
| deleteIrrecoverableChunks | When enabled, irrecoverable chunks will be deleted during runtime, and any other irrecoverable chunk located in the configured storage path directory will be deleted when Fluent-Bit starts. | string |

[Back to TOC](#table-of-contents)
