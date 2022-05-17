# Kubernetes

Kubernetes filter allows to enrich your log files with Kubernetes metadata.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| bufferSize | Set the buffer size for HTTP client when reading responses from Kubernetes API server. | string |
| kubeURL | API Server end-point | string |
| kubeCAFile | CA certificate file | string |
| kubeCAPath | Absolute path to scan for certificate files | string |
| kubeTokenFile | Token file | string |
| kubeTagPrefix | When the source records comes from Tail input plugin, this option allows to specify what's the prefix used in Tail configuration. | string |
| mergeLog | When enabled, it checks if the log field content is a JSON string map, if so, it append the map fields as part of the log structure. | *bool |
| mergeLogKey | When Merge_Log is enabled, the filter tries to assume the log field from the incoming message is a JSON string message and make a structured representation of it at the same level of the log field in the map. Now if Merge_Log_Key is set (a string name), all the new structured fields taken from the original log content are inserted under the new key. | string |
| mergeLogTrim | When Merge_Log is enabled, trim (remove possible \n or \r) field values. | *bool |
| mergeParser | Optional parser name to specify how to parse the data contained in the log key. Recommended use is for developers or testing only. | string |
| keepLog | When Keep_Log is disabled, the log field is removed from the incoming message once it has been successfully merged (Merge_Log must be enabled as well). | *bool |
| tlsDebug | Debug level between 0 (nothing) and 4 (every detail). | *int32 |
| tlsVerify | When enabled, turns on certificate validation when connecting to the Kubernetes API server. | *bool |
| useJournal | When enabled, the filter reads logs coming in Journald format. | *bool |
| regexParser | Set an alternative Parser to process record Tag and extract pod_name, namespace_name, container_name and docker_id. The parser must be registered in a parsers file (refer to parser filter-kube-test as an example). | string |
| k8sLoggingParser | Allow Kubernetes Pods to suggest a pre-defined Parser (read more about it in Kubernetes Annotations section) | *bool |
| k8sLoggingExclude | Allow Kubernetes Pods to exclude their logs from the log processor (read more about it in Kubernetes Annotations section). | *bool |
| labels | Include Kubernetes resource labels in the extra metadata. | *bool |
| annotations | Include Kubernetes resource annotations in the extra metadata. | *bool |
| kubeMetaPreloadCacheDir | If set, Kubernetes meta-data can be cached/pre-loaded from files in JSON format in this directory, named as namespace-pod.meta | string |
| dummyMeta | If set, use dummy-meta data (for test/dev purposes) | *bool |
