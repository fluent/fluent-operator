# PrometheusExporter

PrometheusExporter An output plugin to expose Prometheus Metrics. <br /> The prometheus exporter allows you to take metrics from Fluent Bit and expose them such that a Prometheus instance can scrape them. <br /> **Important Note: The prometheus exporter only works with metric plugins, such as Node Exporter Metrics** <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/outputs/prometheus-exporter**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| host | IP address or hostname of the target HTTP Server, default: 0.0.0.0 | string |
| port | This is the port Fluent Bit will bind to when hosting prometheus metrics. | *int32 |
| addLabels | This allows you to add custom labels to all metrics exposed through the prometheus exporter. You may have multiple of these fields | map[string]string |
