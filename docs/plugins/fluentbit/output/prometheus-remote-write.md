# Prometheus Remote Write

The prometheus remote write plugin allows you to take metrics from Fluent Bit and submit them to a Prometheus server through the remote write mechanism.

> Note: The prometheus exporter only works with metric plugins, such as Node Exporter Metrics


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target HTTP Server.            | string |
| port | TCP port of the target HTTP Server. | *int32 |
| httpUser | Basic Auth Username. | *[plugins.Secret](../secret.md) |
| httpPassword | Basic Auth Password. Requires HTTP_user to be set. | *[plugins.Secret](../secret.md) |
| proxy | Specify an HTTP Proxy. The expected format of this value is http://HOST:PORT. Note that HTTPS is not currently supported. It is recommended not to set this and to configure the  instead as they support both HTTP and HTTPS. | sting |
| uri | Specify an optional HTTP URI for the target web server, e.g: /something. | string |
| header | Add a HTTP header key/value pair. Multiple headers can be set. | map[string]string |
| log_response_payload | Log the response payload within the Fluent Bit log. | string |
| add_label | This allows you to add custom labels to all metrics exposed through the prometheus exporter. You may have multiple of these fields. | map[string]string |
| Workers | Enables dedicated thread(s) for this output. | *int32 |
| tls | | *[plugins.TLS](../tls.md) |
