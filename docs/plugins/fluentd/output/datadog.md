# Datadog

Datadog defines the parameters for out_datadog plugin


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| apiKey | This parameter is required in order to authenticate your fluent agent. | *[plugins.Secret](../secret.md) |
| useJson | Event format, if true, the event is sent in json format. Othwerwise, in plain text. | *bool |
| includeTagKey | Automatically include the Fluentd tag in the record. | *bool |
| tagKey | Where to store the Fluentd tag. | *string |
| timestampKey | Name of the attribute which will contain timestamp of the log event. If nil, timestamp attribute is not added. | *string |
| useSSL | If true, the agent initializes a secure connection to Datadog. In clear TCP otherwise. | *bool |
| noSSLValidation | Disable SSL validation (useful for proxy forwarding) | *bool |
| sslPort | Port used to send logs over a SSL encrypted connection to Datadog. If use_http is disabled, use 10516 for the US region and 443 for the EU region. | *uint32 |
| maxRetries | The number of retries before the output plugin stops. Set to -1 for unlimited retries | *uint32 |
| maxBackoff | The maximum time waited between each retry in seconds | *uint32 |
| useHTTP | Enable HTTP forwarding. If you disable it, make sure to change the port to 10514 or ssl_port to 10516 | *bool |
| useCompression | Enable log compression for HTTP | *bool |
| compressionLevel | Set the log compression level for HTTP (1 to 9, 9 being the best ratio) | *uint32 |
| ddSource | This tells Datadog what integration it is | *string |
| ddSourcecategory | Multiple value attribute. Can be used to refine the source attribute | *string |
| ddTags | Custom tags with the following format \"key1:value1, key2:value2\" | *string |
| ddHostname | Used by Datadog to identify the host submitting the logs. | *string |
| service | Used by Datadog to correlate between logs, traces and metrics. | *string |
| port | Proxy port when logs are not directly forwarded to Datadog and ssl is not used | *uint32 |
| host | Proxy endpoint when logs are not directly forwarded to Datadog | *string |
| httpProxy | HTTP proxy, only takes effect if HTTP forwarding is enabled (use_http). Defaults to HTTP_PROXY/http_proxy env vars. | *string |
