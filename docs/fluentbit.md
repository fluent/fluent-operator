# API Docs
This Document documents the types introduced by the fluentbit Operator to be consumed by users.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.
## Table of Contents
* [ClusterFilter](#clusterfilter)
* [ClusterFluentBitConfig](#clusterfluentbitconfig)
* [ClusterInput](#clusterinput)
* [ClusterOutput](#clusteroutput)
* [ClusterParser](#clusterparser)
* [Decorder](#decorder)
* [FilterItem](#filteritem)
* [FilterList](#filterlist)
* [FilterSpec](#filterspec)
* [FluentBit](#fluentbit)
* [FluentBitConfigList](#fluentbitconfiglist)
* [FluentBitConfigSpec](#fluentbitconfigspec)
* [FluentBitList](#fluentbitlist)
* [FluentBitSpec](#fluentbitspec)
* [InputList](#inputlist)
* [InputSpec](#inputspec)
* [OutputList](#outputlist)
* [OutputSpec](#outputspec)
* [ParserList](#parserlist)
* [ParserSpec](#parserspec)
* [Script](#script)
* [Service](#service)
# ClusterFilter

ClusterFilter defines a cluster-level Filter configuration.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec | Specification of desired Filter configuration. | FilterSpec |

[Back to TOC](#table-of-contents)
# ClusterFluentBitConfig

ClusterFluentBitConfig is the Schema for the cluster-level fluentbitconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FluentBitConfigSpec |

[Back to TOC](#table-of-contents)
# ClusterInput

ClusterInput is the Schema for the inputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | InputSpec |

[Back to TOC](#table-of-contents)
# ClusterOutput

ClusterOutput is the Schema for the cluster-level outputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | OutputSpec |

[Back to TOC](#table-of-contents)
# ClusterParser

ClusterParser is the Schema for the cluster-level parsers API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | ParserSpec |

[Back to TOC](#table-of-contents)
# Decorder




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| decodeField | If the content can be decoded in a structured message, append that structure message (keys and values) to the original log message. | string |
| decodeFieldAs | Any content decoded (unstructured or structured) will be replaced in the same key/value, no extra keys are added. | string |

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

[Back to TOC](#table-of-contents)
# FilterList

FilterList contains a list of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterFilter |

[Back to TOC](#table-of-contents)
# FilterSpec

FilterSpec defines the desired state of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | A pattern to match against the tags of incoming records. It's case-sensitive and support the star (*) character as a wildcard. | string |
| matchRegex | A regular expression to match against the tags of incoming records. Use this option if you want to use the full regex syntax. | string |
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
# FluentBitConfigList

FluentBitConfigList contains a list of ClusterFluentBitConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterFluentBitConfig |

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
| namespace | If namespace is defined, then the configmap and secret for fluent-bit is in this namespace. If it is not defined, it is in the namespace of the fluentd-operator | *string |

[Back to TOC](#table-of-contents)
# FluentBitList

FluentBitList contains a list of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []FluentBit |

[Back to TOC](#table-of-contents)
# FluentBitSpec

FluentBitSpec defines the desired state of FluentBit


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| image | Fluent Bit image. | string |
| args | Fluent Bit Watcher command line arguments. | []string |
| imagePullPolicy | Fluent Bit image pull policy. | corev1.PullPolicy |
| imagePullSecrets | Fluent Bit image pull secret | []corev1.LocalObjectReference |
| positionDB | Storage for position db. You will use it if tail input is enabled. | [corev1.VolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volume-v1-core) |
| containerLogRealPath | Container log path | string |
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

[Back to TOC](#table-of-contents)
# InputList

ClusterInputList contains a list of Input


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterInput |

[Back to TOC](#table-of-contents)
# InputSpec

InputSpec defines the desired state of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| alias | A user friendly alias name for this input plugin. Used in metrics for distinction of each configured input. | string |
| dummy | Dummy defines Dummy Input configuration. | *[input.Dummy](plugins/input/dummy.md) |
| tail | Tail defines Tail Input configuration. | *[input.Tail](plugins/input/tail.md) |
| systemd | Systemd defines Systemd Input configuration. | *[input.Systemd](plugins/input/systemd.md) |

[Back to TOC](#table-of-contents)
# OutputList

OutputList contains a list of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterOutput |

[Back to TOC](#table-of-contents)
# OutputSpec

OutputSpec defines the desired state of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| match | A pattern to match against the tags of incoming records. It's case sensitive and support the star (*) character as a wildcard. | string |
| matchRegex | A regular expression to match against the tags of incoming records. Use this option if you want to use the full regex syntax. | string |
| alias | A user friendly alias name for this output plugin. Used in metrics for distinction of each configured output. | string |
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
| datadog | DataDog defines DataDog Output configuration. | *[output.DataDog](plugins/output/datadog.md) |
| firehose | Firehose defines Firehose Output configuration. | *[output.Firehose](plugins/output/firehose.md) |

[Back to TOC](#table-of-contents)
# ParserList

ParserList contains a list of ClusterParser


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterParser |

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
| httpListen | Address to listen | string |
| httpPort | Port to listen | *int32 |
| httpServer | If true enable statistics HTTP server | *bool |
| logFile | File to log diagnostic output | string |
| logLevel | Diagnostic level (error/warning/info/debug/trace) | string |
| parsersFile | Optional 'parsers' config file (can be multiple) | string |

[Back to TOC](#table-of-contents)
