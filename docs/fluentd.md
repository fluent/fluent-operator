# API Docs
This Document documents the types introduced by the fluentd Operator.
> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.
## Table of Contents
* [BufferVolume](#buffervolume)
* [ClusterFilter](#clusterfilter)
* [ClusterFilterList](#clusterfilterlist)
* [ClusterFilterSpec](#clusterfilterspec)
* [ClusterFluentdConfig](#clusterfluentdconfig)
* [ClusterFluentdConfigList](#clusterfluentdconfiglist)
* [ClusterFluentdConfigSpec](#clusterfluentdconfigspec)
* [ClusterFluentdConfigStatus](#clusterfluentdconfigstatus)
* [ClusterInput](#clusterinput)
* [ClusterInputList](#clusterinputlist)
* [ClusterInputSpec](#clusterinputspec)
* [ClusterOutput](#clusteroutput)
* [ClusterOutputList](#clusteroutputlist)
* [ClusterOutputSpec](#clusteroutputspec)
* [Filter](#filter)
* [FilterList](#filterlist)
* [FilterSpec](#filterspec)
* [FluentDService](#fluentdservice)
* [Fluentd](#fluentd)
* [FluentdConfig](#fluentdconfig)
* [FluentdConfigList](#fluentdconfiglist)
* [FluentdConfigSpec](#fluentdconfigspec)
* [FluentdConfigStatus](#fluentdconfigstatus)
* [FluentdList](#fluentdlist)
* [FluentdSpec](#fluentdspec)
* [FluentdStatus](#fluentdstatus)
* [Input](#input)
* [InputList](#inputlist)
* [InputSpec](#inputspec)
* [Output](#output)
* [OutputList](#outputlist)
* [OutputSpec](#outputspec)
# BufferVolume




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| disableBufferVolume | Enabled buffer pvc by default. | bool |
| hostPath | Volume definition. | *[corev1.HostPathVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#hostpathvolumesource-v1-core) |
| emptyDir |  | *[corev1.EmptyDirVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#emptydirvolumesource-v1-core) |
| pvc | PVC definition | *[corev1.PersistentVolumeClaim](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaim-v1-core) |

[Back to TOC](#table-of-contents)
# ClusterFilter

ClusterFilter is the Schema for the clusterfilters API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [ClusterFilterSpec](#clusterfilterspec) |
| status |  | [ClusterFilterStatus](#clusterfilterstatus) |

[Back to TOC](#table-of-contents)
# ClusterFilterList

ClusterFilterList contains a list of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterFilter](#clusterfilter) |

[Back to TOC](#table-of-contents)
# ClusterFilterSpec

ClusterFilterSpec defines the desired state of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| filters |  | [][filter.Filter](plugins/fluentd/filter/filter.md) |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfig

ClusterFluentdConfig is the Schema for the clusterfluentdconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [ClusterFluentdConfigSpec](#clusterfluentdconfigspec) |
| status |  | [ClusterFluentdConfigStatus](#clusterfluentdconfigstatus) |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfigList

ClusterFluentdConfigList contains a list of ClusterFluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterFluentdConfig](#clusterfluentdconfig) |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfigSpec

ClusterFluentdConfigSpec defines the desired state of ClusterFluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| emit_mode | Emit mode. If batch, the plugin will emit events per labels matched. Enum: record, batch. will make no effect if EnableFilterKubernetes is set false. | string |
| stickyTags | Sticky tags will match only one record from an event stream. The same tag will be treated the same way. will make no effect if EnableFilterKubernetes is set false. | string |
| watchedNamespaces | A set of namespaces. The whole namespaces would be watched if left empty. | []string |
| watchedHosts | A set of hosts. Ignored if left empty. | []string |
| watchedConstainers | A set of container names. Ignored if left empty. | []string |
| watchedLabels | Use this field to filter the logs, will make no effect if EnableFilterKubernetes is set false. | map[string]string |
| clusterFilterSelector | Select cluster filter plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| clusterOutputSelector | Select cluster output plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| clusterInputSelector | Select cluster input plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfigStatus

ClusterFluentdConfigStatus defines the observed state of ClusterFluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| messages | Messages defines the plugin errors which is selected by this fluentdconfig | string |
| state | The state of this fluentd config | [StatusState](#statusstate) |

[Back to TOC](#table-of-contents)
# ClusterInput

ClusterInput is the Schema for the clusterinputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [ClusterInputSpec](#clusterinputspec) |
| status |  | [ClusterInputStatus](#clusterinputstatus) |

[Back to TOC](#table-of-contents)
# ClusterInputList

ClusterInputList contains a list of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterInput](#clusterinput) |

[Back to TOC](#table-of-contents)
# ClusterInputSpec

ClusterInputSpec defines the desired state of ClusterInput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| inputs |  | [][input.Input](plugins/fluentd/input/input.md) |

[Back to TOC](#table-of-contents)
# ClusterOutput

ClusterOutput is the Schema for the clusteroutputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [ClusterOutputSpec](#clusteroutputspec) |
| status |  | [ClusterOutputStatus](#clusteroutputstatus) |

[Back to TOC](#table-of-contents)
# ClusterOutputList

ClusterOutputList contains a list of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][ClusterOutput](#clusteroutput) |

[Back to TOC](#table-of-contents)
# ClusterOutputSpec

ClusterOutputSpec defines the desired state of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| outputs |  | [][output.Output](plugins/fluentd/output/output.md) |

[Back to TOC](#table-of-contents)
# Filter

Filter is the Schema for the filters API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [FilterSpec](#filterspec) |
| status |  | [FilterStatus](#filterstatus) |

[Back to TOC](#table-of-contents)
# FilterList

FilterList contains a list of Filter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Filter](#filter) |

[Back to TOC](#table-of-contents)
# FilterSpec

FilterSpec defines the desired state of Filter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| filters |  | [][filter.Filter](plugins/fluentd/filter/filter.md) |

[Back to TOC](#table-of-contents)
# FluentDService

FluentDService the service of the FluentD


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| name | Name is the name of the FluentD service. | string |
| annotations | Annotations to add to each FluentD service. | map[string]string |
| labels | Labels to add to each FluentD service | map[string]string |
| type | Type is the service type to deploy. | *[corev1.ServiceType](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#servicetype-v1-core) |

[Back to TOC](#table-of-contents)
# Fluentd

Fluentd is the Schema for the fluentds API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [FluentdSpec](#fluentdspec) |
| status |  | [FluentdStatus](#fluentdstatus) |

[Back to TOC](#table-of-contents)
# FluentdConfig

FluentdConfig is the Schema for the fluentdconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [FluentdConfigSpec](#fluentdconfigspec) |
| status |  | [FluentdConfigStatus](#fluentdconfigstatus) |

[Back to TOC](#table-of-contents)
# FluentdConfigList

FluentdConfigList contains a list of FluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][FluentdConfig](#fluentdconfig) |

[Back to TOC](#table-of-contents)
# FluentdConfigSpec

FluentdConfigSpec defines the desired state of FluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| emit_mode | Emit mode. If batch, the plugin will emit events per labels matched. Enum: record, batch. will make no effect if EnableFilterKubernetes is set false. | string |
| stickyTags | Sticky tags will match only one record from an event stream. The same tag will be treated the same way. will make no effect if EnableFilterKubernetes is set false. | string |
| watchedHosts | A set of hosts. Ignored if left empty. | []string |
| watchedConstainers | A set of container names. Ignored if left empty. | []string |
| watchedLabels | Use this field to filter the logs, will make no effect if EnableFilterKubernetes is set false. | map[string]string |
| filterSelector | Select namespaced filter plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| outputSelector | Select namespaced output plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| inputSelector | Select cluster input plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| clusterFilterSelector | Select cluster filter plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| clusterOutputSelector | Select cluster output plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| clusterInputSelector | Select cluster input plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |

[Back to TOC](#table-of-contents)
# FluentdConfigStatus

FluentdConfigStatus defines the observed state of FluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| messages | Messages defines the plugin errors which is selected by this fluentdconfig | string |
| state | The state of this fluentd config | [StatusState](#statusstate) |

[Back to TOC](#table-of-contents)
# FluentdList

FluentdList contains a list of Fluentd


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Fluentd](#fluentd) |

[Back to TOC](#table-of-contents)
# FluentdSpec

FluentdSpec defines the desired state of Fluentd


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| globalInputs | Fluentd global inputs. | [][input.Input](plugins/fluentd/input/input.md) |
| defaultInputSelector | Select cluster input plugins used to gather the default cluster output | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| defaultFilterSelector | Select cluster filter plugins used to filter for the default cluster output | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| defaultOutputSelector | Select cluster output plugins used to send all logs that did not match any route to the matching outputs | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| disableService | By default will build the related service according to the globalinputs definition. | bool |
| replicas | Numbers of the Fluentd instance Applicable when the mode is \"collector\", and will be ignored when the mode is \"agent\" | *int32 |
| workers | Numbers of the workers in Fluentd instance | *int32 |
| logLevel | Global logging verbosity | string |
| image | Fluentd image. | string |
| args | Fluentd Watcher command line arguments. | []string |
| envVars | EnvVars represent environment variables that can be passed to fluentd pods. | [][corev1.EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#envvar-v1-core) |
| envFrom | EnvFrom represent environment variables that can be passed to fluentd pods directly from secret or configmap | [][corev1.EnvFromSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#envfromsource-v1-core) |
| fluentdCfgSelector | FluentdCfgSelector defines the selectors to select the fluentd config CRs. | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#labelselector-v1-meta) |
| buffer | Buffer definition | *[BufferVolume](#buffervolume) |
| imagePullPolicy | Fluentd image pull policy. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#pullpolicy-v1-core) |
| imagePullSecrets | Fluentd image pull secret | [][corev1.LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#localobjectreference-v1-core) |
| resources | Compute Resources required by container. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#resourcerequirements-v1-core) |
| nodeSelector | NodeSelector | map[string]string |
| annotations | Annotations to add to each Fluentd pod. | map[string]string |
| serviceAccountAnnotations | Annotations to add to the Fluentd service account | map[string]string |
| affinity | Pod's scheduling constraints. | *[corev1.Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#affinity-v1-core) |
| tolerations | Tolerations | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#toleration-v1-core) |
| runtimeClassName | RuntimeClassName represents the container runtime configuration. | string |
| priorityClassName | PriorityClassName represents the pod's priority class. | string |
| rbacRules | RBACRules represents additional rbac rules which will be applied to the fluentd clusterrole. | [][rbacv1.PolicyRule](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#policyrule-v1-rbac-authorization-k8s-io) |
| volumes | List of volumes that can be mounted by containers belonging to the pod. | [][corev1.Volume](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volume-v1-core) |
| volumeMounts | Pod volumes to mount into the container's filesystem. Cannot be updated. | [][corev1.VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volumemount-v1-core) |
| volumeClaimTemplates | volumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. Applicable when the mode is \"collector\", and will be ignored when the mode is \"agent\" | [][corev1.PersistentVolumeClaim](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#persistentvolumeclaim-v1-core) |
| service | Service represents configurations on the fluentd service. | [FluentDService](#fluentdservice) |
| securityContext | PodSecurityContext represents the security context for the fluentd pods. | *[corev1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#podsecuritycontext-v1-core) |
| schedulerName | SchedulerName represents the desired scheduler for fluentd pods. | string |
| mode | Mode to determine whether to run Fluentd as collector or agent. | string |
| containerSecurityContext | ContainerSecurityContext represents the security context for the fluentd container. | *[corev1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#securitycontext-v1-core) |
| positionDB | Storage for position db. You will use it if tail input is enabled. Applicable when the mode is \"agent\", and will be ignored when the mode is \"collector\" | [corev1.VolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#volumesource-v1-core) |
| livenessProbe | LivenessProbe represents the liveness probe for the fluentd container. | *[corev1.Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#probe-v1-core) |
| readinessProbe | ReadinessProbe represents the readiness probe for the fluentd container. | *[corev1.Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#probe-v1-core) |
| hostAliases | HostAliases is an optional list of IPs and hostnames that will be injected into the pod's hosts file if specified. | [][corev1.HostAlias](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#hostalias-v1-core) |

[Back to TOC](#table-of-contents)
# FluentdStatus

FluentdStatus defines the observed state of Fluentd


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| messages | Messages defines the plugin errors which is selected by this fluentdconfig | string |
| state | The state of this fluentd | [StatusState](#statusstate) |

[Back to TOC](#table-of-contents)
# Input

Input is the Schema for the inputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [InputSpec](#inputspec) |
| status |  | [InputStatus](#inputstatus) |

[Back to TOC](#table-of-contents)
# InputList

InputList contains a list of Input


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Input](#input) |

[Back to TOC](#table-of-contents)
# InputSpec

InputSpec defines the desired state of Input


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| inputs |  | [][input.Input](plugins/fluentd/input/input.md) |

[Back to TOC](#table-of-contents)
# Output

Output is the Schema for the outputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#objectmeta-v1-meta) |
| spec |  | [OutputSpec](#outputspec) |
| status |  | [OutputStatus](#outputstatus) |

[Back to TOC](#table-of-contents)
# OutputList

OutputList contains a list of Output


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.31/#listmeta-v1-meta) |
| items |  | [][Output](#output) |

[Back to TOC](#table-of-contents)
# OutputSpec

OutputSpec defines the desired state of Output


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| outputs |  | [][output.Output](plugins/fluentd/output/output.md) |

[Back to TOC](#table-of-contents)
