# PrometheusRemoteWrite

An output plugin to submit Prometheus Metrics using the remote write protocol. <br /> The prometheus remote write plugin allows you to take metrics from Fluent Bit and submit them to a Prometheus server through the remote write mechanism. <br /> **Important Note: The prometheus exporter only works with metric plugins, such as Node Exporter Metrics** <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/prometheus-remote-write**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target HTTP Server, default: 127.0.0.1 | string |
| httpUser | Basic Auth Username | *[plugins.Secret](../secret.md) |
| httpPasswd | Basic Auth Password. Requires HTTP_user to be se | *[plugins.Secret](../secret.md) |
| port | TCP port of the target HTTP Serveri, default:80 | *int32 |
| proxy | Specify an HTTP Proxy. The expected format of this value is http://HOST:PORT. | string |
| uri | Specify an optional HTTP URI for the target web server, e.g: /something ,default: / | string |
| headers | Add a HTTP header key/value pair. Multiple headers can be set. | map[string]string |
| logResponsePayload | Log the response payload within the Fluent Bit log,default: false | *bool |
| addLabels | This allows you to add custom labels to all metrics exposed through the prometheus exporter. You may have multiple of these fields | map[string]string |
| workers | Enables dedicated thread(s) for this output. Default value is set since version 1.8.13. For previous versions is 0,default : 2 | *int32 |
| tls |  | *[plugins.TLS](../tls.md) |
