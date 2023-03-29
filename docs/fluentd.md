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
* [ClusterOutput](#clusteroutput)
* [ClusterOutputList](#clusteroutputlist)
* [ClusterOutputSpec](#clusteroutputspec)
* [Filter](#filter)
* [FilterList](#filterlist)
* [FilterSpec](#filterspec)
* [Fluentd](#fluentd)
* [FluentdConfig](#fluentdconfig)
* [FluentdConfigList](#fluentdconfiglist)
* [FluentdConfigSpec](#fluentdconfigspec)
* [FluentdConfigStatus](#fluentdconfigstatus)
* [FluentdList](#fluentdlist)
* [FluentdSpec](#fluentdspec)
* [FluentdStatus](#fluentdstatus)
* [Output](#output)
* [OutputList](#outputlist)
* [OutputSpec](#outputspec)
# BufferVolume




| Field | Description | Scheme |
| ----- | ----------- | ------ |
| disableBufferVolume | Enabled buffer pvc by default. | bool |
| hostPath | Volume definition. | *corev1.HostPathVolumeSource |
| emptyDir |  | *corev1.EmptyDirVolumeSource |
| pvc | PVC definition | *corev1.PersistentVolumeClaim |

[Back to TOC](#table-of-contents)
# ClusterFilter

ClusterFilter is the Schema for the clusterfilters API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | ClusterFilterSpec |
| status |  | ClusterFilterStatus |

[Back to TOC](#table-of-contents)
# ClusterFilterList

ClusterFilterList contains a list of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterFilter |

[Back to TOC](#table-of-contents)
# ClusterFilterSpec

ClusterFilterSpec defines the desired state of ClusterFilter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| filters |  | [][[filter.Filter](plugins/filter/filter.md)](plugins/[filter/filter](plugins/filter/filter/md).md) |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfig

ClusterFluentdConfig is the Schema for the clusterfluentdconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | ClusterFluentdConfigSpec |
| status |  | ClusterFluentdConfigStatus |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfigList

ClusterFluentdConfigList contains a list of ClusterFluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterFluentdConfig |

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
| clusterFilterSelector | Select cluster filter plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| clusterOutputSelector | Select cluster output plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |

[Back to TOC](#table-of-contents)
# ClusterFluentdConfigStatus

ClusterFluentdConfigStatus defines the observed state of ClusterFluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| messages | Messages defines the plugin errors which is selected by this fluentdconfig | string |
| state | The state of this fluentd config | StatusState |

[Back to TOC](#table-of-contents)
# ClusterOutput

ClusterOutput is the Schema for the clusteroutputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | ClusterOutputSpec |
| status |  | ClusterOutputStatus |

[Back to TOC](#table-of-contents)
# ClusterOutputList

ClusterOutputList contains a list of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []ClusterOutput |

[Back to TOC](#table-of-contents)
# ClusterOutputSpec

ClusterOutputSpec defines the desired state of ClusterOutput


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| outputs |  | [][[output.Output](plugins/output/output.md)](plugins/[output/output](plugins/output/output/md).md) |

[Back to TOC](#table-of-contents)
# Filter

Filter is the Schema for the filters API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FilterSpec |
| status |  | FilterStatus |

[Back to TOC](#table-of-contents)
# FilterList

FilterList contains a list of Filter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Filter |

[Back to TOC](#table-of-contents)
# FilterSpec

FilterSpec defines the desired state of Filter


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| filters |  | [][[filter.Filter](plugins/filter/filter.md)](plugins/[filter/filter](plugins/filter/filter/md).md) |

[Back to TOC](#table-of-contents)
# Fluentd

Fluentd is the Schema for the fluentds API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FluentdSpec |
| status |  | FluentdStatus |

[Back to TOC](#table-of-contents)
# FluentdConfig

FluentdConfig is the Schema for the fluentdconfigs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | FluentdConfigSpec |
| status |  | FluentdConfigStatus |

[Back to TOC](#table-of-contents)
# FluentdConfigList

FluentdConfigList contains a list of FluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []FluentdConfig |

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
| filterSelector | Select namespaced filter plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| outputSelector | Select namespaced output plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| clusterFilterSelector | Select cluster filter plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| clusterOutputSelector | Select cluster output plugins | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |

[Back to TOC](#table-of-contents)
# FluentdConfigStatus

FluentdConfigStatus defines the observed state of FluentdConfig


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| messages | Messages defines the plugin errors which is selected by this fluentdconfig | string |
| state | The state of this fluentd config | StatusState |

[Back to TOC](#table-of-contents)
# FluentdList

FluentdList contains a list of Fluentd


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Fluentd |

[Back to TOC](#table-of-contents)
# FluentdSpec

FluentdSpec defines the desired state of Fluentd


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| globalInputs | Fluentd global inputs. | [][[input.Input](plugins/input/input.md)](plugins/[input/input](plugins/input/input/md).md) |
| disableService | By default will build the related service according to the globalinputs definition. | bool |
| replicas | Numbers of the Fluentd instance | *int32 |
| workers | Numbers of the workers in Fluentd instance | *int32 |
| image | Fluentd image. | string |
| args | Fluentd Watcher command line arguments. | []string |
| fluentdCfgSelector | FluentdCfgSelector defines the selectors to select the fluentd config CRs. | [metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) |
| buffer | Buffer definition | *BufferVolume |
| imagePullPolicy | Fluentd image pull policy. | corev1.PullPolicy |
| imagePullSecrets | Fluentd image pull secret | []corev1.LocalObjectReference |
| resources | Compute Resources required by container. | corev1.ResourceRequirements |
| nodeSelector | NodeSelector | map[string]string |
| serviceAccountAnnotations | Annotations to add to the Fluentd service account | map[string]string |
| affinity | Pod's scheduling constraints. | *corev1.Affinity |
| tolerations | Tolerations | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#toleration-v1-core) |
| runtimeClassName | RuntimeClassName represents the container runtime configuration. | string |
| priorityClassName | PriorityClassName represents the pod's priority class. | string |
| rbacRules | RBACRules represents additional rbac rules which will be applied to the fluentd clusterrole. | []rbacv1.PolicyRule |
| volumes | List of volumes that can be mounted by containers belonging to the pod. | []corev1.Volume |
| volumeMounts | Pod volumes to mount into the container's filesystem. Cannot be updated. | []corev1.VolumeMount |
| volumeClaimTemplates | volumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. | []corev1.PersistentVolumeClaim |

[Back to TOC](#table-of-contents)
# FluentdStatus

FluentdStatus defines the observed state of Fluentd


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| messages | Messages defines the plugin errors which is selected by this fluentdconfig | string |
| state | The state of this fluentd | StatusState |

[Back to TOC](#table-of-contents)
# Output

Output is the Schema for the outputs API


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) |
| spec |  | OutputSpec |
| status |  | OutputStatus |

[Back to TOC](#table-of-contents)
# OutputList

OutputList contains a list of Output


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) |
| items |  | []Output |

[Back to TOC](#table-of-contents)
# OutputSpec

OutputSpec defines the desired state of Output


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| outputs |  | [][[output.Output](plugins/output/output.md)](plugins/[output/output](plugins/output/output/md).md) |

[Back to TOC](#table-of-contents)
