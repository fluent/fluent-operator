# Loki

The loki output plugin, allows to ingest your records into a Loki service.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| url | Loki URL. | *string |
| httpUser | Set HTTP basic authentication user name. | *[plugins.Secret](../secret.md) |
| httpPassword | Password for user defined in HTTP_User Set HTTP basic authentication password | *[plugins.Secret](../secret.md) |
| tenantID | Tenant ID used by default to push logs to Loki. If omitted or empty it assumes Loki is running in single-tenant mode and no X-Scope-OrgID header is sent. | *[plugins.Secret](../secret.md) |
| labels | Stream labels for API request. It can be multiple comma separated of strings specifying  key=value pairs. In addition to fixed parameters, it also allows to add custom record keys (similar to label_keys property). | []string |
| labelKeys | Optional list of record keys that will be placed as stream labels. This configuration property is for records key only. | []string |
| removeKeys | Optional list of record keys that will be removed from stream labels. This configuration property is for records key only. | []string |
| lineFormat | Format to use when flattening the record to a log line. Valid values are json or key_value. If set to json,  the log line sent to Loki will be the Fluentd record dumped as JSON. If set to key_value, the log line will be each item in the record concatenated together (separated by a single space) in the format. | string |
| extractKubernetesLabels | If set to true, it will add all Kubernetes labels to the Stream labels. | *bool |
| dropSingleKey | If a record only has 1 key, then just set the log line to the value and discard the key. | *bool |
| includeThreadLabel | Whether or not to include the fluentd_thread label when multiple threads are used for flushing | *bool |
| insecure | Disable certificate validation | *bool |
| tlsCaCertFile | TlsCaCert defines the CA certificate file for TLS. | *string |
| tlsClientCertFile | TlsClientCert defines the client certificate file for TLS. | *string |
| tlsPrivateKeyFile | TlsPrivateKey defines the client private key file for TLS. | *string |
