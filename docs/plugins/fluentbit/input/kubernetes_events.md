# KubernetesEvents

The KubernetesEvents input plugin allows you to collect kubernetes cluster events from kube-api server **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/kubernetes-events*


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Tag name associated to all records coming from this plugin. | string |
| db | Set a database file to keep track of recorded Kubernetes events | string |
| dbSync | Set a database sync method. values: extra, full, normal and off | string |
| intervalSec | Set the polling interval for each channel. | *int32 |
| intervalNsec | Set the polling interval for each channel (sub seconds: nanoseconds). | *int64 |
| kubeURL | API Server end-point | string |
| kubeCAFile | CA certificate file | string |
| kubeCAPath | Absolute path to scan for certificate files | string |
| kubeTokenFile | Token file | string |
| kubeTokenTTL | configurable 'time to live' for the K8s token. By default, it is set to 600 seconds. After this time, the token is reloaded from Kube_Token_File or the Kube_Token_Command. | string |
| kubeRequestLimit | kubernetes limit parameter for events query, no limit applied when set to 0. | *int32 |
| kubeRetentionTime | Kubernetes retention time for events. | string |
| kubeNamespace | Kubernetes namespace to query events from. Gets events from all namespaces by default | string |
| tlsDebug | Debug level between 0 (nothing) and 4 (every detail). | *int32 |
| tlsVerify | When enabled, turns on certificate validation when connecting to the Kubernetes API server. | *bool |
| tlsVhost | Set optional TLS virtual host. | string |
| storageType | Set storage type for buffering can be filesystem or memory. | string |
