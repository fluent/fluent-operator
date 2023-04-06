# Stackdriver

Stackdriver is the Stackdriver output plugin, allows you to ingest your records into GCP Stackdriver. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/stackdriver**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| googleServiceCredentials | Path to GCP Credentials JSON file | string |
| serviceAccountEmail | Email associated with the service | *[plugins.Secret](../secret.md) |
| serviceAccountSecret | Private Key associated with the service | *[plugins.Secret](../secret.md) |
| metadataServer | Metadata Server Prefix | string |
| location | GCP/AWS region to store data. Required if Resource is generic_node or generic_task | string |
| namespace | Namespace identifier. Required if Resource is generic_node or generic_task | string |
| nodeID | Node identifier within the namespace. Required if Resource is generic_node or generic_task | string |
| job | Identifier for a grouping of tasks. Required if Resource is generic_task | string |
| taskID | Identifier for a task within a namespace. Required if Resource is generic_task | string |
| exportToProjectID | The GCP Project that should receive the logs | string |
| resource | Set resource types of data | string |
| k8sClusterName | Name of the cluster that the pod is running in. Required if Resource is k8s_container, k8s_node, or k8s_pod | string |
| k8sClusterLocation | Location of the cluster that contains the pods/nodes. Required if Resource is k8s_container, k8s_node, or k8s_pod | string |
| labelsKey | Used by Stackdriver to find related labels and extract them to LogEntry Labels | string |
| labels | Optional list of comma separated of strings for key/value pairs | []string |
| logNameKey | The value of this field is set as the logName field in Stackdriver | string |
| tagPrefix | Used to validate the tags of logs that when the Resource is k8s_container, k8s_node, or k8s_pod | string |
| severityKey | Specify the key that contains the severity information for the logs | string |
| autoformatStackdriverTrace | Rewrite the trace field to be formatted for use with GCP Cloud Trace | *bool |
| workers | Number of dedicated threads for the Stackdriver Output Plugin | *int32 |
| customK8sRegex | A custom regex to extract fields from the local_resource_id of the logs | string |
| resourceLabels | Optional list of comma seperated strings. Setting these fields overrides the Stackdriver monitored resource API values | []string |
