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
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec | Specification of desired Filter configuration. | [FilterSpec](#filterspec) |

[Back to TOC](#table-of-contents)
# ClusterFilterList

ClusterFilterList contains a list of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterFilter](#clusterfilter) |

[Back to TOC](#table-of-contents)
# ClusterFluentBitConfig

ClusterFluentBitConfig is the Schema for the cluster-level fluentbitconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [FluentBitConfigSpec](#fluentbitconfigspec) |

[Back to TOC](#table-of-contents)
# ClusterFluentBitConfigList

ClusterFluentBitConfigList contains a list of ClusterFluentBitConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterFluentBitConfig](#clusterfluentbitconfig) |

[Back to TOC](#table-of-contents)
# ClusterInput

ClusterInput is the Schema for the inputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [InputSpec](#inputspec) |

[Back to TOC](#table-of-contents)
# ClusterInputList

ClusterInputList contains a list of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterInput](#clusterinput) |

[Back to TOC](#table-of-contents)
# ClusterMultilineParser

ClusterMultilineParser is the Schema for the cluster-level multiline parser API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [MultilineParserSpec](#multilineparserspec) |

[Back to TOC](#table-of-contents)
# ClusterMultilineParserList

ClusterMultilineParserList contains a list of ClusterMultilineParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterMultilineParser](#clustermultilineparser) |

[Back to TOC](#table-of-contents)
# ClusterOutput

ClusterOutput is the Schema for the cluster-level outputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [OutputSpec](#outputspec) |

[Back to TOC](#table-of-contents)
# ClusterOutputList

ClusterOutputList contains a list of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterOutput](#clusteroutput) |

[Back to TOC](#table-of-contents)
# ClusterParser

ClusterParser is the Schema for the cluster-level parsers API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [ParserSpec](#parserspec) |

[Back to TOC](#table-of-contents)
# ClusterParserList

ClusterParserList contains a list of ClusterParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterParser](#clusterparser) |

[Back to TOC](#table-of-contents)
# Collector

Collector is the Schema for the fluentbits API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [CollectorSpec](#collectorspec) |
| status |  | [CollectorStatus](#collectorstatus) |

[Back to TOC](#table-of-contents)
# CollectorList

CollectorList contains a list of Collector


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Collector](#collector) |

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
| imagePullPolicy | Fluent Bit image pull policy. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#pullpolicy-v1-core) |
| imagePullSecrets | Fluent Bit image pull secret | [][corev1.LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#localobjectreference-v1-core) |
| resources | Compute Resources required by container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#resourcerequirements-v1-core) |
| nodeSelector | NodeSelector | map[string]string |
| affinity | Pod's scheduling constraints. | *[corev1.Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#affinity-v1-core) |
| tolerations | Tolerations | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#toleration-v1-core) |
| fluentBitConfigName | Fluentbitconfig object associated with this Fluentbit | string |
| secrets | The Secrets are mounted into /fluent-bit/secrets/<secret-name>. | []string |
| runtimeClassName | RuntimeClassName represents the container runtime configuration. | string |
| priorityClassName | PriorityClassName represents the pod's priority class. | string |
| volumes | List of volumes that can be mounted by containers belonging to the pod. | [][corev1.Volume](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volume-v1-core) |
| volumesMounts | Pod volumes to mount into the container's filesystem. | [][corev1.VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volumemount-v1-core) |
| annotations | Annotations to add to each Fluentbit pod. | map[string]string |
| serviceAccountAnnotations | Annotations to add to the Fluentbit service account | map[string]string |
| securityContext | SecurityContext holds pod-level security attributes and common container settings. | *[corev1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#podsecuritycontext-v1-core) |
| hostNetwork | Host networking is requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false. | bool |
| pvc | PVC definition | *[corev1.PersistentVolumeClaim](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaim-v1-core) |
| rbacRules | RBACRules represents additional rbac rules which will be applied to the fluent-bit clusterrole. | [][rbacv1.PolicyRule](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#policyrule-v1-rbac-authorization-k8s-io) |
| disableService | By default will build the related service according to the globalinputs definition. | bool |
| bufferPath | The path where buffer chunks are stored. | *string |
| ports | Ports represents the pod's ports. | [][corev1.ContainerPort](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#containerport-v1-core) |
| service | Service represents configurations on the fluent-bit service. | [CollectorService](#collectorservice) |
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
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [FilterSpec](#filterspec) |

[Back to TOC](#table-of-contents)
# FilterItem




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| grep | Grep defines Grep Filter configuration. | *[filter.Grep](plugins/fluentbit/filter/grep.md) |
| recordModifier | RecordModifier defines Record Modifier Filter configuration. | *[filter.RecordModifier](plugins/fluentbit/filter/record_modifier.md) |
| kubernetes | Kubernetes defines Kubernetes Filter configuration. | *[filter.Kubernetes](plugins/fluentbit/filter/kubernetes.md) |
| modify | Modify defines Modify Filter configuration. | *[filter.Modify](plugins/fluentbit/filter/modify.md) |
| nest | Nest defines Nest Filter configuration. | *[filter.Nest](plugins/fluentbit/filter/nest.md) |
| parser | Parser defines Parser Filter configuration. | *[filter.Parser](plugins/fluentbit/filter/parser.md) |
| lua | Lua defines Lua Filter configuration. | *[filter.Lua](plugins/fluentbit/filter/lua.md) |
| throttle | Throttle defines a Throttle configuration. | *[filter.Throttle](plugins/fluentbit/filter/throttle.md) |
| rewriteTag | RewriteTag defines a RewriteTag configuration. | *[filter.RewriteTag](plugins/fluentbit/filter/rewrite_tag.md) |
| aws | Aws defines a Aws configuration. | *[filter.AWS](plugins/fluentbit/filter/aws.md) |
| multiline | Multiline defines a Multiline configuration. | *[filter.Multiline](plugins/fluentbit/filter/multiline.md) |
| logToMetrics | LogToMetrics defines a Log to Metrics Filter configuration. | *[filter.LogToMetrics](plugins/fluentbit/filter/log_to_metrics.md) |
| wasm | Wasm defines a Wasm configuration. | *[filter.Wasm](plugins/fluentbit/filter/wasm.md) |
| customPlugin | CustomPlugin defines a Custom plugin configuration. | *[custom.CustomPlugin](plugins/fluentbit/custom/custom_plugin.md) |

[Back to TOC](#table-of-contents)
# FilterList

FilterList contains a list of Filters


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Filter](#filter) |

[Back to TOC](#table-of-contents)
# FilterSpec

FilterSpec defines the desired state of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | A pattern to match against the tags of incoming records. It's case-sensitive and support the star (*) character as a wildcard. | string |
| matchRegex | A regular expression to match against the tags of incoming records. Use this option if you want to use the full regex syntax. | string |
| logLevel |  | string |
| filters | A set of filter plugins in order. | [][FilterItem](#filteritem) |

[Back to TOC](#table-of-contents)
# FluentBit

FluentBit is the Schema for the fluentbits API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [FluentBitSpec](#fluentbitspec) |
| status |  | [FluentBitStatus](#fluentbitstatus) |

[Back to TOC](#table-of-contents)
# FluentBitConfig

FluentBitConfig is the Schema for the API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [NamespacedFluentBitCfgSpec](#namespacedfluentbitcfgspec) |

[Back to TOC](#table-of-contents)
# FluentBitConfigList

FluentBitConfigList contains a list of Collector


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][FluentBitConfig](#fluentbitconfig) |

[Back to TOC](#table-of-contents)
# FluentBitConfigSpec

FluentBitConfigSpec defines the desired state of ClusterFluentBitConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| service | Service defines the global behaviour of the Fluent Bit engine. | *[Service](#service) |
| inputSelector | Select input plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| filterSelector | Select filter plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| outputSelector | Select output plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| parserSelector | Select parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| multilineParserSelector | Select multiline parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| namespace | If namespace is defined, then the configmap and secret for fluent-bit is in this namespace. If it is not defined, it is in the namespace of the fluentd-operator | *string |
| configFileFormat | ConfigFileFormat defines the format of the config file, default is \"classic\", available options are \"classic\" and \"yaml\" | *string |

[Back to TOC](#table-of-contents)
# FluentBitList

FluentBitList contains a list of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][FluentBit](#fluentbit) |

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
| imagePullPolicy | Fluent Bit image pull policy. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#pullpolicy-v1-core) |
| imagePullSecrets | Fluent Bit image pull secret | [][corev1.LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#localobjectreference-v1-core) |
| internalMountPropagation | MountPropagation option for internal mounts | *[corev1.MountPropagationMode](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#mountpropagationmode-v1-core) |
| positionDB | Storage for position db. You will use it if tail input is enabled. | [corev1.VolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volumesource-v1-core) |
| containerLogRealPath | Container log path | string |
| resources | Compute Resources required by container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#resourcerequirements-v1-core) |
| nodeSelector | NodeSelector | map[string]string |
| affinity | Pod's scheduling constraints. | *[corev1.Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#affinity-v1-core) |
| tolerations | Tolerations | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#toleration-v1-core) |
| fluentBitConfigName | Fluentbitconfig object associated with this Fluentbit | string |
| namespaceFluentBitCfgSelector | NamespacedFluentBitCfgSelector selects the namespace FluentBitConfig associated with this FluentBit | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| secrets | The Secrets are mounted into /fluent-bit/secrets/<secret-name>. | []string |
| runtimeClassName | RuntimeClassName represents the container runtime configuration. | string |
| priorityClassName | PriorityClassName represents the pod's priority class. | string |
| volumes | List of volumes that can be mounted by containers belonging to the pod. | [][corev1.Volume](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volume-v1-core) |
| volumesMounts | Pod volumes to mount into the container's filesystem. | [][corev1.VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volumemount-v1-core) |
| disableLogVolumes | DisableLogVolumes removes the hostPath mounts for varlibcontainers, varlogs and systemd. | bool |
| annotations | Annotations to add to each Fluentbit pod. | map[string]string |
| serviceAccountAnnotations | Annotations to add to the Fluentbit service account | map[string]string |
| labels | Labels to add to each FluentBit pod | map[string]string |
| securityContext | SecurityContext holds pod-level security attributes and common container settings. | *[corev1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#podsecuritycontext-v1-core) |
| containerSecurityContext | ContainerSecurityContext holds container-level security attributes. | *[corev1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#securitycontext-v1-core) |
| hostNetwork | Host networking is requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false. | bool |
| envVars | EnvVars represent environment variables that can be passed to fluentbit pods. | [][corev1.EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#envvar-v1-core) |
| livenessProbe | LivenessProbe represents the pod's liveness probe. | *[corev1.Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#probe-v1-core) |
| readinessProbe | ReadinessProbe represents the pod's readiness probe. | *[corev1.Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#probe-v1-core) |
| initContainers | InitContainers represents the pod's init containers. | [][corev1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#container-v1-core) |
| ports | Ports represents the pod's ports. | [][corev1.ContainerPort](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#containerport-v1-core) |
| rbacRules | RBACRules represents additional rbac rules which will be applied to the fluent-bit clusterrole. | [][rbacv1.PolicyRule](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#policyrule-v1-rbac-authorization-k8s-io) |
| dnsPolicy | Set DNS policy for the pod. Defaults to \"ClusterFirst\". Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. | [corev1.DNSPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#dnspolicy-v1-core) |
| metricsPort | MetricsPort is the port used by the metrics server. If this option is set, HttpPort from ClusterFluentBitConfig needs to match this value. Default is 2020. | int32 |
| service | Service represents configurations on the fluent-bit service. | [FluentBitService](#fluentbitservice) |
| schedulerName | SchedulerName represents the desired scheduler for fluent-bit pods. | string |
| terminationGracePeriodSeconds | Optional duration in seconds the pod needs to terminate gracefully. Value must be non-negative integer. | *int64 |

[Back to TOC](#table-of-contents)
# InputSpec

InputSpec defines the desired state of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| alias | A user friendly alias name for this input plugin. Used in metrics for distinction of each configured input. | string |
| logLevel |  | string |
| dummy | Dummy defines Dummy Input configuration. | *[input.Dummy](plugins/fluentbit/input/dummy.md) |
| tail | Tail defines Tail Input configuration. | *[input.Tail](plugins/fluentbit/input/tail.md) |
| systemd | Systemd defines Systemd Input configuration. | *[input.Systemd](plugins/fluentbit/input/systemd.md) |
| nodeExporterMetrics | NodeExporterMetrics defines Node Exporter Metrics Input configuration. | *[input.NodeExporterMetrics](plugins/fluentbit/input/node_exporter_metrics.md) |
| prometheusScrapeMetrics | PrometheusScrapeMetrics  defines Prometheus Scrape Metrics Input configuration. | *[input.PrometheusScrapeMetrics](plugins/fluentbit/input/prometheus_scrape_metrics.md) |
| fluentBitMetrics | FluentBitMetrics defines Fluent Bit Metrics Input configuration. | *[input.FluentbitMetrics](plugins/fluentbit/input/fluentbit_metrics.md) |
| customPlugin | CustomPlugin defines Custom Input configuration. | *[custom.CustomPlugin](plugins/fluentbit/custom/custom_plugin.md) |
| forward | Forward defines forward  input plugin configuration | *[input.Forward](plugins/fluentbit/input/forward.md) |
| openTelemetry | OpenTelemetry defines the OpenTelemetry input plugin configuration | *[input.OpenTelemetry](plugins/fluentbit/input/open_telemetry.md) |
| http | HTTP defines the HTTP input plugin configuration | *[input.HTTP](plugins/fluentbit/input/http.md) |
| mqtt | MQTT defines the MQTT input plugin configuration | *[input.MQTT](plugins/fluentbit/input/mqtt.md) |
| collectd | Collectd defines the Collectd input plugin configuration | *[input.Collectd](plugins/fluentbit/input/collectd.md) |
| statsd | StatsD defines the StatsD input plugin configuration | *[input.StatsD](plugins/fluentbit/input/stats_d.md) |
| nginx | Nginx defines the Nginx input plugin configuration | *[input.Nginx](plugins/fluentbit/input/nginx.md) |
| syslog | Syslog defines the Syslog input plugin configuration | *[input.Syslog](plugins/fluentbit/input/syslog.md) |
| tcp | TCP defines the TCP input plugin configuration | *[input.TCP](plugins/fluentbit/input/tcp.md) |
| udp | UDP defines the UDP input plugin configuration | *[input.UDP](plugins/fluentbit/input/udp.md) |
| kubernetesEvents | KubernetesEvents defines the KubernetesEvents input plugin configuration | *[input.KubernetesEvents](plugins/fluentbit/input/kubernetes_events.md) |
| execWasi | ExecWasi defines the exec wasi input plugin configuration | *[input.ExecWasi](plugins/fluentbit/input/exec_wasi.md) |
| processors | Processors defines the processors configuration | *plugins.Config |

[Back to TOC](#table-of-contents)
# MultilineParser

MultilineParser is the Schema of namespace-level multiline parser API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [MultilineParserSpec](#multilineparserspec) |

[Back to TOC](#table-of-contents)
# MultilineParserList

MultilineParserList contains a list of MultilineParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][MultilineParser](#multilineparser) |

[Back to TOC](#table-of-contents)
# NamespacedFluentBitCfgSpec

NamespacedFluentBitCfgSpec defines the desired state of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| filterSelector | Select filter plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| outputSelector | Select output plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| parserSelector | Select parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| service | Service defines the global behaviour of the Fluent Bit engine. | *[Service](#service) |
| clusterParserSelector | Select cluster level parser config | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| multilineParserSelector | Select multiline parser plugins | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| clusterMultilineParserSelector | Select cluster level multiline parser config | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |

[Back to TOC](#table-of-contents)
# Output

Output is the schema for namespace level output API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [OutputSpec](#outputspec) |

[Back to TOC](#table-of-contents)
# OutputList

OutputList contains a list of Outputs


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Output](#output) |

[Back to TOC](#table-of-contents)
# OutputSpec

OutputSpec defines the desired state of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | A pattern to match against the tags of incoming records. It's case sensitive and support the star (*) character as a wildcard. | string |
| matchRegex | A regular expression to match against the tags of incoming records. Use this option if you want to use the full regex syntax. | string |
| alias | A user friendly alias name for this output plugin. Used in metrics for distinction of each configured output. | string |
| logLevel | Set the plugin's logging verbosity level. Allowed values are: off, error, warn, info, debug and trace, Defaults to the SERVICE section's Log_Level | string |
| azureBlob | AzureBlob defines AzureBlob Output Configuration | *[output.AzureBlob](plugins/fluentbit/output/azure_blob.md) |
| azureLogAnalytics | AzureLogAnalytics defines AzureLogAnalytics Output Configuration | *[output.AzureLogAnalytics](plugins/fluentbit/output/azure_log_analytics.md) |
| cloudWatch | CloudWatch defines CloudWatch Output Configuration | *[output.CloudWatch](plugins/fluentbit/output/cloud_watch.md) |
| retry_limit | RetryLimit represents configuration for the scheduler which can be set independently on each output section. This option allows to disable retries or impose a limit to try N times and then discard the data after reaching that limit. | string |
| es | Elasticsearch defines Elasticsearch Output configuration. | *[output.Elasticsearch](plugins/fluentbit/output/elasticsearch.md) |
| file | File defines File Output configuration. | *[output.File](plugins/fluentbit/output/file.md) |
| forward | Forward defines Forward Output configuration. | *[output.Forward](plugins/fluentbit/output/forward.md) |
| http | HTTP defines HTTP Output configuration. | *[output.HTTP](plugins/fluentbit/output/http.md) |
| kafka | Kafka defines Kafka Output configuration. | *[output.Kafka](plugins/fluentbit/output/kafka.md) |
| null | Null defines Null Output configuration. | *[output.Null](plugins/fluentbit/output/null.md) |
| stdout | Stdout defines Stdout Output configuration. | *[output.Stdout](plugins/fluentbit/output/stdout.md) |
| tcp | TCP defines TCP Output configuration. | *[output.TCP](plugins/fluentbit/output/tcp.md) |
| loki | Loki defines Loki Output configuration. | *[output.Loki](plugins/fluentbit/output/loki.md) |
| syslog | Syslog defines Syslog Output configuration. | *[output.Syslog](plugins/fluentbit/output/syslog.md) |
| influxDB | InfluxDB defines InfluxDB Output configuration. | *[output.InfluxDB](plugins/fluentbit/output/influx_db.md) |
| datadog | DataDog defines DataDog Output configuration. | *[output.DataDog](plugins/fluentbit/output/data_dog.md) |
| firehose | Firehose defines Firehose Output configuration. | *[output.Firehose](plugins/fluentbit/output/firehose.md) |
| kinesis | Kinesis defines Kinesis Output configuration. | *[output.Kinesis](plugins/fluentbit/output/kinesis.md) |
| stackdriver | Stackdriver defines Stackdriver Output Configuration | *[output.Stackdriver](plugins/fluentbit/output/stackdriver.md) |
| splunk | Splunk defines Splunk Output Configuration | *[output.Splunk](plugins/fluentbit/output/splunk.md) |
| opensearch | OpenSearch defines OpenSearch Output configuration. | *[output.OpenSearch](plugins/fluentbit/output/open_search.md) |
| opentelemetry | OpenTelemetry defines OpenTelemetry Output configuration. | *[output.OpenTelemetry](plugins/fluentbit/output/open_telemetry.md) |
| prometheusExporter | PrometheusExporter_types defines Prometheus exporter configuration to expose metrics from Fluent Bit. | *[output.PrometheusExporter](plugins/fluentbit/output/prometheus_exporter.md) |
| prometheusRemoteWrite | PrometheusRemoteWrite_types defines Prometheus Remote Write configuration. | *[output.PrometheusRemoteWrite](plugins/fluentbit/output/prometheus_remote_write.md) |
| s3 | S3 defines S3 Output configuration. | *[output.S3](plugins/fluentbit/output/s3.md) |
| gelf | Gelf defines GELF Output configuration. | *[output.Gelf](plugins/fluentbit/output/gelf.md) |
| customPlugin | CustomPlugin defines Custom Output configuration. | *[custom.CustomPlugin](plugins/fluentbit/custom/custom_plugin.md) |
| processors | Processors defines the processors configuration | *plugins.Config |

[Back to TOC](#table-of-contents)
# Parser

Parser is the Schema for namespace level parser API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [ParserSpec](#parserspec) |

[Back to TOC](#table-of-contents)
# ParserList

ParserList contains a list of Parsers


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Parser](#parser) |

[Back to TOC](#table-of-contents)
# ParserSpec

ParserSpec defines the desired state of ClusterParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| json | JSON defines json parser configuration. | *[parser.JSON](plugins/fluentbit/parser/json.md) |
| regex | Regex defines regex parser configuration. | *[parser.Regex](plugins/fluentbit/parser/regex.md) |
| ltsv | LTSV defines ltsv parser configuration. | *[parser.LSTV](plugins/fluentbit/parser/lstv.md) |
| logfmt | Logfmt defines logfmt parser configuration. | *[parser.Logfmt](plugins/fluentbit/parser/logfmt.md) |
| decoders | Decoders are a built-in feature available through the Parsers file, each Parser definition can optionally set one or multiple decoders. There are two type of decoders type: Decode_Field and Decode_Field_As. | [][Decorder](#decorder) |

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
| flushSeconds | Interval to flush output | *float64 |
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
| storage | Configure a global environment for the storage layer in Service. It is recommended to configure the volume and volumeMount separately for this storage. The hostPath type should be used for that Volume in Fluentbit daemon set. | *[Storage](#storage) |
| emitterName | Per-namespace re-emitter configuration | string |
| emitterMemBufLimit |  | string |
| emitterStorageType |  | string |
| hotReload | If true enable reloading via HTTP | *bool |

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
