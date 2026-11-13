# PrometheusRemoteWrite

The Prometheus Remote Write input plugin receives metrics sent with the Prometheus remote write protocol, <br /> for example by Prometheus, Grafana Alloy or Fluent Bit's Prometheus Remote Write output. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/prometheus-remote-write**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag | Tag name associated to all metrics coming from this plugin. | string |
| listen | The address to listen on, defaults to 0.0.0.0. | string |
| port | The port for Fluent Bit to listen on, defaults to 8080. Also add it to the ports of the FluentBit resource so that the Fluent Bit Service exposes it. | *int32 |
| uri | The HTTP path that accepts remote write requests, e.g. /api/prom/push. If not set, requests are accepted on any path. | string |
| bufferMaxSize | Maximum size of a remote write request, defaults to 4M. | string |
| bufferChunkSize | Size of the chunks allocated to read a request, up to bufferMaxSize, defaults to 512K. | string |
| successfulResponseCode | HTTP status code returned for a successfully ingested request, defaults to 201. | *int32 |
| threaded | Run the input in its own thread instead of the main event loop. | *bool |
| tls | TLS settings of the listener. Set crtFile and keyFile to accept remote write requests over HTTPS. | *[plugins.TLS](../tls.md) |
