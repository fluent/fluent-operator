# FluentbitMetrics

Fluent Bit exposes its own metrics to allow you to monitor the internals of your pipeline. 
The collected metrics can be processed similarly to those from the Prometheus Node Exporter input plugin. 
They can be sent to output plugins including Prometheus Exporter, Prometheus Remote Write or OpenTelemetry.
> Important note: Metrics collected with Node Exporter Metrics flow through a separate pipeline from logs and current filters do not operate on top of metrics.

| Field | Description | Scheme |
| ----- | ----------- |--------|
| tag | Tag name associated to all records comming from this plugin. | string |
| scrape_interval | The rate at which metrics are collected from the host operating system. | string |
| scrape_on_start | Scrape metrics upon start, useful to avoid waiting for 'scrape_interval' for the first round of metrics. | *bool  |
