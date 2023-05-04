# OpenTelemetry

The OpenTelemetry plugin allows you to take logs, metrics, and traces from Fluent Bit and submit them to an OpenTelemetry HTTP endpoint. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/opentelemetry**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target HTTP Server, default `127.0.0.1` | string |
| port | TCP port of the target OpenSearch instance, default `80` | *int32 |
| httpUser | Optional username credential for access | *[plugins.Secret](../secret.md) |
| httpPassword | Password for user defined in HTTP_User | *[plugins.Secret](../secret.md) |
| proxy | Specify an HTTP Proxy. The expected format of this value is http://HOST:PORT. Note that HTTPS is not currently supported. It is recommended not to set this and to configure the HTTP proxy environment variables instead as they support both HTTP and HTTPS. | string |
| metricsUri | Specify an optional HTTP URI for the target web server listening for metrics, e.g: /v1/metrics | string |
| logsUri | Specify an optional HTTP URI for the target web server listening for logs, e.g: /v1/logs | string |
| tracesUri | Specify an optional HTTP URI for the target web server listening for traces, e.g: /v1/traces | string |
| header | Add a HTTP header key/value pair. Multiple headers can be set. | map[string]string |
| logResponsePayload | Log the response payload within the Fluent Bit log. | *bool |
| addLabel | This allows you to add custom labels to all metrics exposed through the OpenTelemetry exporter. You may have multiple of these fields. | map[string]string |
| tls |  | *[plugins.TLS](../tls.md) |
