# Prometheus

Prometheus is the in_prometheus plugin from fluent-plugin-prometheus. It serves the metrics collected by the prometheus plugins over HTTP so that Prometheus can scrape them.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| bind | The bind address to listen to. Defaults to 127.0.0.1, so set it to 0.0.0.0 to scrape from outside the pod. | *string |
| port | The port to listen to. Defaults to 24231. Use 2021 to serve the metrics on the \"metrics\" port that the operator already exposes on the Fluentd pods. | *int64 |
| metricsPath | The HTTP path of the metrics endpoint. Defaults to /metrics. | *string |
| aggregatedMetricsPath | The HTTP path of the endpoint that aggregates the metrics of all workers. Defaults to /aggregated_metrics. | *string |
| contentEncoding | The encoding of the exposed metrics. | *string |
# PrometheusMonitor

PrometheusMonitor is the in_prometheus_monitor plugin from fluent-plugin-prometheus. It collects Fluentd's buffer and retry metrics, such as fluentd_status_buffer_queue_length.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| interval | The interval in seconds between metric updates. Defaults to 5. | *uint32 |
| labels | Additional labels added to the metrics. | map[string]string |
# PrometheusOutputMonitor

PrometheusOutputMonitor is the in_prometheus_output_monitor plugin from fluent-plugin-prometheus. It collects metrics about the output plugins, such as fluentd_output_status_retry_count.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| interval | The interval in seconds between metric updates. Defaults to 5. | *uint32 |
| gaugeAll | Whether to expose the counters as gauges. Defaults to true. | *bool |
| labels | Additional labels added to the metrics. | map[string]string |
