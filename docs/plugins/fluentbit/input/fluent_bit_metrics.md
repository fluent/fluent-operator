# FluentbitMetrics

Fluent Bit exposes its own metrics to allow you to monitor the internals of your pipeline. <br /> The collected metrics can be processed similarly to those from the Prometheus Node Exporter input plugin. <br /> They can be sent to output plugins including Prometheus Exporter, Prometheus Remote Write or OpenTelemetry. <br /> **Important note: Metrics collected with Node Exporter Metrics flow through a separate pipeline from logs and current filters do not operate on top of metrics.** <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/inputs/fluentbit-metrics**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| tag |  | string |
| scrapeInterval | The rate at which metrics are collected from the host operating system. default is 2 seconds. | string |
| scrapeOnStart | Scrape metrics upon start, useful to avoid waiting for 'scrape_interval' for the first round of metrics. | *bool |
