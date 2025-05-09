# Loki

The loki output plugin, allows to ingest your records into a Loki service. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/loki**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | Loki hostname or IP address. | string |
| port | Loki TCP port | *int32 |
| uri | Specify a custom HTTP URI. It must start with forward slash. | string |
| httpUser | Set HTTP basic authentication user name. | *[plugins.Secret](../secret.md) |
| httpPassword | Password for user defined in HTTP_User Set HTTP basic authentication password | *[plugins.Secret](../secret.md) |
| bearerToken | Set bearer token authentication token value. Can be used as alterntative to HTTP basic authentication | *[plugins.Secret](../secret.md) |
| tenantID | Tenant ID used by default to push logs to Loki. If omitted or empty it assumes Loki is running in single-tenant mode and no X-Scope-OrgID header is sent. | *[plugins.Secret](../secret.md) |
| labels | Stream labels for API request. It can be multiple comma separated of strings specifying  key=value pairs. In addition to fixed parameters, it also allows to add custom record keys (similar to label_keys property). | []string |
| labelKeys | Optional list of record keys that will be placed as stream labels. This configuration property is for records key only. | []string |
| labelMapPath | Specify the label map file path. The file defines how to extract labels from each record. | string |
| removeKeys | Optional list of keys to remove. | []string |
| dropSingleKey | If set to true and after extracting labels only a single key remains, the log line sent to Loki will be the value of that key in line_format. | string |
| lineFormat | Format to use when flattening the record to a log line. Valid values are json or key_value. If set to json,  the log line sent to Loki will be the Fluent Bit record dumped as JSON. If set to key_value, the log line will be each item in the record concatenated together (separated by a single space) in the format. | string |
| autoKubernetesLabels | If set to true, it will add all Kubernetes labels to the Stream labels. | string |
| tenantIDKey | Specify the name of the key from the original record that contains the Tenant ID. The value of the key is set as X-Scope-OrgID of HTTP header. It is useful to set Tenant ID dynamically. | string |
| structuredMetadata | Stream structured metadata for API request. It can be multiple comma separated key=value pairs. This is used for high cardinality data that isn't suited for using labels. Only supported in Loki 3.0+ with schema v13 and TSDB storage. | map[string]string |
| structuredMetadataKeys | Optional list of record keys that will be placed as structured metadata. This allows using record accessor patterns (e.g. $kubernetes['pod_name']) to reference record keys. | []string |
| tls |  | *[plugins.TLS](../tls.md) |
| networking | Include fluentbit networking options for this output-plugin | *plugins.Networking |
| totalLimitSize | Limit the maximum number of Chunks in the filesystem for the current output logical destination. | string |
| workers | Enables dedicated thread(s) for this output. Default value is set since version 1.8.13. For previous versions is 0. | *int32 |
